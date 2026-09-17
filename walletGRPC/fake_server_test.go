package walletGRPC

import (
	"context"
	"net"
	"testing"

	"github.com/Snipa22/go-tari-grpc-lib/v3/tari_generated"
	"google.golang.org/grpc"
)

// startFakeWalletServer starts an in-process GRPC server on a loopback TCP port implementing the
// tari_generated.WalletServer interface via srv, and returns its listen address. The server (and its
// listener) are stopped automatically via t.Cleanup, so callers don't need to manage teardown themselves.
func startFakeWalletServer(t *testing.T, srv tari_generated.WalletServer) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake wallet listener: %v", err)
	}
	grpcServer := grpc.NewServer()
	tari_generated.RegisterWalletServer(grpcServer, srv)
	go func() {
		_ = grpcServer.Serve(lis)
	}()
	t.Cleanup(func() {
		grpcServer.Stop()
		_ = lis.Close()
	})
	return lis.Addr().String()
}

// fakeWalletServer is an in-process, per-test fake implementation of the tari_generated.WalletServer
// interface. Only the RPCs under test in a given test need their corresponding *Fn field set; every other
// method falls back to the embedded UnimplementedWalletServer (returning a codes.Unimplemented error),
// which is fine since a test only ever drives the wrapper(s) it's exercising.
type fakeWalletServer struct {
	tari_generated.UnimplementedWalletServer
	getVersionFn                           func(context.Context, *tari_generated.GetVersionRequest) (*tari_generated.GetVersionResponse, error)
	checkForUpdatesFn                      func(context.Context, *tari_generated.Empty) (*tari_generated.SoftwareUpdate, error)
	identifyFn                             func(context.Context, *tari_generated.GetIdentityRequest) (*tari_generated.GetIdentityResponse, error)
	getAddressFn                           func(context.Context, *tari_generated.Empty) (*tari_generated.GetAddressResponse, error)
	getPaymentIdAddressFn                  func(context.Context, *tari_generated.GetPaymentIdAddressRequest) (*tari_generated.GetCompleteAddressResponse, error)
	prepareOneSidedTransactionForSigningFn func(context.Context, *tari_generated.PrepareOneSidedTransactionForSigningRequest) (*tari_generated.PrepareOneSidedTransactionForSigningResponse, error)
	broadcastSignedOneSidedTransactionFn   func(context.Context, *tari_generated.BroadcastSignedOneSidedTransactionRequest) (*tari_generated.BroadcastSignedOneSidedTransactionResponse, error)
	getTransactionPayRefsFn                func(context.Context, *tari_generated.GetTransactionPayRefsRequest) (*tari_generated.GetTransactionPayRefsResponse, error)
	getUnspentAmountsFn                    func(context.Context, *tari_generated.Empty) (*tari_generated.GetUnspentAmountsResponse, error)
	importUtxosFn                          func(context.Context, *tari_generated.ImportUtxosRequest) (*tari_generated.ImportUtxosResponse, error)
	getNetworkStatusFn                     func(context.Context, *tari_generated.Empty) (*tari_generated.NetworkStatusResponse, error)
	getConnectedHttpPeerFn                 func(context.Context, *tari_generated.Empty) (*tari_generated.GetConnectedHttpPeerResponse, error)
	cancelTransactionFn                    func(context.Context, *tari_generated.CancelTransactionRequest) (*tari_generated.CancelTransactionResponse, error)
	sendShaAtomicSwapTransactionFn         func(context.Context, *tari_generated.SendShaAtomicSwapRequest) (*tari_generated.SendShaAtomicSwapResponse, error)
	createBurnTransactionFn                func(context.Context, *tari_generated.CreateBurnTransactionRequest) (*tari_generated.CreateBurnTransactionResponse, error)
	claimShaAtomicSwapTransactionFn        func(context.Context, *tari_generated.ClaimShaAtomicSwapRequest) (*tari_generated.ClaimShaAtomicSwapResponse, error)
	claimHtlcRefundTransactionFn           func(context.Context, *tari_generated.ClaimHtlcRefundRequest) (*tari_generated.ClaimHtlcRefundResponse, error)
	createTemplateRegistrationFn           func(context.Context, *tari_generated.CreateTemplateRegistrationRequest) (*tari_generated.CreateTemplateRegistrationResponse, error)
	signMessageFn                          func(context.Context, *tari_generated.SignMessageRequest) (*tari_generated.SignMessageResponse, error)
	importTransactionsFn                   func(context.Context, *tari_generated.ImportTransactionsRequest) (*tari_generated.ImportTransactionsResponse, error)
	getAllCompletedTransactionsFn          func(context.Context, *tari_generated.GetAllCompletedTransactionsRequest) (*tari_generated.GetAllCompletedTransactionsResponse, error)
	getPaymentByReferenceFn                func(context.Context, *tari_generated.GetPaymentByReferenceRequest) (*tari_generated.GetPaymentByReferenceResponse, error)
	getFeeEstimateFn                       func(context.Context, *tari_generated.GetFeeEstimateRequest) (*tari_generated.GetFeeEstimateResponse, error)
	getFeePerGramStatsFn                   func(context.Context, *tari_generated.GetFeePerGramStatsRequest) (*tari_generated.GetFeePerGramStatsResponse, error)
	replaceByFeeFn                         func(context.Context, *tari_generated.ReplaceByFeeRequest) (*tari_generated.ReplaceByFeeResponse, error)
	userPayForFeeFn                        func(context.Context, *tari_generated.UserPayForFeeRequest) (*tari_generated.UserPayForFeeResponse, error)
	registerValidatorNodeFn                func(context.Context, *tari_generated.RegisterValidatorNodeRequest) (*tari_generated.RegisterValidatorNodeResponse, error)
	submitValidatorEvictionProofFn         func(context.Context, *tari_generated.SubmitValidatorEvictionProofRequest) (*tari_generated.SubmitValidatorEvictionProofResponse, error)
	submitValidatorNodeExitFn              func(context.Context, *tari_generated.SubmitValidatorNodeExitRequest) (*tari_generated.SubmitValidatorNodeExitResponse, error)
	streamTransactionEventsFn              func(*tari_generated.TransactionEventRequest, tari_generated.Wallet_StreamTransactionEventsServer) error
	getAllCompletedTransactionsStreamFn    func(*tari_generated.GetAllCompletedTransactionsRequest, tari_generated.Wallet_GetAllCompletedTransactionsStreamServer) error
}

func (f *fakeWalletServer) GetVersion(ctx context.Context, req *tari_generated.GetVersionRequest) (*tari_generated.GetVersionResponse, error) {
	return f.getVersionFn(ctx, req)
}

func (f *fakeWalletServer) CheckForUpdates(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SoftwareUpdate, error) {
	return f.checkForUpdatesFn(ctx, req)
}

func (f *fakeWalletServer) Identify(ctx context.Context, req *tari_generated.GetIdentityRequest) (*tari_generated.GetIdentityResponse, error) {
	return f.identifyFn(ctx, req)
}

func (f *fakeWalletServer) GetAddress(ctx context.Context, req *tari_generated.Empty) (*tari_generated.GetAddressResponse, error) {
	return f.getAddressFn(ctx, req)
}

func (f *fakeWalletServer) GetPaymentIdAddress(ctx context.Context, req *tari_generated.GetPaymentIdAddressRequest) (*tari_generated.GetCompleteAddressResponse, error) {
	return f.getPaymentIdAddressFn(ctx, req)
}

func (f *fakeWalletServer) PrepareOneSidedTransactionForSigning(ctx context.Context, req *tari_generated.PrepareOneSidedTransactionForSigningRequest) (*tari_generated.PrepareOneSidedTransactionForSigningResponse, error) {
	return f.prepareOneSidedTransactionForSigningFn(ctx, req)
}

func (f *fakeWalletServer) BroadcastSignedOneSidedTransaction(ctx context.Context, req *tari_generated.BroadcastSignedOneSidedTransactionRequest) (*tari_generated.BroadcastSignedOneSidedTransactionResponse, error) {
	return f.broadcastSignedOneSidedTransactionFn(ctx, req)
}

func (f *fakeWalletServer) GetTransactionPayRefs(ctx context.Context, req *tari_generated.GetTransactionPayRefsRequest) (*tari_generated.GetTransactionPayRefsResponse, error) {
	return f.getTransactionPayRefsFn(ctx, req)
}

func (f *fakeWalletServer) GetUnspentAmounts(ctx context.Context, req *tari_generated.Empty) (*tari_generated.GetUnspentAmountsResponse, error) {
	return f.getUnspentAmountsFn(ctx, req)
}

func (f *fakeWalletServer) ImportUtxos(ctx context.Context, req *tari_generated.ImportUtxosRequest) (*tari_generated.ImportUtxosResponse, error) {
	return f.importUtxosFn(ctx, req)
}

func (f *fakeWalletServer) GetNetworkStatus(ctx context.Context, req *tari_generated.Empty) (*tari_generated.NetworkStatusResponse, error) {
	return f.getNetworkStatusFn(ctx, req)
}

func (f *fakeWalletServer) GetConnectedHttpPeer(ctx context.Context, req *tari_generated.Empty) (*tari_generated.GetConnectedHttpPeerResponse, error) {
	return f.getConnectedHttpPeerFn(ctx, req)
}

func (f *fakeWalletServer) CancelTransaction(ctx context.Context, req *tari_generated.CancelTransactionRequest) (*tari_generated.CancelTransactionResponse, error) {
	return f.cancelTransactionFn(ctx, req)
}

func (f *fakeWalletServer) SendShaAtomicSwapTransaction(ctx context.Context, req *tari_generated.SendShaAtomicSwapRequest) (*tari_generated.SendShaAtomicSwapResponse, error) {
	return f.sendShaAtomicSwapTransactionFn(ctx, req)
}

func (f *fakeWalletServer) CreateBurnTransaction(ctx context.Context, req *tari_generated.CreateBurnTransactionRequest) (*tari_generated.CreateBurnTransactionResponse, error) {
	return f.createBurnTransactionFn(ctx, req)
}

func (f *fakeWalletServer) ClaimShaAtomicSwapTransaction(ctx context.Context, req *tari_generated.ClaimShaAtomicSwapRequest) (*tari_generated.ClaimShaAtomicSwapResponse, error) {
	return f.claimShaAtomicSwapTransactionFn(ctx, req)
}

func (f *fakeWalletServer) ClaimHtlcRefundTransaction(ctx context.Context, req *tari_generated.ClaimHtlcRefundRequest) (*tari_generated.ClaimHtlcRefundResponse, error) {
	return f.claimHtlcRefundTransactionFn(ctx, req)
}

func (f *fakeWalletServer) CreateTemplateRegistration(ctx context.Context, req *tari_generated.CreateTemplateRegistrationRequest) (*tari_generated.CreateTemplateRegistrationResponse, error) {
	return f.createTemplateRegistrationFn(ctx, req)
}

func (f *fakeWalletServer) SignMessage(ctx context.Context, req *tari_generated.SignMessageRequest) (*tari_generated.SignMessageResponse, error) {
	return f.signMessageFn(ctx, req)
}

func (f *fakeWalletServer) ImportTransactions(ctx context.Context, req *tari_generated.ImportTransactionsRequest) (*tari_generated.ImportTransactionsResponse, error) {
	return f.importTransactionsFn(ctx, req)
}

func (f *fakeWalletServer) GetAllCompletedTransactions(ctx context.Context, req *tari_generated.GetAllCompletedTransactionsRequest) (*tari_generated.GetAllCompletedTransactionsResponse, error) {
	return f.getAllCompletedTransactionsFn(ctx, req)
}

func (f *fakeWalletServer) GetPaymentByReference(ctx context.Context, req *tari_generated.GetPaymentByReferenceRequest) (*tari_generated.GetPaymentByReferenceResponse, error) {
	return f.getPaymentByReferenceFn(ctx, req)
}

func (f *fakeWalletServer) GetFeeEstimate(ctx context.Context, req *tari_generated.GetFeeEstimateRequest) (*tari_generated.GetFeeEstimateResponse, error) {
	return f.getFeeEstimateFn(ctx, req)
}

func (f *fakeWalletServer) GetFeePerGramStats(ctx context.Context, req *tari_generated.GetFeePerGramStatsRequest) (*tari_generated.GetFeePerGramStatsResponse, error) {
	return f.getFeePerGramStatsFn(ctx, req)
}

func (f *fakeWalletServer) ReplaceByFee(ctx context.Context, req *tari_generated.ReplaceByFeeRequest) (*tari_generated.ReplaceByFeeResponse, error) {
	return f.replaceByFeeFn(ctx, req)
}

func (f *fakeWalletServer) UserPayForFee(ctx context.Context, req *tari_generated.UserPayForFeeRequest) (*tari_generated.UserPayForFeeResponse, error) {
	return f.userPayForFeeFn(ctx, req)
}

func (f *fakeWalletServer) RegisterValidatorNode(ctx context.Context, req *tari_generated.RegisterValidatorNodeRequest) (*tari_generated.RegisterValidatorNodeResponse, error) {
	return f.registerValidatorNodeFn(ctx, req)
}

func (f *fakeWalletServer) SubmitValidatorEvictionProof(ctx context.Context, req *tari_generated.SubmitValidatorEvictionProofRequest) (*tari_generated.SubmitValidatorEvictionProofResponse, error) {
	return f.submitValidatorEvictionProofFn(ctx, req)
}

func (f *fakeWalletServer) SubmitValidatorNodeExit(ctx context.Context, req *tari_generated.SubmitValidatorNodeExitRequest) (*tari_generated.SubmitValidatorNodeExitResponse, error) {
	return f.submitValidatorNodeExitFn(ctx, req)
}

func (f *fakeWalletServer) StreamTransactionEvents(req *tari_generated.TransactionEventRequest, stream tari_generated.Wallet_StreamTransactionEventsServer) error {
	return f.streamTransactionEventsFn(req, stream)
}

func (f *fakeWalletServer) GetAllCompletedTransactionsStream(req *tari_generated.GetAllCompletedTransactionsRequest, stream tari_generated.Wallet_GetAllCompletedTransactionsStreamServer) error {
	return f.getAllCompletedTransactionsStreamFn(req, stream)
}
