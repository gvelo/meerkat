// Code generated from /Users/sdominguez/desa/workspace_go/eventdb/internal/query/MqlParser.g4 by ANTLR 4.7.2. DO NOT EDIT.

package mql_parser // MqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseMqlParserVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseMqlParserVisitor) VisitStart(ctx *StartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitIdentifierList(ctx *IdentifierListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitAgrupTypes(ctx *AgrupTypesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitAgrupCall(ctx *AgrupCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitStringLiteral(ctx *StringLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitDecimalLiteral(ctx *DecimalLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitFloatLiteral(ctx *FloatLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitBoolLiteral(ctx *BoolLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitIdentifier(ctx *IdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitBinaryExpression(ctx *BinaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitParenExpression(ctx *ParenExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitComparatorExpression(ctx *ComparatorExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitComparator(ctx *ComparatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitBinary(ctx *BinaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitCommands(ctx *CommandsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitWhereCommand(ctx *WhereCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitSelectCommand(ctx *SelectCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitRenameCommand(ctx *RenameCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitStatCommand(ctx *StatCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitBinCommand(ctx *BinCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitFieldCommand(ctx *FieldCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitDedupCommand(ctx *DedupCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitSortCommand(ctx *SortCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitHeadCommand(ctx *HeadCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseMqlParserVisitor) VisitCompleteCommand(ctx *CompleteCommandContext) interface{} {
	return v.VisitChildren(ctx)
}
