package nodeGRPC

import (
	"context"
	"flag"
	"go-tari-grpc-lib/tari_generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"os"
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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getBaseConnection builds the connection so we can init the BaseNodeClient, it does NOT close the connection, so we
// need to close the connection down stream.
func getBaseConnection() (*grpc.ClientConn, error) {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(getEnv("basenode_address", "node-pool.tari.jagtech.io:18102"), opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// GetTipInfo wraps the GetTipInfo GRPC call and handles the response from the upstream
func GetTipInfo() (*tari_generated.TipInfoResponse, error) {
	conn, err := getBaseConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := tari_generated.NewBaseNodeClient(conn)
	return client.GetTipInfo(context.Background(), &tari_generated.Empty{})
}

// GetBlockTemplate wraps the GetNewBlockTemplate call, requires the type of blockTemplate to generate
func GetBlockTemplate(algo *tari_generated.PowAlgo) (*tari_generated.NewBlockTemplateResponse, error) {
	conn, err := getBaseConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := tari_generated.NewBaseNodeClient(conn)
	return client.GetNewBlockTemplate(context.Background(), &tari_generated.NewBlockTemplateRequest{Algo: algo})
}

// GetBlockWithCoinbases wraps the GetNewBlockWithCoinbases, requires all data for the GRPC request
func GetBlockWithCoinbases(requestData *tari_generated.GetNewBlockTemplateWithCoinbasesRequest) (*tari_generated.GetNewBlockResult, error) {
	conn, err := getBaseConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := tari_generated.NewBaseNodeClient(conn)
	return client.GetNewBlockTemplateWithCoinbases(context.Background(), requestData)
}

// GetBlockByHeight retrieves blocks, handles the streaming data, then returns the blocks as a slice
func GetBlockByHeight(blockIDs []uint64) ([]*tari_generated.Block, error) {
	conn, err := getBaseConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := tari_generated.NewBaseNodeClient(conn)
	active_client, err := client.GetBlocks(context.Background(), &tari_generated.GetBlocksRequest{Heights: blockIDs})
	if err != nil {
		return nil, err
	}
	resp := make([]*tari_generated.Block, len(blockIDs))
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

// SubmitBlock sends blocks to the daemon for processing
func SubmitBlock(requestData *tari_generated.Block) (*tari_generated.SubmitBlockResponse, error) {
	conn, err := getBaseConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := tari_generated.NewBaseNodeClient(conn)
	return client.SubmitBlock(context.Background(), requestData)
}
