package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/Snipa22/core-go-lib/helpers"
	core "github.com/Snipa22/core-go-lib/milieu"
	"github.com/Snipa22/go-tari-grpc-lib/v3/walletGRPC"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron/v3"
)

var utxoMinCount uint64 = 400
var sqliteDbPath = ""
var milieu *core.Milieu

// splitCallCount is how many times a single checkAndSplitUTXOs run should invoke the split
// submitter when a split is needed. This is a regression guard for a historical copy-paste bug
// that called SubmitCoinSplitRequest twice in a row with the identical computed amount, silently
// doubling the real split amount despite the log line above it reporting only splitValue once.
const splitCallCount = 1

// splitUTXOsIfNeeded checks whether the existing unspent-output count is below minCount, and if
// so, invokes submit exactly splitCallCount times with the computed split amount. It returns the
// number of times submit was attempted, which lets callers/tests assert exactly-once behavior
// without needing a real database or a live GRPC connection.
func splitUTXOsIfNeeded(existingCount uint64, minCount uint64, availableBalance uint64, submit func(splitAmt int, numSplits int) error) (int, error) {
	if existingCount >= minCount {
		return 0, nil
	}
	splitValue := availableBalance / 10 / (minCount * 2)
	calls := 0
	for i := 0; i < splitCallCount; i++ {
		calls++
		if err := submit(int(splitValue), int(minCount)); err != nil {
			return calls, err
		}
	}
	return calls, nil
}

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
		balances, err := walletGRPC.GetBalances(context.Background())
		if err != nil {
			milieu.CaptureException(err)
			milieu.Error(err.Error())
			return
		}
		calls, err := splitUTXOsIfNeeded(uint64(count), utxoMinCount, balances.AvailableBalance, func(splitAmt int, numSplits int) error {
			milieu.Info(fmt.Sprintf("Making %v UTXOs for %v XTM", splitAmt, numSplits*2))
			_, err := walletGRPC.SubmitCoinSplitRequest(context.Background(), splitAmt, numSplits, walletGRPC.DefaultFeePerGram)
			return err
		})
		if err != nil {
			milieu.CaptureException(err)
			milieu.Error(err.Error())
			return
		}
		if calls > 0 {
			milieu.Info(fmt.Sprintf("Submitted %v split request(s)", calls))
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

	milieu, _ = core.NewMilieu(nil, nil, &sentryURI)

	if err := walletGRPC.InitWalletGRPC(*walletGRPCAddressPtr); err != nil {
		milieu.CaptureException(err)
		milieu.Fatal(err.Error())
	}

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
