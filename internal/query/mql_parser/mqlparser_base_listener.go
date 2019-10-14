// Code generated from /Users/sebad/desa/workspace_go/eventdb/internal/query/MqlParser.g4 by ANTLR 4.7.2. DO NOT EDIT.

package mql_parser // MqlParser
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseMqlParserListener is a complete listener for a parse tree produced by MqlParser.
type BaseMqlParserListener struct{}

var _ MqlParserListener = &BaseMqlParserListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseMqlParserListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseMqlParserListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseMqlParserListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseMqlParserListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart is called when production start is entered.
func (s *BaseMqlParserListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BaseMqlParserListener) ExitStart(ctx *StartContext) {}

// EnterAgrupTypes is called when production agrupTypes is entered.
func (s *BaseMqlParserListener) EnterAgrupTypes(ctx *AgrupTypesContext) {}

// ExitAgrupTypes is called when production agrupTypes is exited.
func (s *BaseMqlParserListener) ExitAgrupTypes(ctx *AgrupTypesContext) {}

// EnterAgrupCall is called when production agrupCall is entered.
func (s *BaseMqlParserListener) EnterAgrupCall(ctx *AgrupCallContext) {}

// ExitAgrupCall is called when production agrupCall is exited.
func (s *BaseMqlParserListener) ExitAgrupCall(ctx *AgrupCallContext) {}

// EnterStringLiteral is called when production stringLiteral is entered.
func (s *BaseMqlParserListener) EnterStringLiteral(ctx *StringLiteralContext) {}

// ExitStringLiteral is called when production stringLiteral is exited.
func (s *BaseMqlParserListener) ExitStringLiteral(ctx *StringLiteralContext) {}

// EnterDecimalLiteral is called when production decimalLiteral is entered.
func (s *BaseMqlParserListener) EnterDecimalLiteral(ctx *DecimalLiteralContext) {}

// ExitDecimalLiteral is called when production decimalLiteral is exited.
func (s *BaseMqlParserListener) ExitDecimalLiteral(ctx *DecimalLiteralContext) {}

// EnterFloatLiteral is called when production floatLiteral is entered.
func (s *BaseMqlParserListener) EnterFloatLiteral(ctx *FloatLiteralContext) {}

// ExitFloatLiteral is called when production floatLiteral is exited.
func (s *BaseMqlParserListener) ExitFloatLiteral(ctx *FloatLiteralContext) {}

// EnterBoolLiteral is called when production boolLiteral is entered.
func (s *BaseMqlParserListener) EnterBoolLiteral(ctx *BoolLiteralContext) {}

// ExitBoolLiteral is called when production boolLiteral is exited.
func (s *BaseMqlParserListener) ExitBoolLiteral(ctx *BoolLiteralContext) {}

// EnterIdentifier is called when production identifier is entered.
func (s *BaseMqlParserListener) EnterIdentifier(ctx *IdentifierContext) {}

// ExitIdentifier is called when production identifier is exited.
func (s *BaseMqlParserListener) ExitIdentifier(ctx *IdentifierContext) {}

// EnterIdentifierList is called when production identifierList is entered.
func (s *BaseMqlParserListener) EnterIdentifierList(ctx *IdentifierListContext) {}

// ExitIdentifierList is called when production identifierList is exited.
func (s *BaseMqlParserListener) ExitIdentifierList(ctx *IdentifierListContext) {}

// EnterSort is called when production sort is entered.
func (s *BaseMqlParserListener) EnterSort(ctx *SortContext) {}

// ExitSort is called when production sort is exited.
func (s *BaseMqlParserListener) ExitSort(ctx *SortContext) {}

// EnterSortList is called when production sortList is entered.
func (s *BaseMqlParserListener) EnterSortList(ctx *SortListContext) {}

// ExitSortList is called when production sortList is exited.
func (s *BaseMqlParserListener) ExitSortList(ctx *SortListContext) {}

// EnterIndexExpression is called when production indexExpression is entered.
func (s *BaseMqlParserListener) EnterIndexExpression(ctx *IndexExpressionContext) {}

// ExitIndexExpression is called when production indexExpression is exited.
func (s *BaseMqlParserListener) ExitIndexExpression(ctx *IndexExpressionContext) {}

// EnterBinaryExpression is called when production binaryExpression is entered.
func (s *BaseMqlParserListener) EnterBinaryExpression(ctx *BinaryExpressionContext) {}

// ExitBinaryExpression is called when production binaryExpression is exited.
func (s *BaseMqlParserListener) ExitBinaryExpression(ctx *BinaryExpressionContext) {}

// EnterParenExpression is called when production parenExpression is entered.
func (s *BaseMqlParserListener) EnterParenExpression(ctx *ParenExpressionContext) {}

// ExitParenExpression is called when production parenExpression is exited.
func (s *BaseMqlParserListener) ExitParenExpression(ctx *ParenExpressionContext) {}

// EnterRegexExpression is called when production regexExpression is entered.
func (s *BaseMqlParserListener) EnterRegexExpression(ctx *RegexExpressionContext) {}

// ExitRegexExpression is called when production regexExpression is exited.
func (s *BaseMqlParserListener) ExitRegexExpression(ctx *RegexExpressionContext) {}

// EnterComparatorExpression is called when production comparatorExpression is entered.
func (s *BaseMqlParserListener) EnterComparatorExpression(ctx *ComparatorExpressionContext) {}

// ExitComparatorExpression is called when production comparatorExpression is exited.
func (s *BaseMqlParserListener) ExitComparatorExpression(ctx *ComparatorExpressionContext) {}

// EnterTimeExpression is called when production timeExpression is entered.
func (s *BaseMqlParserListener) EnterTimeExpression(ctx *TimeExpressionContext) {}

// ExitTimeExpression is called when production timeExpression is exited.
func (s *BaseMqlParserListener) ExitTimeExpression(ctx *TimeExpressionContext) {}

// EnterComparator is called when production comparator is entered.
func (s *BaseMqlParserListener) EnterComparator(ctx *ComparatorContext) {}

// ExitComparator is called when production comparator is exited.
func (s *BaseMqlParserListener) ExitComparator(ctx *ComparatorContext) {}

// EnterBinary is called when production binary is entered.
func (s *BaseMqlParserListener) EnterBinary(ctx *BinaryContext) {}

// ExitBinary is called when production binary is exited.
func (s *BaseMqlParserListener) ExitBinary(ctx *BinaryContext) {}

// EnterFieldList is called when production fieldList is entered.
func (s *BaseMqlParserListener) EnterFieldList(ctx *FieldListContext) {}

// ExitFieldList is called when production fieldList is exited.
func (s *BaseMqlParserListener) ExitFieldList(ctx *FieldListContext) {}

// EnterCommands is called when production commands is entered.
func (s *BaseMqlParserListener) EnterCommands(ctx *CommandsContext) {}

// ExitCommands is called when production commands is exited.
func (s *BaseMqlParserListener) ExitCommands(ctx *CommandsContext) {}

// EnterWhereCommand is called when production whereCommand is entered.
func (s *BaseMqlParserListener) EnterWhereCommand(ctx *WhereCommandContext) {}

// ExitWhereCommand is called when production whereCommand is exited.
func (s *BaseMqlParserListener) ExitWhereCommand(ctx *WhereCommandContext) {}

// EnterSelectCommand is called when production selectCommand is entered.
func (s *BaseMqlParserListener) EnterSelectCommand(ctx *SelectCommandContext) {}

// ExitSelectCommand is called when production selectCommand is exited.
func (s *BaseMqlParserListener) ExitSelectCommand(ctx *SelectCommandContext) {}

// EnterRenameCommand is called when production renameCommand is entered.
func (s *BaseMqlParserListener) EnterRenameCommand(ctx *RenameCommandContext) {}

// ExitRenameCommand is called when production renameCommand is exited.
func (s *BaseMqlParserListener) ExitRenameCommand(ctx *RenameCommandContext) {}

// EnterStatCommand is called when production statCommand is entered.
func (s *BaseMqlParserListener) EnterStatCommand(ctx *StatCommandContext) {}

// ExitStatCommand is called when production statCommand is exited.
func (s *BaseMqlParserListener) ExitStatCommand(ctx *StatCommandContext) {}

// EnterBucketCommand is called when production bucketCommand is entered.
func (s *BaseMqlParserListener) EnterBucketCommand(ctx *BucketCommandContext) {}

// ExitBucketCommand is called when production bucketCommand is exited.
func (s *BaseMqlParserListener) ExitBucketCommand(ctx *BucketCommandContext) {}

// EnterFieldCommand is called when production fieldCommand is entered.
func (s *BaseMqlParserListener) EnterFieldCommand(ctx *FieldCommandContext) {}

// ExitFieldCommand is called when production fieldCommand is exited.
func (s *BaseMqlParserListener) ExitFieldCommand(ctx *FieldCommandContext) {}

// EnterDedupCommand is called when production dedupCommand is entered.
func (s *BaseMqlParserListener) EnterDedupCommand(ctx *DedupCommandContext) {}

// ExitDedupCommand is called when production dedupCommand is exited.
func (s *BaseMqlParserListener) ExitDedupCommand(ctx *DedupCommandContext) {}

// EnterRexCommand is called when production rexCommand is entered.
func (s *BaseMqlParserListener) EnterRexCommand(ctx *RexCommandContext) {}

// ExitRexCommand is called when production rexCommand is exited.
func (s *BaseMqlParserListener) ExitRexCommand(ctx *RexCommandContext) {}

// EnterSortCommand is called when production sortCommand is entered.
func (s *BaseMqlParserListener) EnterSortCommand(ctx *SortCommandContext) {}

// ExitSortCommand is called when production sortCommand is exited.
func (s *BaseMqlParserListener) ExitSortCommand(ctx *SortCommandContext) {}

// EnterTopCommand is called when production topCommand is entered.
func (s *BaseMqlParserListener) EnterTopCommand(ctx *TopCommandContext) {}

// ExitTopCommand is called when production topCommand is exited.
func (s *BaseMqlParserListener) ExitTopCommand(ctx *TopCommandContext) {}

// EnterCompleteCommand is called when production completeCommand is entered.
func (s *BaseMqlParserListener) EnterCompleteCommand(ctx *CompleteCommandContext) {}

// ExitCompleteCommand is called when production completeCommand is exited.
func (s *BaseMqlParserListener) ExitCompleteCommand(ctx *CompleteCommandContext) {}
