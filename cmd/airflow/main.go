package main

import (
	"flag"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/joho/godotenv"
	"github.com/tikivn/clickhousectl/internal/utils/mem_stats"
	"github.com/tikivn/ultrago/u_graceful"
	"github.com/tikivn/ultrago/u_logger"
	"golang.org/x/sync/errgroup"
)

func init() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
}

func main() {
	ctx, logger := u_logger.GetLogger(u_graceful.NewCtx())

	eg, gctx := errgroup.WithContext(ctx)
	airFlowApp, cleanup, err := initAirFlow(gctx)
	if err != nil {
		panic(err)
	}

	defer cleanup()

	actionPtr := flag.String("action", "", "actions: InsertData")
	tablePtr := flag.String("table", "", "DDL of interacted table")
	queryPtr := flag.String("query", "", "detail for action")
	flag.Parse()

	go mem_stats.Monitor(gctx, 5*time.Second)

	if flag.Parsed() {
		if *actionPtr == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
		logger.Infof("actionPtr: %s, tablePtr: %s, queryPtr: %s\n",
			*actionPtr, *tablePtr, *queryPtr)

		results := reflect.ValueOf(airFlowApp).MethodByName(*actionPtr).Call([]reflect.Value{
			reflect.ValueOf(gctx),
			reflect.ValueOf(*tablePtr),
			reflect.ValueOf(*queryPtr),
		})
		if len(results) >= 1 {
			result := results[0]
			if result.IsNil() {
				os.Exit(0)
			} else {
				logger.Fatalf("%s failed: %v", *actionPtr, result.Interface())
				os.Exit(1)
			}
		}
	}

	if err := eg.Wait(); err != nil {
		logger.Errorln(err)
	}
}
