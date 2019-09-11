package plan

import (
	"meerkat/internal/query/rel"
	"testing"
)

func buildIndex() *rel.ParsedTree {

	parser := rel.NewMqlParser()
	sql := "indexname=name  c1=1 and  ( c2>2  or c3=3 ) "

	return parser.Parse(sql)

}

func TestMeerkatExecutor_ExecuteQuery(t *testing.T) {

	exe := NewMeerkatExecutor()

	exe.ExecuteQuery(buildIndex())

}
