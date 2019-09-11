// Copyright 2019 The Meerkat Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rel

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"meerkat/internal/query/mql_parser"
	"meerkat/internal/tools"
)

type MQLListener struct {
	*antlr.BaseParseTreeListener
	builder Builder
	lexer   *mql_parser.MqlLexer
}

func (l *MQLListener) EnterStart(c *mql_parser.StartContext) {

}

func (l *MQLListener) EnterIdentifierList(c *mql_parser.IdentifierListContext) {

}

func (l *MQLListener) EnterAgrupTypes(c *mql_parser.AgrupTypesContext) {

}

func (l *MQLListener) EnterAgrupCall(c *mql_parser.AgrupCallContext) {

}

func (l *MQLListener) EnterStringLiteral(c *mql_parser.StringLiteralContext) {

}

func (l *MQLListener) EnterDecimalLiteral(c *mql_parser.DecimalLiteralContext) {

}

func (l *MQLListener) EnterFloatLiteral(c *mql_parser.FloatLiteralContext) {

}

func (l *MQLListener) EnterBoolLiteral(c *mql_parser.BoolLiteralContext) {

}

func (l *MQLListener) EnterIdentifier(c *mql_parser.IdentifierContext) {

}

func (l *MQLListener) EnterBinaryExpression(c *mql_parser.BinaryExpressionContext) {

}

func (l *MQLListener) EnterParenExpression(c *mql_parser.ParenExpressionContext) {

}

func (l *MQLListener) EnterComparatorExpression(c *mql_parser.ComparatorExpressionContext) {

}

func (l *MQLListener) EnterComparator(c *mql_parser.ComparatorContext) {

}

func (l *MQLListener) EnterBinary(c *mql_parser.BinaryContext) {

}

func (l *MQLListener) EnterCommands(c *mql_parser.CommandsContext) {

}

func (l *MQLListener) EnterWhereCommand(c *mql_parser.WhereCommandContext) {

}

func (l *MQLListener) EnterSelectCommand(c *mql_parser.SelectCommandContext) {

}

func (l *MQLListener) EnterRenameCommand(c *mql_parser.RenameCommandContext) {

}

func (l *MQLListener) EnterStatCommand(c *mql_parser.StatCommandContext) {

}

func (l *MQLListener) EnterBinCommand(c *mql_parser.BinCommandContext) {

}

func (l *MQLListener) EnterFieldCommand(c *mql_parser.FieldCommandContext) {

}

func (l *MQLListener) EnterDedupCommand(c *mql_parser.DedupCommandContext) {

}

func (l *MQLListener) EnterSortCommand(c *mql_parser.SortCommandContext) {

}

func (l *MQLListener) EnterHeadCommand(c *mql_parser.HeadCommandContext) {

}

func (l *MQLListener) EnterCompleteCommand(c *mql_parser.CompleteCommandContext) {

}

func (l *MQLListener) ExitStart(c *mql_parser.StartContext) {

}

func (l *MQLListener) ExitIdentifierList(c *mql_parser.IdentifierListContext) {

}

func (l *MQLListener) ExitAgrupTypes(c *mql_parser.AgrupTypesContext) {

}

func (l *MQLListener) ExitAgrupCall(c *mql_parser.AgrupCallContext) {

}

func (l *MQLListener) ExitStringLiteral(c *mql_parser.StringLiteralContext) {

}

func (l *MQLListener) ExitDecimalLiteral(c *mql_parser.DecimalLiteralContext) {

}

func (l *MQLListener) ExitFloatLiteral(c *mql_parser.FloatLiteralContext) {

}

func (l *MQLListener) ExitBoolLiteral(c *mql_parser.BoolLiteralContext) {

}

func (l *MQLListener) ExitIdentifier(c *mql_parser.IdentifierContext) {

}

func (l *MQLListener) ExitBinaryExpression(c *mql_parser.BinaryExpressionContext) {

}

func (l *MQLListener) ExitParenExpression(c *mql_parser.ParenExpressionContext) {

}

func (l *MQLListener) ExitComparatorExpression(c *mql_parser.ComparatorExpressionContext) {

}

func (l *MQLListener) ExitComparator(c *mql_parser.ComparatorContext) {

}

func (l *MQLListener) ExitBinary(c *mql_parser.BinaryContext) {

}

func (l *MQLListener) ExitCommands(c *mql_parser.CommandsContext) {

}

func (l *MQLListener) ExitWhereCommand(c *mql_parser.WhereCommandContext) {

}

func (l *MQLListener) ExitRenameCommand(c *mql_parser.RenameCommandContext) {

}

func (l *MQLListener) ExitStatCommand(c *mql_parser.StatCommandContext) {

}

func (l *MQLListener) ExitBinCommand(c *mql_parser.BinCommandContext) {

}

func (l *MQLListener) ExitFieldCommand(c *mql_parser.FieldCommandContext) {

}

func (l *MQLListener) ExitDedupCommand(c *mql_parser.DedupCommandContext) {

}

func (l *MQLListener) ExitSortCommand(c *mql_parser.SortCommandContext) {

}

func (l *MQLListener) ExitHeadCommand(c *mql_parser.HeadCommandContext) {

}

func (l *MQLListener) ExitCompleteCommand(c *mql_parser.CompleteCommandContext) {

}

func NewListener(lexer *mql_parser.MqlLexer) *MQLListener {
	l := new(MQLListener)
	l.builder = NewRelBuilder()
	l.lexer = lexer
	return l
}

func (l *MQLListener) ExitSelectCommand(c *mql_parser.SelectCommandContext) {
	i := c.GetChildren()[2]
	l.builder.Scan(fmt.Sprintf("%v", i))
	es := c.GetChildren()[3]
	if es != nil {
		f := l.buildFilters(es.(antlr.ParserRuleContext))
		l.builder.Filter(f)
	}
}

func (l *MQLListener) buildFilters(ctx antlr.ParserRuleContext) *Filter {
	if ctx == nil {
		tools.Log("Empty")
		return nil
	}
	switch ctx.(type) {
	// something AND | OR something
	case *mql_parser.BinaryExpressionContext:

		lf := ctx.(*mql_parser.BinaryExpressionContext).GetLeft()
		rg := ctx.(*mql_parser.BinaryExpressionContext).GetRight()

		op := ctx.(*mql_parser.BinaryExpressionContext).GetOp()

		tools.Logf("Bin %v" , op.GetText())


		leftFilter := l.buildFilters(lf)
		rightFilter := l.buildFilters(rg)


		f := NewFilter(leftFilter, parseOperator(op.GetText()), rightFilter)

		return f

	//  ( binary expression )
	case *mql_parser.ParenExpressionContext:

		p := ctx.(*mql_parser.ParenExpressionContext)
		c := p.GetChild(1)
		f := l.buildFilters(c.(antlr.ParserRuleContext))
		f.Group = true
		return f

	// something > = != ...  something
	case *mql_parser.ComparatorExpressionContext:

		lf := ctx.(*mql_parser.ComparatorExpressionContext).GetLeft()

		rg := ctx.(*mql_parser.ComparatorExpressionContext).GetRight()

		op := ctx.(*mql_parser.ComparatorExpressionContext).GetOp()

		tools.Logf("Comp %v" , op.GetText())

		f := NewFilter(l.builder.CreateExpresion(lf), parseOperator(op.GetText()), l.builder.CreateExpresion(rg))
		return f

	default:
		tools.Logf("type %v ignored", ctx)

	}

	return nil
}

func (l *MQLListener) GetTree() *ParsedTree {
	return l.builder.Build()
}

func parseOperator(s string) Operator {
	switch s {
	case "=":
		return EQ
	case ">":
		return GT
	case ">=":
		return GEQT
	case "<":
		return LT
	case "<=":
		return LEQT
	case "!=":
		return DST
	case "or":
		return OR
	case "and":
		return AND
	default:
		return -1
	}
}
