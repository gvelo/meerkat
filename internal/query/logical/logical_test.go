package logical

import (
	"github.com/stretchr/testify/assert"
	"meerkat/internal/query/parser"
	"testing"
)

func TestTransform(t *testing.T) {
	ast, err := parser.Parse("T|where A>10|summarize avg=avg(C) by colA=bin(d)")
	if err != nil {
		t.Fatal(err)
	}
	actual := ToLogical(ast)[0]

	expected := &SummarizeOp{
		Agg: []*AggExpr{{
			ColName: "avg",
			Expr: &CallExpr{
				FuncName: "avg",
				ArgList:  []Node{&ColRefExpr{Name: "C"}},
			},
		}},
		By: []*ColumnExpr{{
			ColName: "colA",
			Expr: &CallExpr{
				FuncName: "bin",
				ArgList:  []Node{&ColRefExpr{Name: "d"}},
			},
		}},
		Child: &FilterOp{
			Predicate: &BinaryExpr{
				LeftExpr:  &ColRefExpr{Name: "A"},
				Op:        GTR,
				RightExpr: &LiteralExpr{Value: 10},
			},
			Child: &SourceOp{TableName: "T"},
		},
	}

	assert.Equal(t, expected, actual)
}
