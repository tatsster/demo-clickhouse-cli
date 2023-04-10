package operator

import (
	"strings"

	"github.com/tikivn/clickhousectl/internal/pkg/entity"
	"github.com/tikivn/clickhousectl/internal/utils/wildcard"
	"github.com/tikivn/ultrago/u_validator"
)

type OperatorInsertData struct {
	OrgId     string `json:"org_id" validate:"required"`

	Data [][]RecordData `json:"data" validate:"required"`
}

type RecordData struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}

func (op *OperatorInsertData) Validate() error {
	return u_validator.Struct(op)
}

func (op *OperatorInsertData) OperatorName() string {
	return "OperatorInsertData"
}

func (op *OperatorInsertData) ServerId() string {
	return op.OrgId
}

func (op *OperatorInsertData) columnName(recordIdx int) string {
	names := make([]string, 0)

	for _, data := range op.Data[recordIdx] {
		names = append(names, data.Name)
	}
	return strings.Join(names, ", ")
}

func (op *OperatorInsertData) ValueField(tableDDL entity.TableDDL, recordIdx int) []interface{} {
	values := []interface{}{}

	for _, data := range op.Data[recordIdx] {
		typeValue := tableDDL.TypeByName(data.Name)
		value, _ := tableDDL.ConvertFromStr(typeValue, data.Value)
		// value := data.Value

		values = append(values, value)
	}
	return values
}

// Return with list of ? for prepate stage
func (op *OperatorInsertData) prepareQuery(recordIdx int) string {
	names := make([]string, 0)
	for range op.Data[recordIdx] {
		names = append(names, "?")
	}
	return strings.Join(names, ", ")
}

func (op *OperatorInsertData) wildCard(server *entity.ClickHouseServer, tableDDL entity.TableDDL, recordIdx int) map[string]string {
	wildcard := make(map[string]string)
	wildcard["database"] = tableDDL.DBName
	wildcard["table"] = tableDDL.TableName
	wildcard["columns"] = op.columnName(recordIdx)
	wildcard["values"] = op.prepareQuery(recordIdx)
	return wildcard
}

func (op *OperatorInsertData) Stmt(server *entity.ClickHouseServer, tableDDL entity.TableDDL, recordIdx int) (string, error) {
	wildcardMap := op.wildCard(server, tableDDL, recordIdx)

	distributeSql := `INSERT INTO {{.database}}.{{.table}} ({{.columns}}) VALUES ({{.values}})`
	// Need to run on every server once

	executeDistribute, err := wildcard.ParseWildCard(distributeSql, wildcardMap)
	return executeDistribute, err
}
