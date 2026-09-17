package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Snipa22/go-tari-grpc-lib/v3/nodeGRPC"
	"github.com/Snipa22/go-tari-grpc-lib/v3/walletGRPC"
	"github.com/coreos/go-systemd/v22/dbus"
)

// Check against a list of nodes, if our local height is 5 under the highest, then reboot it via dbus

// defaultOraclePeers is the built-in fallback list of third-party peers used to determine the
// network's best known height. Overridable at runtime via the -oracle-peers flag.
var defaultOraclePeers = []string{
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

// heightLagThreshold is how many blocks behind the oracle-derived best height the wallet's
// scanned height can be before it's considered "behind" and eligible for a reboot.
const heightLagThreshold uint64 = 10

// peerCallTimeout bounds how long a single oracle-peer GRPC call may take, so one dead/slow peer
// in the list can't stall the whole probing loop.
const peerCallTimeout = 5 * time.Second

// shouldCheckHeightBasedReboot reports whether the wallet's scanned height is far enough behind
// the oracle-derived bestHeight to warrant a reboot.
//
// If hasResponded is false (i.e. every oracle peer was unreachable), this always returns false.
// Previously, bestHeight stayed at its zero value in that scenario, and scannedHeight <
// bestHeight-heightLagThreshold underflowed uint64(0)-heightLagThreshold to a value near
// math.MaxUint64, which made the comparison unconditionally true and triggered a restart even
// though we genuinely had no information about the network. This function is the single source
// of truth for that decision so main() doesn't duplicate the underflow-prone arithmetic.
func shouldCheckHeightBasedReboot(scannedHeight uint64, bestHeight uint64, hasResponded bool) bool {
	if !hasResponded {
		return false
	}
	return scannedHeight < bestHeight-heightLagThreshold
}

func main() {
	oraclePeersPtr := flag.String("oracle-peers", strings.Join(defaultOraclePeers, ","), "Comma-separated list of host:port oracle peers used to determine the network's best known height")
	walletGRPCAddressPtr := flag.String("wallet-grpc-address", "127.0.0.1:18143", "Tari wallet GRPC address")
	flag.Parse()

	nodeList := strings.Split(*oraclePeersPtr, ",")

	var bestHeight uint64
	hasResponded := false
	for _, node := range nodeList {
		node = strings.TrimSpace(node)
		if node == "" {
			continue
		}
		if err := nodeGRPC.InitNodeGRPC(node); err != nil {
			log.Printf("Failed to connect to oracle peer %s: %v", node, err)
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), peerCallTimeout)
		val, err := nodeGRPC.GetTipInfo(ctx)
		cancel()
		if err != nil || val == nil {
			log.Printf("Failed to get tip info from oracle peer %s: %v", node, err)
			continue
		}
		hasResponded = true
		if val.Metadata.BestBlockHeight > bestHeight {
			bestHeight = val.Metadata.BestBlockHeight
		}
	}

	if !hasResponded {
		log.Printf("WARNING: all %d oracle peers were unreachable; skipping the height-based auto-reboot check this run", len(nodeList))
	}

	if err := walletGRPC.InitWalletGRPC(*walletGRPCAddressPtr); err != nil {
		log.Printf("Failed to connect to wallet: %v", err)
		return
	}

	shouldReboot := false

	// The height-based check is independent of the walletConnectivity check below; it's only
	// skipped (via shouldCheckHeightBasedReboot) when zero oracle peers responded.
	stateCtx, cancel := context.WithTimeout(context.Background(), peerCallTimeout)
	walletState, err := walletGRPC.GetWalletState(stateCtx)
	cancel()
	if err != nil {
		log.Printf("Failed to get wallet state: %v", err)
	} else {
		fmt.Printf("Wallet scanned height is %v/%v\n", walletState.ScannedHeight, bestHeight)
		if shouldCheckHeightBasedReboot(walletState.ScannedHeight, bestHeight, hasResponded) {
			shouldReboot = true
		}
	}

	// walletConnectivity.Status == 2 is an independent reboot condition and always runs,
	// regardless of oracle peer availability.
	connCtx, cancel := context.WithTimeout(context.Background(), peerCallTimeout)
	walletConnectivity, err := walletGRPC.GetWalletConnectivity(connCtx)
	cancel()
	if err != nil {
		log.Printf("Failed to get wallet connectivity: %v", err)
	} else {
		fmt.Printf("Wallet connectivity is %v\n", walletConnectivity)
		if walletConnectivity.Status == 2 {
			shouldReboot = true
		}
	}

	if shouldReboot {
		// Perform reboot va dbus
		rebootCtx := context.Background()
		// Connect to systemd
		// Specifically this will look DBUS_SYSTEM_BUS_ADDRESS environment variable
		// For example: `unix:path=/run/dbus/system_bus_socket`
		systemdConnection, err := dbus.NewSystemConnectionContext(rebootCtx)
		if err != nil {
			fmt.Printf("Failed to connect to systemd: %v\n", err)
			panic(err)
		}
		defer systemdConnection.Close()
		holderChan := make(chan string)
		systemdConnection.RestartUnitContext(rebootCtx, "tari_wallet.service", "replace", holderChan)
	}
}
