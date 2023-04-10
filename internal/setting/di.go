package setting

import "github.com/google/wire"

var GraphSet = wire.NewSet(
	NewMySqlConfig,
)
