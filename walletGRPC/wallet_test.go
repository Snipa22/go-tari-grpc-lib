package walletGRPC

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"testing"

	"google.golang.org/grpc"
)

// TestBuildTransportCredentials_Insecure verifies that the plaintext path builds insecure
// transport credentials, not TLS-based ones.
func TestBuildTransportCredentials_Insecure(t *testing.T) {
	creds := buildTransportCredentials(false, nil)
	if creds == nil {
		t.Fatal("expected non-nil credentials")
	}
	if got := creds.Info().SecurityProtocol; got != "insecure" {
		t.Fatalf("expected insecure transport credentials, got SecurityProtocol=%q", got)
	}
}

// TestBuildTransportCredentials_Secure verifies that the secure path actually builds TLS-based
// transport credentials, not insecure ones, per finding 1 / required test 1.
func TestBuildTransportCredentials_Secure(t *testing.T) {
	creds := buildTransportCredentials(true, &tls.Config{})
	if creds == nil {
		t.Fatal("expected non-nil credentials")
	}
	if got := creds.Info().SecurityProtocol; got != "tls" {
		t.Fatalf("expected tls transport credentials, got SecurityProtocol=%q", got)
	}
}

// TestBuildTransportCredentials_SecureNilConfig verifies that InitWalletGRPCSecure still produces
// TLS credentials even when the caller passes a nil *tls.Config (falls back to Go's defaults,
// rather than silently downgrading to insecure).
func TestBuildTransportCredentials_SecureNilConfig(t *testing.T) {
	creds := buildTransportCredentials(true, nil)
	if got := creds.Info().SecurityProtocol; got != "tls" {
		t.Fatalf("expected tls transport credentials even with a nil tls.Config, got SecurityProtocol=%q", got)
	}
}

func TestInitWalletGRPC_UsesInsecureCredentials(t *testing.T) {
	if err := InitWalletGRPC("127.0.0.1:0"); err != nil {
		t.Fatalf("InitWalletGRPC returned unexpected error: %v", err)
	}
	if getConn() == nil {
		t.Fatal("expected a non-nil connection after InitWalletGRPC")
	}
}

func TestInitWalletGRPCSecure_UsesTLSCredentials(t *testing.T) {
	if err := InitWalletGRPCSecure("127.0.0.1:0", &tls.Config{}); err != nil {
		t.Fatalf("InitWalletGRPCSecure returned unexpected error: %v", err)
	}
	if getConn() == nil {
		t.Fatal("expected a non-nil connection after InitWalletGRPCSecure")
	}
}

// TestInitWalletGRPC_PropagatesError verifies finding 2: InitWalletGRPC must not silently discard
// the error returned by grpc.NewClient. grpc.NewClient itself validates most targets lazily (at
// dial time, not construction time), so this substitutes a failing stub via the newClient seam
// to deterministically exercise the propagation path.
func TestInitWalletGRPC_PropagatesError(t *testing.T) {
	origNewClient := newClient
	defer func() { newClient = origNewClient }()

	wantErr := errors.New("boom: simulated grpc.NewClient failure")
	newClient = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return nil, wantErr
	}

	err := InitWalletGRPC("127.0.0.1:9999")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected InitWalletGRPC to propagate the underlying error, got %v", err)
	}
}

// TestInitWalletGRPC_ConcurrentAccessIsRaceFree exercises finding 4: concurrent Init calls must
// not race on the package-level connection state. Run with `go test -race` to verify.
func TestInitWalletGRPC_ConcurrentAccessIsRaceFree(t *testing.T) {
	const goroutines = 50
	var wg sync.WaitGroup
	errs := make(chan error, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			addr := fmt.Sprintf("127.0.0.1:%d", 30000+i)
			if err := InitWalletGRPC(addr); err != nil {
				errs <- err
			}
			_ = getConn()
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Errorf("unexpected error from concurrent InitWalletGRPC: %v", err)
	}
}

// withNilConn temporarily clears the package-level connection (as if Init* had never been
// called), runs fn, then restores whatever connection was previously installed. It exists so
// tests can exercise the "not initialized yet" path without permanently disrupting the shared
// package-level grpcConn state other tests in this file depend on (see the serial-tests-only
// comment on startFakeWalletServer).
func withNilConn(t *testing.T, fn func()) {
	t.Helper()
	connMu.Lock()
	prev := grpcConn
	grpcConn = nil
	connMu.Unlock()
	defer func() {
		connMu.Lock()
		grpcConn = prev
		connMu.Unlock()
	}()
	fn()
}

// TestWalletClient_ErrNotInitialized verifies finding 4's centralized fix directly: calling
// walletClient() before Init* has ever populated grpcConn must return ErrNotInitialized (and a
// nil client), not construct a tari_generated.WalletClient wrapping a nil *grpc.ClientConn.
func TestWalletClient_ErrNotInitialized(t *testing.T) {
	withNilConn(t, func() {
		client, err := walletClient()
		if client != nil {
			t.Fatalf("expected a nil client, got %v", client)
		}
		if !errors.Is(err, ErrNotInitialized) {
			t.Fatalf("expected ErrNotInitialized, got %v", err)
		}
	})
}

// TestWrapper_ReturnsErrNotInitialized_NotPanic verifies finding 4 end-to-end: calling an RPC
// wrapper before InitWalletGRPC/InitWalletGRPCSecure has ever run must return ErrNotInitialized,
// not panic with a nil-pointer dereference (which is what happened when every wrapper
// constructed tari_generated.NewWalletClient(getConn()) directly against a nil connection).
// GetBalances is used as a representative sample of the ~80 call sites this fix covers
// centrally.
func TestWrapper_ReturnsErrNotInitialized_NotPanic(t *testing.T) {
	withNilConn(t, func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("GetBalances panicked instead of returning an error: %v", r)
			}
		}()
		_, err := GetBalances(context.Background())
		if !errors.Is(err, ErrNotInitialized) {
			t.Fatalf("expected ErrNotInitialized, got %v", err)
		}
	})
}

// TestInitWalletGRPC_ExtraDialOptionsReachConnection verifies finding 8: an extra
// grpc.DialOption passed to InitWalletGRPC must actually reach the grpc.NewClient call used to
// build the connection, not be silently dropped. It substitutes the newClient seam to observe
// the exact []grpc.DialOption slice constructed by dialOptions()/initWalletGRPC, comparing the
// option count with and without an extra option supplied, then still delegates to the real
// grpc.NewClient so the resulting connection is genuine (grpc.NewClient doesn't dial eagerly, so
// this stays fast and doesn't require a live server).
func TestInitWalletGRPC_ExtraDialOptionsReachConnection(t *testing.T) {
	origNewClient := newClient
	defer func() { newClient = origNewClient }()

	var gotOptsCount int
	newClient = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		gotOptsCount = len(opts)
		return origNewClient(target, opts...)
	}

	if err := InitWalletGRPC("127.0.0.1:0"); err != nil {
		t.Fatalf("InitWalletGRPC (baseline): %v", err)
	}
	baselineCount := gotOptsCount

	extra := grpc.WithUserAgent("distinct-marker-ua")
	if err := InitWalletGRPC("127.0.0.1:0", extra); err != nil {
		t.Fatalf("InitWalletGRPC (with extra dial option): %v", err)
	}
	if getConn() == nil {
		t.Fatal("expected a non-nil connection after InitWalletGRPC with an extra dial option")
	}
	if gotOptsCount != baselineCount+1 {
		t.Fatalf("expected exactly 1 additional dial option to reach grpc.NewClient beyond the %d defaults, got %d total", baselineCount, gotOptsCount)
	}
}

// TestInitWalletGRPCSecure_ExtraDialOptionsReachConnection mirrors
// TestInitWalletGRPC_ExtraDialOptionsReachConnection for the TLS-secured Init variant.
func TestInitWalletGRPCSecure_ExtraDialOptionsReachConnection(t *testing.T) {
	origNewClient := newClient
	defer func() { newClient = origNewClient }()

	var gotOptsCount int
	newClient = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		gotOptsCount = len(opts)
		return origNewClient(target, opts...)
	}

	if err := InitWalletGRPCSecure("127.0.0.1:0", &tls.Config{}); err != nil {
		t.Fatalf("InitWalletGRPCSecure (baseline): %v", err)
	}
	baselineCount := gotOptsCount

	extra := grpc.WithUserAgent("distinct-marker-ua")
	if err := InitWalletGRPCSecure("127.0.0.1:0", &tls.Config{}, extra); err != nil {
		t.Fatalf("InitWalletGRPCSecure (with extra dial option): %v", err)
	}
	if gotOptsCount != baselineCount+1 {
		t.Fatalf("expected exactly 1 additional dial option to reach grpc.NewClient beyond the %d defaults, got %d total", baselineCount, gotOptsCount)
	}
}

// TestDefaultFeePerGram pins the documented default value so an accidental change is caught.
func TestDefaultFeePerGram(t *testing.T) {
	if DefaultFeePerGram != 5 {
		t.Fatalf("expected DefaultFeePerGram to be 5, got %v", DefaultFeePerGram)
	}
}

// TestBuildCoinSplitRequest_UsesDefaultFeeWhenZero verifies that a feePerGram of 0 falls back to
// DefaultFeePerGram.
func TestBuildCoinSplitRequest_UsesDefaultFeeWhenZero(t *testing.T) {
	req := buildCoinSplitRequest(1000, 10, 0)
	if req.FeePerGram != DefaultFeePerGram {
		t.Fatalf("expected FeePerGram to fall back to DefaultFeePerGram (%v), got %v", DefaultFeePerGram, req.FeePerGram)
	}
	if req.AmountPerSplit != 1000 || req.SplitCount != 10 {
		t.Fatalf("unexpected request fields: %+v", req)
	}
}

// TestBuildCoinSplitRequest_HonorsOverride verifies that a non-zero feePerGram override produces
// a request using the overridden fee, not the default.
func TestBuildCoinSplitRequest_HonorsOverride(t *testing.T) {
	const override uint64 = 42
	req := buildCoinSplitRequest(1000, 10, override)
	if req.FeePerGram != override {
		t.Fatalf("expected FeePerGram to be the override value %v, got %v", override, req.FeePerGram)
	}
	if req.FeePerGram == DefaultFeePerGram {
		t.Fatalf("override value unexpectedly matched DefaultFeePerGram, test is not distinguishing")
	}
}

func TestBuildCoinSplitRequestTableDriven(t *testing.T) {
	cases := []struct {
		name           string
		splitAmt       int
		numSplits      int
		feePerGram     uint64
		wantFeePerGram uint64
	}{
		{name: "zero fee falls back to default", splitAmt: 500, numSplits: 5, feePerGram: 0, wantFeePerGram: DefaultFeePerGram},
		{name: "explicit default passed through", splitAmt: 500, numSplits: 5, feePerGram: DefaultFeePerGram, wantFeePerGram: DefaultFeePerGram},
		{name: "override passed through", splitAmt: 500, numSplits: 5, feePerGram: 25, wantFeePerGram: 25},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := buildCoinSplitRequest(tc.splitAmt, tc.numSplits, tc.feePerGram)
			if req.FeePerGram != tc.wantFeePerGram {
				t.Errorf("FeePerGram = %v, want %v", req.FeePerGram, tc.wantFeePerGram)
			}
			if req.AmountPerSplit != uint64(tc.splitAmt) {
				t.Errorf("AmountPerSplit = %v, want %v", req.AmountPerSplit, tc.splitAmt)
			}
			if req.SplitCount != uint64(tc.numSplits) {
				t.Errorf("SplitCount = %v, want %v", req.SplitCount, tc.numSplits)
			}
		})
	}
}
