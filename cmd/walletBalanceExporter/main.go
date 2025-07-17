package main

import (
	"flag"
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/v3/walletGRPC"
	"net/http"
)

func rootCall(w http.ResponseWriter, req *http.Request) {
	balances, err := walletGRPC.GetBalances()
	if err != nil {
		fmt.Fprintf(w, `{"error": "%v"}`, err)
		return
	}
	fmt.Fprintf(w, `{"available_balance": %v, "pending_incoming_balance": %v, "pending_outgoing_balance": %v, "timelocked_balance": %v}`, balances.AvailableBalance, balances.PendingIncomingBalance, balances.PendingOutgoingBalance, balances.TimelockedBalance)
	return
}

func main() {
	walletGRPCAddressPtr := flag.String("wallet-grpc-address", "127.0.0.1:18143", "Tari wallet GRPC address")
	flag.Parse()
	walletGRPC.InitWalletGRPC(*walletGRPCAddressPtr)
	http.HandleFunc("/", rootCall)
	http.ListenAndServe("127.0.0.1:2049", nil)
}
