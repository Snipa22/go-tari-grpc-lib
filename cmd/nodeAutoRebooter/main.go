package main

import (
	"context"
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/v2/nodeGRPC"
	"github.com/coreos/go-systemd/v22/dbus"
)

// Check against a list of nodes, if our local height is 5 under the highest, then reboot it via dbus

var nodeList = []string{
	"51.91.215.198:18102",
	"51.210.222.91:18102",
	"141.94.99.110:18102",
	"162.218.117.106:18102",
	"162.218.117.98:18102",
	"184.164.76.210:18102",
	"15.235.227.47:18102",
	"15.235.227.59:18102",
	"15.235.228.36:18102",
}

var bestHeight uint64 = 0

func main() {
	for _, node := range nodeList {
		nodeGRPC.InitNodeGRPC(node)
		if val, _ := nodeGRPC.GetTipInfo(); val != nil {
			if val.Metadata.BestBlockHeight > bestHeight {
				bestHeight = val.Metadata.BestBlockHeight
			}
		}
	}
	nodeGRPC.InitNodeGRPC("127.0.0.1:18102")
	if val, _ := nodeGRPC.GetTipInfo(); val != nil {
		if val.Metadata.BestBlockHeight < bestHeight-5 {
			// Perform reboot va dbus
			ctx := context.Background()
			// Connect to systemd
			// Specifically this will look DBUS_SYSTEM_BUS_ADDRESS environment variable
			// For example: `unix:path=/run/dbus/system_bus_socket`
			systemdConnection, err := dbus.NewSystemConnectionContext(ctx)
			if err != nil {
				fmt.Printf("Failed to connect to systemd: %v\n", err)
				panic(err)
			}
			defer systemdConnection.Close()
			holderChan := make(chan string)
			systemdConnection.RestartUnitContext(ctx, "tari.service", "replace", holderChan)
		}
	}
}
