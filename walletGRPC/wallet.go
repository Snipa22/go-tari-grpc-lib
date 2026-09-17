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
