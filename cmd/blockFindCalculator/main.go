package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/v3/nodeGRPC"
	"log"
)

func main() {
	nodeGRPCPtr := flag.String("base-node-grpc-address", "node-pool.tari.jagtech.io:18102", "Address for the base-node, defaults to Impala's public pool")
	flag.Parse()
	if err := nodeGRPC.InitNodeGRPC(*nodeGRPCPtr); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Getting Diff Data")
	diffData, err := nodeGRPC.GetNetworkDiff(context.Background(), 3133)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(diffData)
	hashesPerSecond := 4000000000
	fmt.Printf("%v seconds per block for %v hashes/second or %v hours", diffData.Difficulty/uint64(hashesPerSecond), hashesPerSecond, (diffData.Difficulty/uint64(hashesPerSecond))/3600)
}
