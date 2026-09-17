package walletGRPC

import (
	"context"
	"crypto/tls"
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

// connMu guards grpcWalletAddress/grpcConn below. Without it, concurrent Init calls race on both
// reads and writes of the package-level connection state.
var connMu sync.RWMutex
var grpcWalletAddress string
var grpcConn *grpc.ClientConn

// buildTransportCredentials returns the transport credentials used to dial the wallet.
// Exposed as its own function so it's directly unit-testable without needing a live server.
func buildTransportCredentials(secure bool, tlsConfig *tls.Config) credentials.TransportCredentials {
	if secure {
		return credentials.NewTLS(tlsConfig)
	}
	return insecure.NewCredentials()
}

// dialOptions returns the common set of grpc.DialOption values used for every Init* variant,
// given the transport credentials to use.
func dialOptions(transportCreds credentials.TransportCredentials) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: dialTimeout}),
	}
}

// InitWalletGRPC opens a PLAINTEXT (insecure) GRPC connection to a Tari wallet at the given
// address. This transport is intended for loopback/trusted-network use only (e.g. a wallet
// running on the same host or within a private network you control) — for anything crossing an
// untrusted network, use InitWalletGRPCSecure instead.
func InitWalletGRPC(walletAddress string) error {
	return initWalletGRPC(walletAddress, buildTransportCredentials(false, nil))
}

// InitWalletGRPCSecure opens a TLS-secured GRPC connection to a Tari wallet at the given
// address, using the supplied tls.Config for the transport credentials. Use this for any
// connection that crosses an untrusted network. Pass nil for tlsConfig to use Go's default TLS
// configuration.
func InitWalletGRPCSecure(walletAddress string, tlsConfig *tls.Config) error {
	return initWalletGRPC(walletAddress, buildTransportCredentials(true, tlsConfig))
}

// newClient is the function used to construct the underlying GRPC connection. It's a
// package-level var (defaulting to grpc.NewClient) so tests can substitute a failing stub to
// verify that initWalletGRPC actually propagates the error instead of discarding it (finding 2) —
// grpc.NewClient itself validates most targets lazily at dial time, so it can't reliably be made
// to fail synchronously from a test using only real target strings.
var newClient = grpc.NewClient

func initWalletGRPC(walletAddress string, transportCreds credentials.TransportCredentials) error {
	conn, err := newClient(walletAddress, dialOptions(transportCreds)...)
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

// SendTransactions sends the transactions to the wallet
func SendTransactions(ctx context.Context, transactions []*tari_generated.PaymentRecipient) (*tari_generated.TransferResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.Transfer(ctx, &tari_generated.TransferRequest{
		Recipients: transactions,
	})
}

// GetTransactionsInBlock will return the top of the wallet if it's called with 0, otherwise it pushes the height to
// the GRPC call, though this doesn't seem to actually do anything.  No sorting/order/etc is guaranteed, so callers
// need to parse, cache etc.
func GetTransactionsInBlock(ctx context.Context, blockHeight uint64) ([]*tari_generated.TransactionInfo, error) {
	client := tari_generated.NewWalletClient(getConn())
	var completedTxnsClient tari_generated.Wallet_GetCompletedTransactionsClient
	var err error
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
	client := tari_generated.NewWalletClient(getConn())
	return client.CoinSplit(ctx, buildCoinSplitRequest(splitAmt, numSplits, feePerGram))
}

// GetBalances wraps the GetBalances GRPC call
func GetBalances(ctx context.Context) (*tari_generated.GetBalanceResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetBalance(ctx, &tari_generated.GetBalanceRequest{})
}

// GetTransactionInfoByID wraps the GetTransactionInfo call in GRPC, one at a time
func GetTransactionInfoByID(ctx context.Context, transactionID uint64) (*tari_generated.TransactionInfo, error) {
	client := tari_generated.NewWalletClient(getConn())
	txns, err := client.GetTransactionInfo(ctx, &tari_generated.GetTransactionInfoRequest{
		TransactionIds: []uint64{transactionID},
	})
	if err != nil || len(txns.Transactions) == 0 {
		return nil, err
	}
	return txns.Transactions[0], nil
}

func RevalidateAllTransactions(ctx context.Context) (*tari_generated.RevalidateResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.RevalidateAllTransactions(ctx, &tari_generated.RevalidateRequest{})
}

func ValidateAllTransactions(ctx context.Context) (*tari_generated.ValidateResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.ValidateAllTransactions(ctx, &tari_generated.ValidateRequest{})
}

// GetWalletState wraps the GetState call on the wallet GRPC
func GetWalletState(ctx context.Context) (*tari_generated.GetStateResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetState(ctx, &tari_generated.GetStateRequest{})
}

// GetWalletConnectivity wraps the CheckConnectivity call
func GetWalletConnectivity(ctx context.Context) (*tari_generated.CheckConnectivityResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.CheckConnectivity(ctx, &tari_generated.GetConnectivityRequest{})
}

// GetAddresses gets the addresses for the wallet
func GetAddresses(ctx context.Context) (*tari_generated.GetCompleteAddressResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetCompleteAddress(ctx, nil)
}

// ---- Additional wrapper functions covering the remaining TariWallet service RPCs ----

// GetVersion wraps the GetVersion GRPC call, returning the wallet's build/version info.
func GetVersion(ctx context.Context, req *tari_generated.GetVersionRequest) (*tari_generated.GetVersionResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetVersion(ctx, req)
}

// CheckForUpdates wraps the CheckForUpdates GRPC call, asking the wallet to check for a newer software release.
func CheckForUpdates(ctx context.Context) (*tari_generated.SoftwareUpdate, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.CheckForUpdates(ctx, &tari_generated.Empty{})
}

// Identify wraps the Identify GRPC call, returning the wallet's node identity.
func Identify(ctx context.Context, req *tari_generated.GetIdentityRequest) (*tari_generated.GetIdentityResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.Identify(ctx, req)
}

// GetAddress wraps the GetAddress GRPC call, returning the wallet's interactive address.
func GetAddress(ctx context.Context) (*tari_generated.GetAddressResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetAddress(ctx, &tari_generated.Empty{})
}

// GetPaymentIdAddress wraps the GetPaymentIdAddress GRPC call, deriving a complete address for the supplied payment ID.
func GetPaymentIdAddress(ctx context.Context, req *tari_generated.GetPaymentIdAddressRequest) (*tari_generated.GetCompleteAddressResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetPaymentIdAddress(ctx, req)
}

// PrepareOneSidedTransactionForSigning wraps the PrepareOneSidedTransactionForSigning GRPC call. It builds an unsigned one-sided transaction for external/offline signing; it does not broadcast anything.
func PrepareOneSidedTransactionForSigning(ctx context.Context, req *tari_generated.PrepareOneSidedTransactionForSigningRequest) (*tari_generated.PrepareOneSidedTransactionForSigningResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.PrepareOneSidedTransactionForSigning(ctx, req)
}

// BroadcastSignedOneSidedTransaction wraps the BroadcastSignedOneSidedTransaction GRPC call, submitting a previously-signed one-sided transaction (see PrepareOneSidedTransactionForSigning) to the network.
func BroadcastSignedOneSidedTransaction(ctx context.Context, req *tari_generated.BroadcastSignedOneSidedTransactionRequest) (*tari_generated.BroadcastSignedOneSidedTransactionResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.BroadcastSignedOneSidedTransaction(ctx, req)
}

// GetTransactionPayRefs wraps the GetTransactionPayRefs GRPC call, returning payment reference data for the requested transactions.
func GetTransactionPayRefs(ctx context.Context, req *tari_generated.GetTransactionPayRefsRequest) (*tari_generated.GetTransactionPayRefsResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetTransactionPayRefs(ctx, req)
}

// GetUnspentAmounts wraps the GetUnspentAmounts GRPC call, returning the amounts of all unspent outputs in the wallet.
func GetUnspentAmounts(ctx context.Context) (*tari_generated.GetUnspentAmountsResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetUnspentAmounts(ctx, &tari_generated.Empty{})
}

// ImportUtxos wraps the ImportUtxos GRPC call, importing externally-supplied UTXOs into the wallet's output manager.
func ImportUtxos(ctx context.Context, req *tari_generated.ImportUtxosRequest) (*tari_generated.ImportUtxosResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.ImportUtxos(ctx, req)
}

// GetNetworkStatus wraps the GetNetworkStatus GRPC call, returning the wallet's view of its network/connectivity status.
func GetNetworkStatus(ctx context.Context) (*tari_generated.NetworkStatusResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetNetworkStatus(ctx, &tari_generated.Empty{})
}

// GetConnectedHttpPeer wraps the GetConnectedHttpPeer GRPC call, returning info about the wallet's connected HTTP peer, if any.
func GetConnectedHttpPeer(ctx context.Context) (*tari_generated.GetConnectedHttpPeerResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetConnectedHttpPeer(ctx, &tari_generated.Empty{})
}

// CancelTransaction wraps the CancelTransaction GRPC call, cancelling a pending transaction by ID.
func CancelTransaction(ctx context.Context, req *tari_generated.CancelTransactionRequest) (*tari_generated.CancelTransactionResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.CancelTransaction(ctx, req)
}

// SendShaAtomicSwapTransaction wraps the SendShaAtomicSwapTransaction GRPC call, initiating a SHA atomic swap (HTLC-style) transaction.
func SendShaAtomicSwapTransaction(ctx context.Context, req *tari_generated.SendShaAtomicSwapRequest) (*tari_generated.SendShaAtomicSwapResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.SendShaAtomicSwapTransaction(ctx, req)
}

// CreateBurnTransaction wraps the CreateBurnTransaction GRPC call, creating and broadcasting a coin-burn transaction.
func CreateBurnTransaction(ctx context.Context, req *tari_generated.CreateBurnTransactionRequest) (*tari_generated.CreateBurnTransactionResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.CreateBurnTransaction(ctx, req)
}

// ClaimShaAtomicSwapTransaction wraps the ClaimShaAtomicSwapTransaction GRPC call, claiming the output of a previously-sent SHA atomic swap transaction.
func ClaimShaAtomicSwapTransaction(ctx context.Context, req *tari_generated.ClaimShaAtomicSwapRequest) (*tari_generated.ClaimShaAtomicSwapResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.ClaimShaAtomicSwapTransaction(ctx, req)
}

// ClaimHtlcRefundTransaction wraps the ClaimHtlcRefundTransaction GRPC call, reclaiming funds from an expired/unclaimed HTLC output.
func ClaimHtlcRefundTransaction(ctx context.Context, req *tari_generated.ClaimHtlcRefundRequest) (*tari_generated.ClaimHtlcRefundResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.ClaimHtlcRefundTransaction(ctx, req)
}

// CreateTemplateRegistration wraps the CreateTemplateRegistration GRPC call, registering a validator-node template on-chain.
func CreateTemplateRegistration(ctx context.Context, req *tari_generated.CreateTemplateRegistrationRequest) (*tari_generated.CreateTemplateRegistrationResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.CreateTemplateRegistration(ctx, req)
}

// SignMessage wraps the SignMessage GRPC call, signing an arbitrary message with the wallet's private key.
func SignMessage(ctx context.Context, req *tari_generated.SignMessageRequest) (*tari_generated.SignMessageResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.SignMessage(ctx, req)
}

// ImportTransactions wraps the ImportTransactions GRPC call, importing externally-supplied transaction records into the wallet's transaction history.
func ImportTransactions(ctx context.Context, req *tari_generated.ImportTransactionsRequest) (*tari_generated.ImportTransactionsResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.ImportTransactions(ctx, req)
}

// GetAllCompletedTransactions wraps the GetAllCompletedTransactions GRPC call, returning every completed transaction in a single (non-streamed) response.
func GetAllCompletedTransactions(ctx context.Context, req *tari_generated.GetAllCompletedTransactionsRequest) (*tari_generated.GetAllCompletedTransactionsResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetAllCompletedTransactions(ctx, req)
}

// GetPaymentByReference wraps the GetPaymentByReference GRPC call, looking up payment details by payment reference.
func GetPaymentByReference(ctx context.Context, req *tari_generated.GetPaymentByReferenceRequest) (*tari_generated.GetPaymentByReferenceResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetPaymentByReference(ctx, req)
}

// GetFeeEstimate wraps the GetFeeEstimate GRPC call, estimating the fee for a prospective transaction without sending it.
func GetFeeEstimate(ctx context.Context, req *tari_generated.GetFeeEstimateRequest) (*tari_generated.GetFeeEstimateResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetFeeEstimate(ctx, req)
}

// GetFeePerGramStats wraps the GetFeePerGramStats GRPC call, returning current network fee-per-gram statistics.
func GetFeePerGramStats(ctx context.Context, req *tari_generated.GetFeePerGramStatsRequest) (*tari_generated.GetFeePerGramStatsResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.GetFeePerGramStats(ctx, req)
}

// ReplaceByFee wraps the ReplaceByFee GRPC call, resubmitting a pending transaction with a higher fee (RBF).
func ReplaceByFee(ctx context.Context, req *tari_generated.ReplaceByFeeRequest) (*tari_generated.ReplaceByFeeResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.ReplaceByFee(ctx, req)
}

// UserPayForFee wraps the UserPayForFee GRPC call, having the wallet user cover the fee for a transaction on behalf of another party.
func UserPayForFee(ctx context.Context, req *tari_generated.UserPayForFeeRequest) (*tari_generated.UserPayForFeeResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.UserPayForFee(ctx, req)
}

// RegisterValidatorNode wraps the RegisterValidatorNode GRPC call, registering this node as a validator node on-chain.
func RegisterValidatorNode(ctx context.Context, req *tari_generated.RegisterValidatorNodeRequest) (*tari_generated.RegisterValidatorNodeResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.RegisterValidatorNode(ctx, req)
}

// SubmitValidatorEvictionProof wraps the SubmitValidatorEvictionProof GRPC call, submitting proof to evict a misbehaving validator node.
func SubmitValidatorEvictionProof(ctx context.Context, req *tari_generated.SubmitValidatorEvictionProofRequest) (*tari_generated.SubmitValidatorEvictionProofResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.SubmitValidatorEvictionProof(ctx, req)
}

// SubmitValidatorNodeExit wraps the SubmitValidatorNodeExit GRPC call, voluntarily exiting this node from the validator node set.
func SubmitValidatorNodeExit(ctx context.Context, req *tari_generated.SubmitValidatorNodeExitRequest) (*tari_generated.SubmitValidatorNodeExitResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	return client.SubmitValidatorNodeExit(ctx, req)
}

// StreamTransactionEvents wraps the StreamTransactionEvents streaming GRPC call, draining every transaction event pushed by the wallet into a slice. It blocks until the server closes the stream, so it will not return for wallets that keep the stream open indefinitely without a context deadline/cancellation.
func StreamTransactionEvents(ctx context.Context, req *tari_generated.TransactionEventRequest) ([]*tari_generated.TransactionEventResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
	streamClient, err := client.StreamTransactionEvents(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.TransactionEventResponse, 0)
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

// GetAllCompletedTransactionsStream wraps the GetAllCompletedTransactionsStream streaming GRPC call, draining every completed transaction pushed by the wallet into a slice.
func GetAllCompletedTransactionsStream(ctx context.Context, req *tari_generated.GetAllCompletedTransactionsRequest) ([]*tari_generated.GetCompletedTransactionsResponse, error) {
	client := tari_generated.NewWalletClient(getConn())
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
