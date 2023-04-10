package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func NewCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		shutCh := make(chan os.Signal, 1)
		signal.Notify(shutCh, os.Interrupt, syscall.SIGTERM)
		<-shutCh
		log.Println("Graceful shutdown")
		cancel()
	}()
	return ctx
}

func BlockListen(ctx context.Context, fn func() error) error {
	lisErr := make(chan error, 1)
	go func() {
		if e := fn(); e != nil {
			lisErr <- e
		} else {
			close(lisErr)
		}
	}()
	for {
		select {
		case err, _ := <-lisErr:
			return err
		case <-ctx.Done():
			return nil
		}
	}
}
