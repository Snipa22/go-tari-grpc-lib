package nodeGRPC

import (
	"context"
	"net"
	"testing"

	"github.com/Snipa22/go-tari-grpc-lib/v3/tari_generated"
	"google.golang.org/grpc"
)

// startFakeBaseNodeServer starts an in-process GRPC server on a loopback TCP port implementing the
// tari_generated.BaseNodeServer interface via srv, and returns its listen address. The server (and its
// listener) are stopped automatically via t.Cleanup, so callers don't need to manage teardown themselves.
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
// method falls back to the embedded UnimplementedBaseNodeServer (returning a codes.Unimplemented error),
// which is fine since a test only ever drives the wrapper(s) it's exercising.
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
	return f.getBlockTimingFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetConstants(ctx context.Context, req *tari_generated.BlockHeight) (*tari_generated.ConsensusConstants, error) {
	return f.getConstantsFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetBlockSize(ctx context.Context, req *tari_generated.BlockGroupRequest) (*tari_generated.BlockGroupResponse, error) {
	return f.getBlockSizeFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetBlockFees(ctx context.Context, req *tari_generated.BlockGroupRequest) (*tari_generated.BlockGroupResponse, error) {
	return f.getBlockFeesFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetVersion(ctx context.Context, req *tari_generated.Empty) (*tari_generated.BaseNodeGetVersionResponse, error) {
	return f.getVersionFn(ctx, req)
}

func (f *fakeBaseNodeServer) CheckForUpdates(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SoftwareUpdate, error) {
	return f.checkForUpdatesFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetNewBlockBlob(ctx context.Context, req *tari_generated.NewBlockTemplate) (*tari_generated.GetNewBlockBlobResult, error) {
	return f.getNewBlockBlobFn(ctx, req)
}

func (f *fakeBaseNodeServer) SubmitBlockBlob(ctx context.Context, req *tari_generated.BlockBlobRequest) (*tari_generated.SubmitBlockResponse, error) {
	return f.submitBlockBlobFn(ctx, req)
}

func (f *fakeBaseNodeServer) SubmitTransaction(ctx context.Context, req *tari_generated.SubmitTransactionRequest) (*tari_generated.SubmitTransactionResponse, error) {
	return f.submitTransactionFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetSyncInfo(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SyncInfoResponse, error) {
	return f.getSyncInfoFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetSyncProgress(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SyncProgressResponse, error) {
	return f.getSyncProgressFn(ctx, req)
}

func (f *fakeBaseNodeServer) TransactionState(ctx context.Context, req *tari_generated.TransactionStateRequest) (*tari_generated.TransactionStateResponse, error) {
	return f.transactionStateFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetNetworkStatus(ctx context.Context, req *tari_generated.Empty) (*tari_generated.NetworkStatusResponse, error) {
	return f.getNetworkStatusFn(ctx, req)
}

func (f *fakeBaseNodeServer) ListConnectedPeers(ctx context.Context, req *tari_generated.Empty) (*tari_generated.ListConnectedPeersResponse, error) {
	return f.listConnectedPeersFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetMempoolStats(ctx context.Context, req *tari_generated.Empty) (*tari_generated.MempoolStatsResponse, error) {
	return f.getMempoolStatsFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetValidatorNodeChanges(ctx context.Context, req *tari_generated.GetValidatorNodeChangesRequest) (*tari_generated.GetValidatorNodeChangesResponse, error) {
	return f.getValidatorNodeChangesFn(ctx, req)
}

func (f *fakeBaseNodeServer) GetShardKey(ctx context.Context, req *tari_generated.GetShardKeyRequest) (*tari_generated.GetShardKeyResponse, error) {
	return f.getShardKeyFn(ctx, req)
}

func (f *fakeBaseNodeServer) ListHeaders(req *tari_generated.ListHeadersRequest, stream tari_generated.BaseNode_ListHeadersServer) error {
	return f.listHeadersFn(req, stream)
}

func (f *fakeBaseNodeServer) GetTokensInCirculation(req *tari_generated.GetBlocksRequest, stream tari_generated.BaseNode_GetTokensInCirculationServer) error {
	return f.getTokensInCirculationFn(req, stream)
}

func (f *fakeBaseNodeServer) SearchKernels(req *tari_generated.SearchKernelsRequest, stream tari_generated.BaseNode_SearchKernelsServer) error {
	return f.searchKernelsFn(req, stream)
}

func (f *fakeBaseNodeServer) SearchUtxos(req *tari_generated.SearchUtxosRequest, stream tari_generated.BaseNode_SearchUtxosServer) error {
	return f.searchUtxosFn(req, stream)
}

func (f *fakeBaseNodeServer) FetchMatchingUtxos(req *tari_generated.FetchMatchingUtxosRequest, stream tari_generated.BaseNode_FetchMatchingUtxosServer) error {
	return f.fetchMatchingUtxosFn(req, stream)
}

func (f *fakeBaseNodeServer) GetPeers(req *tari_generated.GetPeersRequest, stream tari_generated.BaseNode_GetPeersServer) error {
	return f.getPeersFn(req, stream)
}

func (f *fakeBaseNodeServer) GetMempoolTransactions(req *tari_generated.GetMempoolTransactionsRequest, stream tari_generated.BaseNode_GetMempoolTransactionsServer) error {
	return f.getMempoolTransactionsFn(req, stream)
}

func (f *fakeBaseNodeServer) GetActiveValidatorNodes(req *tari_generated.GetActiveValidatorNodesRequest, stream tari_generated.BaseNode_GetActiveValidatorNodesServer) error {
	return f.getActiveValidatorNodesFn(req, stream)
}

func (f *fakeBaseNodeServer) GetTemplateRegistrations(req *tari_generated.GetTemplateRegistrationsRequest, stream tari_generated.BaseNode_GetTemplateRegistrationsServer) error {
	return f.getTemplateRegistrationsFn(req, stream)
}

func (f *fakeBaseNodeServer) GetSideChainUtxos(req *tari_generated.GetSideChainUtxosRequest, stream tari_generated.BaseNode_GetSideChainUtxosServer) error {
	return f.getSideChainUtxosFn(req, stream)
}

func (f *fakeBaseNodeServer) SearchPaymentReferences(req *tari_generated.SearchPaymentReferencesRequest, stream tari_generated.BaseNode_SearchPaymentReferencesServer) error {
	return f.searchPaymentReferencesFn(req, stream)
}

func (f *fakeBaseNodeServer) SearchPaymentReferencesViaOutputHash(req *tari_generated.FetchMatchingUtxosRequest, stream tari_generated.BaseNode_SearchPaymentReferencesViaOutputHashServer) error {
	return f.searchPaymentReferencesViaOutputHashFn(req, stream)
}
