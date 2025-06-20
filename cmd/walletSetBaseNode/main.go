package main

import (
	"flag"
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/v2/walletGRPC"
	"log"
)

func main() {
	walletGRPCAddressPtr := flag.String("wallet-grpc-address", "127.0.0.1:18143", "Tari wallet GRPC address")
	nodePublicKeyPtr := flag.String("node-public-key", "", "Public key of the node to make the base node")
	nodeNetAddrPtr := flag.String("node-net-address", "", "Network address of the node, normally something like `/ip4/192.168.4.254/tcp/9998`")
	flag.Parse()
	if *nodePublicKeyPtr == "" {
		log.Fatal("Node public key is required")
	}
	if *nodeNetAddrPtr == "" {
		log.Fatal("Node network address is required")
	}
	walletGRPC.InitWalletGRPC(*walletGRPCAddressPtr)
	_, err := walletGRPC.SetBaseNode(*nodePublicKeyPtr, *nodeNetAddrPtr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Set node public key and address")
}
