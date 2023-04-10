package operator

import (
	"github.com/tikivn/clickhousectl/internal/pkg/entity"
	"github.com/tikivn/clickhousectl/internal/utils/wildcard"
	"github.com/tikivn/ultrago/u_validator"
)

type OperatorDropTable struct {
	OrgId     string `json:"org_id" validate:"required"`
	DBName    string `json:"database" validate:"required"`
	TableName string `json:"table" validate:"required"`
}

func (op *OperatorDropTable) Validate() error {
	return u_validator.Struct(op)
}

func (op *OperatorDropTable) wildcard(server *entity.ClickHouseServer) map[string]string {
	wildCard := make(map[string]string)
	wildCard["database"] = op.DBName
	wildCard["table"] = op.TableName
	wildCard["cluster"] = server.Cluster
	return wildCard
}

func (op *OperatorDropTable) Stmts(servers []*entity.ClickHouseServer) (map[string][]string, error) {

	creates := make(map[string][]string)

	for _, server := range servers {
		var dropReplicaSql, dropDistributedSql string
		wildcardMap := op.wildcard(server)
		creates[server.Id] = make([]string, 0)

		for _, shard := range server.Shards {
			wildcardMap["shard"] = shard

			dropReplicaSql = `DROP TABLE IF EXISTS {{.shard}}.{{.cluster}}_{{.table}} SYNC;`

			// Drop  table
			dropReplica, err := wildcard.ParseWildCard(dropReplicaSql, wildcardMap)
			if err != nil {
				return nil, err
			}
			creates[server.Id] = append(creates[server.Id], dropReplica)
		}
		// Skip sync time between clickhouse and zookeeper
		dropDistributedSql = `DROP TABLE IF EXISTS {{.database}}.{{.table}} NO DELAY;`

		// Execute drop distributed
		dropDistributed, err := wildcard.ParseWildCard(dropDistributedSql, wildcardMap)
		if err != nil {
			return nil, err
		}
		creates[server.Id] = append(creates[server.Id], dropDistributed)
	}

	return creates, nil
}

func (op *OperatorDropTable) OperatorName() string {
	return "OperatorDropTable"
}

func (op *OperatorDropTable) ServerId() string {
	return op.OrgId
}
