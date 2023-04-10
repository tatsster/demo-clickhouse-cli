//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"github.com/tikivn/clickhousectl/internal/app/cli"
)

func initGoApp(ctx context.Context) (cli.App, func(), error) {
	wire.Build(cli.GraphSet)
	return nil, nil, nil
}
