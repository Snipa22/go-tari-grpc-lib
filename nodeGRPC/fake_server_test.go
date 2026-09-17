package nodeGRPC

import (
	"context"
	"net"
	"testing"

	"github.com/Snipa22/go-tari-grpc-lib/v3/tari_generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// startFakeBaseNodeServer starts an in-process GRPC server on a loopback TCP port implementing the
// tari_generated.BaseNodeServer interface via srv, and returns its listen address. The server (and its
// listener) are stopped automatically via t.Cleanup, so callers don't need to manage teardown themselves.
//
// Tests using this helper MUST run serially (no t.Parallel()): InitNodeGRPC/initNodeGRPC stores
// the resulting *grpc.ClientConn in this package's single package-level grpcConn variable, and
// nothing here (or in initNodeGRPC) closes the previous connection before installing a new one
// when a test calls InitNodeGRPC again. Each test still passes and `-race` stays clean because
// tests run sequentially by default and every fake server's address is only used by the test that
// started it, but running two of these tests concurrently would leak/overwrite grpcConn out from
// under whichever test loses the race, hence no t.Parallel() anywhere in this file.
func startFakeBaseNodeServer(t *testing.T, srv tari_generated.BaseNodeServer) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake base node listener: %v", err)
	}
	grpcServer := grpc.NewServer()
	tari_generated.RegisterBaseNodeServer(grpcServer, srv)
	go func() {
		_ = grpcServer.Serve(lis)
	}()
	t.Cleanup(func() {
		grpcServer.Stop()
		_ = lis.Close()
	})
	return lis.Addr().String()
}

// fakeBaseNodeServer is an in-process, per-test fake implementation of the tari_generated.BaseNodeServer
// interface. Only the RPCs under test in a given test need their corresponding *Fn field set; every other
// method call on that RPC falls back to the embedded UnimplementedBaseNodeServer (returning a
// codes.Unimplemented error), which is fine since a test only ever drives the wrapper(s) it's
// exercising. Every overridden method below guards its *Fn field with a nil check before calling
// it, precisely so this fallback actually happens on an unset field instead of a nil func call
// panicking inside a grpc-go handler goroutine (which grpc-go does not recover from, crashing the
// whole test binary rather than failing one test cleanly).
type fakeBaseNodeServer struct {
	tari_generated.UnimplementedBaseNodeServer
	getBlockTimingFn                       func(context.Context, *tari_generated.HeightRequest) (*tari_generated.BlockTimingResponse, error)
	getConstantsFn                         func(context.Context, *tari_generated.BlockHeight) (*tari_generated.ConsensusConstants, error)
	getBlockSizeFn                         func(context.Context, *tari_generated.BlockGroupRequest) (*tari_generated.BlockGroupResponse, error)
	getBlockFeesFn                         func(context.Context, *tari_generated.BlockGroupRequest) (*tari_generated.BlockGroupResponse, error)
	getVersionFn                           func(context.Context, *tari_generated.Empty) (*tari_generated.BaseNodeGetVersionResponse, error)
	checkForUpdatesFn                      func(context.Context, *tari_generated.Empty) (*tari_generated.SoftwareUpdate, error)
	getNewBlockBlobFn                      func(context.Context, *tari_generated.NewBlockTemplate) (*tari_generated.GetNewBlockBlobResult, error)
	submitBlockBlobFn                      func(context.Context, *tari_generated.BlockBlobRequest) (*tari_generated.SubmitBlockResponse, error)
	submitTransactionFn                    func(context.Context, *tari_generated.SubmitTransactionRequest) (*tari_generated.SubmitTransactionResponse, error)
	getSyncInfoFn                          func(context.Context, *tari_generated.Empty) (*tari_generated.SyncInfoResponse, error)
	getSyncProgressFn                      func(context.Context, *tari_generated.Empty) (*tari_generated.SyncProgressResponse, error)
	transactionStateFn                     func(context.Context, *tari_generated.TransactionStateRequest) (*tari_generated.TransactionStateResponse, error)
	getNetworkStatusFn                     func(context.Context, *tari_generated.Empty) (*tari_generated.NetworkStatusResponse, error)
	listConnectedPeersFn                   func(context.Context, *tari_generated.Empty) (*tari_generated.ListConnectedPeersResponse, error)
	getMempoolStatsFn                      func(context.Context, *tari_generated.Empty) (*tari_generated.MempoolStatsResponse, error)
	getValidatorNodeChangesFn              func(context.Context, *tari_generated.GetValidatorNodeChangesRequest) (*tari_generated.GetValidatorNodeChangesResponse, error)
	getShardKeyFn                          func(context.Context, *tari_generated.GetShardKeyRequest) (*tari_generated.GetShardKeyResponse, error)
	listHeadersFn                          func(*tari_generated.ListHeadersRequest, tari_generated.BaseNode_ListHeadersServer) error
	getTokensInCirculationFn               func(*tari_generated.GetBlocksRequest, tari_generated.BaseNode_GetTokensInCirculationServer) error
	searchKernelsFn                        func(*tari_generated.SearchKernelsRequest, tari_generated.BaseNode_SearchKernelsServer) error
	searchUtxosFn                          func(*tari_generated.SearchUtxosRequest, tari_generated.BaseNode_SearchUtxosServer) error
	fetchMatchingUtxosFn                   func(*tari_generated.FetchMatchingUtxosRequest, tari_generated.BaseNode_FetchMatchingUtxosServer) error
	getPeersFn                             func(*tari_generated.GetPeersRequest, tari_generated.BaseNode_GetPeersServer) error
	getMempoolTransactionsFn               func(*tari_generated.GetMempoolTransactionsRequest, tari_generated.BaseNode_GetMempoolTransactionsServer) error
	getActiveValidatorNodesFn              func(*tari_generated.GetActiveValidatorNodesRequest, tari_generated.BaseNode_GetActiveValidatorNodesServer) error
	getTemplateRegistrationsFn             func(*tari_generated.GetTemplateRegistrationsRequest, tari_generated.BaseNode_GetTemplateRegistrationsServer) error
	getSideChainUtxosFn                    func(*tari_generated.GetSideChainUtxosRequest, tari_generated.BaseNode_GetSideChainUtxosServer) error
	searchPaymentReferencesFn              func(*tari_generated.SearchPaymentReferencesRequest, tari_generated.BaseNode_SearchPaymentReferencesServer) error
	searchPaymentReferencesViaOutputHashFn func(*tari_generated.FetchMatchingUtxosRequest, tari_generated.BaseNode_SearchPaymentReferencesViaOutputHashServer) error
}

func (f *fakeBaseNodeServer) GetBlockTiming(ctx context.Context, req *tari_generated.HeightRequest) (*tari_generated.BlockTimingResponse, error) {
	if f.getBlockTimingFn == nil {
		return f.UnimplementedBaseNodeServer.GetBlockTiming(ctx, req)
	}
	return f.getBlockTimingFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetConstants(ctx context.Context, req *tari_generated.BlockHeight) (*tari_generated.ConsensusConstants, error) {
	if f.getConstantsFn == nil {
		return f.UnimplementedBaseNodeServer.GetConstants(ctx, req)
	}
	return f.getConstantsFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetBlockSize(ctx context.Context, req *tari_generated.BlockGroupRequest) (*tari_generated.BlockGroupResponse, error) {
	if f.getBlockSizeFn == nil {
		return f.UnimplementedBaseNodeServer.GetBlockSize(ctx, req)
	}
	return f.getBlockSizeFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetBlockFees(ctx context.Context, req *tari_generated.BlockGroupRequest) (*tari_generated.BlockGroupResponse, error) {
	if f.getBlockFeesFn == nil {
		return f.UnimplementedBaseNodeServer.GetBlockFees(ctx, req)
	}
	return f.getBlockFeesFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetVersion(ctx context.Context, req *tari_generated.Empty) (*tari_generated.BaseNodeGetVersionResponse, error) {
	if f.getVersionFn == nil {
		return f.UnimplementedBaseNodeServer.GetVersion(ctx, req)
	}
	return f.getVersionFn(ctx, req)
}

func (f *fakeBaseNodeServer) CheckForUpdates(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SoftwareUpdate, error) {
	if f.checkForUpdatesFn == nil {
		return f.UnimplementedBaseNodeServer.CheckForUpdates(ctx, req)
	}
	return f.checkForUpdatesFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetNewBlockBlob(ctx context.Context, req *tari_generated.NewBlockTemplate) (*tari_generated.GetNewBlockBlobResult, error) {
	if f.getNewBlockBlobFn == nil {
		return f.UnimplementedBaseNodeServer.GetNewBlockBlob(ctx, req)
	}
	return f.getNewBlockBlobFn(ctx, req)
}

func (f *fakeBaseNodeServer) SubmitBlockBlob(ctx context.Context, req *tari_generated.BlockBlobRequest) (*tari_generated.SubmitBlockResponse, error) {
	if f.submitBlockBlobFn == nil {
		return f.UnimplementedBaseNodeServer.SubmitBlockBlob(ctx, req)
	}
	return f.submitBlockBlobFn(ctx, req)
}

func (f *fakeBaseNodeServer) SubmitTransaction(ctx context.Context, req *tari_generated.SubmitTransactionRequest) (*tari_generated.SubmitTransactionResponse, error) {
	if f.submitTransactionFn == nil {
		return f.UnimplementedBaseNodeServer.SubmitTransaction(ctx, req)
	}
	return f.submitTransactionFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetSyncInfo(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SyncInfoResponse, error) {
	if f.getSyncInfoFn == nil {
		return f.UnimplementedBaseNodeServer.GetSyncInfo(ctx, req)
	}
	return f.getSyncInfoFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetSyncProgress(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SyncProgressResponse, error) {
	if f.getSyncProgressFn == nil {
		return f.UnimplementedBaseNodeServer.GetSyncProgress(ctx, req)
	}
	return f.getSyncProgressFn(ctx, req)
}

func (f *fakeBaseNodeServer) TransactionState(ctx context.Context, req *tari_generated.TransactionStateRequest) (*tari_generated.TransactionStateResponse, error) {
	if f.transactionStateFn == nil {
		return f.UnimplementedBaseNodeServer.TransactionState(ctx, req)
	}
	return f.transactionStateFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetNetworkStatus(ctx context.Context, req *tari_generated.Empty) (*tari_generated.NetworkStatusResponse, error) {
	if f.getNetworkStatusFn == nil {
		return f.UnimplementedBaseNodeServer.GetNetworkStatus(ctx, req)
	}
	return f.getNetworkStatusFn(ctx, req)
}

func (f *fakeBaseNodeServer) ListConnectedPeers(ctx context.Context, req *tari_generated.Empty) (*tari_generated.ListConnectedPeersResponse, error) {
	if f.listConnectedPeersFn == nil {
		return f.UnimplementedBaseNodeServer.ListConnectedPeers(ctx, req)
	}
	return f.listConnectedPeersFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetMempoolStats(ctx context.Context, req *tari_generated.Empty) (*tari_generated.MempoolStatsResponse, error) {
	if f.getMempoolStatsFn == nil {
		return f.UnimplementedBaseNodeServer.GetMempoolStats(ctx, req)
	}
	return f.getMempoolStatsFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetValidatorNodeChanges(ctx context.Context, req *tari_generated.GetValidatorNodeChangesRequest) (*tari_generated.GetValidatorNodeChangesResponse, error) {
	if f.getValidatorNodeChangesFn == nil {
		return f.UnimplementedBaseNodeServer.GetValidatorNodeChanges(ctx, req)
	}
	return f.getValidatorNodeChangesFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetShardKey(ctx context.Context, req *tari_generated.GetShardKeyRequest) (*tari_generated.GetShardKeyResponse, error) {
	if f.getShardKeyFn == nil {
		return f.UnimplementedBaseNodeServer.GetShardKey(ctx, req)
	}
	return f.getShardKeyFn(ctx, req)
}

func (f *fakeBaseNodeServer) ListHeaders(req *tari_generated.ListHeadersRequest, stream tari_generated.BaseNode_ListHeadersServer) error {
	if f.listHeadersFn == nil {
		return f.UnimplementedBaseNodeServer.ListHeaders(req, stream)
	}
	return f.listHeadersFn(req, stream)
}

func (f *fakeBaseNodeServer) GetTokensInCirculation(req *tari_generated.GetBlocksRequest, stream tari_generated.BaseNode_GetTokensInCirculationServer) error {
	if f.getTokensInCirculationFn == nil {
		return f.UnimplementedBaseNodeServer.GetTokensInCirculation(req, stream)
	}
	return f.getTokensInCirculationFn(req, stream)
}

func (f *fakeBaseNodeServer) SearchKernels(req *tari_generated.SearchKernelsRequest, stream tari_generated.BaseNode_SearchKernelsServer) error {
	if f.searchKernelsFn == nil {
		return f.UnimplementedBaseNodeServer.SearchKernels(req, stream)
	}
	return f.searchKernelsFn(req, stream)
}

func (f *fakeBaseNodeServer) SearchUtxos(req *tari_generated.SearchUtxosRequest, stream tari_generated.BaseNode_SearchUtxosServer) error {
	if f.searchUtxosFn == nil {
		return f.UnimplementedBaseNodeServer.SearchUtxos(req, stream)
	}
	return f.searchUtxosFn(req, stream)
}

func (f *fakeBaseNodeServer) FetchMatchingUtxos(req *tari_generated.FetchMatchingUtxosRequest, stream tari_generated.BaseNode_FetchMatchingUtxosServer) error {
	if f.fetchMatchingUtxosFn == nil {
		return f.UnimplementedBaseNodeServer.FetchMatchingUtxos(req, stream)
	}
	return f.fetchMatchingUtxosFn(req, stream)
}

func (f *fakeBaseNodeServer) GetPeers(req *tari_generated.GetPeersRequest, stream tari_generated.BaseNode_GetPeersServer) error {
	if f.getPeersFn == nil {
		return f.UnimplementedBaseNodeServer.GetPeers(req, stream)
	}
	return f.getPeersFn(req, stream)
}

func (f *fakeBaseNodeServer) GetMempoolTransactions(req *tari_generated.GetMempoolTransactionsRequest, stream tari_generated.BaseNode_GetMempoolTransactionsServer) error {
	if f.getMempoolTransactionsFn == nil {
		return f.UnimplementedBaseNodeServer.GetMempoolTransactions(req, stream)
	}
	return f.getMempoolTransactionsFn(req, stream)
}

func (f *fakeBaseNodeServer) GetActiveValidatorNodes(req *tari_generated.GetActiveValidatorNodesRequest, stream tari_generated.BaseNode_GetActiveValidatorNodesServer) error {
	if f.getActiveValidatorNodesFn == nil {
		return f.UnimplementedBaseNodeServer.GetActiveValidatorNodes(req, stream)
	}
	return f.getActiveValidatorNodesFn(req, stream)
}

func (f *fakeBaseNodeServer) GetTemplateRegistrations(req *tari_generated.GetTemplateRegistrationsRequest, stream tari_generated.BaseNode_GetTemplateRegistrationsServer) error {
	if f.getTemplateRegistrationsFn == nil {
		return f.UnimplementedBaseNodeServer.GetTemplateRegistrations(req, stream)
	}
	return f.getTemplateRegistrationsFn(req, stream)
}

func (f *fakeBaseNodeServer) GetSideChainUtxos(req *tari_generated.GetSideChainUtxosRequest, stream tari_generated.BaseNode_GetSideChainUtxosServer) error {
	if f.getSideChainUtxosFn == nil {
		return f.UnimplementedBaseNodeServer.GetSideChainUtxos(req, stream)
	}
	return f.getSideChainUtxosFn(req, stream)
}

func (f *fakeBaseNodeServer) SearchPaymentReferences(req *tari_generated.SearchPaymentReferencesRequest, stream tari_generated.BaseNode_SearchPaymentReferencesServer) error {
	if f.searchPaymentReferencesFn == nil {
		return f.UnimplementedBaseNodeServer.SearchPaymentReferences(req, stream)
	}
	return f.searchPaymentReferencesFn(req, stream)
}

func (f *fakeBaseNodeServer) SearchPaymentReferencesViaOutputHash(req *tari_generated.FetchMatchingUtxosRequest, stream tari_generated.BaseNode_SearchPaymentReferencesViaOutputHashServer) error {
	if f.searchPaymentReferencesViaOutputHashFn == nil {
		return f.UnimplementedBaseNodeServer.SearchPaymentReferencesViaOutputHash(req, stream)
	}
	return f.searchPaymentReferencesViaOutputHashFn(req, stream)
}

// TestFakeBaseNodeServer_UnsetFnFallsBackToUnimplemented is a regression test for
// production-readiness finding 6: an overridden method on fakeBaseNodeServer whose corresponding
// *Fn field is left nil (e.g. a test exercising a different RPC than the one it happens to share
// a fake server instance with) must degrade to the embedded UnimplementedBaseNodeServer's
// codes.Unimplemented error, not panic on a nil func call inside the grpc-go handler goroutine
// (which would crash this entire test binary rather than failing one test).
func TestFakeBaseNodeServer_UnsetFnFallsBackToUnimplemented(t *testing.T) {
	srv := &fakeBaseNodeServer{} // deliberately: every *Fn field left nil.
	addr := startFakeBaseNodeServer(t, srv)
	if err := InitNodeGRPC(addr); err != nil {
		t.Fatalf("InitNodeGRPC: %v", err)
	}
	_, err := GetBlockTiming(context.Background(), &tari_generated.HeightRequest{})
	if err == nil {
		t.Fatal("expected an Unimplemented error, got nil")
	}
	if status.Code(err) != codes.Unimplemented {
		t.Fatalf("expected codes.Unimplemented, got %v (%v)", status.Code(err), err)
	}
}
