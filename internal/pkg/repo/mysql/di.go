package mysql

import "github.com/google/wire"

var GraphSet = wire.NewSet(
	NewUserRepo,
	NewClickHouseServerRepo,
)
