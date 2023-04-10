package operator

import "github.com/tikivn/clickhousectl/internal/pkg/entity"

type OperatorStatus string

const (
	OperatorStatus_EXECUTING       OperatorStatus = "executing"
	OperatorStatus_SUCCESS         OperatorStatus = "success"
	OperatorStatus_ROLLBACK        OperatorStatus = "rollback"
	OperatorStatus_ROLLBACK_FAILED OperatorStatus = "rollback_failed"
)

type OperatorTable interface {
	Stmts(servers []*entity.ClickHouseServer) (map[string][]string, error)
	OperatorName() string
	ServerId() string
}

type OperatorColumn interface {
	Stmts(servers []*entity.ClickHouseServer, serverId string) (map[string][]string, error)
	OperatorName() string
	ServerId() string
}

type OperatorArgs interface {
	ValueField(entity.TableDDL, int) []interface{}
	Stmt(server *entity.ClickHouseServer, tableDDL entity.TableDDL, recordIdx int) (string, error)
	OperatorName() string
	ServerId() string
}

type Operator struct {
	Id           string              `json:"id"`
	Owner        string              `json:"owner"`
	Name         string              `json:"name"`
	ExecuteStmts map[string][]string `json:"execute_stmts"`
	Status       OperatorStatus      `json:"status"`
	ExecutedAt   int64               `json:"executed_at"`
}

type ArgumentOperator struct {
	Id          string         `json:"id"`
	Owner       string         `json:"owner"`
	Name        string         `json:"name"`
	ExecuteStmt string         `json:"execute_stmt"` // Almost always is Prepare stmt
	Argument    []interface{}  `json:"args"`
	Status      OperatorStatus `json:"status"`
	ExecutedAt  int64          `json:"executed_at"`
}
