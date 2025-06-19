package main

import (
	"flag"
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/v2/tari_generated"
	"github.com/Snipa22/go-tari-grpc-lib/v2/walletGRPC"
	"log"
)

func main() {
	walletGRPCAddressPtr := flag.String("wallet-grpc-address", "127.0.0.1:18143", "Tari wallet GRPC address")
	destAddressPtr := flag.String("dest-address", "", "Destination address for sweeping to")
	flag.Parse()
	if *destAddressPtr == "" {
		log.Fatal("Destination address is required")
	}
	walletGRPC.InitWalletGRPC(*walletGRPCAddressPtr)
	balances, err := walletGRPC.GetBalances()
	if err != nil {
		log.Fatal(err)
	}
	remainder := balances.AvailableBalance % 100000
	toSend := balances.AvailableBalance - remainder
	if toSend <= 1000000 {
		fmt.Println("Balance to send is less than 1 XTM, not sweeping.")
		return
	}
	fmt.Printf("Sending %v uT to %v\n", toSend, *destAddressPtr)
	resp, err := walletGRPC.SendTransactions([]*tari_generated.PaymentRecipient{
		{
			Address:     *destAddressPtr,
			Amount:      toSend,
			FeePerGram:  5,
			PaymentType: 2,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transactions sent successfully: %v, error if exists: %v\n", resp.Results[0].IsSuccess, resp.Results[0].FailureMessage)
}
