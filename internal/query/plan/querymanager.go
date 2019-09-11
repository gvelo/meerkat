package plan

import (
	"meerkat/internal/query/rel"
)

type QueryManager struct {
	dbPath    string
	parser    rel.Parser
	optimizer Optimizer
	executor  Executor
}

//TODO(sebad): pagination ?
type ResultSet struct {
	rowAffected int
	rowScanned  int
	colsName    [] string
	cols        []interface{}
}

func NewQueryManager(path string, p rel.Parser, o Optimizer, e Executor) *QueryManager {

	return &QueryManager{dbPath: path,
		parser:    p,
		executor:  e,
		optimizer: o,
	}

}

func (qm *QueryManager) Query(q string) *ResultSet {

	t := qm.parser.Parse(q)

	p := qm.optimizer.OptimizeQuery(t)

	return qm.executor.ExecuteQuery(p)
}
