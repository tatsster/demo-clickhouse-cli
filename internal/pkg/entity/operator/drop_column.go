package operator

import (
	"github.com/tikivn/clickhousectl/internal/pkg/entity"
	"github.com/tikivn/clickhousectl/internal/utils/wildcard"
	"github.com/tikivn/ultrago/u_validator"
)

type OperatorDropColumn struct {
	OrgId     string `json:"org_id" validate:"required"`
	DBName    string `json:"database" validate:"required"`
	TableName string `json:"table" validate:"required"`
	Column    string `json:"column" validate:"required"`
}

func (op *OperatorDropColumn) Validate() error {
	return u_validator.Struct(op)
}

func (op *OperatorDropColumn) wildCard(server *entity.ClickHouseServer) map[string]string {
	wildcard := make(map[string]string)
	wildcard["database"] = op.DBName
	wildcard["table"] = op.TableName
	wildcard["cluster"] = server.Cluster
	wildcard["column"] = op.Column
	return wildcard
}

func (op *OperatorDropColumn) Stmts(servers []*entity.ClickHouseServer, serverId string) (map[string][]string, error) {
	execute := make(map[string][]string)

	for _, server := range servers {
		var replicaSql, distributeSql string
		wildcardMap := op.wildCard(server)
		// Replica Sql only need to run once
		execute[server.Id] = make([]string, 0)

		for _, shard := range server.Shards {
			wildcardMap["shard"] = shard

			replicaSql = `ALTER TABLE {{.shard}}.{{.cluster}}_{{.table}} DROP COLUMN IF EXISTS {{.column}}`
			// Execute create table
			executeReplica, err := wildcard.ParseWildCard(replicaSql, wildcardMap)
			if err != nil {
				return nil, err
			}
			execute[server.Id] = append(execute[server.Id], executeReplica)
		}

		distributeSql = `ALTER TABLE {{.database}}.{{.table}} DROP COLUMN IF EXISTS {{.column}}`
		// Need to run on every server
		executeDistribute, err := wildcard.ParseWildCard(distributeSql, wildcardMap)
		if err != nil {
			return nil, err
		}
		execute[server.Id] = append(execute[server.Id], executeDistribute)
	}

	return execute, nil
}

func (op *OperatorDropColumn) OperatorName() string {
	return "OperatorDropColumn"
}

func (op *OperatorDropColumn) ServerId() string {
	return op.OrgId
}
