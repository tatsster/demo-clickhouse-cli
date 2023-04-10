package operator

import (
	"github.com/tikivn/clickhousectl/internal/pkg/entity"
	"github.com/tikivn/clickhousectl/internal/utils/wildcard"
	"github.com/tikivn/ultrago/u_validator"
)

type OperatorDeleteData struct {
	OrgId     string  `json:"org_id" validate:"required"`
	DBName    string  `json:"database" validate:"required"`
	TableName string  `json:"table" validate:"required"`
	Queries   []Query `json:"query" validate:"required"`
}

type Query struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}

func (op *OperatorDeleteData) Validate() error {
	return u_validator.Struct(op)
}

func (op *OperatorDeleteData) OperatorName() string {
	return "OperatorDeleteData"
}

func (op *OperatorDeleteData) ServerId() string {
	return op.OrgId
}

func (op *OperatorDeleteData) wildCard(server *entity.ClickHouseServer) map[string]string {
	wildcard := make(map[string]string)
	wildcard["database"] = op.DBName
	wildcard["table"] = op.TableName
	wildcard["cluster"] = server.Cluster
	// TODO: Add query string here
	// wildcard["query"] = op.queries
	return wildcard
}

func (op *OperatorDeleteData) Stmt(server *entity.ClickHouseServer) (string, error) {
	wildcardMap := op.wildCard(server)

	// Clickhouse doesnt support DELETE on distributed table yet
	distributeSql := `ALTER TABLE {{.database}}.{{.table}} ON CLUSTER {{.cluster}} DELETE WHERE {.query}`

	executeDistribute, err := wildcard.ParseWildCard(distributeSql, wildcardMap)
	return executeDistribute, err
}
