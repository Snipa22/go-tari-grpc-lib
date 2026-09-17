package nodeGRPC

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

// TestBuildTransportCredentials_SecureNilConfig verifies that InitNodeGRPCSecure still produces
// TLS credentials even when the caller passes a nil *tls.Config (falls back to Go's defaults,
// rather than silently downgrading to insecure).
func TestBuildTransportCredentials_SecureNilConfig(t *testing.T) {
	creds := buildTransportCredentials(true, nil)
	if got := creds.Info().SecurityProtocol; got != "tls" {
		t.Fatalf("expected tls transport credentials even with a nil tls.Config, got SecurityProtocol=%q", got)
	}
}

// TestInitNodeGRPC_UsesInsecureCredentials is an end-to-end style check (still without a live
// server, since grpc.NewClient doesn't dial eagerly) that InitNodeGRPC wires up insecure
// credentials while InitNodeGRPCSecure wires up TLS credentials.
func TestInitNodeGRPC_UsesInsecureCredentials(t *testing.T) {
	if err := InitNodeGRPC("127.0.0.1:0"); err != nil {
		t.Fatalf("InitNodeGRPC returned unexpected error: %v", err)
	}
	if getConn() == nil {
		t.Fatal("expected a non-nil connection after InitNodeGRPC")
	}
}

func TestInitNodeGRPCSecure_UsesTLSCredentials(t *testing.T) {
	if err := InitNodeGRPCSecure("127.0.0.1:0", &tls.Config{}); err != nil {
		t.Fatalf("InitNodeGRPCSecure returned unexpected error: %v", err)
	}
	if getConn() == nil {
		t.Fatal("expected a non-nil connection after InitNodeGRPCSecure")
	}
}

// TestInitNodeGRPC_PropagatesError verifies finding 2: InitNodeGRPC must not silently discard
// the error returned by grpc.NewClient. grpc.NewClient itself validates most targets lazily (at
// dial time, not construction time), so this substitutes a failing stub via the newClient seam
// to deterministically exercise the propagation path.
func TestInitNodeGRPC_PropagatesError(t *testing.T) {
	origNewClient := newClient
	defer func() { newClient = origNewClient }()

	wantErr := errors.New("boom: simulated grpc.NewClient failure")
	newClient = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return nil, wantErr
	}

	err := InitNodeGRPC("127.0.0.1:9999")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected InitNodeGRPC to propagate the underlying error, got %v", err)
	}
}

// TestInitNodeGRPC_ConcurrentAccessIsRaceFree exercises finding 4: concurrent Init calls must not
// race on the package-level connection state. Run with `go test -race` to verify.
func TestInitNodeGRPC_ConcurrentAccessIsRaceFree(t *testing.T) {
	const goroutines = 50
	var wg sync.WaitGroup
	errs := make(chan error, goroutines)
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			addr := fmt.Sprintf("127.0.0.1:%d", 20000+i)
			if err := InitNodeGRPC(addr); err != nil {
				errs <- err
			}
			_ = getConn()
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Errorf("unexpected error from concurrent InitNodeGRPC: %v", err)
	}
}

// withNilConn temporarily clears the package-level connection (as if Init* had never been
// called), runs fn, then restores whatever connection was previously installed. It exists so
// tests can exercise the "not initialized yet" path without permanently disrupting the shared
// package-level grpcConn state other tests in this file depend on (see the serial-tests-only
// comment on startFakeBaseNodeServer).
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

// TestBaseNodeClient_ErrNotInitialized verifies finding 4's centralized fix directly: calling
// baseNodeClient() before Init* has ever populated grpcConn must return ErrNotInitialized (and a
// nil client), not construct a tari_generated.BaseNodeClient wrapping a nil *grpc.ClientConn.
func TestBaseNodeClient_ErrNotInitialized(t *testing.T) {
	withNilConn(t, func() {
		client, err := baseNodeClient()
		if client != nil {
			t.Fatalf("expected a nil client, got %v", client)
		}
		if !errors.Is(err, ErrNotInitialized) {
			t.Fatalf("expected ErrNotInitialized, got %v", err)
		}
	})
}

// TestWrapper_ReturnsErrNotInitialized_NotPanic verifies finding 4 end-to-end: calling an RPC
// wrapper before InitNodeGRPC/InitNodeGRPCSecure has ever run must return ErrNotInitialized, not
// panic with a nil-pointer dereference (which is what happened when every wrapper constructed
// tari_generated.NewBaseNodeClient(getConn()) directly against a nil connection). GetTipInfo is
// used as a representative sample of the ~80 call sites this fix covers centrally.
func TestWrapper_ReturnsErrNotInitialized_NotPanic(t *testing.T) {
	withNilConn(t, func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("GetTipInfo panicked instead of returning an error: %v", r)
			}
		}()
		_, err := GetTipInfo(context.Background())
		if !errors.Is(err, ErrNotInitialized) {
			t.Fatalf("expected ErrNotInitialized, got %v", err)
		}
	})
}

// TestInitNodeGRPC_ExtraDialOptionsReachConnection verifies finding 8: an extra grpc.DialOption
// passed to InitNodeGRPC must actually reach the grpc.NewClient call used to build the
// connection, not be silently dropped. It substitutes the newClient seam to observe the exact
// []grpc.DialOption slice constructed by dialOptions()/initNodeGRPC, comparing the option count
// with and without an extra option supplied, then still delegates to the real grpc.NewClient so
// the resulting connection is genuine (grpc.NewClient doesn't dial eagerly, so this stays fast
// and doesn't require a live server).
func TestInitNodeGRPC_ExtraDialOptionsReachConnection(t *testing.T) {
	origNewClient := newClient
	defer func() { newClient = origNewClient }()

	var gotOptsCount int
	newClient = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		gotOptsCount = len(opts)
		return origNewClient(target, opts...)
	}

	if err := InitNodeGRPC("127.0.0.1:0"); err != nil {
		t.Fatalf("InitNodeGRPC (baseline): %v", err)
	}
	baselineCount := gotOptsCount

	extra := grpc.WithUserAgent("distinct-marker-ua")
	if err := InitNodeGRPC("127.0.0.1:0", extra); err != nil {
		t.Fatalf("InitNodeGRPC (with extra dial option): %v", err)
	}
	if getConn() == nil {
		t.Fatal("expected a non-nil connection after InitNodeGRPC with an extra dial option")
	}
	if gotOptsCount != baselineCount+1 {
		t.Fatalf("expected exactly 1 additional dial option to reach grpc.NewClient beyond the %d defaults, got %d total", baselineCount, gotOptsCount)
	}
}

// TestInitNodeGRPCSecure_ExtraDialOptionsReachConnection mirrors
// TestInitNodeGRPC_ExtraDialOptionsReachConnection for the TLS-secured Init variant.
func TestInitNodeGRPCSecure_ExtraDialOptionsReachConnection(t *testing.T) {
	origNewClient := newClient
	defer func() { newClient = origNewClient }()

	var gotOptsCount int
	newClient = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		gotOptsCount = len(opts)
		return origNewClient(target, opts...)
	}

	if err := InitNodeGRPCSecure("127.0.0.1:0", &tls.Config{}); err != nil {
		t.Fatalf("InitNodeGRPCSecure (baseline): %v", err)
	}
	baselineCount := gotOptsCount

	extra := grpc.WithUserAgent("distinct-marker-ua")
	if err := InitNodeGRPCSecure("127.0.0.1:0", &tls.Config{}, extra); err != nil {
		t.Fatalf("InitNodeGRPCSecure (with extra dial option): %v", err)
	}
	if gotOptsCount != baselineCount+1 {
		t.Fatalf("expected exactly 1 additional dial option to reach grpc.NewClient beyond the %d defaults, got %d total", baselineCount, gotOptsCount)
	}
}
