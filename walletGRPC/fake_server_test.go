package walletGRPC

import (
	"context"
	"net"
	"testing"

	"github.com/Snipa22/go-tari-grpc-lib/v3/tari_generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// startFakeWalletServer starts an in-process GRPC server on a loopback TCP port implementing the
// tari_generated.WalletServer interface via srv, and returns its listen address. The server (and its
// listener) are stopped automatically via t.Cleanup, so callers don't need to manage teardown themselves.
//
// Tests using this helper MUST run serially (no t.Parallel()): InitWalletGRPC/initWalletGRPC
// stores the resulting *grpc.ClientConn in this package's single package-level grpcConn variable,
// and nothing here (or in initWalletGRPC) closes the previous connection before installing a new
// one when a test calls InitWalletGRPC again. Each test still passes and `-race` stays clean
// because tests run sequentially by default and every fake server's address is only used by the
// test that started it, but running two of these tests concurrently would leak/overwrite grpcConn
// out from under whichever test loses the race, hence no t.Parallel() anywhere in this file.
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
// method call on that RPC falls back to the embedded UnimplementedWalletServer (returning a
// codes.Unimplemented error), which is fine since a test only ever drives the wrapper(s) it's
// exercising. Every overridden method below guards its *Fn field with a nil check before calling
// it, precisely so this fallback actually happens on an unset field instead of a nil func call
// panicking inside a grpc-go handler goroutine (which grpc-go does not recover from, crashing the
// whole test binary rather than failing one test cleanly).
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
	if f.getVersionFn == nil {
		return f.UnimplementedWalletServer.GetVersion(ctx, req)
	}
	return f.getVersionFn(ctx, req)
}

func (f *fakeWalletServer) CheckForUpdates(ctx context.Context, req *tari_generated.Empty) (*tari_generated.SoftwareUpdate, error) {
	if f.checkForUpdatesFn == nil {
		return f.UnimplementedWalletServer.CheckForUpdates(ctx, req)
	}
	return f.checkForUpdatesFn(ctx, req)
}

func (f *fakeWalletServer) Identify(ctx context.Context, req *tari_generated.GetIdentityRequest) (*tari_generated.GetIdentityResponse, error) {
	if f.identifyFn == nil {
		return f.UnimplementedWalletServer.Identify(ctx, req)
	}
	return f.identifyFn(ctx, req)
}

func (f *fakeWalletServer) GetAddress(ctx context.Context, req *tari_generated.Empty) (*tari_generated.GetAddressResponse, error) {
	if f.getAddressFn == nil {
		return f.UnimplementedWalletServer.GetAddress(ctx, req)
	}
	return f.getAddressFn(ctx, req)
}

func (f *fakeWalletServer) GetPaymentIdAddress(ctx context.Context, req *tari_generated.GetPaymentIdAddressRequest) (*tari_generated.GetCompleteAddressResponse, error) {
	if f.getPaymentIdAddressFn == nil {
		return f.UnimplementedWalletServer.GetPaymentIdAddress(ctx, req)
	}
	return f.getPaymentIdAddressFn(ctx, req)
}

func (f *fakeWalletServer) PrepareOneSidedTransactionForSigning(ctx context.Context, req *tari_generated.PrepareOneSidedTransactionForSigningRequest) (*tari_generated.PrepareOneSidedTransactionForSigningResponse, error) {
	if f.prepareOneSidedTransactionForSigningFn == nil {
		return f.UnimplementedWalletServer.PrepareOneSidedTransactionForSigning(ctx, req)
	}
	return f.prepareOneSidedTransactionForSigningFn(ctx, req)
}

func (f *fakeWalletServer) BroadcastSignedOneSidedTransaction(ctx context.Context, req *tari_generated.BroadcastSignedOneSidedTransactionRequest) (*tari_generated.BroadcastSignedOneSidedTransactionResponse, error) {
	if f.broadcastSignedOneSidedTransactionFn == nil {
		return f.UnimplementedWalletServer.BroadcastSignedOneSidedTransaction(ctx, req)
	}
	return f.broadcastSignedOneSidedTransactionFn(ctx, req)
}

func (f *fakeWalletServer) GetTransactionPayRefs(ctx context.Context, req *tari_generated.GetTransactionPayRefsRequest) (*tari_generated.GetTransactionPayRefsResponse, error) {
	if f.getTransactionPayRefsFn == nil {
		return f.UnimplementedWalletServer.GetTransactionPayRefs(ctx, req)
	}
	return f.getTransactionPayRefsFn(ctx, req)
}

func (f *fakeWalletServer) GetUnspentAmounts(ctx context.Context, req *tari_generated.Empty) (*tari_generated.GetUnspentAmountsResponse, error) {
	if f.getUnspentAmountsFn == nil {
		return f.UnimplementedWalletServer.GetUnspentAmounts(ctx, req)
	}
	return f.getUnspentAmountsFn(ctx, req)
}

func (f *fakeWalletServer) ImportUtxos(ctx context.Context, req *tari_generated.ImportUtxosRequest) (*tari_generated.ImportUtxosResponse, error) {
	if f.importUtxosFn == nil {
		return f.UnimplementedWalletServer.ImportUtxos(ctx, req)
	}
	return f.importUtxosFn(ctx, req)
}

func (f *fakeWalletServer) GetNetworkStatus(ctx context.Context, req *tari_generated.Empty) (*tari_generated.NetworkStatusResponse, error) {
	if f.getNetworkStatusFn == nil {
		return f.UnimplementedWalletServer.GetNetworkStatus(ctx, req)
	}
	return f.getNetworkStatusFn(ctx, req)
}

func (f *fakeWalletServer) GetConnectedHttpPeer(ctx context.Context, req *tari_generated.Empty) (*tari_generated.GetConnectedHttpPeerResponse, error) {
	if f.getConnectedHttpPeerFn == nil {
		return f.UnimplementedWalletServer.GetConnectedHttpPeer(ctx, req)
	}
	return f.getConnectedHttpPeerFn(ctx, req)
}

func (f *fakeWalletServer) CancelTransaction(ctx context.Context, req *tari_generated.CancelTransactionRequest) (*tari_generated.CancelTransactionResponse, error) {
	if f.cancelTransactionFn == nil {
		return f.UnimplementedWalletServer.CancelTransaction(ctx, req)
	}
	return f.cancelTransactionFn(ctx, req)
}

func (f *fakeWalletServer) SendShaAtomicSwapTransaction(ctx context.Context, req *tari_generated.SendShaAtomicSwapRequest) (*tari_generated.SendShaAtomicSwapResponse, error) {
	if f.sendShaAtomicSwapTransactionFn == nil {
		return f.UnimplementedWalletServer.SendShaAtomicSwapTransaction(ctx, req)
	}
	return f.sendShaAtomicSwapTransactionFn(ctx, req)
}

func (f *fakeWalletServer) CreateBurnTransaction(ctx context.Context, req *tari_generated.CreateBurnTransactionRequest) (*tari_generated.CreateBurnTransactionResponse, error) {
	if f.createBurnTransactionFn == nil {
		return f.UnimplementedWalletServer.CreateBurnTransaction(ctx, req)
	}
	return f.createBurnTransactionFn(ctx, req)
}

func (f *fakeWalletServer) ClaimShaAtomicSwapTransaction(ctx context.Context, req *tari_generated.ClaimShaAtomicSwapRequest) (*tari_generated.ClaimShaAtomicSwapResponse, error) {
	if f.claimShaAtomicSwapTransactionFn == nil {
		return f.UnimplementedWalletServer.ClaimShaAtomicSwapTransaction(ctx, req)
	}
	return f.claimShaAtomicSwapTransactionFn(ctx, req)
}

func (f *fakeWalletServer) ClaimHtlcRefundTransaction(ctx context.Context, req *tari_generated.ClaimHtlcRefundRequest) (*tari_generated.ClaimHtlcRefundResponse, error) {
	if f.claimHtlcRefundTransactionFn == nil {
		return f.UnimplementedWalletServer.ClaimHtlcRefundTransaction(ctx, req)
	}
	return f.claimHtlcRefundTransactionFn(ctx, req)
}

func (f *fakeWalletServer) CreateTemplateRegistration(ctx context.Context, req *tari_generated.CreateTemplateRegistrationRequest) (*tari_generated.CreateTemplateRegistrationResponse, error) {
	if f.createTemplateRegistrationFn == nil {
		return f.UnimplementedWalletServer.CreateTemplateRegistration(ctx, req)
	}
	return f.createTemplateRegistrationFn(ctx, req)
}

func (f *fakeWalletServer) SignMessage(ctx context.Context, req *tari_generated.SignMessageRequest) (*tari_generated.SignMessageResponse, error) {
	if f.signMessageFn == nil {
		return f.UnimplementedWalletServer.SignMessage(ctx, req)
	}
	return f.signMessageFn(ctx, req)
}

func (f *fakeWalletServer) ImportTransactions(ctx context.Context, req *tari_generated.ImportTransactionsRequest) (*tari_generated.ImportTransactionsResponse, error) {
	if f.importTransactionsFn == nil {
		return f.UnimplementedWalletServer.ImportTransactions(ctx, req)
	}
	return f.importTransactionsFn(ctx, req)
}

func (f *fakeWalletServer) GetAllCompletedTransactions(ctx context.Context, req *tari_generated.GetAllCompletedTransactionsRequest) (*tari_generated.GetAllCompletedTransactionsResponse, error) {
	if f.getAllCompletedTransactionsFn == nil {
		return f.UnimplementedWalletServer.GetAllCompletedTransactions(ctx, req)
	}
	return f.getAllCompletedTransactionsFn(ctx, req)
}

func (f *fakeWalletServer) GetPaymentByReference(ctx context.Context, req *tari_generated.GetPaymentByReferenceRequest) (*tari_generated.GetPaymentByReferenceResponse, error) {
	if f.getPaymentByReferenceFn == nil {
		return f.UnimplementedWalletServer.GetPaymentByReference(ctx, req)
	}
	return f.getPaymentByReferenceFn(ctx, req)
}

func (f *fakeWalletServer) GetFeeEstimate(ctx context.Context, req *tari_generated.GetFeeEstimateRequest) (*tari_generated.GetFeeEstimateResponse, error) {
	if f.getFeeEstimateFn == nil {
		return f.UnimplementedWalletServer.GetFeeEstimate(ctx, req)
	}
	return f.getFeeEstimateFn(ctx, req)
}

func (f *fakeWalletServer) GetFeePerGramStats(ctx context.Context, req *tari_generated.GetFeePerGramStatsRequest) (*tari_generated.GetFeePerGramStatsResponse, error) {
	if f.getFeePerGramStatsFn == nil {
		return f.UnimplementedWalletServer.GetFeePerGramStats(ctx, req)
	}
	return f.getFeePerGramStatsFn(ctx, req)
}

func (f *fakeWalletServer) ReplaceByFee(ctx context.Context, req *tari_generated.ReplaceByFeeRequest) (*tari_generated.ReplaceByFeeResponse, error) {
	if f.replaceByFeeFn == nil {
		return f.UnimplementedWalletServer.ReplaceByFee(ctx, req)
	}
	return f.replaceByFeeFn(ctx, req)
}

func (f *fakeWalletServer) UserPayForFee(ctx context.Context, req *tari_generated.UserPayForFeeRequest) (*tari_generated.UserPayForFeeResponse, error) {
	if f.userPayForFeeFn == nil {
		return f.UnimplementedWalletServer.UserPayForFee(ctx, req)
	}
	return f.userPayForFeeFn(ctx, req)
}

func (f *fakeWalletServer) RegisterValidatorNode(ctx context.Context, req *tari_generated.RegisterValidatorNodeRequest) (*tari_generated.RegisterValidatorNodeResponse, error) {
	if f.registerValidatorNodeFn == nil {
		return f.UnimplementedWalletServer.RegisterValidatorNode(ctx, req)
	}
	return f.registerValidatorNodeFn(ctx, req)
}

func (f *fakeWalletServer) SubmitValidatorEvictionProof(ctx context.Context, req *tari_generated.SubmitValidatorEvictionProofRequest) (*tari_generated.SubmitValidatorEvictionProofResponse, error) {
	if f.submitValidatorEvictionProofFn == nil {
		return f.UnimplementedWalletServer.SubmitValidatorEvictionProof(ctx, req)
	}
	return f.submitValidatorEvictionProofFn(ctx, req)
}

func (f *fakeWalletServer) SubmitValidatorNodeExit(ctx context.Context, req *tari_generated.SubmitValidatorNodeExitRequest) (*tari_generated.SubmitValidatorNodeExitResponse, error) {
	if f.submitValidatorNodeExitFn == nil {
		return f.UnimplementedWalletServer.SubmitValidatorNodeExit(ctx, req)
	}
	return f.submitValidatorNodeExitFn(ctx, req)
}

func (f *fakeWalletServer) StreamTransactionEvents(req *tari_generated.TransactionEventRequest, stream tari_generated.Wallet_StreamTransactionEventsServer) error {
	if f.streamTransactionEventsFn == nil {
		return f.UnimplementedWalletServer.StreamTransactionEvents(req, stream)
	}
	return f.streamTransactionEventsFn(req, stream)
}

func (f *fakeWalletServer) GetAllCompletedTransactionsStream(req *tari_generated.GetAllCompletedTransactionsRequest, stream tari_generated.Wallet_GetAllCompletedTransactionsStreamServer) error {
	if f.getAllCompletedTransactionsStreamFn == nil {
		return f.UnimplementedWalletServer.GetAllCompletedTransactionsStream(req, stream)
	}
	return f.getAllCompletedTransactionsStreamFn(req, stream)
}

// TestFakeWalletServer_UnsetFnFallsBackToUnimplemented is a regression test for
// production-readiness finding 6: an overridden method on fakeWalletServer whose corresponding
// *Fn field is left nil (e.g. a test exercising a different RPC than the one it happens to share
// a fake server instance with) must degrade to the embedded UnimplementedWalletServer's
// codes.Unimplemented error, not panic on a nil func call inside the grpc-go handler goroutine
// (which would crash this entire test binary rather than failing one test).
func TestFakeWalletServer_UnsetFnFallsBackToUnimplemented(t *testing.T) {
	srv := &fakeWalletServer{} // deliberately: every *Fn field left nil.
	addr := startFakeWalletServer(t, srv)
	if err := InitWalletGRPC(addr); err != nil {
		t.Fatalf("InitWalletGRPC: %v", err)
	}
	_, err := Identify(context.Background(), &tari_generated.GetIdentityRequest{})
	if err == nil {
		t.Fatal("expected an Unimplemented error, got nil")
	}
	if status.Code(err) != codes.Unimplemented {
		t.Fatalf("expected codes.Unimplemented, got %v (%v)", status.Code(err), err)
	}
}
