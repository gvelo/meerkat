package plan

import (
	"meerkat/internal/query/exec"
	"meerkat/internal/query/logical"
	"meerkat/internal/query/parse"
	"meerkat/internal/schema"
)

type QueryManager struct {
	schema    schema.Schema
	optimizer exec.Optimizer
	executor  exec.Executor
	parser    parse.Parser
}

func NewQueryManager(schema schema.Schema, p parse.Parser, o exec.Optimizer, e exec.Executor) (*QueryManager, error) {
	return &QueryManager{schema: schema, optimizer: o, executor: e}, nil
}

func (qm *QueryManager) Query(q string) (*exec.ResultSet, error) {
	var err error
	var nodes []logical.Node
	var eNode exec.ExNode

	nodes, err = qm.parser.Parse(qm.schema, q)
	if err != nil {
		return nil, err
	}

	eNode, err = qm.optimizer.OptimizeQuery(nodes)
	if err != nil {
		return nil, err
	}

	//TODO(sebad): VER que devuelvo.
	_, err = qm.executor.ExecuteQuery(eNode)

	return nil, nil
}
