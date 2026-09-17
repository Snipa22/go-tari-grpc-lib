package nodeGRPC

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

/*
RPC management notes, because of the way this works, we need to take a request, open a connection to the client, then
pass the response back to the client, everything is one-shot and every call is responsible for closing it's own conn

use minotari_app_grpc::tari_rpc::{
    GetNewBlockTemplateWithCoinbasesRequest,
    SubmitBlockRequest,
    SubmitBlockResponse,
};
*/

// dialTimeout bounds how long a single connection attempt to the base node GRPC endpoint may take
// before it's abandoned, via grpc.ConnectParams' MinConnectTimeout.
const dialTimeout = 5 * time.Second

// maxRecvMsgSize is the maximum message size (in bytes) this package will accept on ANY BaseNode
// RPC response, unary or streaming. It's set well above grpc-go's 4MiB default because real Tari
// blocks (as returned by, e.g., GetBlocks/HistoricalBlock-shaped RPCs such as GetBlockByHeight,
// SearchKernels, SearchUtxos) routinely exceed that default and would otherwise fail with
// codes.ResourceExhausted. It's applied once, here, as a default call option in dialOptions() so
// it automatically covers every existing and future BaseNode RPC — see finding 1 in the
// production-readiness review that added this constant: two new streaming RPCs (SearchKernels,
// SearchUtxos) had reintroduced the exact bug GetBlockByHeight's call-site-specific
// grpc.MaxCallRecvMsgSize option was originally added to fix, because that fix wasn't centralized.
const maxRecvMsgSize = 16 * 1024 * 1024

// connMu guards grpcNodeAddress/grpcConn below. Without it, concurrent Init calls race on both
// reads and writes of the package-level connection state.
var connMu sync.RWMutex
var grpcNodeAddress string
var grpcConn *grpc.ClientConn

// ErrNotInitialized is returned by every BaseNode RPC wrapper in this package when it's called
// before InitNodeGRPC/InitNodeGRPCSecure has successfully established a connection. Without this,
// tari_generated.NewBaseNodeClient(getConn()) would be constructed from a nil *grpc.ClientConn and
// the subsequent RPC call would panic with a nil-pointer dereference instead of returning a
// handleable error (production-readiness finding 4).
var ErrNotInitialized = errors.New("nodeGRPC: not initialized, call InitNodeGRPC or InitNodeGRPCSecure first")

// buildTransportCredentials returns the transport credentials used to dial the base node.
// Exposed as its own function so it's directly unit-testable without needing a live server.
func buildTransportCredentials(secure bool, tlsConfig *tls.Config) credentials.TransportCredentials {
	if secure {
		return credentials.NewTLS(tlsConfig)
	}
	return insecure.NewCredentials()
}

// dialOptions returns the common set of grpc.DialOption values used for every Init* variant,
// given the transport credentials to use, plus any caller-supplied extra options appended after
// the defaults. extra is exposed all the way up through InitNodeGRPC/InitNodeGRPCSecure so
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

// InitNodeGRPC opens a PLAINTEXT (insecure) GRPC connection to a Tari base node at the given
// address. This transport is intended for loopback/trusted-network use only (e.g. a base node
// running on the same host or within a private network you control) — for anything crossing an
// untrusted network, use InitNodeGRPCSecure instead.
//
// extra, if supplied, is appended to the default dial options (after them, so an extra option of
// the same kind takes precedence per grpc-go's normal last-one-wins semantics) — e.g. to attach a
// unary/stream interceptor or stats handler for instrumentation. Existing callers passing no extra
// arguments are unaffected.
func InitNodeGRPC(nodeAddress string, extra ...grpc.DialOption) error {
	return initNodeGRPC(nodeAddress, buildTransportCredentials(false, nil), extra...)
}

// InitNodeGRPCSecure opens a TLS-secured GRPC connection to a Tari base node at the given
// address, using the supplied tls.Config for the transport credentials. Use this for any
// connection that crosses an untrusted network. Pass nil for tlsConfig to use Go's default TLS
// configuration.
//
// extra, if supplied, is appended to the default dial options (after them, so an extra option of
// the same kind takes precedence per grpc-go's normal last-one-wins semantics) — e.g. to attach a
// unary/stream interceptor or stats handler for instrumentation. Existing callers passing no extra
// arguments are unaffected.
func InitNodeGRPCSecure(nodeAddress string, tlsConfig *tls.Config, extra ...grpc.DialOption) error {
	return initNodeGRPC(nodeAddress, buildTransportCredentials(true, tlsConfig), extra...)
}

// newClient is the function used to construct the underlying GRPC connection. It's a
// package-level var (defaulting to grpc.NewClient) so tests can substitute a failing stub to
// verify that initNodeGRPC actually propagates the error instead of discarding it (finding 2) —
// grpc.NewClient itself validates most targets lazily at dial time, so it can't reliably be made
// to fail synchronously from a test using only real target strings.
var newClient = grpc.NewClient

func initNodeGRPC(nodeAddress string, transportCreds credentials.TransportCredentials, extra ...grpc.DialOption) error {
	conn, err := newClient(nodeAddress, dialOptions(transportCreds, extra...)...)
	if err != nil {
		return err
	}
	connMu.Lock()
	defer connMu.Unlock()
	grpcNodeAddress = nodeAddress
	grpcConn = conn
	return nil
}

// getConn returns the current GRPC connection in a race-free way.
func getConn() *grpc.ClientConn {
	connMu.RLock()
	defer connMu.RUnlock()
	return grpcConn
}

// baseNodeClient returns a tari_generated.BaseNodeClient bound to the current connection, or
// ErrNotInitialized if InitNodeGRPC/InitNodeGRPCSecure hasn't successfully run yet. Every RPC
// wrapper in this file goes through this helper instead of calling
// tari_generated.NewBaseNodeClient(getConn()) directly, so the nil-conn-panic bug (finding 4)
// is fixed once, centrally, for all ~80 call sites instead of needing a nil check hand-added to
// each wrapper individually.
func baseNodeClient() (tari_generated.BaseNodeClient, error) {
	conn := getConn()
	if conn == nil {
		return nil, ErrNotInitialized
	}
	return tari_generated.NewBaseNodeClient(conn), nil
}

// GetTipInfo wraps the GetTipInfo GRPC call and handles the response from the upstream
func GetTipInfo(ctx context.Context) (*tari_generated.TipInfoResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetTipInfo(ctx, &tari_generated.Empty{})
}

// GetBlockTemplate wraps the GetNewBlockTemplate call, requires the type of blockTemplate to generate
func GetBlockTemplate(ctx context.Context, algo *tari_generated.PowAlgo) (*tari_generated.NewBlockTemplateResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetNewBlockTemplate(ctx, &tari_generated.NewBlockTemplateRequest{Algo: algo})
}

// GetBlockWithCoinbases wraps the GetNewBlockWithCoinbases, requires all data for the GRPC request
func GetBlockWithCoinbases(ctx context.Context, requestData *tari_generated.GetNewBlockWithCoinbasesRequest) (*tari_generated.GetNewBlockResult, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetNewBlockWithCoinbases(ctx, requestData)
}

// GetNewBlockTemplateWithCoinbases This incorrectly tells you that you're getting a template, but the response is a full block
func GetNewBlockTemplateWithCoinbases(ctx context.Context, requestData *tari_generated.GetNewBlockTemplateWithCoinbasesRequest) (*tari_generated.GetNewBlockResult, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetNewBlockTemplateWithCoinbases(ctx, requestData)
}

// GetNetworkState wraps the GetNetworkState RPC call
func GetNetworkState(ctx context.Context) (*tari_generated.GetNetworkStateResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetNetworkState(ctx, nil)
}

// GetNewBlock wraps the GetNewBlock GRPC call
func GetNewBlock(ctx context.Context, requestData *tari_generated.NewBlockTemplate) (*tari_generated.GetNewBlockResult, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetNewBlock(ctx, requestData)
}

// GetBlockByHeight retrieves blocks, handles the streaming data, then returns the blocks as a slice
func GetBlockByHeight(ctx context.Context, blockIDs []uint64) ([]*tari_generated.Block, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	// No per-call grpc.MaxCallRecvMsgSize override needed here: dialOptions() now sets it as a
	// connection-wide default (finding 1), which covers this call too.
	active_client, err := client.GetBlocks(ctx, &tari_generated.GetBlocksRequest{Heights: blockIDs})
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.Block, 0)
	for {
		blockResp, err := active_client.Recv()
		if err != nil {
			if err == io.EOF {
				return resp, nil
			}
			return nil, err
		}
		resp = append(resp, blockResp.GetBlock())
	}
}

// GetHeaderByHash wraps the GRPC call of the same name.
func GetHeaderByHash(ctx context.Context, blockHash []byte) (*tari_generated.BlockHeaderResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetHeaderByHash(ctx, &tari_generated.GetHeaderByHashRequest{Hash: blockHash})
}

// SubmitBlock sends blocks to the daemon for processing
func SubmitBlock(ctx context.Context, requestData *tari_generated.Block) (*tari_generated.SubmitBlockResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.SubmitBlock(ctx, requestData)
}

// GetNetworkDiff pulls the network diff of a given block, or it will just use tip if you give it a 0
func GetNetworkDiff(ctx context.Context, height uint64) (*tari_generated.NetworkDifficultyResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	var diffClient tari_generated.BaseNode_GetNetworkDifficultyClient
	if height == 0 {
		diffClient, err = client.GetNetworkDifficulty(ctx, &tari_generated.HeightRequest{FromTip: 1})
	} else {
		diffClient, err = client.GetNetworkDifficulty(ctx, &tari_generated.HeightRequest{StartHeight: height, EndHeight: height})
	}
	if err != nil {
		return nil, err
	}
	return diffClient.Recv()
}

// GetNodeIdentity returns a list of valid rust identities for an opened GRPC node
func GetNodeIdentity(ctx context.Context) (*tari_generated.NodeIdentity, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.Identify(ctx, nil)
}

// ---- Additional wrapper functions covering the remaining BaseNode service RPCs ----

// GetBlockTiming wraps the GetBlockTiming GRPC call, returning block timing statistics for the requested height range.
func GetBlockTiming(ctx context.Context, req *tari_generated.HeightRequest) (*tari_generated.BlockTimingResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetBlockTiming(ctx, req)
}

// GetConstants wraps the GetConstants GRPC call, returning the consensus constants in effect at the given block height.
func GetConstants(ctx context.Context, req *tari_generated.BlockHeight) (*tari_generated.ConsensusConstants, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetConstants(ctx, req)
}

// GetBlockSize wraps the GetBlockSize GRPC call, returning block size statistics for the requested height range.
func GetBlockSize(ctx context.Context, req *tari_generated.BlockGroupRequest) (*tari_generated.BlockGroupResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetBlockSize(ctx, req)
}

// GetBlockFees wraps the GetBlockFees GRPC call, returning block fee statistics for the requested height range.
func GetBlockFees(ctx context.Context, req *tari_generated.BlockGroupRequest) (*tari_generated.BlockGroupResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetBlockFees(ctx, req)
}

// GetVersion wraps the GetVersion GRPC call, returning the base node's build/version info.
func GetVersion(ctx context.Context) (*tari_generated.BaseNodeGetVersionResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetVersion(ctx, &tari_generated.Empty{})
}

// CheckForUpdates wraps the CheckForUpdates GRPC call, asking the base node to check for a newer software release.
func CheckForUpdates(ctx context.Context) (*tari_generated.SoftwareUpdate, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.CheckForUpdates(ctx, &tari_generated.Empty{})
}

// GetNewBlockBlob wraps the GetNewBlockBlob GRPC call, returning a mined-block template as an opaque blob (used by miners that work on serialized blobs instead of structured blocks).
func GetNewBlockBlob(ctx context.Context, req *tari_generated.NewBlockTemplate) (*tari_generated.GetNewBlockBlobResult, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetNewBlockBlob(ctx, req)
}

// SubmitBlockBlob wraps the SubmitBlockBlob GRPC call, submitting a mined block supplied as an opaque blob (see GetNewBlockBlob) for processing.
func SubmitBlockBlob(ctx context.Context, req *tari_generated.BlockBlobRequest) (*tari_generated.SubmitBlockResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.SubmitBlockBlob(ctx, req)
}

// SubmitTransaction wraps the SubmitTransaction GRPC call, submitting a transaction directly to the base node's mempool.
func SubmitTransaction(ctx context.Context, req *tari_generated.SubmitTransactionRequest) (*tari_generated.SubmitTransactionResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.SubmitTransaction(ctx, req)
}

// GetSyncInfo wraps the GetSyncInfo GRPC call, returning the base node's current sync peer/state info.
func GetSyncInfo(ctx context.Context) (*tari_generated.SyncInfoResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetSyncInfo(ctx, &tari_generated.Empty{})
}

// GetSyncProgress wraps the GetSyncProgress GRPC call, returning the base node's current sync progress.
func GetSyncProgress(ctx context.Context) (*tari_generated.SyncProgressResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetSyncProgress(ctx, &tari_generated.Empty{})
}

// TransactionState wraps the TransactionState GRPC call, returning the mempool/chain state of a specific transaction (by excess signature).
func TransactionState(ctx context.Context, req *tari_generated.TransactionStateRequest) (*tari_generated.TransactionStateResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.TransactionState(ctx, req)
}

// GetNetworkStatus wraps the GetNetworkStatus GRPC call, returning the base node's view of its network/connectivity status.
func GetNetworkStatus(ctx context.Context) (*tari_generated.NetworkStatusResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetNetworkStatus(ctx, &tari_generated.Empty{})
}

// ListConnectedPeers wraps the ListConnectedPeers GRPC call, returning the peers this base node currently has an active connection to.
func ListConnectedPeers(ctx context.Context) (*tari_generated.ListConnectedPeersResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.ListConnectedPeers(ctx, &tari_generated.Empty{})
}

// GetMempoolStats wraps the GetMempoolStats GRPC call, returning aggregate statistics about the base node's mempool.
func GetMempoolStats(ctx context.Context) (*tari_generated.MempoolStatsResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetMempoolStats(ctx, &tari_generated.Empty{})
}

// GetValidatorNodeChanges wraps the GetValidatorNodeChanges GRPC call, returning validator-node set changes for the requested height range.
func GetValidatorNodeChanges(ctx context.Context, req *tari_generated.GetValidatorNodeChangesRequest) (*tari_generated.GetValidatorNodeChangesResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetValidatorNodeChanges(ctx, req)
}

// GetShardKey wraps the GetShardKey GRPC call, returning the shard key for a public key at a given height.
func GetShardKey(ctx context.Context, req *tari_generated.GetShardKeyRequest) (*tari_generated.GetShardKeyResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	return client.GetShardKey(ctx, req)
}

// ListHeaders wraps the ListHeaders streaming GRPC call, draining every header pushed by the base node into a slice.
func ListHeaders(ctx context.Context, req *tari_generated.ListHeadersRequest) ([]*tari_generated.BlockHeaderResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.ListHeaders(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.BlockHeaderResponse, 0)
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

// GetTokensInCirculation wraps the GetTokensInCirculation streaming GRPC call, draining the circulating-supply value for every requested height into a slice.
func GetTokensInCirculation(ctx context.Context, req *tari_generated.GetBlocksRequest) ([]*tari_generated.ValueAtHeightResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.GetTokensInCirculation(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.ValueAtHeightResponse, 0)
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

// SearchKernels wraps the SearchKernels streaming GRPC call, draining every historical block containing a matching kernel into a slice.
func SearchKernels(ctx context.Context, req *tari_generated.SearchKernelsRequest) ([]*tari_generated.HistoricalBlock, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.SearchKernels(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.HistoricalBlock, 0)
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

// SearchUtxos wraps the SearchUtxos streaming GRPC call, draining every historical block containing a matching UTXO commitment into a slice.
func SearchUtxos(ctx context.Context, req *tari_generated.SearchUtxosRequest) ([]*tari_generated.HistoricalBlock, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.SearchUtxos(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.HistoricalBlock, 0)
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

// FetchMatchingUtxos wraps the FetchMatchingUtxos streaming GRPC call, draining every matching UTXO into a slice.
func FetchMatchingUtxos(ctx context.Context, req *tari_generated.FetchMatchingUtxosRequest) ([]*tari_generated.FetchMatchingUtxosResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.FetchMatchingUtxos(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.FetchMatchingUtxosResponse, 0)
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

// GetPeers wraps the GetPeers streaming GRPC call, draining every known peer pushed by the base node into a slice.
func GetPeers(ctx context.Context, req *tari_generated.GetPeersRequest) ([]*tari_generated.GetPeersResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.GetPeers(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.GetPeersResponse, 0)
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

// GetMempoolTransactions wraps the GetMempoolTransactions streaming GRPC call, draining every mempool transaction pushed by the base node into a slice.
func GetMempoolTransactions(ctx context.Context, req *tari_generated.GetMempoolTransactionsRequest) ([]*tari_generated.GetMempoolTransactionsResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.GetMempoolTransactions(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.GetMempoolTransactionsResponse, 0)
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

// GetActiveValidatorNodes wraps the GetActiveValidatorNodes streaming GRPC call, draining every active validator node pushed by the base node into a slice.
func GetActiveValidatorNodes(ctx context.Context, req *tari_generated.GetActiveValidatorNodesRequest) ([]*tari_generated.GetActiveValidatorNodesResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.GetActiveValidatorNodes(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.GetActiveValidatorNodesResponse, 0)
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

// GetTemplateRegistrations wraps the GetTemplateRegistrations streaming GRPC call, draining every validator-node template registration into a slice.
func GetTemplateRegistrations(ctx context.Context, req *tari_generated.GetTemplateRegistrationsRequest) ([]*tari_generated.GetTemplateRegistrationResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.GetTemplateRegistrations(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.GetTemplateRegistrationResponse, 0)
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

// GetSideChainUtxos wraps the GetSideChainUtxos streaming GRPC call, draining every side-chain UTXO into a slice.
func GetSideChainUtxos(ctx context.Context, req *tari_generated.GetSideChainUtxosRequest) ([]*tari_generated.GetSideChainUtxosResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.GetSideChainUtxos(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.GetSideChainUtxosResponse, 0)
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

// SearchPaymentReferences wraps the SearchPaymentReferences streaming GRPC call, draining every matching payment reference into a slice.
func SearchPaymentReferences(ctx context.Context, req *tari_generated.SearchPaymentReferencesRequest) ([]*tari_generated.PaymentReferenceResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.SearchPaymentReferences(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.PaymentReferenceResponse, 0)
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

// SearchPaymentReferencesViaOutputHash wraps the SearchPaymentReferencesViaOutputHash streaming GRPC call (looked up by output hash instead of commitment), draining every matching payment reference into a slice.
func SearchPaymentReferencesViaOutputHash(ctx context.Context, req *tari_generated.FetchMatchingUtxosRequest) ([]*tari_generated.PaymentReferenceResponse, error) {
	client, err := baseNodeClient()
	if err != nil {
		return nil, err
	}
	streamClient, err := client.SearchPaymentReferencesViaOutputHash(ctx, req)
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.PaymentReferenceResponse, 0)
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
