package main

import (
	"flag"
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/v2/tari_generated"
	"github.com/Snipa22/go-tari-grpc-lib/v2/walletGRPC"
	"log"
)

func main() {
	walletAddress := flag.String("walletAddress", "", "Wallet address to send to")
	amount := flag.Int("amount", 0, "Amount of uT to send, if you want to send 1XTM, this should be 1000000")
	walletGRPCAddressPtr := flag.String("wallet-grpc-address", "127.0.0.1:18143", "Tari wallet GRPC address")
	flag.Parse()
	walletGRPC.InitWalletGRPC(*walletGRPCAddressPtr)
	if *walletAddress == "" {
		log.Fatalln("No valid wallet address passed")
	}
	fmt.Printf("Sending %v uT to %v\n", *amount, *walletAddress)
	txns := make([]*tari_generated.PaymentRecipient, 0)
	txns = append(txns, &tari_generated.PaymentRecipient{
		Address:     *walletAddress,
		Amount:      uint64(*amount),
		FeePerGram:  uint64(5),
		PaymentType: tari_generated.PaymentRecipient_ONE_SIDED,
	})
	resp, err := walletGRPC.SendTransactions(txns)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Submitted request to wallet, txid: %v or error: %v\n", resp.Results[0].TransactionId, resp.Results[0].FailureMessage)
}
