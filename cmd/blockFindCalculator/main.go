package main

import (
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/v3/nodeGRPC"
	"log"
)

func main() {
	fmt.Println("Getting Diff Data")
	diffData, err := nodeGRPC.GetNetworkDiff(3133)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(diffData)
	hashesPerSecond := 4000000000
	fmt.Printf("%v seconds per block for %v hashes/second or %v hours", diffData.Difficulty/uint64(hashesPerSecond), hashesPerSecond, (diffData.Difficulty/uint64(hashesPerSecond))/3600)
}
