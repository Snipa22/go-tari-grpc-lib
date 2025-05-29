package main

import (
	"flag"
	"github.com/Snipa22/go-tari-grpc-lib/v2/walletGRPC"
	"log"
)

func main() {
	revalidateAllTxns := flag.Bool("revalidate-all-txns", false, "Revalidate all transactions in the wallet")
	validateAllTxns := flag.Bool("validate-all-txns", false, "Validate all transactions in the wallet")
	walletGRPCAddressPtr := flag.String("wallet-grpc-address", "127.0.0.1:18143", "Tari wallet GRPC address")
	flag.Parse()
	walletGRPC.InitWalletGRPC(*walletGRPCAddressPtr)
	if *revalidateAllTxns {
		_, err := walletGRPC.RevalidateAllTransactions()
		if err != nil {
			log.Fatal(err)
		}
	}
	if *validateAllTxns {
		_, err := walletGRPC.ValidateAllTransactions()
		if err != nil {
			log.Fatal(err)
		}
	}
}
