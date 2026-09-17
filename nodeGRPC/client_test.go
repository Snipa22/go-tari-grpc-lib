package nodeGRPC

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
