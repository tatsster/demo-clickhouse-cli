package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"

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
	apiServer, cleanup, err := initGoApp(gctx)
	if err != nil {
		panic(err)
	}

	eg.Go(func() error {
		return apiServer.Start(gctx)
	})

	defer func() {
		shutDownErr := apiServer.Stop(context.Background())
		logger.Infof("The API Server is shutting down with err=%v", shutDownErr)
		cleanup()
	}()

	if err := eg.Wait(); err != nil {
		logger.Errorln(err)
	}
}
