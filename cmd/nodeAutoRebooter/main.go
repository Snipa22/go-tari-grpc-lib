package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Snipa22/go-tari-grpc-lib/v3/nodeGRPC"
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

// heightLagThreshold is how many blocks behind the oracle-derived best height the local node can
// be before it's considered "behind" and eligible for a reboot.
const heightLagThreshold uint64 = 5

// peerCallTimeout bounds how long a single oracle-peer GRPC call may take, so one dead/slow peer
// in the list can't stall the whole probing loop.
const peerCallTimeout = 5 * time.Second

// shouldCheckHeightBasedReboot reports whether the local node is far enough behind the
// oracle-derived bestHeight to warrant a reboot.
//
// If hasResponded is false (i.e. every oracle peer was unreachable), this always returns false.
// Previously, bestHeight stayed at its zero value in that scenario, and localHeight <
// bestHeight-heightLagThreshold underflowed uint64(0)-heightLagThreshold to a value near
// math.MaxUint64, which made the comparison unconditionally true and triggered a restart even
// though we genuinely had no information about the network. This function is the single source
// of truth for that decision so main() doesn't duplicate the underflow-prone arithmetic.
func shouldCheckHeightBasedReboot(localHeight uint64, bestHeight uint64, hasResponded bool) bool {
	if !hasResponded {
		return false
	}
	return localHeight < bestHeight-heightLagThreshold
}

func main() {
	oraclePeersPtr := flag.String("oracle-peers", strings.Join(defaultOraclePeers, ","), "Comma-separated list of host:port oracle peers used to determine the network's best known height")
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

	if err := nodeGRPC.InitNodeGRPC("127.0.0.1:18102"); err != nil {
		log.Printf("Failed to connect to local node: %v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), peerCallTimeout)
	val, err := nodeGRPC.GetTipInfo(ctx)
	cancel()
	if err != nil || val == nil {
		log.Printf("Failed to get tip info from local node: %v", err)
		return
	}

	if shouldCheckHeightBasedReboot(val.Metadata.BestBlockHeight, bestHeight, hasResponded) {
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
		systemdConnection.RestartUnitContext(rebootCtx, "tari.service", "replace", holderChan)
	}
}
