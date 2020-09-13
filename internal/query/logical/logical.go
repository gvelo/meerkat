package logical

import (
	"meerkat/internal/query/parser"
)

//go:generate protoc -I . -I ../../../build/proto/ --plugin ../../../build/protoc-gen-gogofaster --gogofaster_out=plugins=grpc,paths=source_relative:.  ./logical.proto

func ToLogical(query *parser.TabularStmt) []Node {
	t := &transform{}
	t.transform(query)
	return t.outputNodes
}

type transform struct {
	child       Node
	outputNodes []Node
}

func (t *transform) transform(n parser.Node) {

	switch node := n.(type) {
	case *parser.TabularStmt:
		t.transform(node.TabularExpr)
	case *parser.TabularExpr:
		t.child = &SourceOp{TableName: node.Source.Value.(string)}
		for _, tabOp := range node.TabularOp {
			switch op := tabOp.(type) {
			case *parser.WhereOp:
				t.child = t.transformWhereOp(t.child, op)
			case *parser.SummarizeOp:
				t.child = t.transformSummarizeOp(t.child, op)
			case *parser.LimitOp:
				t.child = t.transformLimitOp(t.child, op)
			default:
				panic("unknown operator")
			}
		}
		t.outputNodes = append(t.outputNodes, t.child)
	}

}

func (t *transform) transformWhereOp(child Node, whereOp *parser.WhereOp) *FilterOp {
	return &FilterOp{
		Predicate: t.transformExpr(whereOp.Predicate),
		Child:     child,
	}
}

func (t *transform) transformSummarizeOp(child Node, op *parser.SummarizeOp) Node {

	return &SummarizeOp{
		Agg:   t.transformAggExprList(op.Agg),
		By:    t.transformColumnExprList(op.By),
		Child: child,
	}

}

func (t *transform) transformLimitOp(child Node, op *parser.LimitOp) *LimitOp {
	return &LimitOp{
		NumOfRows: op.NumberOfRows.Value.(int),
		Child:     child,
	}
}

func (t *transform) transformExpr(expr parser.Node) Node {

	switch e := expr.(type) {
	case *parser.BinaryExpr:
		return t.transformBinaryExpr(e)
	case *parser.UnaryExpr:
		return t.transformUnaryExpr(e)
	case *parser.CallExpr:
		return t.transformCallExpr(e)
	case *parser.LitExpr:
		return t.transformLitExpr(e)
	default:
		panic("unknown expr node")
	}

}

func (t *transform) transformBinaryExpr(expr *parser.BinaryExpr) Node {
	return &BinaryExpr{
		LeftExpr:  t.transformExpr(expr.LeftExpr),
		Op:        t.transformOp(expr.Op.Type),
		RightExpr: t.transformExpr(expr.RightExpr),
	}
}

func (t *transform) transformUnaryExpr(expr *parser.UnaryExpr) Node {
	return &UnaryExpr{
		Op:   t.transformOp(expr.Op.Type),
		Expr: t.transformExpr(expr),
	}
}

func (t *transform) transformCallExpr(expr *parser.CallExpr) *CallExpr {

	argList := make([]Node, len(expr.ArgList))

	for i, arg := range expr.ArgList {
		argList[i] = t.transformExpr(arg)
	}

	return &CallExpr{
		FuncName: expr.FuncName.Value.(string),
		ArgList:  argList,
	}

}

func (t *transform) transformOp(token parser.TokenType) Operator {
	if op, found := tokenToOp[token]; found {
		return op
	}
	panic("unknown operator")
}

func (t *transform) transformLitExpr(expr *parser.LitExpr) Node {
	switch expr.Token.Type {
	case parser.IDENT:
		return &ColRefExpr{Name: expr.Value.(string)}
	case parser.INT, parser.FLOAT, parser.STRING, parser.TIME, parser.DATETIME, parser.BOOL:
		return &LiteralExpr{Value: expr.Value}
	default:
		panic("invalid literal type")
	}
}

func (t *transform) transformAggExpr(expr *parser.AggExpr) *AggExpr {
	return &AggExpr{
		ColName: expr.ColName.Value.(string),
		Expr:    t.transformCallExpr(expr.Expr),
	}
}

func (t *transform) transformColumnExpr(expr *parser.ColumnExpr) *ColumnExpr {
	return &ColumnExpr{
		ColName: expr.ColName.Value.(string),
		Expr:    t.transformExpr(expr.Expr),
	}
}

func (t *transform) transformAggExprList(exprList []*parser.AggExpr) []*AggExpr {
	agg := make([]*AggExpr, len(exprList))
	for i, expr := range exprList {
		agg[i] = t.transformAggExpr(expr)
	}
	return agg
}

func (t *transform) transformColumnExprList(exprList []*parser.ColumnExpr) []*ColumnExpr {
	colExpr := make([]*ColumnExpr, len(exprList))
	for i, expr := range exprList {
		colExpr[i] = t.transformColumnExpr(expr)
	}
	return colExpr
}
