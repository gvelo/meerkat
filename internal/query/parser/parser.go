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
	"strconv"
	"strings"
	"time"
	"unicode"
)

var u2ns = map[string]int{
	"ns":           int(time.Nanosecond),
	"nanosecond":   int(time.Nanosecond),
	"nanoseconds":  int(time.Nanosecond),
	"microsecond":  int(time.Microsecond),
	"microseconds": int(time.Microsecond),
	"ms":           int(time.Millisecond),
	"millisecond":  int(time.Millisecond),
	"milliseconds": int(time.Millisecond),
	"s":            int(time.Second),
	"second":       int(time.Second),
	"seconds":      int(time.Second),
	"m":            int(time.Minute),
	"minute":       int(time.Minute),
	"minutes":      int(time.Minute),
	"h":            int(time.Hour),
	"hour":         int(time.Hour),
	"hours":        int(time.Hour),
	"d":            int(time.Hour * 24),
	"day":          int(time.Hour * 24),
	"days":         int(time.Hour * 24),
}

type Parser struct {
	scanner *Scanner
	token   Token
}

func NewParser(s string) (p *Parser, err error) {

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

	p = &Parser{
		scanner: NewScanner(s),
	}

	p.next()

	return
}

func (p *Parser) Parse() (stmt *TabularStmt, err error) {

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

	stmt = p.parseTabularStmt()

	return stmt, nil
}

func (p *Parser) expect(tokenTypes ...TokenType) {

	for _, tokenType := range tokenTypes {
		if p.token.Type == tokenType {
			return
		}
	}

	p.errorf("expected %v found %v", tokenTypes, p.token)

}

func (p *Parser) error(msg string) {

	err := &ParseError{
		msg:    msg,
		offset: p.token.Offset,
		line:   p.token.Line,
		column: p.token.Column,
	}

	panic(err)

}

func (p *Parser) errorf(msg string, a ...interface{}) {
	p.error(fmt.Sprintf(msg, a...))
}

func (p *Parser) next() {

	p.token = p.scanner.Scan()

	// skip comments.
	for p.token.Type == COMMENT {
		p.scanner.next()
	}

}

func (p *Parser) parseLit(expect ...TokenType) *LitExpr {

	if expect != nil {
		p.expect(expect...)
	} else {
		if !p.token.IsLiteral() {
			p.errorf("expect literal found %v", p.token.Type)
		}
	}

	lit := &LitExpr{
		Token: p.token,
		Value: p.parseLitValue(p.token),
	}

	p.next()

	return lit

}

func (p *Parser) parseLitValue(t Token) interface{} {
	switch t.Type {
	case IDENT:
		return t.Literal
	case INT:
		return p.parseIntLit()
	case FLOAT:
		return p.parseFloatLit()
	case STRING:
		return p.parseStringLit()
	case TIME:
		return p.parseTimeLit()
	case DATETIME:
		return p.parseDateTimeLit()
	default:
		panic(fmt.Sprintf("%s is not a literal token", t.Type))
	}
	return nil
}

func (p *Parser) parseIntLit() int {
	i, err := strconv.Atoi(p.token.Literal)
	if err != nil {
		p.errorf("invalid INT format %s", p.token.Literal)
	}
	return i
}

func (p *Parser) parseFloatLit() float64 {
	f, err := strconv.ParseFloat(p.token.Literal, 64)
	if err != nil {
		p.errorf("invalid FLOAT format %s", p.token.Literal)
	}
	return f
}

func (p *Parser) parseStringLit() string {
	s, err := strconv.Unquote(p.token.Literal)
	if err != nil {
		p.error(err.Error())
	}
	return s
}

func (p *Parser) parseTimeLit() int {

	// find the postfix start

	i := len(p.token.Literal) - 1

	for i > 0 {

		ch := rune(p.token.Literal[i])

		if unicode.IsLetter(ch) {
			i--
			continue
		}

		break

	}

	postfix := p.token.Literal[i+1:]

	if postfix == "" {
		postfix = "day"
	}

	numLit := strings.TrimSpace(p.token.Literal[:i+1])
	num, err := strconv.Atoi(numLit)

	if err != nil {
		p.error(err.Error())
	}

	if m, ok := u2ns[postfix]; ok {
		return num * m
	} else {
		p.errorf("invalid time unit: %s", postfix)
	}

	return -1

}

func (p *Parser) parseDateTimeLit() int {

	// TODO(gvelo): currently this time representation is not compatible
	//  with Kusto date representation which is measured in ticks
	//  ( 100 nanosecond ) since 12:00 midnight, January 1, 0001 A.D.
	//  (C.E.) in the GregorianCalendar calendar.
	//  We need to provide a date parser ( and a datetime type )
	//  compatible with Kusto datetime type.

	// TODO(gvelo): refactor this brittle parsing.

	// TODO(gvelo) this is a truncated ISO
	t, err := time.Parse("2006-01-02", p.token.Literal)

	if err == nil {
		return int(t.UnixNano())
	}

	t, err = time.Parse(time.RFC3339, p.token.Literal)

	if err == nil {
		return int(t.UnixNano())
	}

	t, err = time.Parse(time.RFC3339Nano, p.token.Literal)

	if err == nil {
		return int(t.UnixNano())
	}

	t, err = time.Parse(time.RFC822, p.token.Literal)

	if err == nil {
		return int(t.UnixNano())
	}

	t, err = time.Parse(time.RFC850, p.token.Literal)

	if err == nil {
		return int(t.UnixNano())
	}

	p.errorf("cannot parse datetime %s", p.token.Literal)

	return -1

}

func (p *Parser) parseTabularStmt() *TabularStmt {
	return &TabularStmt{p.parseTabularExpr()}
}

func (p *Parser) parseTabularExpr() *TabularExpr {

	tExpr := &TabularExpr{
		Source:    p.parseLit(IDENT),
		TabularOp: p.parseTabularOperatorList(),
	}

	return tExpr

}

func (p *Parser) parseTabularOperatorList() []Node {

	var tOps []Node

	for p.token.Type == PIPE {
		p.next()
		op := p.parseTabularOperator()
		tOps = append(tOps, op)
	}

	p.expect(PIPE, EOF)

	return tOps

}

func (p *Parser) parseTabularOperator() Node {

	p.expect(IDENT)
	t := p.token
	p.next()

	switch t.Literal {
	case "where":
		return p.parseWhereOp()
	case "take":
		return p.parseLimitOp()
	case "limit":
		return p.parseLimitOp()
	case "count":
		return p.parseCountOp()
	case "summarize":
		return p.parseSummarizeOp()
	case "sort":
		return p.parseSortOp()
	case "order":
		return p.parseSortOp()
	case "extend":
		return p.parseExtendOp()
	case "project":
		return p.parseProjectOp()
	default:
		p.errorf("unknown tabular operator %q", t.Literal)
	}

	return nil

}

func (p *Parser) parseWhereOp() *WhereOp {

	return &WhereOp{
		Predicate: p.parseExpr(),
	}

}

func (p *Parser) parseLimitOp() *LimitOp {

	lit := p.parseLit(INT)

	return &LimitOp{
		NumberOfRows: lit,
	}

}

func (p *Parser) parseCountOp() *CountOp {
	return &CountOp{}
}

// exp = unaryExp || binaryExp
func (p *Parser) parseExpr() Node {
	return p.parseBinaryExpr(LowestPrec + 1)
}

// call = literal(IDENT) "(" { [ expr [","] } ]  ")"
func (p *Parser) parseCallExpr(funcName *LitExpr) *CallExpr {

	callExpr := &CallExpr{
		FuncName: funcName,
	}

	p.expect(LPAREN)
	p.next()

	// parse argument list
	for p.token.Type != RPAREN {
		expr := p.parseExpr()
		callExpr.ArgList = append(callExpr.ArgList, expr)
		if p.token.Type == COMMA {
			p.next()
			continue
		}
	}

	p.expect(RPAREN)
	p.next()

	return callExpr

}

// primaryExpr = literal(IDEN/INT/FLOAT/DATE/DATETIME) | callExpr
func (p *Parser) parsePrimaryExpr() Node {

	switch p.token.Type {

	case LPAREN:
		p.next()
		expr := p.parseExpr()
		p.expect(RPAREN)
		p.next()
		return expr

	case IDENT:

		lit := p.parseLit()

		if p.token.Type == LPAREN {
			return p.parseCallExpr(lit)
		}

		return lit

	}

	if p.token.IsLiteral() {
		return p.parseLit()
	}

	p.errorf("expected expr found %v", p.token.Type)

	return nil
}

// unaryExpr = primaryExpr | unary_op primaryExpr
func (p *Parser) parseUnaryExpr() Node {

	if p.token.Type == ADD || p.token.Type == SUB {
		t := p.token
		p.next()
		return &UnaryExpr{
			Op:   t,
			Expr: p.parseUnaryExpr(),
		}
	}

	return p.parsePrimaryExpr()

}

func (p *Parser) parseBinaryExpr(prec int) Node {

	l := p.parseUnaryExpr()

	for {

		if p.token.Precedence() < prec {
			return l
		}

		if !p.token.IsOperator() {
			p.errorf("expect operator got %v", p.token)
		}

		op := p.token
		p.next()

		r := p.parseBinaryExpr(op.Precedence() + 1)

		l = &BinaryExpr{
			LeftExpr:  l,
			Op:        op,
			RightExpr: r,
		}

	}

}

func (p *Parser) parseExtendOp() Node {

	return &ExtendOp{
		Columns: p.parseColumnExprList(),
	}

}

func (p *Parser) parseColumnExprList() []*ColumnExpr {

	var l []*ColumnExpr

	for {
		expr := p.parseColumnExpr()
		l = append(l, expr)
		if p.token.Type != COMMA {
			break
		}
		p.next()
	}

	return l

}

func (p *Parser) parseColumnExpr() *ColumnExpr {

	nameOrExp := p.parseExpr()

	if p.token.Type == ASSIGN {

		name, ok := nameOrExp.(*LitExpr)

		if !ok {
			p.error("invalid column name")
			return nil
		}

		if name.Token.Type != IDENT {
			p.errorf("expect IDENT found %v", name.Token)
			return nil
		}

		p.next()

		return &ColumnExpr{
			ColName: name,
			Expr:    p.parseExpr(),
		}

	}

	// TODO(gvelo) we need to provide a synthetic column
	//  name here ( or maybe create one in the resolver
	//  given that in the parser we don't have enough
	//  information to avoid name coalitions  )

	return &ColumnExpr{
		ColName: nil,
		Expr:    nameOrExp,
	}

}

func (p *Parser) parseSortOp() *SortOp {

	if p.token.Type != IDENT {
		p.errorf("expect IDENT got %v", p.token)
	}

	if p.token.Literal != "by" {
		p.errorf("expect \"by\" keyword got %v", p.token)
	}

	p.next()

	sortOp := &SortOp{}

	for {

		expr := p.parseSortExpr()

		sortOp.SortExpr = append(sortOp.SortExpr, expr)

		if p.token.Type != COMMA {
			break
		}

		p.next()

	}

	return sortOp

}

func (p *Parser) parseSortExpr() *SortExpr {

	sortExpr := &SortExpr{
		Expr:      p.parseExpr(),
		Asc:       false, // default to desc
		NullFirst: false, // default to nulls Last
	}

	if p.token.Type == IDENT && p.token.Literal == "asc" {
		sortExpr.Asc = true
		p.next()

	} else if p.token.Type == IDENT && p.token.Literal == "desc" {
		sortExpr.Asc = false
		p.next()
	}

	if p.token.Type == IDENT && p.token.Literal == "nulls" {

		p.next()

		if p.token.Type == IDENT && p.token.Literal == "first" {
			sortExpr.NullFirst = true
			p.next()

		} else if p.token.Type == IDENT && p.token.Literal == "last" {
			sortExpr.NullFirst = false
			p.next()
		} else {
			p.errorf("expected [\"first\"|\"last\"] got %q", p.token.Literal)
		}

	}

	return sortExpr

}

func (p *Parser) parseSummarizeOp() *SummarizeOp {

	op := &SummarizeOp{
		Agg: p.parseAggExprList(),
	}

	if p.token.Type == IDENT && p.token.Literal == "by" {
		p.next()
		op.By = p.parseColumnExprList()
	}

	return op

}

func (p *Parser) parseAggExprList() []*AggExpr {

	var l []*AggExpr

	for {
		expr := p.parseAggExpr()
		l = append(l, expr)
		if p.token.Type != COMMA {
			break
		}
		p.next()
	}

	return l

}

func (p *Parser) parseAggExpr() *AggExpr {

	lit := p.parseLit(IDENT)

	if p.token.Type == ASSIGN {

		p.next()

		funcName := p.parseLit(IDENT)

		callExpr := p.parseCallExpr(funcName)

		return &AggExpr{
			ColName: lit,
			Expr:    callExpr,
		}

	}

	return &AggExpr{
		ColName: nil,
		Expr:    p.parseCallExpr(lit),
	}

}

func (p *Parser) parseProjectOp() *ProjectOp {

	return &ProjectOp{
		Columns: p.parseColumnExprList(),
	}

}
