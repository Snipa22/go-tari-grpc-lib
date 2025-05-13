package walletGRPC

import (
	"context"
	"flag"
	"github.com/Snipa22/go-tari-grpc-lib/tari_generated"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var grpcWalletAddress string

func InitWalletGRPC(walletAddress string) {
	grpcWalletAddress = walletAddress
}

// getWalletConnection builds the connection so we can init the BaseWalletClient, it does NOT close the connection, so
// we need to close the connection down stream.
func getWalletConnection() (*grpc.ClientConn, error) {
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(grpcWalletAddress, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// SendTransactions sends the transactions to the wallet
func SendTransactions(transactions []*tari_generated.PaymentRecipient) (*tari_generated.TransferResponse, error) {
	conn, err := getWalletConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := tari_generated.NewWalletClient(conn)
	return client.Transfer(context.Background(), &tari_generated.TransferRequest{
		Recipients: transactions,
	})
}
