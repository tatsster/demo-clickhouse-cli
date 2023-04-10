//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/tikivn/clickhousectl/internal/app/airflow"
)

func initAirFlow(ctx context.Context) (airflow.App, func(), error) {
	wire.Build(airflow.GraphSet)
	return nil, nil, nil
}
