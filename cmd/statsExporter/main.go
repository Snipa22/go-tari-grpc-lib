package main

import (
	"flag"
	"github.com/Snipa22/core-go-lib/helpers"
	core "github.com/Snipa22/core-go-lib/milieu"
	"github.com/Snipa22/core-go-lib/milieu/middleware"
	v1 "github.com/Snipa22/go-tari-grpc-lib/v2/cmd/statsExporter/routes/v1"
	"github.com/Snipa22/go-tari-grpc-lib/v2/nodeGRPC"
	"github.com/Snipa22/go-tari-grpc-lib/v2/walletGRPC"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func main() {
	// Build Milieu
	// psqlURL := helpers.GetEnv("PSQL_SERVER", "postgres://postgres@localhost/postgres?sslmode=disable")
	sentry := helpers.GetEnv("SENTRY_SERVER", "")
	milieu, err := core.NewMilieu(nil, nil, &sentry)
	if err != nil {
		milieu.CaptureException(err)
		milieu.Fatal(err.Error())
	}
	// Milieu initialized

	// Load config flags
	debugEnabledPtr := flag.Bool("debug-enabled", false, "Enable Debug Logging")
	nodeGRPCPtr := flag.String("base-node-grpc-address", "node-pool.tari.jagtech.io:18102", "Address for the base-node, defaults to Impala's public pool")
	walletGRPCAddressPtr := flag.String("wallet-grpc-address", "127.0.0.1:18143", "Tari wallet GRPC address")
	flag.Parse()

	nodeGRPC.InitNodeGRPC(*nodeGRPCPtr)
	walletGRPC.InitWalletGRPC(*walletGRPCAddressPtr)

	if *debugEnabledPtr {
		milieu.SetLogLevel(logrus.DebugLevel)
	}

	r := gin.Default()
	r.Use(middleware.SetupMilieu(milieu))

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	v1Router := r.Group("/v1")
	{
		v1Router.GET("/networkStats", v1.GetNetworkStats)
	}

	err = r.Run("127.0.0.1:2050") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		milieu.CaptureException(err)
		log.Fatalf("Unable to initalize gin")
	}
}
