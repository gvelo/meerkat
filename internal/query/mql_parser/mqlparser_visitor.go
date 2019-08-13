// Code generated from /Users/sdominguez/desa/workspace_go/eventdb/internal/query/MqlParser.g4 by ANTLR 4.7.2. DO NOT EDIT.

package mql_parser // MqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by MqlParser.
type MqlParserVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by MqlParser#start.
	VisitStart(ctx *StartContext) interface{}

	// Visit a parse tree produced by MqlParser#identifierList.
	VisitIdentifierList(ctx *IdentifierListContext) interface{}

	// Visit a parse tree produced by MqlParser#agrupTypes.
	VisitAgrupTypes(ctx *AgrupTypesContext) interface{}

	// Visit a parse tree produced by MqlParser#agrupCall.
	VisitAgrupCall(ctx *AgrupCallContext) interface{}

	// Visit a parse tree produced by MqlParser#stringLiteral.
	VisitStringLiteral(ctx *StringLiteralContext) interface{}

	// Visit a parse tree produced by MqlParser#decimalLiteral.
	VisitDecimalLiteral(ctx *DecimalLiteralContext) interface{}

	// Visit a parse tree produced by MqlParser#floatLiteral.
	VisitFloatLiteral(ctx *FloatLiteralContext) interface{}

	// Visit a parse tree produced by MqlParser#boolLiteral.
	VisitBoolLiteral(ctx *BoolLiteralContext) interface{}

	// Visit a parse tree produced by MqlParser#identifier.
	VisitIdentifier(ctx *IdentifierContext) interface{}

	// Visit a parse tree produced by MqlParser#binaryExpression.
	VisitBinaryExpression(ctx *BinaryExpressionContext) interface{}

	// Visit a parse tree produced by MqlParser#parenExpression.
	VisitParenExpression(ctx *ParenExpressionContext) interface{}

	// Visit a parse tree produced by MqlParser#comparatorExpression.
	VisitComparatorExpression(ctx *ComparatorExpressionContext) interface{}

	// Visit a parse tree produced by MqlParser#comparator.
	VisitComparator(ctx *ComparatorContext) interface{}

	// Visit a parse tree produced by MqlParser#binary.
	VisitBinary(ctx *BinaryContext) interface{}

	// Visit a parse tree produced by MqlParser#commands.
	VisitCommands(ctx *CommandsContext) interface{}

	// Visit a parse tree produced by MqlParser#whereCommand.
	VisitWhereCommand(ctx *WhereCommandContext) interface{}

	// Visit a parse tree produced by MqlParser#selectCommand.
	VisitSelectCommand(ctx *SelectCommandContext) interface{}

	// Visit a parse tree produced by MqlParser#renameCommand.
	VisitRenameCommand(ctx *RenameCommandContext) interface{}

	// Visit a parse tree produced by MqlParser#statCommand.
	VisitStatCommand(ctx *StatCommandContext) interface{}

	// Visit a parse tree produced by MqlParser#binCommand.
	VisitBinCommand(ctx *BinCommandContext) interface{}

	// Visit a parse tree produced by MqlParser#fieldCommand.
	VisitFieldCommand(ctx *FieldCommandContext) interface{}

	// Visit a parse tree produced by MqlParser#dedupCommand.
	VisitDedupCommand(ctx *DedupCommandContext) interface{}

	// Visit a parse tree produced by MqlParser#sortCommand.
	VisitSortCommand(ctx *SortCommandContext) interface{}

	// Visit a parse tree produced by MqlParser#headCommand.
	VisitHeadCommand(ctx *HeadCommandContext) interface{}

	// Visit a parse tree produced by MqlParser#completeCommand.
	VisitCompleteCommand(ctx *CompleteCommandContext) interface{}
}
