package plan

import (
	"meerkat/internal/query/mql_parser"
	"meerkat/internal/schema"
)

type QueryManager struct {
	schema    schema.Schema
	optimizer Optimizer
	executor  Executor
}

type ResultSet struct {
	rowAffected int
	rowScanned  int
	colsName    []string
	cols        []interface{}
}

func NewQueryManager(schema schema.Schema) (*QueryManager, error) {

	return &QueryManager{schema: schema}, nil

}

func (qm *QueryManager) Query(q string) (*ResultSet, error) {

	_, err := mql_parser.Parse(qm.schema, q)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
