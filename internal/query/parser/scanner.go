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
//
// some functions on this lexer were copied from go/scanner/scanner.go
// published under the license below.
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parser

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type ParseError struct {
	msg    string
	offset int
	line   int
	column int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("ParseError: %v line:%v column:%v", e.msg, e.line, e.column)
}

type Scanner struct {
	src        []byte // source
	ch         rune   // current character
	start      int    // token start offset
	offset     int    // character offset
	rdOffset   int    // reading offset (position after current character)
	lineOffset int    // current line offset
	line       int    // line number
}

func NewScanner(s string) *Scanner {

	return &Scanner{
		src: []byte(s),
		ch:  ' ',
	}

}

// Read the next Unicode char into s.ch.
// s.ch < 0 means end-of-file.
//
func (s *Scanner) next() {

	if s.rdOffset < len(s.src) {

		s.offset = s.rdOffset

		if s.ch == '\n' {
			s.newLine()
		}

		r, w := rune(s.src[s.rdOffset]), 1

		switch {

		case r == 0:
			s.error(s.offset, "illegal character NUL")

		case r >= utf8.RuneSelf:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				s.error(s.offset, "illegal UTF-8 encoding")
			}

		}

		s.rdOffset += w
		s.ch = r

	} else {

		s.offset = len(s.src)

		if s.ch == '\n' {
			s.newLine()
		}

		s.ch = -1 // eof

	}
}

func (s *Scanner) newLine() {
	s.lineOffset = s.offset
	s.line++
}

func (s *Scanner) error(offs int, msg string) {

	err := &ParseError{
		msg:    msg,
		offset: offs,
		line:   s.line,
		column: offs - s.lineOffset,
	}

	panic(err)

}

func (s *Scanner) errorf(offs int, format string, args ...interface{}) {
	s.error(offs, fmt.Sprintf(format, args...))
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_' || ch >= utf8.RuneSelf && unicode.IsLetter(ch)
}

func isDigit(ch rune) bool {
	return isDecimal(ch) || ch >= utf8.RuneSelf && unicode.IsDigit(ch)
}

func (s *Scanner) scanIdentifier() string {
	start := s.offset
	for isLetter(s.ch) || isDigit(s.ch) || s.ch == '~' {
		s.next()
	}
	return string(s.src[start:s.offset])
}

func (s *Scanner) scanOp() string {
	// ! already consumed.
	start := s.offset - 1
	for isLetter(s.ch) || s.ch == '~' {
		s.next()
	}
	return string(s.src[start:s.offset])
}

func digitVal(ch rune) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch - '0')
	case 'a' <= lower(ch) && lower(ch) <= 'f':
		return int(lower(ch) - 'a' + 10)
	}
	return 16 // larger than any legal digit val
}

func lower(ch rune) rune     { return ('a' - 'A') | ch } // returns lower-case ch iff ch is ASCII letter
func isDecimal(ch rune) bool { return '0' <= ch && ch <= '9' }

// scanEscape parses an escape sequence where rune is the accepted
// escaped quote. In case of a syntax error, it stops at the offending
// character (without consuming it) and returns false. Otherwise
// it returns true.
func (s *Scanner) scanEscape(quote rune) bool {
	offs := s.offset

	var n int
	var base, max uint32
	switch s.ch {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', quote:
		s.next()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		s.next()
		n, base, max = 2, 16, 255
	case 'u':
		s.next()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		s.next()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		msg := "unknown escape sequence"
		if s.ch < 0 {
			msg = "escape sequence not terminated"
		}
		s.error(offs, msg)
		return false
	}

	var x uint32
	for n > 0 {
		d := uint32(digitVal(s.ch))
		if d >= base {
			msg := fmt.Sprintf("illegal character %#U in escape sequence", s.ch)
			if s.ch < 0 {
				msg = "escape sequence not terminated"
			}
			s.error(s.offset, msg)
			return false
		}
		x = x*base + d
		s.next()
		n--
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		s.error(offs, "escape sequence is invalid Unicode code point")
		return false
	}

	return true
}

func (s *Scanner) scanString(quote rune) string {

	start := s.offset - 1

	for {

		ch := s.ch

		if ch == '\n' || ch < 0 {
			s.error(s.offset, "string literal not terminated")
			break
		}

		s.next()

		if ch == quote {
			break
		}

		if ch == '\\' {
			s.scanEscape(quote)
		}

	}

	return string(s.src[start:s.offset])

}

func (s *Scanner) scanNumber() (string, TokenType) {

	tt := INT
	lastDec := s.offset
	start := s.offset

	for {

		if isDecimal(s.ch) {
			lastDec = s.offset
			s.next()
			continue
		}

		if s.ch == '.' {
			tt = FLOAT
			s.next()
			continue
		}

		if isLetter(s.ch) {
			s.next()
			continue
		}

		break

	}

	suffix := string(s.src[lastDec+1 : s.offset])

	if len(suffix) > 0 {
		if _, ok := timeSuffix[suffix]; ok {
			tt = TIME
		} else {
			s.error(lastDec+1, "invalid suffix")
		}
	}

	return string(s.src[start:s.offset]), tt
}

func (s *Scanner) scanComment() string {

	start := s.offset + 1

	//-style comment
	if s.ch == '/' {
		// consume until EOL or EOF
		for s.ch != '\n' && s.ch > 0 {
			s.next()
		}
		return string(s.src[start:s.offset])
	}

	/*-style comment */
	for s.ch > 0 {
		ch := s.ch
		s.next()
		if ch == '*' && s.ch == '/' {
			return string(s.src[start : s.offset-2])
		}
	}

	s.error(s.offset, "comment not terminated")

	return ""

}

// Scan scans the next token and returns the token
// The source end is indicated by token.EOF.
//
// This is a non-recoverable lexer and will panic at the first error found,
// error information will be provided using a ParseError type
//
func (s *Scanner) Scan() (tok Token) {

	s.skipWhitespace()

	// current token start
	s.start = s.offset

	// determine token value
	switch ch := s.ch; {

	case isLetter(ch):
		lit := s.scanIdentifier()
		tok = s.resolveLiteral(lit)

	case isDecimal(ch) || ch == '.' && isDecimal(rune(s.peek())):
		lit, tt := s.scanNumber()
		tok = s.newToken(tt, lit)

	default:
		s.next() // always make progress
		switch ch {
		case -1:
			tok = s.newToken(EOF, "")
		case '"':
			lit := s.scanString('"')
			tok = s.newToken(STRING, lit)
		case '\'':
			lit := s.scanString('\'')
			tok = s.newToken(STRING, lit)
		case ':':
			tok = s.newToken(COLON, ":")
		case '.':
			// fractions starting with a '.' are handled by outer switch
			if s.ch == '.' {
				s.next() // consulme '.'
				tok = s.newToken(RANGE, "..")
			}
			tok = s.newToken(PERIOD, ".")
		case ',':
			tok = s.newToken(COMMA, ",")
		case ';':
			tok = s.newToken(SEMICOLON, ";")
		case '(':
			tok = s.newToken(LPAREN, "(")
		case ')':
			tok = s.newToken(RPAREN, ")")
		case '[':
			tok = s.newToken(LBRACK, "[")
		case ']':
			tok = s.newToken(RBRACK, "]")
		case '{':
			tok = s.newToken(LBRACE, "{")
		case '}':
			tok = s.newToken(RBRACE, "}")
		case '+':
			tok = s.newToken(ADD, "+")
		case '-':
			tok = s.newToken(SUB, "-")
		case '*':
			tok = s.newToken(MUL, "*")
		case '/':
			if s.ch == '/' || s.ch == '*' {
				lit := s.scanComment()
				tok = s.newToken(COMMENT, lit)
			} else {
				tok = s.newToken(QUO, "/")
			}
		case '%':
			tok = s.newToken(REM, "%")
		case '<':
			if s.ch == '=' {
				s.next()
				tok = s.newToken(LEQ, "<=")
			} else {
				tok = s.newToken(LSS, "<")
			}
		case '>':
			if s.ch == '=' {
				s.next()
				tok = s.newToken(GEQ, ">=")
			} else {
				tok = s.newToken(GTR, ">")
			}
		case '=':
			switch s.ch {
			case '=':
				tok = s.newToken(EQL, "==")
				s.next()
			case '~':
				tok = s.newToken(EQL_CI, "=~")
				s.next()
			default:
				tok = s.newToken(ASSIGN, "=")
			}
		case '!':
			switch s.ch {
			case '=':
				tok = s.newToken(NEQ, "!=")
				s.next()
			case '~':
				tok = s.newToken(NEQ_CI, "!~")
				s.next()
			default:
				if isLetter(s.ch) {
					lit := s.scanOp()
					tok = s.resolveLiteral(lit)
					if tok.Type == IDENT {
						s.errorf(s.start, "unknown op : %v", lit)
					}
				} else {
					s.errorf(s.start, "invalid token: %v", s.ch)
				}
			}
		case '|':
			tok = s.newToken(PIPE, "|")
		default:
			s.errorf(s.start, "invalid token %v", s.ch)
		}
	}

	return

}

func (s *Scanner) scanDatetime() string {

	start := s.offset

	for {

		ch := s.ch

		if ch == '\n' || ch < 0 {
			s.error(s.offset, "datetime or timespan literal not terminated")
		}

		s.next()

		if ch == ')' {
			break
		}

	}

	return string(s.src[start : s.offset-1])

}

func (s *Scanner) resolveLiteral(lit string) Token {

	t := getTokenType(lit)

	if t == DATETIME || t == TIME {

		if s.ch != '(' {
			return s.newToken(IDENT, lit)
		}

		s.next() // consume '('
		lit = s.scanDatetime()

		return s.newToken(t, lit)

	}

	return s.newToken(t, lit)

}

func (s *Scanner) newToken(t TokenType, lit string) Token {
	return Token{
		Type:    t,
		Literal: lit,
		Offset:  s.start,
		Column:  s.start - s.lineOffset,
		Line:    s.line,
	}
}

// peek returns the byte following the most recently read character without
// advancing the scanner. If the scanner is at EOF, peek returns 0.
func (s *Scanner) peek() byte {
	if s.rdOffset < len(s.src) {
		return s.src[s.rdOffset]
	}
	return 0
}
