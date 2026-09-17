package nodeGRPC

import (
	"context"
	"testing"

	"github.com/Snipa22/go-tari-grpc-lib/v3/tari_generated"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// TestGetBlockTiming_Success verifies that GetBlockTiming forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetBlockTiming_Success(t *testing.T) {
	wantReq := &tari_generated.HeightRequest{}
	wantResp := &tari_generated.BlockTimingResponse{}
	var gotReq *tari_generated.HeightRequest
	srv := &fakeBaseNodeServer{
		getBlockTimingFn: func(ctx context.Context, req *tari_generated.HeightRequest) (*tari_generated.BlockTimingResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := GetBlockTiming(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetBlockTiming returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetBlockTiming response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetBlockTiming_Error verifies that GetBlockTiming propagates the underlying GRPC error rather than swallowing it.
func TestGetBlockTiming_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getblocktiming")
	srv := &fakeBaseNodeServer{
		getBlockTimingFn: func(ctx context.Context, req *tari_generated.HeightRequest) (*tari_generated.BlockTimingResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetBlockTiming(context.Background(), &tari_generated.HeightRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getblocktiming" {
		t.Fatalf("GetBlockTiming did not propagate the underlying error, got: %v", err)
	}
}

// TestGetConstants_Success verifies that GetConstants forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetConstants_Success(t *testing.T) {
	wantReq := &tari_generated.BlockHeight{}
	wantResp := &tari_generated.ConsensusConstants{}
	var gotReq *tari_generated.BlockHeight
	srv := &fakeBaseNodeServer{
		getConstantsFn: func(ctx context.Context, req *tari_generated.BlockHeight) (*tari_generated.ConsensusConstants, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := GetConstants(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetConstants returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetConstants response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetConstants_Error verifies that GetConstants propagates the underlying GRPC error rather than swallowing it.
func TestGetConstants_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getconstants")
	srv := &fakeBaseNodeServer{
		getConstantsFn: func(ctx context.Context, req *tari_generated.BlockHeight) (*tari_generated.ConsensusConstants, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetConstants(context.Background(), &tari_generated.BlockHeight{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getconstants" {
		t.Fatalf("GetConstants did not propagate the underlying error, got: %v", err)
	}
}

// TestGetBlockSize_Success verifies that GetBlockSize forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetBlockSize_Success(t *testing.T) {
	wantReq := &tari_generated.BlockGroupRequest{}
	wantResp := &tari_generated.BlockGroupResponse{}
	var gotReq *tari_generated.BlockGroupRequest
	srv := &fakeBaseNodeServer{
		getBlockSizeFn: func(ctx context.Context, req *tari_generated.BlockGroupRequest) (*tari_generated.BlockGroupResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := GetBlockSize(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetBlockSize returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetBlockSize response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetBlockSize_Error verifies that GetBlockSize propagates the underlying GRPC error rather than swallowing it.
func TestGetBlockSize_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getblocksize")
	srv := &fakeBaseNodeServer{
		getBlockSizeFn: func(ctx context.Context, req *tari_generated.BlockGroupRequest) (*tari_generated.BlockGroupResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetBlockSize(context.Background(), &tari_generated.BlockGroupRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getblocksize" {
		t.Fatalf("GetBlockSize did not propagate the underlying error, got: %v", err)
	}
}

// TestGetBlockFees_Success verifies that GetBlockFees forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetBlockFees_Success(t *testing.T) {
	wantReq := &tari_generated.BlockGroupRequest{}
	wantResp := &tari_generated.BlockGroupResponse{}
	var gotReq *tari_generated.BlockGroupRequest
	srv := &fakeBaseNodeServer{
		getBlockFeesFn: func(ctx context.Context, req *tari_generated.BlockGroupRequest) (*tari_generated.BlockGroupResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := GetBlockFees(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetBlockFees returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetBlockFees response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetBlockFees_Error verifies that GetBlockFees propagates the underlying GRPC error rather than swallowing it.
func TestGetBlockFees_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getblockfees")
	srv := &fakeBaseNodeServer{
		getBlockFeesFn: func(ctx context.Context, req *tari_generated.BlockGroupRequest) (*tari_generated.BlockGroupResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetBlockFees(context.Background(), &tari_generated.BlockGroupRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getblockfees" {
		t.Fatalf("GetBlockFees did not propagate the underlying error, got: %v", err)
	}
}

// TestGetVersion_Success verifies that GetVersion returns the response the server sends, and that
// the server actually received a request (arg marshaling for the no-argument/Empty case).
func TestGetVersion_Success(t *testing.T) {
	wantResp := &tari_generated.BaseNodeGetVersionResponse{}
	var gotReq *tari_generated.Empty
	srv := &fakeBaseNodeServer{
		getVersionFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.BaseNodeGetVersionResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := GetVersion(context.Background())
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
	srv := &fakeBaseNodeServer{
		getVersionFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.BaseNodeGetVersionResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetVersion(context.Background())
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
	srv := &fakeBaseNodeServer{
		checkForUpdatesFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SoftwareUpdate, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
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
	srv := &fakeBaseNodeServer{
		checkForUpdatesFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SoftwareUpdate, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := CheckForUpdates(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-checkforupdates" {
		t.Fatalf("CheckForUpdates did not propagate the underlying error, got: %v", err)
	}
}

// TestGetNewBlockBlob_Success verifies that GetNewBlockBlob forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetNewBlockBlob_Success(t *testing.T) {
	wantReq := &tari_generated.NewBlockTemplate{}
	wantResp := &tari_generated.GetNewBlockBlobResult{}
	var gotReq *tari_generated.NewBlockTemplate
	srv := &fakeBaseNodeServer{
		getNewBlockBlobFn: func(ctx context.Context, req *tari_generated.NewBlockTemplate) (*tari_generated.GetNewBlockBlobResult, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := GetNewBlockBlob(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetNewBlockBlob returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetNewBlockBlob response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetNewBlockBlob_Error verifies that GetNewBlockBlob propagates the underlying GRPC error rather than swallowing it.
func TestGetNewBlockBlob_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getnewblockblob")
	srv := &fakeBaseNodeServer{
		getNewBlockBlobFn: func(ctx context.Context, req *tari_generated.NewBlockTemplate) (*tari_generated.GetNewBlockBlobResult, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetNewBlockBlob(context.Background(), &tari_generated.NewBlockTemplate{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getnewblockblob" {
		t.Fatalf("GetNewBlockBlob did not propagate the underlying error, got: %v", err)
	}
}

// TestSubmitBlockBlob_Success verifies that SubmitBlockBlob forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestSubmitBlockBlob_Success(t *testing.T) {
	wantReq := &tari_generated.BlockBlobRequest{}
	wantResp := &tari_generated.SubmitBlockResponse{}
	var gotReq *tari_generated.BlockBlobRequest
	srv := &fakeBaseNodeServer{
		submitBlockBlobFn: func(ctx context.Context, req *tari_generated.BlockBlobRequest) (*tari_generated.SubmitBlockResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := SubmitBlockBlob(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("SubmitBlockBlob returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("SubmitBlockBlob response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestSubmitBlockBlob_Error verifies that SubmitBlockBlob propagates the underlying GRPC error rather than swallowing it.
func TestSubmitBlockBlob_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-submitblockblob")
	srv := &fakeBaseNodeServer{
		submitBlockBlobFn: func(ctx context.Context, req *tari_generated.BlockBlobRequest) (*tari_generated.SubmitBlockResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := SubmitBlockBlob(context.Background(), &tari_generated.BlockBlobRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-submitblockblob" {
		t.Fatalf("SubmitBlockBlob did not propagate the underlying error, got: %v", err)
	}
}

// TestSubmitTransaction_Success verifies that SubmitTransaction forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestSubmitTransaction_Success(t *testing.T) {
	wantReq := &tari_generated.SubmitTransactionRequest{}
	wantResp := &tari_generated.SubmitTransactionResponse{}
	var gotReq *tari_generated.SubmitTransactionRequest
	srv := &fakeBaseNodeServer{
		submitTransactionFn: func(ctx context.Context, req *tari_generated.SubmitTransactionRequest) (*tari_generated.SubmitTransactionResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := SubmitTransaction(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("SubmitTransaction returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("SubmitTransaction response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestSubmitTransaction_Error verifies that SubmitTransaction propagates the underlying GRPC error rather than swallowing it.
func TestSubmitTransaction_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-submittransaction")
	srv := &fakeBaseNodeServer{
		submitTransactionFn: func(ctx context.Context, req *tari_generated.SubmitTransactionRequest) (*tari_generated.SubmitTransactionResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := SubmitTransaction(context.Background(), &tari_generated.SubmitTransactionRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-submittransaction" {
		t.Fatalf("SubmitTransaction did not propagate the underlying error, got: %v", err)
	}
}

// TestGetSyncInfo_Success verifies that GetSyncInfo returns the response the server sends, and that
// the server actually received a request (arg marshaling for the no-argument/Empty case).
func TestGetSyncInfo_Success(t *testing.T) {
	wantResp := &tari_generated.SyncInfoResponse{}
	var gotReq *tari_generated.Empty
	srv := &fakeBaseNodeServer{
		getSyncInfoFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SyncInfoResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := GetSyncInfo(context.Background())
	if err != nil {
		t.Fatalf("GetSyncInfo returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetSyncInfo response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetSyncInfo_Error verifies that GetSyncInfo propagates the underlying GRPC error rather than swallowing it.
func TestGetSyncInfo_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getsyncinfo")
	srv := &fakeBaseNodeServer{
		getSyncInfoFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SyncInfoResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetSyncInfo(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getsyncinfo" {
		t.Fatalf("GetSyncInfo did not propagate the underlying error, got: %v", err)
	}
}

// TestGetSyncProgress_Success verifies that GetSyncProgress returns the response the server sends, and that
// the server actually received a request (arg marshaling for the no-argument/Empty case).
func TestGetSyncProgress_Success(t *testing.T) {
	wantResp := &tari_generated.SyncProgressResponse{}
	var gotReq *tari_generated.Empty
	srv := &fakeBaseNodeServer{
		getSyncProgressFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SyncProgressResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := GetSyncProgress(context.Background())
	if err != nil {
		t.Fatalf("GetSyncProgress returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetSyncProgress response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetSyncProgress_Error verifies that GetSyncProgress propagates the underlying GRPC error rather than swallowing it.
func TestGetSyncProgress_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getsyncprogress")
	srv := &fakeBaseNodeServer{
		getSyncProgressFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SyncProgressResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetSyncProgress(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getsyncprogress" {
		t.Fatalf("GetSyncProgress did not propagate the underlying error, got: %v", err)
	}
}

// TestTransactionState_Success verifies that TransactionState forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestTransactionState_Success(t *testing.T) {
	wantReq := &tari_generated.TransactionStateRequest{}
	wantResp := &tari_generated.TransactionStateResponse{}
	var gotReq *tari_generated.TransactionStateRequest
	srv := &fakeBaseNodeServer{
		transactionStateFn: func(ctx context.Context, req *tari_generated.TransactionStateRequest) (*tari_generated.TransactionStateResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := TransactionState(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("TransactionState returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("TransactionState response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestTransactionState_Error verifies that TransactionState propagates the underlying GRPC error rather than swallowing it.
func TestTransactionState_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-transactionstate")
	srv := &fakeBaseNodeServer{
		transactionStateFn: func(ctx context.Context, req *tari_generated.TransactionStateRequest) (*tari_generated.TransactionStateResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := TransactionState(context.Background(), &tari_generated.TransactionStateRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-transactionstate" {
		t.Fatalf("TransactionState did not propagate the underlying error, got: %v", err)
	}
}

// TestGetNetworkStatus_Success verifies that GetNetworkStatus returns the response the server sends, and that
// the server actually received a request (arg marshaling for the no-argument/Empty case).
func TestGetNetworkStatus_Success(t *testing.T) {
	wantResp := &tari_generated.NetworkStatusResponse{}
	var gotReq *tari_generated.Empty
	srv := &fakeBaseNodeServer{
		getNetworkStatusFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.NetworkStatusResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
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
	srv := &fakeBaseNodeServer{
		getNetworkStatusFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.NetworkStatusResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetNetworkStatus(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getnetworkstatus" {
		t.Fatalf("GetNetworkStatus did not propagate the underlying error, got: %v", err)
	}
}

// TestListConnectedPeers_Success verifies that ListConnectedPeers returns the response the server sends, and that
// the server actually received a request (arg marshaling for the no-argument/Empty case).
func TestListConnectedPeers_Success(t *testing.T) {
	wantResp := &tari_generated.ListConnectedPeersResponse{}
	var gotReq *tari_generated.Empty
	srv := &fakeBaseNodeServer{
		listConnectedPeersFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.ListConnectedPeersResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := ListConnectedPeers(context.Background())
	if err != nil {
		t.Fatalf("ListConnectedPeers returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("ListConnectedPeers response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestListConnectedPeers_Error verifies that ListConnectedPeers propagates the underlying GRPC error rather than swallowing it.
func TestListConnectedPeers_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-listconnectedpeers")
	srv := &fakeBaseNodeServer{
		listConnectedPeersFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.ListConnectedPeersResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := ListConnectedPeers(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-listconnectedpeers" {
		t.Fatalf("ListConnectedPeers did not propagate the underlying error, got: %v", err)
	}
}

// TestGetMempoolStats_Success verifies that GetMempoolStats returns the response the server sends, and that
// the server actually received a request (arg marshaling for the no-argument/Empty case).
func TestGetMempoolStats_Success(t *testing.T) {
	wantResp := &tari_generated.MempoolStatsResponse{}
	var gotReq *tari_generated.Empty
	srv := &fakeBaseNodeServer{
		getMempoolStatsFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.MempoolStatsResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := GetMempoolStats(context.Background())
	if err != nil {
		t.Fatalf("GetMempoolStats returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetMempoolStats response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetMempoolStats_Error verifies that GetMempoolStats propagates the underlying GRPC error rather than swallowing it.
func TestGetMempoolStats_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getmempoolstats")
	srv := &fakeBaseNodeServer{
		getMempoolStatsFn: func(ctx context.Context, req *tari_generated.Empty) (*tari_generated.MempoolStatsResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetMempoolStats(context.Background())
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getmempoolstats" {
		t.Fatalf("GetMempoolStats did not propagate the underlying error, got: %v", err)
	}
}

// TestGetValidatorNodeChanges_Success verifies that GetValidatorNodeChanges forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetValidatorNodeChanges_Success(t *testing.T) {
	wantReq := &tari_generated.GetValidatorNodeChangesRequest{}
	wantResp := &tari_generated.GetValidatorNodeChangesResponse{}
	var gotReq *tari_generated.GetValidatorNodeChangesRequest
	srv := &fakeBaseNodeServer{
		getValidatorNodeChangesFn: func(ctx context.Context, req *tari_generated.GetValidatorNodeChangesRequest) (*tari_generated.GetValidatorNodeChangesResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := GetValidatorNodeChanges(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetValidatorNodeChanges returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetValidatorNodeChanges response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetValidatorNodeChanges_Error verifies that GetValidatorNodeChanges propagates the underlying GRPC error rather than swallowing it.
func TestGetValidatorNodeChanges_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getvalidatornodechanges")
	srv := &fakeBaseNodeServer{
		getValidatorNodeChangesFn: func(ctx context.Context, req *tari_generated.GetValidatorNodeChangesRequest) (*tari_generated.GetValidatorNodeChangesResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetValidatorNodeChanges(context.Background(), &tari_generated.GetValidatorNodeChangesRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getvalidatornodechanges" {
		t.Fatalf("GetValidatorNodeChanges did not propagate the underlying error, got: %v", err)
	}
}

// TestGetShardKey_Success verifies that GetShardKey forwards the request as-is and returns the response
// the server sends (arg marshaling + response passthrough).
func TestGetShardKey_Success(t *testing.T) {
	wantReq := &tari_generated.GetShardKeyRequest{}
	wantResp := &tari_generated.GetShardKeyResponse{}
	var gotReq *tari_generated.GetShardKeyRequest
	srv := &fakeBaseNodeServer{
		getShardKeyFn: func(ctx context.Context, req *tari_generated.GetShardKeyRequest) (*tari_generated.GetShardKeyResponse, error) {
			gotReq = req
			return wantResp, nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	gotResp, err := GetShardKey(context.Background(), wantReq)
	if err != nil {
		t.Fatalf("GetShardKey returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if !proto.Equal(gotResp, wantResp) {
		t.Fatalf("GetShardKey response = %+v, want %+v", gotResp, wantResp)
	}
}

// TestGetShardKey_Error verifies that GetShardKey propagates the underlying GRPC error rather than swallowing it.
func TestGetShardKey_Error(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getshardkey")
	srv := &fakeBaseNodeServer{
		getShardKeyFn: func(ctx context.Context, req *tari_generated.GetShardKeyRequest) (*tari_generated.GetShardKeyResponse, error) {
			return nil, wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetShardKey(context.Background(), &tari_generated.GetShardKeyRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if status.Convert(err).Message() != "boom-getshardkey" {
		t.Fatalf("GetShardKey did not propagate the underlying error, got: %v", err)
	}
}

// TestListHeaders_DrainsStream verifies that ListHeaders forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestListHeaders_DrainsStream(t *testing.T) {
	want := []*tari_generated.BlockHeaderResponse{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.ListHeadersRequest
	srv := &fakeBaseNodeServer{
		listHeadersFn: func(req *tari_generated.ListHeadersRequest, stream tari_generated.BaseNode_ListHeadersServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := ListHeaders(context.Background(), &tari_generated.ListHeadersRequest{})
	if err != nil {
		t.Fatalf("ListHeaders returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("ListHeaders returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("ListHeaders item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestListHeaders_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, ListHeaders propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestListHeaders_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-listheaders-stream")
	srv := &fakeBaseNodeServer{
		listHeadersFn: func(req *tari_generated.ListHeadersRequest, stream tari_generated.BaseNode_ListHeadersServer) error {
			if err := stream.Send(&tari_generated.BlockHeaderResponse{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := ListHeaders(context.Background(), &tari_generated.ListHeadersRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-listheaders-stream" {
		t.Fatalf("ListHeaders did not propagate the underlying error, got: %v", err)
	}
}

// TestGetTokensInCirculation_DrainsStream verifies that GetTokensInCirculation forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestGetTokensInCirculation_DrainsStream(t *testing.T) {
	want := []*tari_generated.ValueAtHeightResponse{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.GetBlocksRequest
	srv := &fakeBaseNodeServer{
		getTokensInCirculationFn: func(req *tari_generated.GetBlocksRequest, stream tari_generated.BaseNode_GetTokensInCirculationServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := GetTokensInCirculation(context.Background(), &tari_generated.GetBlocksRequest{})
	if err != nil {
		t.Fatalf("GetTokensInCirculation returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("GetTokensInCirculation returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("GetTokensInCirculation item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestGetTokensInCirculation_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, GetTokensInCirculation propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestGetTokensInCirculation_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-gettokensincirculation-stream")
	srv := &fakeBaseNodeServer{
		getTokensInCirculationFn: func(req *tari_generated.GetBlocksRequest, stream tari_generated.BaseNode_GetTokensInCirculationServer) error {
			if err := stream.Send(&tari_generated.ValueAtHeightResponse{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := GetTokensInCirculation(context.Background(), &tari_generated.GetBlocksRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-gettokensincirculation-stream" {
		t.Fatalf("GetTokensInCirculation did not propagate the underlying error, got: %v", err)
	}
}

// TestSearchKernels_DrainsStream verifies that SearchKernels forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestSearchKernels_DrainsStream(t *testing.T) {
	want := []*tari_generated.HistoricalBlock{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.SearchKernelsRequest
	srv := &fakeBaseNodeServer{
		searchKernelsFn: func(req *tari_generated.SearchKernelsRequest, stream tari_generated.BaseNode_SearchKernelsServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := SearchKernels(context.Background(), &tari_generated.SearchKernelsRequest{})
	if err != nil {
		t.Fatalf("SearchKernels returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("SearchKernels returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("SearchKernels item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestSearchKernels_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, SearchKernels propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestSearchKernels_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-searchkernels-stream")
	srv := &fakeBaseNodeServer{
		searchKernelsFn: func(req *tari_generated.SearchKernelsRequest, stream tari_generated.BaseNode_SearchKernelsServer) error {
			if err := stream.Send(&tari_generated.HistoricalBlock{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := SearchKernels(context.Background(), &tari_generated.SearchKernelsRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-searchkernels-stream" {
		t.Fatalf("SearchKernels did not propagate the underlying error, got: %v", err)
	}
}

// TestSearchUtxos_DrainsStream verifies that SearchUtxos forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestSearchUtxos_DrainsStream(t *testing.T) {
	want := []*tari_generated.HistoricalBlock{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.SearchUtxosRequest
	srv := &fakeBaseNodeServer{
		searchUtxosFn: func(req *tari_generated.SearchUtxosRequest, stream tari_generated.BaseNode_SearchUtxosServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := SearchUtxos(context.Background(), &tari_generated.SearchUtxosRequest{})
	if err != nil {
		t.Fatalf("SearchUtxos returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("SearchUtxos returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("SearchUtxos item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestSearchUtxos_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, SearchUtxos propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestSearchUtxos_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-searchutxos-stream")
	srv := &fakeBaseNodeServer{
		searchUtxosFn: func(req *tari_generated.SearchUtxosRequest, stream tari_generated.BaseNode_SearchUtxosServer) error {
			if err := stream.Send(&tari_generated.HistoricalBlock{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := SearchUtxos(context.Background(), &tari_generated.SearchUtxosRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-searchutxos-stream" {
		t.Fatalf("SearchUtxos did not propagate the underlying error, got: %v", err)
	}
}

// TestFetchMatchingUtxos_DrainsStream verifies that FetchMatchingUtxos forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestFetchMatchingUtxos_DrainsStream(t *testing.T) {
	want := []*tari_generated.FetchMatchingUtxosResponse{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.FetchMatchingUtxosRequest
	srv := &fakeBaseNodeServer{
		fetchMatchingUtxosFn: func(req *tari_generated.FetchMatchingUtxosRequest, stream tari_generated.BaseNode_FetchMatchingUtxosServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := FetchMatchingUtxos(context.Background(), &tari_generated.FetchMatchingUtxosRequest{})
	if err != nil {
		t.Fatalf("FetchMatchingUtxos returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("FetchMatchingUtxos returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("FetchMatchingUtxos item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestFetchMatchingUtxos_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, FetchMatchingUtxos propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestFetchMatchingUtxos_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-fetchmatchingutxos-stream")
	srv := &fakeBaseNodeServer{
		fetchMatchingUtxosFn: func(req *tari_generated.FetchMatchingUtxosRequest, stream tari_generated.BaseNode_FetchMatchingUtxosServer) error {
			if err := stream.Send(&tari_generated.FetchMatchingUtxosResponse{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := FetchMatchingUtxos(context.Background(), &tari_generated.FetchMatchingUtxosRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-fetchmatchingutxos-stream" {
		t.Fatalf("FetchMatchingUtxos did not propagate the underlying error, got: %v", err)
	}
}

// TestGetPeers_DrainsStream verifies that GetPeers forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestGetPeers_DrainsStream(t *testing.T) {
	want := []*tari_generated.GetPeersResponse{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.GetPeersRequest
	srv := &fakeBaseNodeServer{
		getPeersFn: func(req *tari_generated.GetPeersRequest, stream tari_generated.BaseNode_GetPeersServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := GetPeers(context.Background(), &tari_generated.GetPeersRequest{})
	if err != nil {
		t.Fatalf("GetPeers returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("GetPeers returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("GetPeers item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestGetPeers_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, GetPeers propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestGetPeers_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getpeers-stream")
	srv := &fakeBaseNodeServer{
		getPeersFn: func(req *tari_generated.GetPeersRequest, stream tari_generated.BaseNode_GetPeersServer) error {
			if err := stream.Send(&tari_generated.GetPeersResponse{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := GetPeers(context.Background(), &tari_generated.GetPeersRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-getpeers-stream" {
		t.Fatalf("GetPeers did not propagate the underlying error, got: %v", err)
	}
}

// TestGetMempoolTransactions_DrainsStream verifies that GetMempoolTransactions forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestGetMempoolTransactions_DrainsStream(t *testing.T) {
	want := []*tari_generated.GetMempoolTransactionsResponse{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.GetMempoolTransactionsRequest
	srv := &fakeBaseNodeServer{
		getMempoolTransactionsFn: func(req *tari_generated.GetMempoolTransactionsRequest, stream tari_generated.BaseNode_GetMempoolTransactionsServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := GetMempoolTransactions(context.Background(), &tari_generated.GetMempoolTransactionsRequest{})
	if err != nil {
		t.Fatalf("GetMempoolTransactions returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("GetMempoolTransactions returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("GetMempoolTransactions item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestGetMempoolTransactions_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, GetMempoolTransactions propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestGetMempoolTransactions_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getmempooltransactions-stream")
	srv := &fakeBaseNodeServer{
		getMempoolTransactionsFn: func(req *tari_generated.GetMempoolTransactionsRequest, stream tari_generated.BaseNode_GetMempoolTransactionsServer) error {
			if err := stream.Send(&tari_generated.GetMempoolTransactionsResponse{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := GetMempoolTransactions(context.Background(), &tari_generated.GetMempoolTransactionsRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-getmempooltransactions-stream" {
		t.Fatalf("GetMempoolTransactions did not propagate the underlying error, got: %v", err)
	}
}

// TestGetActiveValidatorNodes_DrainsStream verifies that GetActiveValidatorNodes forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestGetActiveValidatorNodes_DrainsStream(t *testing.T) {
	want := []*tari_generated.GetActiveValidatorNodesResponse{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.GetActiveValidatorNodesRequest
	srv := &fakeBaseNodeServer{
		getActiveValidatorNodesFn: func(req *tari_generated.GetActiveValidatorNodesRequest, stream tari_generated.BaseNode_GetActiveValidatorNodesServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := GetActiveValidatorNodes(context.Background(), &tari_generated.GetActiveValidatorNodesRequest{})
	if err != nil {
		t.Fatalf("GetActiveValidatorNodes returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("GetActiveValidatorNodes returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("GetActiveValidatorNodes item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestGetActiveValidatorNodes_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, GetActiveValidatorNodes propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestGetActiveValidatorNodes_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getactivevalidatornodes-stream")
	srv := &fakeBaseNodeServer{
		getActiveValidatorNodesFn: func(req *tari_generated.GetActiveValidatorNodesRequest, stream tari_generated.BaseNode_GetActiveValidatorNodesServer) error {
			if err := stream.Send(&tari_generated.GetActiveValidatorNodesResponse{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := GetActiveValidatorNodes(context.Background(), &tari_generated.GetActiveValidatorNodesRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-getactivevalidatornodes-stream" {
		t.Fatalf("GetActiveValidatorNodes did not propagate the underlying error, got: %v", err)
	}
}

// TestGetTemplateRegistrations_DrainsStream verifies that GetTemplateRegistrations forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestGetTemplateRegistrations_DrainsStream(t *testing.T) {
	want := []*tari_generated.GetTemplateRegistrationResponse{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.GetTemplateRegistrationsRequest
	srv := &fakeBaseNodeServer{
		getTemplateRegistrationsFn: func(req *tari_generated.GetTemplateRegistrationsRequest, stream tari_generated.BaseNode_GetTemplateRegistrationsServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := GetTemplateRegistrations(context.Background(), &tari_generated.GetTemplateRegistrationsRequest{})
	if err != nil {
		t.Fatalf("GetTemplateRegistrations returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("GetTemplateRegistrations returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("GetTemplateRegistrations item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestGetTemplateRegistrations_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, GetTemplateRegistrations propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestGetTemplateRegistrations_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-gettemplateregistrations-stream")
	srv := &fakeBaseNodeServer{
		getTemplateRegistrationsFn: func(req *tari_generated.GetTemplateRegistrationsRequest, stream tari_generated.BaseNode_GetTemplateRegistrationsServer) error {
			if err := stream.Send(&tari_generated.GetTemplateRegistrationResponse{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := GetTemplateRegistrations(context.Background(), &tari_generated.GetTemplateRegistrationsRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-gettemplateregistrations-stream" {
		t.Fatalf("GetTemplateRegistrations did not propagate the underlying error, got: %v", err)
	}
}

// TestGetSideChainUtxos_DrainsStream verifies that GetSideChainUtxos forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestGetSideChainUtxos_DrainsStream(t *testing.T) {
	want := []*tari_generated.GetSideChainUtxosResponse{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.GetSideChainUtxosRequest
	srv := &fakeBaseNodeServer{
		getSideChainUtxosFn: func(req *tari_generated.GetSideChainUtxosRequest, stream tari_generated.BaseNode_GetSideChainUtxosServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := GetSideChainUtxos(context.Background(), &tari_generated.GetSideChainUtxosRequest{})
	if err != nil {
		t.Fatalf("GetSideChainUtxos returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("GetSideChainUtxos returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("GetSideChainUtxos item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestGetSideChainUtxos_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, GetSideChainUtxos propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestGetSideChainUtxos_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-getsidechainutxos-stream")
	srv := &fakeBaseNodeServer{
		getSideChainUtxosFn: func(req *tari_generated.GetSideChainUtxosRequest, stream tari_generated.BaseNode_GetSideChainUtxosServer) error {
			if err := stream.Send(&tari_generated.GetSideChainUtxosResponse{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := GetSideChainUtxos(context.Background(), &tari_generated.GetSideChainUtxosRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-getsidechainutxos-stream" {
		t.Fatalf("GetSideChainUtxos did not propagate the underlying error, got: %v", err)
	}
}

// TestSearchPaymentReferences_DrainsStream verifies that SearchPaymentReferences forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestSearchPaymentReferences_DrainsStream(t *testing.T) {
	want := []*tari_generated.PaymentReferenceResponse{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.SearchPaymentReferencesRequest
	srv := &fakeBaseNodeServer{
		searchPaymentReferencesFn: func(req *tari_generated.SearchPaymentReferencesRequest, stream tari_generated.BaseNode_SearchPaymentReferencesServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := SearchPaymentReferences(context.Background(), &tari_generated.SearchPaymentReferencesRequest{})
	if err != nil {
		t.Fatalf("SearchPaymentReferences returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("SearchPaymentReferences returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("SearchPaymentReferences item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestSearchPaymentReferences_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, SearchPaymentReferences propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestSearchPaymentReferences_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-searchpaymentreferences-stream")
	srv := &fakeBaseNodeServer{
		searchPaymentReferencesFn: func(req *tari_generated.SearchPaymentReferencesRequest, stream tari_generated.BaseNode_SearchPaymentReferencesServer) error {
			if err := stream.Send(&tari_generated.PaymentReferenceResponse{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := SearchPaymentReferences(context.Background(), &tari_generated.SearchPaymentReferencesRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-searchpaymentreferences-stream" {
		t.Fatalf("SearchPaymentReferences did not propagate the underlying error, got: %v", err)
	}
}

// TestSearchPaymentReferencesViaOutputHash_DrainsStream verifies that SearchPaymentReferencesViaOutputHash forwards the request as-is and drains every
// item the server streams back into a slice, in order.
func TestSearchPaymentReferencesViaOutputHash_DrainsStream(t *testing.T) {
	want := []*tari_generated.PaymentReferenceResponse{
		{},
		{},
		{},
	}
	var gotReq *tari_generated.FetchMatchingUtxosRequest
	srv := &fakeBaseNodeServer{
		searchPaymentReferencesViaOutputHashFn: func(req *tari_generated.FetchMatchingUtxosRequest, stream tari_generated.BaseNode_SearchPaymentReferencesViaOutputHashServer) error {
			gotReq = req
			for _, item := range want {
				if err := stream.Send(item); err != nil {
					return err
				}
			}
			return nil
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := SearchPaymentReferencesViaOutputHash(context.Background(), &tari_generated.FetchMatchingUtxosRequest{})
	if err != nil {
		t.Fatalf("SearchPaymentReferencesViaOutputHash returned unexpected error: %v", err)
	}
	if gotReq == nil {
		t.Fatal("server did not receive a request")
	}
	if len(got) != len(want) {
		t.Fatalf("SearchPaymentReferencesViaOutputHash returned %d items, want %d", len(got), len(want))
	}
	for i := range want {
		if !proto.Equal(got[i], want[i]) {
			t.Fatalf("SearchPaymentReferencesViaOutputHash item %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestSearchPaymentReferencesViaOutputHash_PropagatesMidStreamError verifies that when the server errors out partway through the
// stream, SearchPaymentReferencesViaOutputHash propagates that error and discards any items already received (mirroring the existing
// GetTransactionsInBlock/GetBlockByHeight streaming wrappers' behavior).
func TestSearchPaymentReferencesViaOutputHash_PropagatesMidStreamError(t *testing.T) {
	wantErr := status.Error(codes.Internal, "boom-searchpaymentreferencesviaoutputhash-stream")
	srv := &fakeBaseNodeServer{
		searchPaymentReferencesViaOutputHashFn: func(req *tari_generated.FetchMatchingUtxosRequest, stream tari_generated.BaseNode_SearchPaymentReferencesViaOutputHashServer) error {
			if err := stream.Send(&tari_generated.PaymentReferenceResponse{}); err != nil {
				return err
			}
			return wantErr
		},
	}
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	got, err := SearchPaymentReferencesViaOutputHash(context.Background(), &tari_generated.FetchMatchingUtxosRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if got != nil {
		t.Fatalf("expected a nil result on mid-stream error, got %+v", got)
	}
	if status.Convert(err).Message() != "boom-searchpaymentreferencesviaoutputhash-stream" {
		t.Fatalf("SearchPaymentReferencesViaOutputHash did not propagate the underlying error, got: %v", err)
	}
}
