package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/v3/nodeGRPC"
	"github.com/Snipa22/go-tari-grpc-lib/v3/walletGRPC"
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
	walletGRPCAddressPtr := flag.String("wallet-grpc-address", "127.0.0.1:18143", "Tari wallet GRPC address")
	flag.Parse()
	for _, node := range nodeList {
		nodeGRPC.InitNodeGRPC(node)
		if val, _ := nodeGRPC.GetTipInfo(); val != nil {
			if val.Metadata.BestBlockHeight > bestHeight {
				bestHeight = val.Metadata.BestBlockHeight
			}
		}
	}
	walletGRPC.InitWalletGRPC(*walletGRPCAddressPtr)
	walletState, err := walletGRPC.GetWalletState()
	if err != nil {
		return
	}
	shouldReboot := false
	fmt.Printf("Wallet scanned height is %v/%v\n", walletState.ScannedHeight, bestHeight)
	if walletState.ScannedHeight < bestHeight-10 {
		shouldReboot = true
	}
	walletConnectivity, _ := walletGRPC.GetWalletConnectivity()
	fmt.Printf("Wallet connectivity is %v\n", walletConnectivity)
	if walletConnectivity.Status == 2 {
		shouldReboot = true
	}
	if shouldReboot {
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
		systemdConnection.RestartUnitContext(ctx, "tari_wallet.service", "replace", holderChan)
	}
}
