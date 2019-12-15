package parse

import (
	"meerkat/internal/query/logical"
	"meerkat/internal/query/mql_parser"
	"meerkat/internal/schema"
)

type Parser interface {
	Parse(schema schema.Schema, qry string) ([]logical.Node, error)
}

type MqlParser struct {
}

func (m *MqlParser) Parse(schema schema.Schema, qry string) ([]logical.Node, error) {
	return mql_parser.Parse(schema, qry)
}
