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

package mql_parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"meerkat/internal/query/logical"
	"meerkat/internal/tools"
	"strconv"
)

type MQLListener struct {
	*antlr.BaseParseTreeListener
	builder Builder
	lexer   *MqlLexer
}

func newListener(lexer *MqlLexer) *MQLListener {
	l := new(MQLListener)
	l.builder = NewRelBuilder()
	l.lexer = lexer
	return l
}

func (l *MQLListener) ExitRexCommand(c *RexCommandContext) {
	r := c.regex.GetText()
	var f = "_raw"
	if c.rexfield != nil {
		f = c.rexfield.GetText()
	}
	l.builder.Regex(f, r)
}

func (l *MQLListener) ExitBucketCommand(c *BucketCommandContext) {
	t := c.GetSpan()
	l.builder.Span(l.builder.CreateExpresion(t))
}

func (l *MQLListener) ExitSortCommand(c *SortCommandContext) {
	children := c.GetSList().GetChildren()
	li := make([]string, 0)
	for i, _ := range children {
		ch, ok := children[i].(*SortContext)
		if ok {
			li = append(li, ch.GetField().GetText())
			if ch.GetDirection() != nil {
				li = append(li, ch.GetDirection().GetText())
			} else {
				li = append(li, "asc")
			}
		}
	}
	l.builder.Sort(li...)
}

func (l *MQLListener) ExitTopCommand(c *TopCommandContext) {
	i, _ := strconv.Atoi(c.GetLimit().GetText())
	l.builder.Limit(i)
}

func (l *MQLListener) ExitSelectCommand(c *SelectCommandContext) {

	f := make([]interface{}, 0)
	for _, ctx := range c.GetChildren() {

		if c.GetIndex() == ctx {
			l.builder.Scan(c.GetIndex().GetName().GetText())
			continue
		}

		f = append(f, ctx)

	}

	rf := &logical.RootFilter{} // root filter
	for i := 0; i < len(f); i++ {
		f := l.buildFilters(f[i].(antlr.ParserRuleContext))
		rf.RootFilter = f
	}
	l.builder.Filter(rf)

}

func (l *MQLListener) buildFilters(ctx antlr.ParserRuleContext) *logical.Filter {
	if ctx == nil {
		tools.Log("Empty")
		return nil
	}
	switch ctx.(type) {
	// something AND | OR something
	case *BinaryExpressionContext:

		lf := ctx.(*BinaryExpressionContext).GetLeft()
		rg := ctx.(*BinaryExpressionContext).GetRight()

		op := ctx.(*BinaryExpressionContext).GetOp()

		tools.Logf("Bin %v", op.GetText())

		leftFilter := l.buildFilters(lf)
		rightFilter := l.buildFilters(rg)

		f := logical.NewFilter(leftFilter, parseOperator(op.GetText()), rightFilter)

		return f

	case *TimeExpressionContext:

		rg := ctx.(*TimeExpressionContext).GetRight()

		op := ctx.(*TimeExpressionContext).GetOp()

		tools.Logf("Time exp %v", op.GetText())
		e := &logical.Exp{
			ExpType: logical.STRING,
			Value:   "_time",
		}
		f := logical.NewFilter(e, parseOperator(op.GetText()), l.builder.CreateExpresion(rg))

		return f

	//  ( binary expression )
	case *ParenExpressionContext:

		p := ctx.(*ParenExpressionContext)
		c := p.GetChild(1)
		f := l.buildFilters(c.(antlr.ParserRuleContext))
		f.Group = true
		return f

	// something > = != ...  something
	case *ComparatorExpressionContext:

		lf := ctx.(*ComparatorExpressionContext).GetLeft()

		rg := ctx.(*ComparatorExpressionContext).GetRight()

		op := ctx.(*ComparatorExpressionContext).GetOp()

		tools.Logf("Comp %v", op.GetText())

		f := logical.NewFilter(l.builder.CreateExpresion(lf), parseOperator(op.GetText()), l.builder.CreateExpresion(rg))
		return f

	default:
		tools.Logf("type %v ignored", ctx)

	}

	return nil
}

func (l *MQLListener) GetTree() []logical.Node {
	return l.builder.Build()
}

func parseOperator(s string) logical.Operator {
	switch s {
	case "=":
		return logical.EQ
	case ">":
		return logical.GT
	case ">=":
		return logical.GEQT
	case "<":
		return logical.LT
	case "<=":
		return logical.LEQT
	case "!=":
		return logical.DST
	case "or":
		return logical.OR
	case "and":
		return logical.AND
	default:
		return -1
	}
}

func (l *MQLListener) EnterSortCommand(c *SortCommandContext) {

}

func (l *MQLListener) EnterBucketCommand(c *BucketCommandContext) {
}

func (l *MQLListener) EnterIndexExpression(c *IndexExpressionContext) {

}

func (l *MQLListener) ExitIndexExpression(c *IndexExpressionContext) {

}

func (l *MQLListener) EnterTimeExpression(c *TimeExpressionContext) {
}

func (l *MQLListener) ExitTimeExpression(c *TimeExpressionContext) {

}

func (l *MQLListener) EnterSort(c *SortContext) {
}

func (l *MQLListener) EnterSortList(c *SortListContext) {
}

func (l *MQLListener) ExitSort(c *SortContext) {
}

func (l *MQLListener) ExitSortList(c *SortListContext) {
}

func (l *MQLListener) EnterTopCommand(c *TopCommandContext) {
}

func (l *MQLListener) EnterStart(c *StartContext) {

}

func (l *MQLListener) EnterIdentifierList(c *IdentifierListContext) {

}

func (l *MQLListener) EnterAgrupTypes(c *AgrupTypesContext) {

}

func (l *MQLListener) EnterAgrupCall(c *AgrupCallContext) {

}

func (l *MQLListener) EnterStringLiteral(c *StringLiteralContext) {

}

func (l *MQLListener) EnterDecimalLiteral(c *DecimalLiteralContext) {

}

func (l *MQLListener) EnterFloatLiteral(c *FloatLiteralContext) {

}

func (l *MQLListener) EnterBoolLiteral(c *BoolLiteralContext) {

}

func (l *MQLListener) EnterIdentifier(c *IdentifierContext) {

}

func (l *MQLListener) EnterBinaryExpression(c *BinaryExpressionContext) {

}

func (l *MQLListener) EnterParenExpression(c *ParenExpressionContext) {

}

func (l *MQLListener) EnterComparatorExpression(c *ComparatorExpressionContext) {

}

func (l *MQLListener) EnterComparator(c *ComparatorContext) {

}

func (l *MQLListener) EnterBinary(c *BinaryContext) {

}

func (l *MQLListener) EnterCommands(c *CommandsContext) {

}

func (l *MQLListener) EnterWhereCommand(c *WhereCommandContext) {

}

func (l *MQLListener) EnterSelectCommand(c *SelectCommandContext) {

}

func (l *MQLListener) EnterRenameCommand(c *RenameCommandContext) {

}

func (l *MQLListener) EnterStatCommand(c *StatCommandContext) {

}

func (l *MQLListener) EnterFieldCommand(c *FieldCommandContext) {

}

func (l *MQLListener) EnterDedupCommand(c *DedupCommandContext) {

}

func (l *MQLListener) EnterCompleteCommand(c *CompleteCommandContext) {

}

func (l *MQLListener) ExitStart(c *StartContext) {

}

func (l *MQLListener) ExitIdentifierList(c *IdentifierListContext) {

}

func (l *MQLListener) ExitAgrupTypes(c *AgrupTypesContext) {

}

func (l *MQLListener) ExitAgrupCall(c *AgrupCallContext) {

}

func (l *MQLListener) ExitStringLiteral(c *StringLiteralContext) {

}

func (l *MQLListener) ExitDecimalLiteral(c *DecimalLiteralContext) {

}

func (l *MQLListener) ExitFloatLiteral(c *FloatLiteralContext) {

}

func (l *MQLListener) ExitBoolLiteral(c *BoolLiteralContext) {

}

func (l *MQLListener) ExitIdentifier(c *IdentifierContext) {

}

func (l *MQLListener) ExitBinaryExpression(c *BinaryExpressionContext) {

}

func (l *MQLListener) ExitParenExpression(c *ParenExpressionContext) {

}

func (l *MQLListener) ExitComparatorExpression(c *ComparatorExpressionContext) {

}

func (l *MQLListener) ExitComparator(c *ComparatorContext) {

}

func (l *MQLListener) ExitBinary(c *BinaryContext) {

}

func (l *MQLListener) ExitCommands(c *CommandsContext) {

}

func (l *MQLListener) ExitWhereCommand(c *WhereCommandContext) {

}

func (l *MQLListener) ExitRenameCommand(c *RenameCommandContext) {

}

func (l *MQLListener) ExitStatCommand(c *StatCommandContext) {
	t := c.GetF()
	list := make([]string, 0)
	for i := c.GetField().GetTokenIndex(); i < c.GetChildCount(); i++ {
		if s := c.GetChildren()[i].(antlr.Token).GetText(); s != "," {
			list = append(list, s)
		}
	}
	l.builder.Aggregate(t.GetText(), list)
}

func (l *MQLListener) ExitFieldCommand(c *FieldCommandContext) {

}

func (l *MQLListener) ExitDedupCommand(c *DedupCommandContext) {

}

func (l *MQLListener) ExitCompleteCommand(c *CompleteCommandContext) {

}

func (l *MQLListener) EnterRexCommand(c *RexCommandContext) {

}
