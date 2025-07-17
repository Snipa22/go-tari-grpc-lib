package nodeGRPC

import (
	"context"
	"github.com/Snipa22/go-tari-grpc-lib/v2/tari_generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
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

var grpcNodeAddress string
var grpcConn *grpc.ClientConn

func InitNodeGRPC(nodeAddress string) {
	grpcNodeAddress = nodeAddress
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	grpcConn, _ = grpc.NewClient(grpcNodeAddress, opts...)
}

// GetTipInfo wraps the GetTipInfo GRPC call and handles the response from the upstream
func GetTipInfo() (*tari_generated.TipInfoResponse, error) {
	client := tari_generated.NewBaseNodeClient(grpcConn)
	return client.GetTipInfo(context.Background(), &tari_generated.Empty{})
}

// GetBlockTemplate wraps the GetNewBlockTemplate call, requires the type of blockTemplate to generate
func GetBlockTemplate(algo *tari_generated.PowAlgo) (*tari_generated.NewBlockTemplateResponse, error) {
	client := tari_generated.NewBaseNodeClient(grpcConn)
	return client.GetNewBlockTemplate(context.Background(), &tari_generated.NewBlockTemplateRequest{Algo: algo})
}

// GetBlockWithCoinbases wraps the GetNewBlockWithCoinbases, requires all data for the GRPC request
func GetBlockWithCoinbases(requestData *tari_generated.GetNewBlockWithCoinbasesRequest) (*tari_generated.GetNewBlockResult, error) {
	client := tari_generated.NewBaseNodeClient(grpcConn)
	return client.GetNewBlockWithCoinbases(context.Background(), requestData)
}

// GetNetworkState wraps the GetNetworkState RPC call
func GetNetworkState() (*tari_generated.GetNetworkStateResponse, error) {
	client := tari_generated.NewBaseNodeClient(grpcConn)
	return client.GetNetworkState(context.Background(), nil)
}

// GetNewBlock wraps the GetNewBlock GRPC call
func GetNewBlock(requestData *tari_generated.NewBlockTemplate) (*tari_generated.GetNewBlockResult, error) {
	client := tari_generated.NewBaseNodeClient(grpcConn)
	return client.GetNewBlock(context.Background(), requestData)
}

// GetBlockByHeight retrieves blocks, handles the streaming data, then returns the blocks as a slice
func GetBlockByHeight(blockIDs []uint64) ([]*tari_generated.Block, error) {
	client := tari_generated.NewBaseNodeClient(grpcConn)
	active_client, err := client.GetBlocks(context.Background(), &tari_generated.GetBlocksRequest{Heights: blockIDs}, grpc.MaxCallRecvMsgSize(16*1024*1024))
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
func GetHeaderByHash(blockHash []byte) (*tari_generated.BlockHeaderResponse, error) {
	client := tari_generated.NewBaseNodeClient(grpcConn)
	return client.GetHeaderByHash(context.Background(), &tari_generated.GetHeaderByHashRequest{Hash: blockHash})
}

// SubmitBlock sends blocks to the daemon for processing
func SubmitBlock(requestData *tari_generated.Block) (*tari_generated.SubmitBlockResponse, error) {
	client := tari_generated.NewBaseNodeClient(grpcConn)
	return client.SubmitBlock(context.Background(), requestData)
}

// GetNetworkDiff pulls the network diff of a given block, or it will just use tip if you give it a 0
func GetNetworkDiff(height uint64) (*tari_generated.NetworkDifficultyResponse, error) {
	client := tari_generated.NewBaseNodeClient(grpcConn)
	var diffClient tari_generated.BaseNode_GetNetworkDifficultyClient
	var err error
	if height == 0 {
		diffClient, err = client.GetNetworkDifficulty(context.Background(), &tari_generated.HeightRequest{FromTip: 1})
	} else {
		diffClient, err = client.GetNetworkDifficulty(context.Background(), &tari_generated.HeightRequest{StartHeight: height, EndHeight: height})
	}
	if err != nil {
		return nil, err
	}
	return diffClient.Recv()
}

// GetNodeIdentity returns a list of valid rust identities for an opened GRPC node
func GetNodeIdentity() (*tari_generated.NodeIdentity, error) {
	client := tari_generated.NewBaseNodeClient(grpcConn)
	return client.Identify(context.Background(), nil)
}
