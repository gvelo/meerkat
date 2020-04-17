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

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(s string) (*TabularStmt, error) {

	// TODO(gvelo): handle error properly

	p.scanner = NewScanner(s)
	p.next()
	stmt := p.parseTabularStmt()
	return stmt, nil
}

func (p *Parser) expect(tokenTypes ...TokenType) {

	for _, tokenType := range tokenTypes {
		if p.token.Type == tokenType {
			return
		}
	}

	p.errorf("expected %v found %v", tokenTypes, p.token.Type)

}

func (p *Parser) error(msg string) {

	err := ParseError{
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

	i := len(p.token.Literal)

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

	t, err := time.Parse(time.RFC3339Nano, p.token.Literal)

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

	// unreachable
	return -1

}

func (p *Parser) parseTabularStmt() *TabularStmt {
	return &TabularStmt{p.parseTabularExpr()}
}

func (p *Parser) parseTabularExpr() *TabularExpr {

	tExpr := &TabularExpr{
		Source:    p.parseLit(IDENT),
		TabularOp: p.parseTabularOperators(),
	}

	return tExpr

}

func (p *Parser) parseTabularOperators() []interface{} {

	var tOps []interface{}

	for p.token.Type == PIPE {
		p.next()
		op := p.parseTabularOperator()
		tOps = append(tOps, op)
	}

	p.expect(EOF)

	return tOps

}

func (p *Parser) parseTabularOperator() interface{} {

	p.expect(IDENT)

	switch p.token.Literal {
	case "where":
	case "take":
		p.next()
		return p.parseLimit()
	case "limit":
		p.next()
		return p.parseLimit()
	case "count":
		p.next()
		return p.parseCount()
	case "summarize":
	case "sort":
	case "extend":
	default:
		p.errorf("unknown tabular operator %q", p.token.Literal)
	}

	return nil

}

func (p *Parser) parseLimit() *LimitOp {

	lit := p.parseLit(INT)

	return &LimitOp{
		NumberOfRows: lit,
	}

}

func (p *Parser) parseCount() *CountOp {
	return &CountOp{}
}
