package entity

import (
	"strconv"
	"strings"
)

type TableDDL struct {
	OrgId     string `json:"org_id" validate:"required"`
	DBName    string `json:"database" validate:"required"`
	TableName string `json:"table" validate:"required"`

	Columns []DataType `json:"columns" validate:"required"`
}

type DataType struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"expression" validate:"required"`
}

// TODO: If have any different type, add more case here

func (table *TableDDL) TypeByName(name string) string {
	col := table.DataDefinition(name)

	switch {
	case strings.Contains(col.Type, "Int"):
		return "Int"
	case strings.Contains(col.Type, "Float"):
		return "Float"
	case strings.Contains(col.Type, "String"):
		return "String"
	default:
		return col.Type
	}
}

func (table *TableDDL) ConvertFromStr(typeValue string, value string) (interface{}, error) {
	switch typeValue {
	case "Int":
		return strconv.Atoi(value)

	case "Float":
		return strconv.ParseFloat(value, 64)
	default:
		return value, nil
	}
}

func (table *TableDDL) DataDefinition(name string) DataType {
	for _, col := range table.Columns {
		if col.Name == name {
			return col
		}
	}
	return DataType{}
}
