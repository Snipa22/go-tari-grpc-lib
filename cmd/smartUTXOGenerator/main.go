package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/Snipa22/core-go-lib/helpers"
	core "github.com/Snipa22/core-go-lib/milieu"
	"github.com/Snipa22/go-tari-grpc-lib/v2/walletGRPC"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
)

var utxoMinCount uint64 = 400
var sqliteDbPath = ""
var milieu *core.Milieu

func checkAndSplitUTXOs() {
	milieu.Info("Starting UTXO Split")
	sqliteDSN := fmt.Sprintf("file:%s?cache=shared&mode=ro", sqliteDbPath)
	db, err := sql.Open("sqlite3", sqliteDSN)
	if err != nil {
		milieu.CaptureException(err)
		milieu.Error(err.Error())
		return
	}
	defer db.Close()
	row := db.QueryRow("select count(1) from outputs where status = 0")
	count := 0
	if err = row.Scan(&count); err != nil {
		milieu.CaptureException(err)
		milieu.Error(err.Error())
		return
	}
	db.Close()
	milieu.Info(fmt.Sprintf("%v/%v UTXOs", count, utxoMinCount))
	if uint64(count) < utxoMinCount {
		balances, err := walletGRPC.GetBalances()
		if err != nil {
			milieu.CaptureException(err)
			milieu.Error(err.Error())
			return
		}
		splitValue := balances.AvailableBalance / 10 / (utxoMinCount * 2)
		milieu.Info(fmt.Sprintf("Making %v UTXOs for %v XTM", splitValue, utxoMinCount*2))
		_, err = walletGRPC.SubmitCoinSplitRequest(int(splitValue), int(utxoMinCount))
		if err != nil {
			milieu.CaptureException(err)
			milieu.Error(err.Error())
			return
		}
		_, err = walletGRPC.SubmitCoinSplitRequest(int(splitValue), int(utxoMinCount))
		if err != nil {
			milieu.CaptureException(err)
			milieu.Error(err.Error())
			return
		}
	}
	milieu.Info("UTXO Split Complete")
}

func main() {
	minUtxoCountFlag := flag.Int("min-utxos", 400, "Minimum number of UTXOs to keep in wallet, serves as generation amount")
	walletGRPCAddressPtr := flag.String("wallet-grpc-address", "127.0.0.1:18143", "Tari wallet GRPC address")
	walletSqliteDBPtr := flag.String("wallet-sqlite-db", "", "Path to the source tari wallet sqlite DB")
	sentryURI := helpers.GetEnv("SENTRY_SERVER", "")
	flag.Parse()

	// Configure all the things
	utxoMinCount = uint64(*minUtxoCountFlag)
	sqliteDbPath = *walletSqliteDBPtr
	walletGRPC.InitWalletGRPC(*walletGRPCAddressPtr)

	milieu, _ = core.NewMilieu(nil, nil, &sentryURI)

	checkAndSplitUTXOs()

	// Build the cron spinner
	c := cron.New()
	_, _ = c.AddFunc("30 * * * *", func() {
		checkAndSplitUTXOs()
	})
	c.Run()

	// Idle loop!
	for {
		select {}
	}
}
