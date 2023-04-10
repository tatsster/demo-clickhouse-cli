package operator

import (
	"fmt"
	"strings"

	"github.com/tikivn/clickhousectl/internal/pkg/entity"
	"github.com/tikivn/clickhousectl/internal/utils/wildcard"
	"github.com/tikivn/ultrago/u_validator"
)

type OperatorCreateTable struct {
	OrgId       string   `json:"org_id" validate:"required"`
	DBName      string   `json:"database" validate:"required"`
	TableName   string   `json:"table" validate:"required"`
	Engine      string   `json:"engine" validate:"required"`
	PartitionBy []string `json:"partition_by" validate:"required"`
	OrderBy     []string `json:"order_by" validate:"required"`

	Columns []ClickHouseColumn `json:"columns" valid:"required"`
}

func (op *OperatorCreateTable) Validate() error {
	return u_validator.Struct(op)
}

func (op *OperatorCreateTable) columns() string {
	var definitions = make([]string, 0, len(op.Columns))
	for _, c := range op.Columns {
		definitions = append(definitions, c.String())
	}
	return strings.Join(definitions, ",\n")
}

func (op *OperatorCreateTable) partitionBy() string {
	return strings.Join(op.PartitionBy, ", ")
}

func (op *OperatorCreateTable) orderBy() string {
	return strings.Join(op.OrderBy, ", ")
}

func (op *OperatorCreateTable) wildcard(server *entity.ClickHouseServer) map[string]string {
	wildCard := make(map[string]string)
	wildCard["database"] = op.DBName
	wildCard["cluster"] = server.Cluster
	wildCard["table"] = op.TableName
	wildCard["engine"] = op.Engine
	wildCard["columns"] = op.columns()
	wildCard["partition_by"] = op.partitionBy()
	wildCard["order_by"] = op.orderBy()
	return wildCard
}

func (op *OperatorCreateTable) Stmts(servers []*entity.ClickHouseServer) (map[string][]string, error) {
	// layer - shard - replica autoload from macros.xml
	creates := make(map[string][]string)

	for _, server := range servers {
		var replicaSql, distributedSql string
		wildcardMap := op.wildcard(server)

		creates[server.Id] = make([]string, 0)

		for i, shard := range server.Shards {
			wildcardMap["shard"] = shard

			createDB := `CREATE DATABASE IF NOT EXISTS {{.shard}}`
			createReplica, err := wildcard.ParseWildCard(createDB, wildcardMap)
			if err != nil {
				return nil, err
			}
			creates[server.Id] = append(creates[server.Id], createReplica)

			replicaSql = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS {{.shard}}.{{.cluster}}_{{.table}} (
			{{.columns}}
			) ENGINE = Replicated{{.engine}}('/clickhouse/{{.database}}/tables/{layer}-{shard%v}/{{.cluster}}_{{.table}}', '{replica%v}') 
			PARTITION BY {{.partition_by}} 
			ORDER BY {{.order_by}};`, i+1, i+1)

			// Execute create table
			createReplica, err = wildcard.ParseWildCard(replicaSql, wildcardMap)
			if err != nil {
				return nil, err
			}
			creates[server.Id] = append(creates[server.Id], createReplica)
		}

		createDB := `CREATE DATABASE IF NOT EXISTS {{.database}}`
		createDistribute, err := wildcard.ParseWildCard(createDB, wildcardMap)
		if err != nil {
			return nil, err
		}
		creates[server.Id] = append(creates[server.Id], createDistribute)

		distributedSql = `CREATE TABLE IF NOT EXISTS {{.database}}.{{.table}} (
			{{.columns}}
		) ENGINE = Distributed({{.cluster}}, '', {{.cluster}}_{{.table}}, rand());`

		createDistribute, err = wildcard.ParseWildCard(distributedSql, wildcardMap)
		if err != nil {
			return nil, err
		}
		creates[server.Id] = append(creates[server.Id], createDistribute)
	}
	return creates, nil
}

func (op *OperatorCreateTable) OperatorName() string {
	return "OperatorCreateTable"
}

func (op *OperatorCreateTable) ServerId() string {
	return op.OrgId
}

type ClickHouseColumn struct {
	Name       string `json:"name"`
	Expression string `json:"expression"`
}

func (c *ClickHouseColumn) String() string {
	return fmt.Sprintf("%s %s", c.Name, c.Expression)
}
