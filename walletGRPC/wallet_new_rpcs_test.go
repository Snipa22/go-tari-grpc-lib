package walletGRPC

import (
	"context"
	"testing"

	"github.com/Snipa22/go-tari-grpc-lib/v3/tari_generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// TestGetVersion_Success verifies that GetVersion forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetVersion_Success(t *testing.T) {
	wantReq := &tari_generated.GetVersionRequest{}
	wantResp := &tari_generated.GetVersionResponse{}
	var gotReq *tari_generated.GetVersionRequest
	srv := &fakeWalletServer{
		getVersionFn: func(ctx context.Context, req *tari_generated.GetVersionRequest) (*tari_generated.GetVersionResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := GetVersion(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetVersion returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetVersion response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetVersion_Error verifies that GetVersion propagates the underlying GRPC error rather than swallowing it.
func TestGetVersion_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getversion")
	srv := &fakeWalletServer{
		getVersionFn: func(ctx context.Context, req *tari_generated.GetVersionRequest) (*tari_generated.GetVersionResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := GetVersion(context.Background(), &tari_generated.GetVersionRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getversion" {
		t.Fatalf("GetVersion did not propagate the underlying error, got: %v", err)
	}
}

// TestCheckForUpdates_Success verifies that CheckForUpdates returns the response the server sends, and that
// the server actually received a request (arg marshaling for the no-argument/Empty case).
func TestCheckForUpdates_Success(t *testing.T) {
	wantResp := &tari_generated.SoftwareUpdate{}
	var gotReq *tari_generated.Empty
	srv := &fakeWalletServer{
		checkForUpdatesFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SoftwareUpdate, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := CheckForUpdates(context.Background())
	if err != nil {
		t.Fatalf("CheckForUpdates returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("CheckForUpdates response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestCheckForUpdates_Error verifies that CheckForUpdates propagates the underlying GRPC error rather than swallowing it.
func TestCheckForUpdates_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-checkforupdates")
	srv := &fakeWalletServer{
		checkForUpdatesFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SoftwareUpdate, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := CheckForUpdates(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-checkforupdates" {
		t.Fatalf("CheckForUpdates did not propagate the underlying error, got: %v", err)
	}
}

// TestIdentify_Success verifies that Identify forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestIdentify_Success(t *testing.T) {
	wantReq := &tari_generated.GetIdentityRequest{}
	wantResp := &tari_generated.GetIdentityResponse{}
	var gotReq *tari_generated.GetIdentityRequest
	srv := &fakeWalletServer{
		identifyFn: func(ctx context.Context, req *tari_generated.GetIdentityRequest) (*tari_generated.GetIdentityResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := Identify(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("Identify returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("Identify response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestIdentify_Error verifies that Identify propagates the underlying GRPC error rather than swallowing it.
func TestIdentify_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-identify")
	srv := &fakeWalletServer{
		identifyFn: func(ctx context.Context, req *tari_generated.GetIdentityRequest) (*tari_generated.GetIdentityResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := Identify(context.Background(), &tari_generated.GetIdentityRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-identify" {
		t.Fatalf("Identify did not propagate the underlying error, got: %v", err)
	}
}

// TestGetAddress_Success verifies that GetAddress returns the response the server sends, and that
// the server actually received a request (arg marshaling for the no-argument/Empty case).
func TestGetAddress_Success(t *testing.T) {
	wantResp := &tari_generated.GetAddressResponse{}
	var gotReq *tari_generated.Empty
	srv := &fakeWalletServer{
		getAddressFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.GetAddressResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := GetAddress(context.Background())
	if err != nil {
		t.Fatalf("GetAddress returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetAddress response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetAddress_Error verifies that GetAddress propagates the underlying GRPC error rather than swallowing it.
func TestGetAddress_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getaddress")
	srv := &fakeWalletServer{
		getAddressFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.GetAddressResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := GetAddress(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getaddress" {
		t.Fatalf("GetAddress did not propagate the underlying error, got: %v", err)
	}
}

// TestGetPaymentIdAddress_Success verifies that GetPaymentIdAddress forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetPaymentIdAddress_Success(t *testing.T) {
	wantReq := &tari_generated.GetPaymentIdAddressRequest{}
	wantResp := &tari_generated.GetCompleteAddressResponse{}
	var gotReq *tari_generated.GetPaymentIdAddressRequest
	srv := &fakeWalletServer{
		getPaymentIdAddressFn: func(ctx context.Context, req *tari_generated.GetPaymentIdAddressRequest) (*tari_generated.GetCompleteAddressResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := GetPaymentIdAddress(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetPaymentIdAddress returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetPaymentIdAddress response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetPaymentIdAddress_Error verifies that GetPaymentIdAddress propagates the underlying GRPC error rather than swallowing it.
func TestGetPaymentIdAddress_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getpaymentidaddress")
	srv := &fakeWalletServer{
		getPaymentIdAddressFn: func(ctx context.Context, req *tari_generated.GetPaymentIdAddressRequest) (*tari_generated.GetCompleteAddressResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := GetPaymentIdAddress(context.Background(), &tari_generated.GetPaymentIdAddressRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getpaymentidaddress" {
		t.Fatalf("GetPaymentIdAddress did not propagate the underlying error, got: %v", err)
	}
}

// TestPrepareOneSidedTransactionForSigning_Success verifies that PrepareOneSidedTransactionForSigning forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestPrepareOneSidedTransactionForSigning_Success(t *testing.T) {
	wantReq := &tari_generated.PrepareOneSidedTransactionForSigningRequest{}
	wantResp := &tari_generated.PrepareOneSidedTransactionForSigningResponse{}
	var gotReq *tari_generated.PrepareOneSidedTransactionForSigningRequest
	srv := &fakeWalletServer{
		prepareOneSidedTransactionForSigningFn: func(ctx context.Context, req *tari_generated.PrepareOneSidedTransactionForSigningRequest) (*tari_generated.PrepareOneSidedTransactionForSigningResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := PrepareOneSidedTransactionForSigning(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("PrepareOneSidedTransactionForSigning returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("PrepareOneSidedTransactionForSigning response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestPrepareOneSidedTransactionForSigning_Error verifies that PrepareOneSidedTransactionForSigning propagates the underlying GRPC error rather than swallowing it.
func TestPrepareOneSidedTransactionForSigning_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-prepareonesidedtransactionforsigning")
	srv := &fakeWalletServer{
		prepareOneSidedTransactionForSigningFn: func(ctx context.Context, req *tari_generated.PrepareOneSidedTransactionForSigningRequest) (*tari_generated.PrepareOneSidedTransactionForSigningResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := PrepareOneSidedTransactionForSigning(context.Background(), &tari_generated.PrepareOneSidedTransactionForSigningRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-prepareonesidedtransactionforsigning" {
		t.Fatalf("PrepareOneSidedTransactionForSigning did not propagate the underlying error, got: %v", err)
	}
}

// TestBroadcastSignedOneSidedTransaction_Success verifies that BroadcastSignedOneSidedTransaction forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestBroadcastSignedOneSidedTransaction_Success(t *testing.T) {
	wantReq := &tari_generated.BroadcastSignedOneSidedTransactionRequest{}
	wantResp := &tari_generated.BroadcastSignedOneSidedTransactionResponse{}
	var gotReq *tari_generated.BroadcastSignedOneSidedTransactionRequest
	srv := &fakeWalletServer{
		broadcastSignedOneSidedTransactionFn: func(ctx context.Context, req *tari_generated.BroadcastSignedOneSidedTransactionRequest) (*tari_generated.BroadcastSignedOneSidedTransactionResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := BroadcastSignedOneSidedTransaction(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("BroadcastSignedOneSidedTransaction returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("BroadcastSignedOneSidedTransaction response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestBroadcastSignedOneSidedTransaction_Error verifies that BroadcastSignedOneSidedTransaction propagates the underlying GRPC error rather than swallowing it.
func TestBroadcastSignedOneSidedTransaction_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-broadcastsignedonesidedtransaction")
	srv := &fakeWalletServer{
		broadcastSignedOneSidedTransactionFn: func(ctx context.Context, req *tari_generated.BroadcastSignedOneSidedTransactionRequest) (*tari_generated.BroadcastSignedOneSidedTransactionResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := BroadcastSignedOneSidedTransaction(context.Background(), &tari_generated.BroadcastSignedOneSidedTransactionRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-broadcastsignedonesidedtransaction" {
		t.Fatalf("BroadcastSignedOneSidedTransaction did not propagate the underlying error, got: %v", err)
	}
}

// TestGetTransactionPayRefs_Success verifies that GetTransactionPayRefs forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetTransactionPayRefs_Success(t *testing.T) {
	wantReq := &tari_generated.GetTransactionPayRefsRequest{}
	wantResp := &tari_generated.GetTransactionPayRefsResponse{}
	var gotReq *tari_generated.GetTransactionPayRefsRequest
	srv := &fakeWalletServer{
		getTransactionPayRefsFn: func(ctx context.Context, req *tari_generated.GetTransactionPayRefsRequest) (*tari_generated.GetTransactionPayRefsResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := GetTransactionPayRefs(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetTransactionPayRefs returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetTransactionPayRefs response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetTransactionPayRefs_Error verifies that GetTransactionPayRefs propagates the underlying GRPC error rather than swallowing it.
func TestGetTransactionPayRefs_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-gettransactionpayrefs")
	srv := &fakeWalletServer{
		getTransactionPayRefsFn: func(ctx context.Context, req *tari_generated.GetTransactionPayRefsRequest) (*tari_generated.GetTransactionPayRefsResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := GetTransactionPayRefs(context.Background(), &tari_generated.GetTransactionPayRefsRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-gettransactionpayrefs" {
		t.Fatalf("GetTransactionPayRefs did not propagate the underlying error, got: %v", err)
	}
}

// TestGetUnspentAmounts_Success verifies that GetUnspentAmounts returns the response the server sends, and that
// the server actually received a request (arg marshaling for the no-argument/Empty case).
func TestGetUnspentAmounts_Success(t *testing.T) {
	wantResp := &tari_generated.GetUnspentAmountsResponse{}
	var gotReq *tari_generated.Empty
	srv := &fakeWalletServer{
		getUnspentAmountsFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.GetUnspentAmountsResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := GetUnspentAmounts(context.Background())
	if err != nil {
		t.Fatalf("GetUnspentAmounts returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetUnspentAmounts response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetUnspentAmounts_Error verifies that GetUnspentAmounts propagates the underlying GRPC error rather than swallowing it.
func TestGetUnspentAmounts_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getunspentamounts")
	srv := &fakeWalletServer{
		getUnspentAmountsFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.GetUnspentAmountsResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := GetUnspentAmounts(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getunspentamounts" {
		t.Fatalf("GetUnspentAmounts did not propagate the underlying error, got: %v", err)
	}
}

// TestImportUtxos_Success verifies that ImportUtxos forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestImportUtxos_Success(t *testing.T) {
	wantReq := &tari_generated.ImportUtxosRequest{}
	wantResp := &tari_generated.ImportUtxosResponse{}
	var gotReq *tari_generated.ImportUtxosRequest
	srv := &fakeWalletServer{
		importUtxosFn: func(ctx context.Context, req *tari_generated.ImportUtxosRequest) (*tari_generated.ImportUtxosResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := ImportUtxos(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("ImportUtxos returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("ImportUtxos response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestImportUtxos_Error verifies that ImportUtxos propagates the underlying GRPC error rather than swallowing it.
func TestImportUtxos_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-importutxos")
	srv := &fakeWalletServer{
		importUtxosFn: func(ctx context.Context, req *tari_generated.ImportUtxosRequest) (*tari_generated.ImportUtxosResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := ImportUtxos(context.Background(), &tari_generated.ImportUtxosRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-importutxos" {
		t.Fatalf("ImportUtxos did not propagate the underlying error, got: %v", err)
	}
}

// TestGetNetworkStatus_Success verifies that GetNetworkStatus returns the response the server sends, and that
// the server actually received a request (arg marshaling for the no-argument/Empty case).
func TestGetNetworkStatus_Success(t *testing.T) {
	wantResp := &tari_generated.NetworkStatusResponse{}
	var gotReq *tari_generated.Empty
	srv := &fakeWalletServer{
		getNetworkStatusFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.NetworkStatusResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := GetNetworkStatus(context.Background())
	if err != nil {
		t.Fatalf("GetNetworkStatus returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetNetworkStatus response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetNetworkStatus_Error verifies that GetNetworkStatus propagates the underlying GRPC error rather than swallowing it.
func TestGetNetworkStatus_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getnetworkstatus")
	srv := &fakeWalletServer{
		getNetworkStatusFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.NetworkStatusResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := GetNetworkStatus(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getnetworkstatus" {
		t.Fatalf("GetNetworkStatus did not propagate the underlying error, got: %v", err)
	}
}

// TestGetConnectedHttpPeer_Success verifies that GetConnectedHttpPeer returns the response the server sends, and that
// the server actually received a request (arg marshaling for the no-argument/Empty case).
func TestGetConnectedHttpPeer_Success(t *testing.T) {
	wantResp := &tari_generated.GetConnectedHttpPeerResponse{}
	var gotReq *tari_generated.Empty
	srv := &fakeWalletServer{
		getConnectedHttpPeerFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.GetConnectedHttpPeerResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := GetConnectedHttpPeer(context.Background())
	if err != nil {
		t.Fatalf("GetConnectedHttpPeer returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetConnectedHttpPeer response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetConnectedHttpPeer_Error verifies that GetConnectedHttpPeer propagates the underlying GRPC error rather than swallowing it.
func TestGetConnectedHttpPeer_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getconnectedhttppeer")
	srv := &fakeWalletServer{
		getConnectedHttpPeerFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.GetConnectedHttpPeerResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := GetConnectedHttpPeer(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getconnectedhttppeer" {
		t.Fatalf("GetConnectedHttpPeer did not propagate the underlying error, got: %v", err)
	}
}

// TestCancelTransaction_Success verifies that CancelTransaction forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestCancelTransaction_Success(t *testing.T) {
	wantReq := &tari_generated.CancelTransactionRequest{}
	wantResp := &tari_generated.CancelTransactionResponse{}
	var gotReq *tari_generated.CancelTransactionRequest
	srv := &fakeWalletServer{
		cancelTransactionFn: func(ctx context.Context, req *tari_generated.CancelTransactionRequest) (*tari_generated.CancelTransactionResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := CancelTransaction(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("CancelTransaction returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("CancelTransaction response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestCancelTransaction_Error verifies that CancelTransaction propagates the underlying GRPC error rather than swallowing it.
func TestCancelTransaction_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-canceltransaction")
	srv := &fakeWalletServer{
		cancelTransactionFn: func(ctx context.Context, req *tari_generated.CancelTransactionRequest) (*tari_generated.CancelTransactionResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := CancelTransaction(context.Background(), &tari_generated.CancelTransactionRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-canceltransaction" {
		t.Fatalf("CancelTransaction did not propagate the underlying error, got: %v", err)
	}
}

// TestSendShaAtomicSwapTransaction_Success verifies that SendShaAtomicSwapTransaction forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestSendShaAtomicSwapTransaction_Success(t *testing.T) {
	wantReq := &tari_generated.SendShaAtomicSwapRequest{}
	wantResp := &tari_generated.SendShaAtomicSwapResponse{}
	var gotReq *tari_generated.SendShaAtomicSwapRequest
	srv := &fakeWalletServer{
		sendShaAtomicSwapTransactionFn: func(ctx context.Context, req *tari_generated.SendShaAtomicSwapRequest) (*tari_generated.SendShaAtomicSwapResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := SendShaAtomicSwapTransaction(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("SendShaAtomicSwapTransaction returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("SendShaAtomicSwapTransaction response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestSendShaAtomicSwapTransaction_Error verifies that SendShaAtomicSwapTransaction propagates the underlying GRPC error rather than swallowing it.
func TestSendShaAtomicSwapTransaction_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-sendshaatomicswaptransaction")
	srv := &fakeWalletServer{
		sendShaAtomicSwapTransactionFn: func(ctx context.Context, req *tari_generated.SendShaAtomicSwapRequest) (*tari_generated.SendShaAtomicSwapResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := SendShaAtomicSwapTransaction(context.Background(), &tari_generated.SendShaAtomicSwapRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-sendshaatomicswaptransaction" {
		t.Fatalf("SendShaAtomicSwapTransaction did not propagate the underlying error, got: %v", err)
	}
}

// TestCreateBurnTransaction_Success verifies that CreateBurnTransaction forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestCreateBurnTransaction_Success(t *testing.T) {
	wantReq := &tari_generated.CreateBurnTransactionRequest{}
	wantResp := &tari_generated.CreateBurnTransactionResponse{}
	var gotReq *tari_generated.CreateBurnTransactionRequest
	srv := &fakeWalletServer{
		createBurnTransactionFn: func(ctx context.Context, req *tari_generated.CreateBurnTransactionRequest) (*tari_generated.CreateBurnTransactionResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := CreateBurnTransaction(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("CreateBurnTransaction returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("CreateBurnTransaction response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestCreateBurnTransaction_Error verifies that CreateBurnTransaction propagates the underlying GRPC error rather than swallowing it.
func TestCreateBurnTransaction_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-createburntransaction")
	srv := &fakeWalletServer{
		createBurnTransactionFn: func(ctx context.Context, req *tari_generated.CreateBurnTransactionRequest) (*tari_generated.CreateBurnTransactionResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := CreateBurnTransaction(context.Background(), &tari_generated.CreateBurnTransactionRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-createburntransaction" {
		t.Fatalf("CreateBurnTransaction did not propagate the underlying error, got: %v", err)
	}
}

// TestClaimShaAtomicSwapTransaction_Success verifies that ClaimShaAtomicSwapTransaction forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestClaimShaAtomicSwapTransaction_Success(t *testing.T) {
	wantReq := &tari_generated.ClaimShaAtomicSwapRequest{}
	wantResp := &tari_generated.ClaimShaAtomicSwapResponse{}
	var gotReq *tari_generated.ClaimShaAtomicSwapRequest
	srv := &fakeWalletServer{
		claimShaAtomicSwapTransactionFn: func(ctx context.Context, req *tari_generated.ClaimShaAtomicSwapRequest) (*tari_generated.ClaimShaAtomicSwapResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := ClaimShaAtomicSwapTransaction(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("ClaimShaAtomicSwapTransaction returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("ClaimShaAtomicSwapTransaction response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestClaimShaAtomicSwapTransaction_Error verifies that ClaimShaAtomicSwapTransaction propagates the underlying GRPC error rather than swallowing it.
func TestClaimShaAtomicSwapTransaction_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-claimshaatomicswaptransaction")
	srv := &fakeWalletServer{
		claimShaAtomicSwapTransactionFn: func(ctx context.Context, req *tari_generated.ClaimShaAtomicSwapRequest) (*tari_generated.ClaimShaAtomicSwapResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := ClaimShaAtomicSwapTransaction(context.Background(), &tari_generated.ClaimShaAtomicSwapRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-claimshaatomicswaptransaction" {
		t.Fatalf("ClaimShaAtomicSwapTransaction did not propagate the underlying error, got: %v", err)
	}
}

// TestClaimHtlcRefundTransaction_Success verifies that ClaimHtlcRefundTransaction forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestClaimHtlcRefundTransaction_Success(t *testing.T) {
	wantReq := &tari_generated.ClaimHtlcRefundRequest{}
	wantResp := &tari_generated.ClaimHtlcRefundResponse{}
	var gotReq *tari_generated.ClaimHtlcRefundRequest
	srv := &fakeWalletServer{
		claimHtlcRefundTransactionFn: func(ctx context.Context, req *tari_generated.ClaimHtlcRefundRequest) (*tari_generated.ClaimHtlcRefundResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := ClaimHtlcRefundTransaction(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("ClaimHtlcRefundTransaction returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("ClaimHtlcRefundTransaction response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestClaimHtlcRefundTransaction_Error verifies that ClaimHtlcRefundTransaction propagates the underlying GRPC error rather than swallowing it.
func TestClaimHtlcRefundTransaction_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-claimhtlcrefundtransaction")
	srv := &fakeWalletServer{
		claimHtlcRefundTransactionFn: func(ctx context.Context, req *tari_generated.ClaimHtlcRefundRequest) (*tari_generated.ClaimHtlcRefundResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := ClaimHtlcRefundTransaction(context.Background(), &tari_generated.ClaimHtlcRefundRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-claimhtlcrefundtransaction" {
		t.Fatalf("ClaimHtlcRefundTransaction did not propagate the underlying error, got: %v", err)
	}
}

// TestCreateTemplateRegistration_Success verifies that CreateTemplateRegistration forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestCreateTemplateRegistration_Success(t *testing.T) {
	wantReq := &tari_generated.CreateTemplateRegistrationRequest{}
	wantResp := &tari_generated.CreateTemplateRegistrationResponse{}
	var gotReq *tari_generated.CreateTemplateRegistrationRequest
	srv := &fakeWalletServer{
		createTemplateRegistrationFn: func(ctx context.Context, req *tari_generated.CreateTemplateRegistrationRequest) (*tari_generated.CreateTemplateRegistrationResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := CreateTemplateRegistration(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("CreateTemplateRegistration returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("CreateTemplateRegistration response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestCreateTemplateRegistration_Error verifies that CreateTemplateRegistration propagates the underlying GRPC error rather than swallowing it.
func TestCreateTemplateRegistration_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-createtemplateregistration")
	srv := &fakeWalletServer{
		createTemplateRegistrationFn: func(ctx context.Context, req *tari_generated.CreateTemplateRegistrationRequest) (*tari_generated.CreateTemplateRegistrationResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := CreateTemplateRegistration(context.Background(), &tari_generated.CreateTemplateRegistrationRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-createtemplateregistration" {
		t.Fatalf("CreateTemplateRegistration did not propagate the underlying error, got: %v", err)
	}
}

// TestSignMessage_Success verifies that SignMessage forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestSignMessage_Success(t *testing.T) {
	wantReq := &tari_generated.SignMessageRequest{}
	wantResp := &tari_generated.SignMessageResponse{}
	var gotReq *tari_generated.SignMessageRequest
	srv := &fakeWalletServer{
		signMessageFn: func(ctx context.Context, req *tari_generated.SignMessageRequest) (*tari_generated.SignMessageResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := SignMessage(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("SignMessage returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("SignMessage response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestSignMessage_Error verifies that SignMessage propagates the underlying GRPC error rather than swallowing it.
func TestSignMessage_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-signmessage")
	srv := &fakeWalletServer{
		signMessageFn: func(ctx context.Context, req *tari_generated.SignMessageRequest) (*tari_generated.SignMessageResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := SignMessage(context.Background(), &tari_generated.SignMessageRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-signmessage" {
		t.Fatalf("SignMessage did not propagate the underlying error, got: %v", err)
	}
}

// TestImportTransactions_Success verifies that ImportTransactions forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestImportTransactions_Success(t *testing.T) {
	wantReq := &tari_generated.ImportTransactionsRequest{}
	wantResp := &tari_generated.ImportTransactionsResponse{}
	var gotReq *tari_generated.ImportTransactionsRequest
	srv := &fakeWalletServer{
		importTransactionsFn: func(ctx context.Context, req *tari_generated.ImportTransactionsRequest) (*tari_generated.ImportTransactionsResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := ImportTransactions(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("ImportTransactions returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("ImportTransactions response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestImportTransactions_Error verifies that ImportTransactions propagates the underlying GRPC error rather than swallowing it.
func TestImportTransactions_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-importtransactions")
	srv := &fakeWalletServer{
		importTransactionsFn: func(ctx context.Context, req *tari_generated.ImportTransactionsRequest) (*tari_generated.ImportTransactionsResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := ImportTransactions(context.Background(), &tari_generated.ImportTransactionsRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-importtransactions" {
		t.Fatalf("ImportTransactions did not propagate the underlying error, got: %v", err)
	}
}

// TestGetAllCompletedTransactions_Success verifies that GetAllCompletedTransactions forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetAllCompletedTransactions_Success(t *testing.T) {
	wantReq := &tari_generated.GetAllCompletedTransactionsRequest{}
	wantResp := &tari_generated.GetAllCompletedTransactionsResponse{}
	var gotReq *tari_generated.GetAllCompletedTransactionsRequest
	srv := &fakeWalletServer{
		getAllCompletedTransactionsFn: func(ctx context.Context, req *tari_generated.GetAllCompletedTransactionsRequest) (*tari_generated.GetAllCompletedTransactionsResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := GetAllCompletedTransactions(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetAllCompletedTransactions returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetAllCompletedTransactions response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetAllCompletedTransactions_Error verifies that GetAllCompletedTransactions propagates the underlying GRPC error rather than swallowing it.
func TestGetAllCompletedTransactions_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getallcompletedtransactions")
	srv := &fakeWalletServer{
		getAllCompletedTransactionsFn: func(ctx context.Context, req *tari_generated.GetAllCompletedTransactionsRequest) (*tari_generated.GetAllCompletedTransactionsResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := GetAllCompletedTransactions(context.Background(), &tari_generated.GetAllCompletedTransactionsRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getallcompletedtransactions" {
		t.Fatalf("GetAllCompletedTransactions did not propagate the underlying error, got: %v", err)
	}
}

// TestGetPaymentByReference_Success verifies that GetPaymentByReference forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetPaymentByReference_Success(t *testing.T) {
	wantReq := &tari_generated.GetPaymentByReferenceRequest{}
	wantResp := &tari_generated.GetPaymentByReferenceResponse{}
	var gotReq *tari_generated.GetPaymentByReferenceRequest
	srv := &fakeWalletServer{
		getPaymentByReferenceFn: func(ctx context.Context, req *tari_generated.GetPaymentByReferenceRequest) (*tari_generated.GetPaymentByReferenceResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := GetPaymentByReference(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetPaymentByReference returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetPaymentByReference response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetPaymentByReference_Error verifies that GetPaymentByReference propagates the underlying GRPC error rather than swallowing it.
func TestGetPaymentByReference_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getpaymentbyreference")
	srv := &fakeWalletServer{
		getPaymentByReferenceFn: func(ctx context.Context, req *tari_generated.GetPaymentByReferenceRequest) (*tari_generated.GetPaymentByReferenceResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := GetPaymentByReference(context.Background(), &tari_generated.GetPaymentByReferenceRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getpaymentbyreference" {
		t.Fatalf("GetPaymentByReference did not propagate the underlying error, got: %v", err)
	}
}

// TestGetFeeEstimate_Success verifies that GetFeeEstimate forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetFeeEstimate_Success(t *testing.T) {
	wantReq := &tari_generated.GetFeeEstimateRequest{}
	wantResp := &tari_generated.GetFeeEstimateResponse{}
	var gotReq *tari_generated.GetFeeEstimateRequest
	srv := &fakeWalletServer{
		getFeeEstimateFn: func(ctx context.Context, req *tari_generated.GetFeeEstimateRequest) (*tari_generated.GetFeeEstimateResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := GetFeeEstimate(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetFeeEstimate returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetFeeEstimate response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetFeeEstimate_Error verifies that GetFeeEstimate propagates the underlying GRPC error rather than swallowing it.
func TestGetFeeEstimate_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getfeeestimate")
	srv := &fakeWalletServer{
		getFeeEstimateFn: func(ctx context.Context, req *tari_generated.GetFeeEstimateRequest) (*tari_generated.GetFeeEstimateResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := GetFeeEstimate(context.Background(), &tari_generated.GetFeeEstimateRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getfeeestimate" {
		t.Fatalf("GetFeeEstimate did not propagate the underlying error, got: %v", err)
	}
}

// TestGetFeePerGramStats_Success verifies that GetFeePerGramStats forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetFeePerGramStats_Success(t *testing.T) {
	wantReq := &tari_generated.GetFeePerGramStatsRequest{}
	wantResp := &tari_generated.GetFeePerGramStatsResponse{}
	var gotReq *tari_generated.GetFeePerGramStatsRequest
	srv := &fakeWalletServer{
		getFeePerGramStatsFn: func(ctx context.Context, req *tari_generated.GetFeePerGramStatsRequest) (*tari_generated.GetFeePerGramStatsResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := GetFeePerGramStats(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetFeePerGramStats returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetFeePerGramStats response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetFeePerGramStats_Error verifies that GetFeePerGramStats propagates the underlying GRPC error rather than swallowing it.
func TestGetFeePerGramStats_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getfeepergramstats")
	srv := &fakeWalletServer{
		getFeePerGramStatsFn: func(ctx context.Context, req *tari_generated.GetFeePerGramStatsRequest) (*tari_generated.GetFeePerGramStatsResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := GetFeePerGramStats(context.Background(), &tari_generated.GetFeePerGramStatsRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getfeepergramstats" {
		t.Fatalf("GetFeePerGramStats did not propagate the underlying error, got: %v", err)
	}
}

// TestReplaceByFee_Success verifies that ReplaceByFee forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestReplaceByFee_Success(t *testing.T) {
	wantReq := &tari_generated.ReplaceByFeeRequest{}
	wantResp := &tari_generated.ReplaceByFeeResponse{}
	var gotReq *tari_generated.ReplaceByFeeRequest
	srv := &fakeWalletServer{
		replaceByFeeFn: func(ctx context.Context, req *tari_generated.ReplaceByFeeRequest) (*tari_generated.ReplaceByFeeResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := ReplaceByFee(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("ReplaceByFee returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("ReplaceByFee response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestReplaceByFee_Error verifies that ReplaceByFee propagates the underlying GRPC error rather than swallowing it.
func TestReplaceByFee_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-replacebyfee")
	srv := &fakeWalletServer{
		replaceByFeeFn: func(ctx context.Context, req *tari_generated.ReplaceByFeeRequest) (*tari_generated.ReplaceByFeeResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := ReplaceByFee(context.Background(), &tari_generated.ReplaceByFeeRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-replacebyfee" {
		t.Fatalf("ReplaceByFee did not propagate the underlying error, got: %v", err)
	}
}

// TestUserPayForFee_Success verifies that UserPayForFee forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestUserPayForFee_Success(t *testing.T) {
	wantReq := &tari_generated.UserPayForFeeRequest{}
	wantResp := &tari_generated.UserPayForFeeResponse{}
	var gotReq *tari_generated.UserPayForFeeRequest
	srv := &fakeWalletServer{
		userPayForFeeFn: func(ctx context.Context, req *tari_generated.UserPayForFeeRequest) (*tari_generated.UserPayForFeeResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := UserPayForFee(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("UserPayForFee returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("UserPayForFee response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestUserPayForFee_Error verifies that UserPayForFee propagates the underlying GRPC error rather than swallowing it.
func TestUserPayForFee_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-userpayforfee")
	srv := &fakeWalletServer{
		userPayForFeeFn: func(ctx context.Context, req *tari_generated.UserPayForFeeRequest) (*tari_generated.UserPayForFeeResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := UserPayForFee(context.Background(), &tari_generated.UserPayForFeeRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-userpayforfee" {
		t.Fatalf("UserPayForFee did not propagate the underlying error, got: %v", err)
	}
}

// TestRegisterValidatorNode_Success verifies that RegisterValidatorNode forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestRegisterValidatorNode_Success(t *testing.T) {
	wantReq := &tari_generated.RegisterValidatorNodeRequest{}
	wantResp := &tari_generated.RegisterValidatorNodeResponse{}
	var gotReq *tari_generated.RegisterValidatorNodeRequest
	srv := &fakeWalletServer{
		registerValidatorNodeFn: func(ctx context.Context, req *tari_generated.RegisterValidatorNodeRequest) (*tari_generated.RegisterValidatorNodeResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := RegisterValidatorNode(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("RegisterValidatorNode returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("RegisterValidatorNode response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestRegisterValidatorNode_Error verifies that RegisterValidatorNode propagates the underlying GRPC error rather than swallowing it.
func TestRegisterValidatorNode_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-registervalidatornode")
	srv := &fakeWalletServer{
		registerValidatorNodeFn: func(ctx context.Context, req *tari_generated.RegisterValidatorNodeRequest) (*tari_generated.RegisterValidatorNodeResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := RegisterValidatorNode(context.Background(), &tari_generated.RegisterValidatorNodeRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-registervalidatornode" {
		t.Fatalf("RegisterValidatorNode did not propagate the underlying error, got: %v", err)
	}
}

// TestSubmitValidatorEvictionProof_Success verifies that SubmitValidatorEvictionProof forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestSubmitValidatorEvictionProof_Success(t *testing.T) {
	wantReq := &tari_generated.SubmitValidatorEvictionProofRequest{}
	wantResp := &tari_generated.SubmitValidatorEvictionProofResponse{}
	var gotReq *tari_generated.SubmitValidatorEvictionProofRequest
	srv := &fakeWalletServer{
		submitValidatorEvictionProofFn: func(ctx context.Context, req *tari_generated.SubmitValidatorEvictionProofRequest) (*tari_generated.SubmitValidatorEvictionProofResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := SubmitValidatorEvictionProof(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("SubmitValidatorEvictionProof returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("SubmitValidatorEvictionProof response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestSubmitValidatorEvictionProof_Error verifies that SubmitValidatorEvictionProof propagates the underlying GRPC error rather than swallowing it.
func TestSubmitValidatorEvictionProof_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-submitvalidatorevictionproof")
	srv := &fakeWalletServer{
		submitValidatorEvictionProofFn: func(ctx context.Context, req *tari_generated.SubmitValidatorEvictionProofRequest) (*tari_generated.SubmitValidatorEvictionProofResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := SubmitValidatorEvictionProof(context.Background(), &tari_generated.SubmitValidatorEvictionProofRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-submitvalidatorevictionproof" {
		t.Fatalf("SubmitValidatorEvictionProof did not propagate the underlying error, got: %v", err)
	}
}

// TestSubmitValidatorNodeExit_Success verifies that SubmitValidatorNodeExit forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestSubmitValidatorNodeExit_Success(t *testing.T) {
	wantReq := &tari_generated.SubmitValidatorNodeExitRequest{}
	wantResp := &tari_generated.SubmitValidatorNodeExitResponse{}
	var gotReq *tari_generated.SubmitValidatorNodeExitRequest
	srv := &fakeWalletServer{
		submitValidatorNodeExitFn: func(ctx context.Context, req *tari_generated.SubmitValidatorNodeExitRequest) (*tari_generated.SubmitValidatorNodeExitResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	gotResp, err := SubmitValidatorNodeExit(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("SubmitValidatorNodeExit returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("SubmitValidatorNodeExit response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestSubmitValidatorNodeExit_Error verifies that SubmitValidatorNodeExit propagates the underlying GRPC error rather than swallowing it.
func TestSubmitValidatorNodeExit_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-submitvalidatornodeexit")
	srv := &fakeWalletServer{
		submitValidatorNodeExitFn: func(ctx context.Context, req *tari_generated.SubmitValidatorNodeExitRequest) (*tari_generated.SubmitValidatorNodeExitResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := SubmitValidatorNodeExit(context.Background(), &tari_generated.SubmitValidatorNodeExitRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-submitvalidatornodeexit" {
		t.Fatalf("SubmitValidatorNodeExit did not propagate the underlying error, got: %v", err)
	}
}

// TestStreamTransactionEvents_DrainsStream verifies that StreamTransactionEvents forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestStreamTransactionEvents_DrainsStream(t *testing.T) {
	want := []*tari_generated.TransactionEventResponse{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.TransactionEventRequest
	srv := &fakeWalletServer{
		streamTransactionEventsFn: func(req *tari_generated.TransactionEventRequest, stream tari_generated.Wallet_StreamTransactionEventsServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	got, err := StreamTransactionEvents(context.Background(), &tari_generated.TransactionEventRequest{})
	if err != nil {
		t.Fatalf("StreamTransactionEvents returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("StreamTransactionEvents returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("StreamTransactionEvents item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestStreamTransactionEvents_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, StreamTransactionEvents propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestStreamTransactionEvents_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-streamtransactionevents-stream")
	srv := &fakeWalletServer{
		streamTransactionEventsFn: func(req *tari_generated.TransactionEventRequest, stream tari_generated.Wallet_StreamTransactionEventsServer) error {
			if err := stream.Send(&tari_generated.TransactionEventResponse{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	got, err := StreamTransactionEvents(context.Background(), &tari_generated.TransactionEventRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-streamtransactionevents-stream" {
		t.Fatalf("StreamTransactionEvents did not propagate the underlying error, got: %v", err)
	}
}

// TestGetAllCompletedTransactionsStream_DrainsStream verifies that GetAllCompletedTransactionsStream forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestGetAllCompletedTransactionsStream_DrainsStream(t *testing.T) {
	want := []*tari_generated.GetCompletedTransactionsResponse{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.GetAllCompletedTransactionsRequest
	srv := &fakeWalletServer{
		getAllCompletedTransactionsStreamFn: func(req *tari_generated.GetAllCompletedTransactionsRequest, stream tari_generated.Wallet_GetAllCompletedTransactionsStreamServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	got, err := GetAllCompletedTransactionsStream(context.Background(), &tari_generated.GetAllCompletedTransactionsRequest{})
	if err != nil {
		t.Fatalf("GetAllCompletedTransactionsStream returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("GetAllCompletedTransactionsStream returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("GetAllCompletedTransactionsStream item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestGetAllCompletedTransactionsStream_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, GetAllCompletedTransactionsStream propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestGetAllCompletedTransactionsStream_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getallcompletedtransactionsstream-stream")
	srv := &fakeWalletServer{
		getAllCompletedTransactionsStreamFn: func(req *tari_generated.GetAllCompletedTransactionsRequest, stream tari_generated.Wallet_GetAllCompletedTransactionsStreamServer) error {
			if err := stream.Send(&tari_generated.GetCompletedTransactionsResponse{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	got, err := GetAllCompletedTransactionsStream(context.Background(), &tari_generated.GetAllCompletedTransactionsRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-getallcompletedtransactionsstream-stream" {
		t.Fatalf("GetAllCompletedTransactionsStream did not propagate the underlying error, got: %v", err)
	}
}
