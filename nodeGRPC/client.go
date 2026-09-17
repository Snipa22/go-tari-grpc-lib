package nodeGRPC

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

// connMu guards grpcNodeAddress/grpcConn below. Without it, concurrent Init calls race on both
// reads and writes of the package-level connection state.
var connMu sync.RWMutex
var grpcNodeAddress string
var grpcConn *grpc.ClientConn

// buildTransportCredentials returns the transport credentials used to dial the base node.
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

// InitNodeGRPC opens a PLAINTEXT (insecure) GRPC connection to a Tari base node at the given
// address. This transport is intended for loopback/trusted-network use only (e.g. a base node
// running on the same host or within a private network you control) — for anything crossing an
// untrusted network, use InitNodeGRPCSecure instead.
func InitNodeGRPC(nodeAddress string) error {
	return initNodeGRPC(nodeAddress, buildTransportCredentials(false, nil))
}

// InitNodeGRPCSecure opens a TLS-secured GRPC connection to a Tari base node at the given
// address, using the supplied tls.Config for the transport credentials. Use this for any
// connection that crosses an untrusted network. Pass nil for tlsConfig to use Go's default TLS
// configuration.
func InitNodeGRPCSecure(nodeAddress string, tlsConfig *tls.Config) error {
	return initNodeGRPC(nodeAddress, buildTransportCredentials(true, tlsConfig))
}

// newClient is the function used to construct the underlying GRPC connection. It's a
// package-level var (defaulting to grpc.NewClient) so tests can substitute a failing stub to
// verify that initNodeGRPC actually propagates the error instead of discarding it (finding 2) —
// grpc.NewClient itself validates most targets lazily at dial time, so it can't reliably be made
// to fail synchronously from a test using only real target strings.
var newClient = grpc.NewClient

func initNodeGRPC(nodeAddress string, transportCreds credentials.TransportCredentials) error {
	conn, err := newClient(nodeAddress, dialOptions(transportCreds)...)
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

// GetTipInfo wraps the GetTipInfo GRPC call and handles the response from the upstream
func GetTipInfo(ctx context.Context) (*tari_generated.TipInfoResponse, error) {
	client := tari_generated.NewBaseNodeClient(getConn())
	return client.GetTipInfo(ctx, &tari_generated.Empty{})
}

// GetBlockTemplate wraps the GetNewBlockTemplate call, requires the type of blockTemplate to generate
func GetBlockTemplate(ctx context.Context, algo *tari_generated.PowAlgo) (*tari_generated.NewBlockTemplateResponse, error) {
	client := tari_generated.NewBaseNodeClient(getConn())
	return client.GetNewBlockTemplate(ctx, &tari_generated.NewBlockTemplateRequest{Algo: algo})
}

// GetBlockWithCoinbases wraps the GetNewBlockWithCoinbases, requires all data for the GRPC request
func GetBlockWithCoinbases(ctx context.Context, requestData *tari_generated.GetNewBlockWithCoinbasesRequest) (*tari_generated.GetNewBlockResult, error) {
	client := tari_generated.NewBaseNodeClient(getConn())
	return client.GetNewBlockWithCoinbases(ctx, requestData)
}

// GetNewBlockTemplateWithCoinbases This incorrectly tells you that you're getting a template, but the response is a full block
func GetNewBlockTemplateWithCoinbases(ctx context.Context, requestData *tari_generated.GetNewBlockTemplateWithCoinbasesRequest) (*tari_generated.GetNewBlockResult, error) {
	client := tari_generated.NewBaseNodeClient(getConn())
	return client.GetNewBlockTemplateWithCoinbases(ctx, requestData)
}

// GetNetworkState wraps the GetNetworkState RPC call
func GetNetworkState(ctx context.Context) (*tari_generated.GetNetworkStateResponse, error) {
	client := tari_generated.NewBaseNodeClient(getConn())
	return client.GetNetworkState(ctx, nil)
}

// GetNewBlock wraps the GetNewBlock GRPC call
func GetNewBlock(ctx context.Context, requestData *tari_generated.NewBlockTemplate) (*tari_generated.GetNewBlockResult, error) {
	client := tari_generated.NewBaseNodeClient(getConn())
	return client.GetNewBlock(ctx, requestData)
}

// GetBlockByHeight retrieves blocks, handles the streaming data, then returns the blocks as a slice
func GetBlockByHeight(ctx context.Context, blockIDs []uint64) ([]*tari_generated.Block, error) {
	client := tari_generated.NewBaseNodeClient(getConn())
	active_client, err := client.GetBlocks(ctx, &tari_generated.GetBlocksRequest{Heights: blockIDs}, grpc.MaxCallRecvMsgSize(16*1024*1024))
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
	client := tari_generated.NewBaseNodeClient(getConn())
	return client.GetHeaderByHash(ctx, &tari_generated.GetHeaderByHashRequest{Hash: blockHash})
}

// SubmitBlock sends blocks to the daemon for processing
func SubmitBlock(ctx context.Context, requestData *tari_generated.Block) (*tari_generated.SubmitBlockResponse, error) {
	client := tari_generated.NewBaseNodeClient(getConn())
	return client.SubmitBlock(ctx, requestData)
}

// GetNetworkDiff pulls the network diff of a given block, or it will just use tip if you give it a 0
func GetNetworkDiff(ctx context.Context, height uint64) (*tari_generated.NetworkDifficultyResponse, error) {
	client := tari_generated.NewBaseNodeClient(getConn())
	var diffClient tari_generated.BaseNode_GetNetworkDifficultyClient
	var err error
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
	client := tari_generated.NewBaseNodeClient(getConn())
	return client.Identify(ctx, nil)
}
