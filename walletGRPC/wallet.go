package walletGRPC

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/Snipa22/go-tari-grpc-lib/v3/tari_generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

// DefaultFeePerGram is the default per-gram transaction fee (in uT) used when callers don't
// supply their own fee rate. This was previously duplicated as a hardcoded literal across
// walletGRPC and several cmd/* binaries; it's now defined once here so it can be referenced
// and, where needed, overridden consistently.
const DefaultFeePerGram uint64 = 5

// dialTimeout bounds how long a single connection attempt to the wallet GRPC endpoint may take
// before it's abandoned, via grpc.ConnectParams' MinConnectTimeout.
const dialTimeout = 5 * time.Second

// maxRecvMsgSize is the maximum message size (in bytes) this package will accept on ANY Wallet
// RPC response, unary or streaming. It's set well above grpc-go's 4MiB default because several
// non-streamed wallet RPCs return an unbounded collection in a single response — e.g.
// GetBlockHeightTransactionsResponse.Transactions is a repeated TransactionInfo with no
// server-side page limit, so a block containing many transactions could plausibly exceed the
// default on its own — and the streamed RPCs return the same TransactionInfo/GetCompletedTransactionsResponse
// shapes one at a time, so applying the raised limit uniformly costs nothing there either. This
// mirrors the identical fix in nodeGRPC's dialOptions() (production-readiness finding 1), applied
// here defensively per the same finding's guidance to audit walletGRPC's own streaming/paginated
// RPCs (GetCompletedTransactions, GetAllCompletedTransactionsStream, GetBlockHeightTransactions)
// for the same class of bug.
const maxRecvMsgSize = 16 * 1024 * 1024

// connMu guards grpcWalletAddress/grpcConn below. Without it, concurrent Init calls race on both
// reads and writes of the package-level connection state.
var connMu sync.RWMutex
var grpcWalletAddress string
var grpcConn *grpc.ClientConn

// ErrNotInitialized is returned by every Wallet RPC wrapper in this package when it's called
// before InitWalletGRPC/InitWalletGRPCSecure has successfully established a connection. Without
// this, tari_generated.NewWalletClient(getConn()) would be constructed from a nil *grpc.ClientConn
// and the subsequent RPC call would panic with a nil-pointer dereference instead of returning a
// handleable error (production-readiness finding 4).
var ErrNotInitialized = errors.New("walletGRPC: not initialized, call InitWalletGRPC or InitWalletGRPCSecure first")

// buildTransportCredentials returns the transport credentials used to dial the wallet.
// Exposed as its own function so it's directly unit-testable without needing a live server.
func buildTransportCredentials(secure bool, tlsConfig *tls.Config) credentials.TransportCredentials {
	if secure {
		return credentials.NewTLS(tlsConfig)
	}
	return insecure.NewCredentials()
}

// dialOptions returns the common set of grpc.DialOption values used for every Init* variant,
// given the transport credentials to use, plus any caller-supplied extra options appended after
// the defaults. extra is exposed all the way up through InitWalletGRPC/InitWalletGRPCSecure so
// consumers can attach their own instrumentation (grpc.WithChainUnaryInterceptor,
// grpc.WithChainStreamInterceptor, grpc.WithStatsHandler — e.g. otelgrpc or a Prometheus
// interceptor) to the package-level connection without this package needing to depend on any of
// those integrations itself (production-readiness finding 8).
func dialOptions(transportCreds credentials.TransportCredentials, extra ...grpc.DialOption) []grpc.DialOption {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: dialTimeout}),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxRecvMsgSize)),
	}
	return append(opts, extra...)
}

// InitWalletGRPC opens a PLAINTEXT (insecure) GRPC connection to a Tari wallet at the given
// address. This transport is intended for loopback/trusted-network use only (e.g. a wallet
// running on the same host or within a private network you control) — for anything crossing an
// untrusted network, use InitWalletGRPCSecure instead.
//
// extra, if supplied, is appended to the default dial options (after them, so an extra option of
// the same kind takes precedence per grpc-go's normal last-one-wins semantics) — e.g. to attach a
// unary/stream interceptor or stats handler for instrumentation. Existing callers passing no extra
// arguments are unaffected.
func InitWalletGRPC(walletAddress string, extra ...grpc.DialOption) error {
	return initWalletGRPC(walletAddress, buildTransportCredentials(false, nil), extra...)
}

// InitWalletGRPCSecure opens a TLS-secured GRPC connection to a Tari wallet at the given
// address, using the supplied tls.Config for the transport credentials. Use this for any
// connection that crosses an untrusted network. Pass nil for tlsConfig to use Go's default TLS
// configuration.
//
// extra, if supplied, is appended to the default dial options (after them, so an extra option of
// the same kind takes precedence per grpc-go's normal last-one-wins semantics) — e.g. to attach a
// unary/stream interceptor or stats handler for instrumentation. Existing callers passing no extra
// arguments are unaffected.
func InitWalletGRPCSecure(walletAddress string, tlsConfig *tls.Config, extra ...grpc.DialOption) error {
	return initWalletGRPC(walletAddress, buildTransportCredentials(true, tlsConfig), extra...)
}

// newClient is the function used to construct the underlying GRPC connection. It's a
// package-level var (defaulting to grpc.NewClient) so tests can substitute a failing stub to
// verify that initWalletGRPC actually propagates the error instead of discarding it (finding 2) —
// grpc.NewClient itself validates most targets lazily at dial time, so it can't reliably be made
// to fail synchronously from a test using only real target strings.
var newClient = grpc.NewClient

func initWalletGRPC(walletAddress string, transportCreds credentials.TransportCredentials, extra ...grpc.DialOption) error {
	conn, err := newClient(walletAddress, dialOptions(transportCreds, extra...)...)
	if err != nil {
		return err
	}
	connMu.Lock()
	defer connMu.Unlock()
	grpcWalletAddress = walletAddress
	grpcConn = conn
	return nil
}

// getConn returns the current GRPC connection in a race-free way.
func getConn() *grpc.ClientConn {
	connMu.RLock()
	defer connMu.RUnlock()
	return grpcConn
}

// walletClient returns a tari_generated.WalletClient bound to the current connection, or
// ErrNotInitialized if InitWalletGRPC/InitWalletGRPCSecure hasn't successfully run yet. Every RPC
// wrapper in this file goes through this helper instead of calling
// tari_generated.NewWalletClient(getConn()) directly, so the nil-conn-panic bug (finding 4) is
// fixed once, centrally, for all ~80 call sites instead of needing a nil check hand-added to each
// wrapper individually.
func walletClient() (tari_generated.WalletClient, error) {
	conn := getConn()
	if conn == nil {
		return nil, ErrNotInitialized
	}
	return tari_generated.NewWalletClient(conn), nil
}

// SendTransactions sends the transactions to the wallet
func SendTransactions(ctx context.Context, transactions []*tari_generated.PaymentRecipient) (*tari_generated.TransferResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.Transfer(ctx, &tari_generated.TransferRequest{
		Recipients: transactions,
	})
}

// GetTransactionsInBlock will return the top of the wallet if it's called with 0, otherwise it pushes the height to
// the GRPC call, though this doesn't seem to actually do anything.  No sorting/order/etc is guaranteed, so callers
// need to parse, cache etc.
func GetTransactionsInBlock(ctx context.Context, blockHeight uint64) ([]*tari_generated.TransactionInfo, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	var completedTxnsClient tari_generated.Wallet_GetCompletedTransactionsClient
	if blockHeight == 0 {
		completedTxnsClient, err = client.GetCompletedTransactions(ctx, nil)
	} else {
		getBlockHeightTxns, err := client.GetBlockHeightTransactions(ctx, &tari_generated.GetBlockHeightTransactionsRequest{
			BlockHeight: blockHeight,
		})
		if err != nil {
			return nil, err
		}
		return getBlockHeightTxns.Transactions, nil
	}
	if err != nil {
		return nil, err
	}

	resp := make([]*tari_generated.TransactionInfo, 0)
	for {
		txnResp, err := completedTxnsClient.Recv()
		if err != nil {
			if err == io.EOF {
				return resp, nil
			}
			return nil, err
		}
		resp = append(resp, txnResp.Transaction)
	}
}

// buildCoinSplitRequest builds the CoinSplitRequest payload used by SubmitCoinSplitRequest. If
// feePerGram is 0, DefaultFeePerGram is used instead. Exposed as its own function so the request
// construction is directly unit-testable without a live GRPC server.
func buildCoinSplitRequest(splitAmt int, numSplits int, feePerGram uint64) *tari_generated.CoinSplitRequest {
	if feePerGram == 0 {
		feePerGram = DefaultFeePerGram
	}
	return &tari_generated.CoinSplitRequest{
		AmountPerSplit: uint64(splitAmt),
		SplitCount:     uint64(numSplits),
		FeePerGram:     feePerGram,
		LockHeight:     0,
		PaymentId:      nil,
	}
}

// SubmitCoinSplitRequest wraps the CoinSplit GRPC call so we can split coins easier. feePerGram
// lets callers override the fee rate used for the split; pass 0 to use DefaultFeePerGram.
func SubmitCoinSplitRequest(ctx context.Context, splitAmt int, numSplits int, feePerGram uint64) (*tari_generated.CoinSplitResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.CoinSplit(ctx, buildCoinSplitRequest(splitAmt, numSplits, feePerGram))
}

// GetBalances wraps the GetBalances GRPC call
func GetBalances(ctx context.Context) (*tari_generated.GetBalanceResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetBalance(ctx, &tari_generated.GetBalanceRequest{})
}

// GetTransactionInfoByID wraps the GetTransactionInfo call in GRPC, one at a time
func GetTransactionInfoByID(ctx context.Context, transactionID uint64) (*tari_generated.TransactionInfo, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	txns, err := client.GetTransactionInfo(ctx, &tari_generated.GetTransactionInfoRequest{
		TransactionIds: []uint64{transactionID},
	})
	if err != nil || len(txns.Transactions) == 0 {
		return nil, err
	}
	return txns.Transactions[0], nil
}

func RevalidateAllTransactions(ctx context.Context) (*tari_generated.RevalidateResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.RevalidateAllTransactions(ctx, &tari_generated.RevalidateRequest{})
}

func ValidateAllTransactions(ctx context.Context) (*tari_generated.ValidateResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.ValidateAllTransactions(ctx, &tari_generated.ValidateRequest{})
}

// GetWalletState wraps the GetState call on the wallet GRPC
func GetWalletState(ctx context.Context) (*tari_generated.GetStateResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetState(ctx, &tari_generated.GetStateRequest{})
}

// GetWalletConnectivity wraps the CheckConnectivity call
func GetWalletConnectivity(ctx context.Context) (*tari_generated.CheckConnectivityResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.CheckConnectivity(ctx, &tari_generated.GetConnectivityRequest{})
}

// GetAddresses gets the addresses for the wallet
func GetAddresses(ctx context.Context) (*tari_generated.GetCompleteAddressResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetCompleteAddress(ctx, nil)
}

// ---- Additional wrapper functions covering the remaining TariWallet service RPCs ----

// GetVersion wraps the GetVersion GRPC call, returning the wallet's build/version info. Like
// every other zero-field-request wrapper in this package (e.g. GetAddress, GetUnspentAmounts) and
// its nodeGRPC.GetVersion sibling, it hides the empty request and takes only ctx — this was
// originally (within this same not-yet-merged PR) mistakenly given a req *GetVersionRequest
// parameter, breaking that convention for no reason (production-readiness finding 9); fixed
// before merge since no released version ever shipped the old signature.
func GetVersion(ctx context.Context) (*tari_generated.GetVersionResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetVersion(ctx, &tari_generated.GetVersionRequest{})
}

// CheckForUpdates wraps the CheckForUpdates GRPC call, asking the wallet to check for a newer software release.
func CheckForUpdates(ctx context.Context) (*tari_generated.SoftwareUpdate, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.CheckForUpdates(ctx, &tari_generated.Empty{})
}

// Identify wraps the Identify GRPC call, returning the wallet's node identity.
func Identify(ctx context.Context, req *tari_generated.GetIdentityRequest) (*tari_generated.GetIdentityResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.Identify(ctx, req)
}

// GetAddress wraps the GetAddress GRPC call, returning the wallet's interactive address.
func GetAddress(ctx context.Context) (*tari_generated.GetAddressResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetAddress(ctx, &tari_generated.Empty{})
}

// GetPaymentIdAddress wraps the GetPaymentIdAddress GRPC call, deriving a complete address for the supplied payment ID.
func GetPaymentIdAddress(ctx context.Context, req *tari_generated.GetPaymentIdAddressRequest) (*tari_generated.GetCompleteAddressResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetPaymentIdAddress(ctx, req)
}

// PrepareOneSidedTransactionForSigning wraps the PrepareOneSidedTransactionForSigning GRPC call. It builds an unsigned one-sided transaction for external/offline signing; it does not broadcast anything.
func PrepareOneSidedTransactionForSigning(ctx context.Context, req *tari_generated.PrepareOneSidedTransactionForSigningRequest) (*tari_generated.PrepareOneSidedTransactionForSigningResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.PrepareOneSidedTransactionForSigning(ctx, req)
}

// BroadcastSignedOneSidedTransaction wraps the BroadcastSignedOneSidedTransaction GRPC call, submitting a previously-signed one-sided transaction (see PrepareOneSidedTransactionForSigning) to the network.
//
// req carries the transaction's signing material (the completed signature produced offline from
// the unsigned transaction PrepareOneSidedTransactionForSigning returned). Use InitWalletGRPCSecure
// (TLS) rather than InitWalletGRPC (plaintext) when the connection crosses any untrusted network.
func BroadcastSignedOneSidedTransaction(ctx context.Context, req *tari_generated.BroadcastSignedOneSidedTransactionRequest) (*tari_generated.BroadcastSignedOneSidedTransactionResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.BroadcastSignedOneSidedTransaction(ctx, req)
}

// GetTransactionPayRefs wraps the GetTransactionPayRefs GRPC call, returning payment reference data for the requested transactions.
func GetTransactionPayRefs(ctx context.Context, req *tari_generated.GetTransactionPayRefsRequest) (*tari_generated.GetTransactionPayRefsResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetTransactionPayRefs(ctx, req)
}

// GetUnspentAmounts wraps the GetUnspentAmounts GRPC call, returning the amounts of all unspent outputs in the wallet.
func GetUnspentAmounts(ctx context.Context) (*tari_generated.GetUnspentAmountsResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetUnspentAmounts(ctx, &tari_generated.Empty{})
}

// ImportUtxos wraps the ImportUtxos GRPC call, importing externally-supplied UTXOs into the wallet's output manager.
//
// req carries the UTXO's spending_key, which is private key material: it lets the wallet spend
// the imported output. Use InitWalletGRPCSecure (TLS) rather than InitWalletGRPC (plaintext) when
// the connection crosses any untrusted network.
func ImportUtxos(ctx context.Context, req *tari_generated.ImportUtxosRequest) (*tari_generated.ImportUtxosResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.ImportUtxos(ctx, req)
}

// GetNetworkStatus wraps the GetNetworkStatus GRPC call, returning the wallet's view of its network/connectivity status.
func GetNetworkStatus(ctx context.Context) (*tari_generated.NetworkStatusResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetNetworkStatus(ctx, &tari_generated.Empty{})
}

// GetConnectedHttpPeer wraps the GetConnectedHttpPeer GRPC call, returning info about the wallet's connected HTTP peer, if any.
func GetConnectedHttpPeer(ctx context.Context) (*tari_generated.GetConnectedHttpPeerResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetConnectedHttpPeer(ctx, &tari_generated.Empty{})
}

// CancelTransaction wraps the CancelTransaction GRPC call, cancelling a pending transaction by ID.
func CancelTransaction(ctx context.Context, req *tari_generated.CancelTransactionRequest) (*tari_generated.CancelTransactionResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.CancelTransaction(ctx, req)
}

// SendShaAtomicSwapTransaction wraps the SendShaAtomicSwapTransaction GRPC call, initiating a SHA atomic swap (HTLC-style) transaction.
func SendShaAtomicSwapTransaction(ctx context.Context, req *tari_generated.SendShaAtomicSwapRequest) (*tari_generated.SendShaAtomicSwapResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.SendShaAtomicSwapTransaction(ctx, req)
}

// CreateBurnTransaction wraps the CreateBurnTransaction GRPC call, creating and broadcasting a coin-burn transaction.
func CreateBurnTransaction(ctx context.Context, req *tari_generated.CreateBurnTransactionRequest) (*tari_generated.CreateBurnTransactionResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.CreateBurnTransaction(ctx, req)
}

// ClaimShaAtomicSwapTransaction wraps the ClaimShaAtomicSwapTransaction GRPC call, claiming the output of a previously-sent SHA atomic swap transaction.
//
// req carries the HTLC pre-image needed to claim the swap output, which is sensitive material:
// anyone who observes it in transit can claim the output themselves before the legitimate
// claimant's request lands. Use InitWalletGRPCSecure (TLS) rather than InitWalletGRPC (plaintext)
// when the connection crosses any untrusted network.
func ClaimShaAtomicSwapTransaction(ctx context.Context, req *tari_generated.ClaimShaAtomicSwapRequest) (*tari_generated.ClaimShaAtomicSwapResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.ClaimShaAtomicSwapTransaction(ctx, req)
}

// ClaimHtlcRefundTransaction wraps the ClaimHtlcRefundTransaction GRPC call, reclaiming funds from an expired/unclaimed HTLC output.
//
// req carries the HTLC pre-image/refund material needed to reclaim the output, which is
// sensitive: anyone who observes it in transit could attempt to race the legitimate refund. Use
// InitWalletGRPCSecure (TLS) rather than InitWalletGRPC (plaintext) when the connection crosses
// any untrusted network.
func ClaimHtlcRefundTransaction(ctx context.Context, req *tari_generated.ClaimHtlcRefundRequest) (*tari_generated.ClaimHtlcRefundResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.ClaimHtlcRefundTransaction(ctx, req)
}

// CreateTemplateRegistration wraps the CreateTemplateRegistration GRPC call, registering a validator-node template on-chain.
func CreateTemplateRegistration(ctx context.Context, req *tari_generated.CreateTemplateRegistrationRequest) (*tari_generated.CreateTemplateRegistrationResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.CreateTemplateRegistration(ctx, req)
}

// SignMessage wraps the SignMessage GRPC call, signing an arbitrary message with the wallet's private key.
//
// This RPC signs with the wallet's node-identity private key: while the private key itself never
// leaves the wallet process, the request/response round-trip is still worth protecting against a
// network-level attacker who could otherwise observe which messages are being signed or attempt a
// man-in-the-middle substitution. Use InitWalletGRPCSecure (TLS) rather than InitWalletGRPC
// (plaintext) when the connection crosses any untrusted network.
func SignMessage(ctx context.Context, req *tari_generated.SignMessageRequest) (*tari_generated.SignMessageResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.SignMessage(ctx, req)
}

// ImportTransactions wraps the ImportTransactions GRPC call, importing externally-supplied transaction records into the wallet's transaction history.
func ImportTransactions(ctx context.Context, req *tari_generated.ImportTransactionsRequest) (*tari_generated.ImportTransactionsResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.ImportTransactions(ctx, req)
}

// GetAllCompletedTransactions wraps the GetAllCompletedTransactions GRPC call, returning a single
// page of completed transactions (sorted by timestamp) selected by req.Offset/req.Limit; the
// underlying RPC caps req.Limit at 50 rows per request server-side regardless of what's
// requested, so this never returns "every" completed transaction in one call — callers needing
// the full history must page through with repeated calls (or use
// GetAllCompletedTransactionsStream instead).
//
// Deprecated: use GetAllCompletedTransactionsStream instead. The underlying proto RPC is marked
// `option deprecated = true` ("Use GetAllCompletedTransactionsStream for better performance and
// memory efficiency"), and the generated tari_generated.WalletClient.GetAllCompletedTransactions
// method it wraps already carries a `// Deprecated: Do not use.` comment for the same reason.
func GetAllCompletedTransactions(ctx context.Context, req *tari_generated.GetAllCompletedTransactionsRequest) (*tari_generated.GetAllCompletedTransactionsResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetAllCompletedTransactions(ctx, req) //nolint:staticcheck // SA1019: intentional passthrough of a deprecated RPC, see doc comment above.
}

// GetPaymentByReference wraps the GetPaymentByReference GRPC call, looking up payment details by payment reference.
func GetPaymentByReference(ctx context.Context, req *tari_generated.GetPaymentByReferenceRequest) (*tari_generated.GetPaymentByReferenceResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetPaymentByReference(ctx, req)
}

// GetFeeEstimate wraps the GetFeeEstimate GRPC call, estimating the fee for a prospective transaction without sending it.
func GetFeeEstimate(ctx context.Context, req *tari_generated.GetFeeEstimateRequest) (*tari_generated.GetFeeEstimateResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetFeeEstimate(ctx, req)
}

// GetFeePerGramStats wraps the GetFeePerGramStats GRPC call, returning current network fee-per-gram statistics.
func GetFeePerGramStats(ctx context.Context, req *tari_generated.GetFeePerGramStatsRequest) (*tari_generated.GetFeePerGramStatsResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.GetFeePerGramStats(ctx, req)
}

// ReplaceByFee wraps the ReplaceByFee GRPC call, resubmitting a pending transaction with a higher fee (RBF).
func ReplaceByFee(ctx context.Context, req *tari_generated.ReplaceByFeeRequest) (*tari_generated.ReplaceByFeeResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.ReplaceByFee(ctx, req)
}

// UserPayForFee wraps the UserPayForFee GRPC call, having the wallet user cover the fee for a transaction on behalf of another party.
func UserPayForFee(ctx context.Context, req *tari_generated.UserPayForFeeRequest) (*tari_generated.UserPayForFeeResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.UserPayForFee(ctx, req)
}

// RegisterValidatorNode wraps the RegisterValidatorNode GRPC call, registering this node as a validator node on-chain.
func RegisterValidatorNode(ctx context.Context, req *tari_generated.RegisterValidatorNodeRequest) (*tari_generated.RegisterValidatorNodeResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.RegisterValidatorNode(ctx, req)
}

// SubmitValidatorEvictionProof wraps the SubmitValidatorEvictionProof GRPC call, submitting proof to evict a misbehaving validator node.
func SubmitValidatorEvictionProof(ctx context.Context, req *tari_generated.SubmitValidatorEvictionProofRequest) (*tari_generated.SubmitValidatorEvictionProofResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.SubmitValidatorEvictionProof(ctx, req)
}

// SubmitValidatorNodeExit wraps the SubmitValidatorNodeExit GRPC call, voluntarily exiting this node from the validator node set.
func SubmitValidatorNodeExit(ctx context.Context, req *tari_generated.SubmitValidatorNodeExitRequest) (*tari_generated.SubmitValidatorNodeExitResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	return client.SubmitValidatorNodeExit(ctx, req)
}

// StreamTransactionEvents subscribes to the wallet's continuous transaction-event stream. Unlike
// this package's other streaming wrappers (ListHeaders-style callers in nodeGRPC,
// GetAllCompletedTransactionsStream below, etc.), which wrap RPCs that stream a FINITE, bounded
// result set and therefore drain cleanly to a slice, the underlying StreamTransactionEvents RPC
// is documented as "a continuous stream of transaction events as they occur ... ideal for
// real-time UI updates or external monitoring" — it never sends io.EOF on its own while the
// wallet is alive. Drain-to-slice would never return against a real wallet, and worse, would
// silently discard every event already received the moment ctx was cancelled/timed out
// (production-readiness finding 2).
//
// Instead, onEvent is invoked once per event, synchronously, in the order received.
// StreamTransactionEvents itself returns when exactly one of the following happens:
//   - onEvent returns a non-nil error: that error is returned to the caller, unwrapped, and no
//     further events are delivered (this is the intended way for a long-lived caller to end the
//     subscription from inside its own callback, e.g. after deciding it's done).
//   - the underlying stream ends on its own (server closes it, e.g. wallet shutdown): returns nil.
//   - ctx is cancelled or its deadline expires: returns ctx.Err().
//
// This shape (an onEvent callback plus ctx for cancellation) was chosen over returning the raw
// tari_generated.Wallet_StreamTransactionEventsClient because it keeps this package's existing
// idiom of not leaking generated GRPC stream types across the wrapper boundary, while still
// giving a long-lived caller a way to actually consume a subscription that may run for the
// lifetime of the process. There's no established precedent for a genuinely continuous stream
// elsewhere in this package (every other streaming RPC wrapped here is finite), so this is a new
// pattern — introduced here rather than reusing the finite-stream drain-to-slice shape, which
// would be actively wrong for this RPC.
func StreamTransactionEvents(ctx context.Context, req *tari_generated.TransactionEventRequest, onEvent func(*tari_generated.TransactionEventResponse) error) error {
	client, err := walletClient()
	if err != nil {
		return err
	}
	streamClient, err := client.StreamTransactionEvents(ctx, req)
	if err != nil {
		return err
	}
	for {
		item, err := streamClient.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}
		if err := onEvent(item); err != nil {
			return err
		}
	}
}

// GetAllCompletedTransactionsStream wraps the GetAllCompletedTransactionsStream streaming GRPC
// call, draining every completed transaction pushed by the wallet into a slice. Unlike
// GetAllCompletedTransactions, this RPC streams its results rather than returning a single
// paginated response, but the underlying data set is still selected by req.Offset/req.Limit —
// the server still caps req.Limit at 50 rows per request — this wrapper simply drains whatever
// single page the server streams back for the given request into a slice; it does not
// auto-paginate across multiple requests.
func GetAllCompletedTransactionsStream(ctx context.Context, req *tari_generated.GetAllCompletedTransactionsRequest) ([]*tari_generated.GetCompletedTransactionsResponse, error) {
	client, err := walletClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.GetAllCompletedTransactionsStream(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.GetCompletedTransactionsResponse, 0)
	for {
		item, err := streamClient.Recv()
		if err != nil {
			if err == io.EOF {
				return resp, nil
			}
			return nil, err
		}
		resp = append(resp, item)
	}
}
