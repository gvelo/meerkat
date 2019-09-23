package plan

import (
	"meerkat/internal/query/mql_parser"
)

type QueryManager struct {
	dbPath    string
	optimizer Optimizer
	executor  Executor
}

//TODO(sebad): pagination ?
type ResultSet struct {
	rowAffected int
	rowScanned  int
	colsName    []string
	cols        []interface{}
}

func NewQueryManager(path string, o Optimizer, e Executor) *QueryManager {

	return &QueryManager{dbPath: path,
		executor:  e,
		optimizer: o,
	}

}

func (qm *QueryManager) Query(q string) *ResultSet {

	t := mql_parser.Parse(q)

	p := qm.optimizer.OptimizeQuery(t)

	return qm.executor.ExecuteQuery(p)
}
