package query

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"meerkat/internal/query/mql_parser"
)

type MQLVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *MQLVisitor) VisitStart(ctx *mql_parser.StartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *MQLVisitor) VisitIdentifier_list(ctx *mql_parser.Identifier_listContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *MQLVisitor) VisitAgrupTypes(ctx *mql_parser.AgrupTypesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *MQLVisitor) VisitAgrupCall(ctx *mql_parser.AgrupCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *MQLVisitor) VisitStat_expresion(ctx *mql_parser.Stat_expresionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *MQLVisitor) VisitComparators(ctx *mql_parser.ComparatorsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *MQLVisitor) VisitExpression(ctx *mql_parser.ExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *MQLVisitor) VisitExpressions(ctx *mql_parser.ExpressionsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *MQLVisitor) VisitWhere_expresion(ctx *mql_parser.Where_expresionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *MQLVisitor) VisitSelect_expression(ctx *mql_parser.Select_expressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *MQLVisitor) VisitFields(ctx *mql_parser.FieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *MQLVisitor) VisitA(ctx *mql_parser.AContext) interface{} {
	return v.VisitChildren(ctx)
}
