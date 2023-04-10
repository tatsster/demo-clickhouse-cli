package repo

import (
	"context"

	"github.com/tikivn/clickhousectl/internal/pkg/entity"
)

type ClickHouseServerRepo interface {
	FindAll(ctx context.Context, serverId string) ([]*entity.ClickHouseServer, error)
}
