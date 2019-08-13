// Code generated from /Users/sdominguez/desa/workspace_go/eventdb/internal/query/MqlParser.g4 by ANTLR 4.7.2. DO NOT EDIT.

package mql_parser // MqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// MqlParserListener is a complete listener for a parse tree produced by MqlParser.
type MqlParserListener interface {
	antlr.ParseTreeListener

	// EnterStart is called when entering the start production.
	EnterStart(c *StartContext)

	// EnterIdentifierList is called when entering the identifierList production.
	EnterIdentifierList(c *IdentifierListContext)

	// EnterAgrupTypes is called when entering the agrupTypes production.
	EnterAgrupTypes(c *AgrupTypesContext)

	// EnterAgrupCall is called when entering the agrupCall production.
	EnterAgrupCall(c *AgrupCallContext)

	// EnterStringLiteral is called when entering the stringLiteral production.
	EnterStringLiteral(c *StringLiteralContext)

	// EnterDecimalLiteral is called when entering the decimalLiteral production.
	EnterDecimalLiteral(c *DecimalLiteralContext)

	// EnterFloatLiteral is called when entering the floatLiteral production.
	EnterFloatLiteral(c *FloatLiteralContext)

	// EnterBoolLiteral is called when entering the boolLiteral production.
	EnterBoolLiteral(c *BoolLiteralContext)

	// EnterIdentifier is called when entering the identifier production.
	EnterIdentifier(c *IdentifierContext)

	// EnterBinaryExpression is called when entering the binaryExpression production.
	EnterBinaryExpression(c *BinaryExpressionContext)

	// EnterParenExpression is called when entering the parenExpression production.
	EnterParenExpression(c *ParenExpressionContext)

	// EnterComparatorExpression is called when entering the comparatorExpression production.
	EnterComparatorExpression(c *ComparatorExpressionContext)

	// EnterComparator is called when entering the comparator production.
	EnterComparator(c *ComparatorContext)

	// EnterBinary is called when entering the binary production.
	EnterBinary(c *BinaryContext)

	// EnterCommands is called when entering the commands production.
	EnterCommands(c *CommandsContext)

	// EnterWhereCommand is called when entering the whereCommand production.
	EnterWhereCommand(c *WhereCommandContext)

	// EnterSelectCommand is called when entering the selectCommand production.
	EnterSelectCommand(c *SelectCommandContext)

	// EnterRenameCommand is called when entering the renameCommand production.
	EnterRenameCommand(c *RenameCommandContext)

	// EnterStatCommand is called when entering the statCommand production.
	EnterStatCommand(c *StatCommandContext)

	// EnterBinCommand is called when entering the binCommand production.
	EnterBinCommand(c *BinCommandContext)

	// EnterFieldCommand is called when entering the fieldCommand production.
	EnterFieldCommand(c *FieldCommandContext)

	// EnterDedupCommand is called when entering the dedupCommand production.
	EnterDedupCommand(c *DedupCommandContext)

	// EnterSortCommand is called when entering the sortCommand production.
	EnterSortCommand(c *SortCommandContext)

	// EnterHeadCommand is called when entering the headCommand production.
	EnterHeadCommand(c *HeadCommandContext)

	// EnterCompleteCommand is called when entering the completeCommand production.
	EnterCompleteCommand(c *CompleteCommandContext)

	// ExitStart is called when exiting the start production.
	ExitStart(c *StartContext)

	// ExitIdentifierList is called when exiting the identifierList production.
	ExitIdentifierList(c *IdentifierListContext)

	// ExitAgrupTypes is called when exiting the agrupTypes production.
	ExitAgrupTypes(c *AgrupTypesContext)

	// ExitAgrupCall is called when exiting the agrupCall production.
	ExitAgrupCall(c *AgrupCallContext)

	// ExitStringLiteral is called when exiting the stringLiteral production.
	ExitStringLiteral(c *StringLiteralContext)

	// ExitDecimalLiteral is called when exiting the decimalLiteral production.
	ExitDecimalLiteral(c *DecimalLiteralContext)

	// ExitFloatLiteral is called when exiting the floatLiteral production.
	ExitFloatLiteral(c *FloatLiteralContext)

	// ExitBoolLiteral is called when exiting the boolLiteral production.
	ExitBoolLiteral(c *BoolLiteralContext)

	// ExitIdentifier is called when exiting the identifier production.
	ExitIdentifier(c *IdentifierContext)

	// ExitBinaryExpression is called when exiting the binaryExpression production.
	ExitBinaryExpression(c *BinaryExpressionContext)

	// ExitParenExpression is called when exiting the parenExpression production.
	ExitParenExpression(c *ParenExpressionContext)

	// ExitComparatorExpression is called when exiting the comparatorExpression production.
	ExitComparatorExpression(c *ComparatorExpressionContext)

	// ExitComparator is called when exiting the comparator production.
	ExitComparator(c *ComparatorContext)

	// ExitBinary is called when exiting the binary production.
	ExitBinary(c *BinaryContext)

	// ExitCommands is called when exiting the commands production.
	ExitCommands(c *CommandsContext)

	// ExitWhereCommand is called when exiting the whereCommand production.
	ExitWhereCommand(c *WhereCommandContext)

	// ExitSelectCommand is called when exiting the selectCommand production.
	ExitSelectCommand(c *SelectCommandContext)

	// ExitRenameCommand is called when exiting the renameCommand production.
	ExitRenameCommand(c *RenameCommandContext)

	// ExitStatCommand is called when exiting the statCommand production.
	ExitStatCommand(c *StatCommandContext)

	// ExitBinCommand is called when exiting the binCommand production.
	ExitBinCommand(c *BinCommandContext)

	// ExitFieldCommand is called when exiting the fieldCommand production.
	ExitFieldCommand(c *FieldCommandContext)

	// ExitDedupCommand is called when exiting the dedupCommand production.
	ExitDedupCommand(c *DedupCommandContext)

	// ExitSortCommand is called when exiting the sortCommand production.
	ExitSortCommand(c *SortCommandContext)

	// ExitHeadCommand is called when exiting the headCommand production.
	ExitHeadCommand(c *HeadCommandContext)

	// ExitCompleteCommand is called when exiting the completeCommand production.
	ExitCompleteCommand(c *CompleteCommandContext)
}
