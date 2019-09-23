package plan

import (
	"meerkat/internal/query/logical"
	"meerkat/internal/query/mql_parser"
	"testing"
)

func buildIndex() *logical.Projection {

	sql := "indexname=name  c1=1 and  ( c2>2  or c3=3 ) "
	return mql_parser.Parse(sql)

}

func TestMeerkatExecutor_ExecuteQuery(t *testing.T) {

	exe := NewMeerkatExecutor()

	exe.ExecuteQuery(buildIndex())

}
