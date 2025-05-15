package main

import (
	"flag"
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/v2/walletGRPC"
)

func main() {
	amtPerSplit := flag.Int("amount-per-split", 1000000, "Amount of uT per split, defaults to 1T or 1000000 uT")
	numSplits := flag.Int("num-splits", 400, "Number of splits to make, wallet must have amount-per-split * num-splits available, defaults to 500")
	walletGRPCAddressPtr := flag.String("wallet-grpc-address", "127.0.0.1:18143", "Tari wallet GRPC address")
	flag.Parse()
	walletGRPC.InitWalletGRPC(*walletGRPCAddressPtr)
	fmt.Println("Prepping split for a wallet to UTXO's")
	resp, err := walletGRPC.SubmitCoinSplitRequest(*amtPerSplit, *numSplits)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Submitted request to wallet, txid: %v", resp.TxId)
}
