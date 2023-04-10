package cli

import (
	"github.com/google/wire"
	"github.com/tikivn/clickhousectl/internal/pkg/api"
	"github.com/tikivn/clickhousectl/internal/pkg/infra"
	"github.com/tikivn/clickhousectl/internal/pkg/repo/mysql"
	"github.com/tikivn/clickhousectl/internal/pkg/service"
	"github.com/tikivn/clickhousectl/internal/setting"
)

var deps = wire.NewSet(
	api.GraphSet,
	infra.GraphSet,
	mysql.GraphSet,
	service.GraphSet,
	setting.GraphSet,
)

var GraphSet = wire.NewSet(
	deps,
	NewApp,
	NewHttpServer,
)
