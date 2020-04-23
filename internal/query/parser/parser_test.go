// Copyright 2020 The Meerkat Authors
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

package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type PrintVisitor struct {
	stack []string
}

func (v *PrintVisitor) push(s string) {
	v.stack = append(v.stack, s)
}

func (v *PrintVisitor) pushf(format string, a ...interface{}) {
	v.push(fmt.Sprintf(format, a...))
}

func (v *PrintVisitor) pop() string {
	i := len(v.stack) - 1
	s := v.stack[i]
	v.stack = v.stack[:i]
	return s
}

func (v *PrintVisitor) VisitPre(n Node) Node {
	return n
}

func (v *PrintVisitor) VisitPost(n Node) Node {

	switch node := n.(type) {
	case *LitExpr:
		v.pushf("( LitExpr %v [%v] %T )", node.Token.Type, node.Value, node.Value)
	case *UnaryExpr:
		v.pushf("( UnaryExpr op %v expr %v )", node.Op.Type, v.pop())
	case *BinaryExpr:
		// pushed left to right
		r := v.pop()
		l := v.pop()
		v.pushf("( BinaryExpr op %v LeftExpr %v RightExpr %v )", node.Op.Type, l, r)
	case *CallExpr:
		children := v.printStack(len(node.ArgList))
		v.pushf("( CallExpr FuncName %v ArgList ( %v ) )", node.FuncName.Value, children)
	case *WhereOp:
		v.pushf("( WhereOp Predicate %v )", v.pop())
	case *CountOp:
		v.push("( CountOp )")
	case *LimitOp:
		v.pushf("( LimitOp NumberOfRows %v )", node.NumberOfRows.Value)
	case *SummarizeOp:
		by := v.printStack(len(node.By))
		agg := v.printStack(len(node.Agg))
		v.pushf("( SummarizeOp Agg %v By %v )", agg, by)
	case *ExtendOp:
		col := v.printStack(len(node.Columns))
		v.pushf("( ExtendOp Columns %v )", col)
	case *ColumnExpr:
		colName := ""
		if node.ColName != nil {
			colName = node.ColName.Value.(string)
		}
		v.pushf("( ColumnExpr ColName %v Expr %v )", colName, v.pop())
	}

	return n

}

func (v *PrintVisitor) printStack(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		if i > 0 {
			s = " " + s
		}
		// pushed left to right.
		s = v.pop() + s
	}
	return s
}

type TestCase struct {
	name     string
	input    string
	expected string
	values   []interface{}
	isError  bool
	fun      func(*Parser) Node
}

var cases = []TestCase{
	{
		name:     "LitExpr IDENT",
		input:    "meerkat",
		expected: "( LitExpr IDENT [meerkat] string )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "LitExpr INT",
		input:    "10",
		expected: "( LitExpr INT [10] int )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "LitExpr FLOAT",
		input:    "0.42",
		expected: "( LitExpr FLOAT [0.42] float64 )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "LitExpr TIME",
		input:    "2d",
		expected: "( LitExpr TIME [%v] int )",
		values:   []interface{}{int(time.Hour * 24 * 2)},
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "LitExpr TIME with prefix",
		input:    "time(2hours)",
		expected: "( LitExpr TIME [%v] int )",
		values:   []interface{}{int(time.Hour * 2)},
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "LitExpr invalid time unit",
		input:    "2dd",
		expected: "",
		isError:  true,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "LitExpr String",
		input:    `"meerkat"`,
		expected: "( LitExpr STRING [meerkat] string )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "LitExpr String escaped",
		input:    `"\"meerkat"`,
		expected: `( LitExpr STRING ["meerkat] string )`,
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "LitExpr String not terminated",
		input:    `"meerkat`,
		expected: "",
		isError:  true,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "LitExpr DATETIME truncated ISO",
		input:    "datetime(2020-01-01)",
		expected: "( LitExpr DATETIME [%v] int )",
		values:   []interface{}{time.Date(2020, time.January, 01, 0, 0, 0, 0, time.UTC).UnixNano()},
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "UnaryExpr",
		input:    "+1",
		expected: "( UnaryExpr op ADD expr ( LitExpr INT [1] int ) )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "CallExpr",
		input:    "ago(1)",
		expected: "( CallExpr FuncName ago ArgList ( ( LitExpr INT [1] int ) ) )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "CallExpr multiple args",
		input:    "bin(1,2)",
		expected: "( CallExpr FuncName bin ArgList ( ( LitExpr INT [1] int ) ( LitExpr INT [2] int ) ) )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "CallExpr with TIME arg",
		input:    "ago(1h)",
		expected: "( CallExpr FuncName ago ArgList ( ( LitExpr TIME [%v] int ) ) )",
		values:   []interface{}{int(time.Hour)},
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "BinaryExpr",
		input:    "1+2",
		expected: "( BinaryExpr op ADD LeftExpr ( LitExpr INT [1] int ) RightExpr ( LitExpr INT [2] int ) )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "BinaryExpr Precedence 1",
		input:    "1+2+3",
		expected: "( BinaryExpr op ADD LeftExpr ( BinaryExpr op ADD LeftExpr ( LitExpr INT [1] int ) RightExpr ( LitExpr INT [2] int ) ) RightExpr ( LitExpr INT [3] int ) )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "BinaryExpr Precedence MUL",
		input:    "1+2*3",
		expected: "( BinaryExpr op ADD LeftExpr ( LitExpr INT [1] int ) RightExpr ( BinaryExpr op MUL LeftExpr ( LitExpr INT [2] int ) RightExpr ( LitExpr INT [3] int ) ) )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "BinaryExpr parenthesized",
		input:    "(1+2)*3",
		expected: "( BinaryExpr op MUL LeftExpr ( BinaryExpr op ADD LeftExpr ( LitExpr INT [1] int ) RightExpr ( LitExpr INT [2] int ) ) RightExpr ( LitExpr INT [3] int ) )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseExpr() },
	},
	{
		name:     "WhereOp",
		input:    "where ColA > ago(2h)",
		expected: "( WhereOp Predicate ( BinaryExpr op GTR LeftExpr ( LitExpr IDENT [ColA] string ) RightExpr ( CallExpr FuncName ago ArgList ( ( LitExpr TIME [%v] int ) ) ) ) )",
		isError:  false,
		values:   []interface{}{int(2 * time.Hour)},
		fun:      func(p *Parser) Node { return p.parseTabularOperator() },
	},
	{
		name:     "CountOp",
		input:    "count",
		expected: "( CountOp )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseTabularOperator() },
	},
	{
		name:     "LimitOp",
		input:    "limit 10",
		expected: "( LimitOp NumberOfRows 10 )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseTabularOperator() },
	},
	{
		name:     "ExtendOp",
		input:    "extend v=ColA + 1,ColB+2",
		expected: "( ExtendOp Columns ( ColumnExpr ColName v Expr ( BinaryExpr op ADD LeftExpr ( LitExpr IDENT [ColA] string ) RightExpr ( LitExpr INT [1] int ) ) ) ( ColumnExpr ColName  Expr ( BinaryExpr op ADD LeftExpr ( LitExpr IDENT [ColB] string ) RightExpr ( LitExpr INT [2] int ) ) ) )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseTabularOperator() },
	},
	{
		name:     "SummarizeOp",
		input:    "summarize avg=Avg(ColA) by ColB",
		expected: "( SummarizeOp Agg ( CallExpr FuncName Avg ArgList ( ( LitExpr IDENT [ColA] string ) ) ) By ( ColumnExpr ColName  Expr ( LitExpr IDENT [ColB] string ) ) )",
		isError:  false,
		fun:      func(p *Parser) Node { return p.parseTabularOperator() },
	},
}

func TestParser(t *testing.T) {

	for _, testCase := range cases {

		t.Run(testCase.name, func(t *testing.T) {

			node, err := parse(testCase.fun, testCase.input)

			if err != nil {
				t.Log(err)
				assert.IsType(t, &ParseError{}, err)
				assert.True(t, testCase.isError, "unexpected ParseError")
				return
			}

			v := &PrintVisitor{}

			Walk(node, v)

			actual := v.pop()

			t.Log(actual)

			expected := fmt.Sprintf(testCase.expected, testCase.values...)

			assert.Equal(t, expected, actual)

		})

	}

}

func parse(fun func(*Parser) Node, s string) (n Node, err error) {

	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(*ParseError)
			if ok {
				err = e
				return
			} else {
				panic(r)
			}
		}
	}()

	p, e := NewParser(s)

	if e != nil {
		return nil, e
	}

	n = fun(p)

	return n, nil

}
