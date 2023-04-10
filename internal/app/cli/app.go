package cli

import (
	"context"
	"log"

	"github.com/tikivn/clickhousectl/internal/utils/shutdown"
	"golang.org/x/sync/errgroup"
)

func NewApp(httpServer *HttpServer) App {
	return &app{httpServer: httpServer}
}

type App interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type app struct {
	httpServer *HttpServer
}

func (a *app) Start(ctx context.Context) error {
	eg, childCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return shutdown.BlockListen(childCtx, func() error {
			return a.httpServer.ListenAndServe()
		})
	})

	log.Println("API Server started")
	return eg.Wait()
}

func (a *app) Stop(ctx context.Context) error {
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
