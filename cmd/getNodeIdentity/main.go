package main

import (
	"flag"
	"fmt"
	"github.com/Snipa22/core-go-lib/helpers"
	core "github.com/Snipa22/core-go-lib/milieu"
	"github.com/Snipa22/go-tari-grpc-lib/v2/nodeGRPC"
)

func main() {
	sentry := helpers.GetEnv("SENTRY_SERVER", "")
	milieu, err := core.NewMilieu(nil, nil, &sentry)
	if err != nil {
		milieu.CaptureException(err)
		milieu.Fatal(err.Error())
	}
	// Milieu initialized

	// Load config flags
	nodeGRPCPtr := flag.String("base-node-grpc-address", "node-pool.tari.jagtech.io:18102", "Address for the base-node, defaults to Impala's public pool")
	flag.Parse()

	nodeGRPC.InitNodeGRPC(*nodeGRPCPtr)

	nodeIdents, err := nodeGRPC.GetNodeIdentity()
	if err != nil {
		milieu.CaptureException(err)
		milieu.Fatal(err.Error())
	}
	for _, v := range nodeIdents.PublicAddresses {
		fmt.Printf(`"%x::%v",`, nodeIdents.PublicKey, v)
		fmt.Println()
	}
}
