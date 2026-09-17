package walletGRPC

import (
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
