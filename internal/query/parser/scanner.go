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
	"go/token"
	"unicode"
	"unicode/utf8"
)

type ParseError struct {
	msg     string
	offset  int
	line    int
	columnt int
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
			} else if r == bom && s.offset > 0 {
				s.error(s.offset, "illegal byte order mark")
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
		msg:     msg,
		offset:  offs,
		line:    s.line,
		columnt: offs - s.lineOffset,
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
	for isLetter(s.ch) || isDigit(s.ch) {
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
func isHex(ch rune) bool     { return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f' }

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

func (s *Scanner) scanString() string {

	start := s.offset

	for {

		ch := s.ch

		if ch == '\n' || ch < 0 {
			s.error(offs, "string literal not terminated")
			break
		}

		s.next()

		if ch == '"' {
			break
		}

		if ch == '\\' {
			s.scanEscape('"')
		}

	}

	return string(s.src[start : s.offset-1])
}

func (s *Scanner) scanNumber() (string, TokenType) {

	tt := INT
	lastDec := -1
	start := s.offset

	for {

		ch := s.ch
		s.next()

		if isDecimal(ch) {
			lastDec = s.offset - 1
		}

		if ch == '.' {
			tt = FLOAT
		}

		if isLetter(ch) {
			continue
		}

		break

	}

	if lastDec > 0 {

		suffix := string(s.src[lastDec:s.offset])

		if _, ok := timeSuffix[suffix]; ok {
			tt = TIME
		} else {
			s.error(lastDec+1, "invalid suffix")
		}

	}

	return string(s.src[start:s.offset]), tt
}

// Scan scans the next token and returns the token position, the token,
// and its literal string if applicable. The source end is indicated by
// token.EOF.
//
// If the returned token is a literal (token.IDENT, token.INT, token.FLOAT,
// token.STRING) or token.COMMENT, the literal string has the
// corresponding value.
//
// If the returned token is a keyword, the literal string is the keyword.
//
// In all other cases, Scan returns an empty literal string.
//
// This is a non-recoverable lexer and will panic at the first error found,
// error information will be provided using a ParseError type
//
func (s *Scanner) Scan() (tok Token) {

scanAgain:

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
			lit := s.scanString()
			tok = s.newToken(STRING, lit)
		case ':':
			tok = s.switch2(token.COLON, token.DEFINE)
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
				// comment
				if s.insertSemi && s.findLineEnd() {
					// reset position to the beginning of the comment
					s.ch = '/'
					s.offset = s.file.Offset(pos)
					s.rdOffset = s.offset + 1
					s.insertSemi = false // newline consumed
					return pos, token.SEMICOLON, "\n"
				}
				comment := s.scanComment()
				if s.mode&ScanComments == 0 {
					// skip comment
					s.insertSemi = false // newline consumed
					goto scanAgain
				}
				tok = token.COMMENT
				lit = comment
			} else {
				tok = s.switch2(token.QUO, token.QUO_ASSIGN)
			}
		case '%':
			tok = s.switch2(token.REM, token.REM_ASSIGN)
		case '^':
			tok = s.switch2(token.XOR, token.XOR_ASSIGN)
		case '<':
			if s.ch == '-' {
				s.next()
				tok = token.ARROW
			} else {
				tok = s.switch4(token.LSS, token.LEQ, '<', token.SHL, token.SHL_ASSIGN)
			}
		case '>':
			tok = s.switch4(token.GTR, token.GEQ, '>', token.SHR, token.SHR_ASSIGN)
		case '=':
			tok = s.switch2(token.ASSIGN, token.EQL)
		case '!':
			tok = s.switch2(token.NOT, token.NEQ)
		case '&':
			if s.ch == '^' {
				s.next()
				tok = s.switch2(token.AND_NOT, token.AND_NOT_ASSIGN)
			} else {
				tok = s.switch3(token.AND, token.AND_ASSIGN, '&', token.LAND)
			}
		case '|':
			tok = s.switch3(token.OR, token.OR_ASSIGN, '|', token.LOR)
		default:
			// next reports unexpected BOMs - don't repeat
			if ch != bom {
				s.errorf(s.file.Offset(pos), "illegal character %#U", ch)
			}
			insertSemi = s.insertSemi // preserve insertSemi info
			tok = token.ILLEGAL
			lit = string(ch)
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

	return string(s.src[start:s.offset])

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

	return s.newToken(t)

}

func (s *Scanner) newToken(t TokenType, lit string) Token {
	return Token{
		Typte:   t,
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

func (s *Scanner) literal() string {
	// since utf8 is not allowed in literals we can use offset instead of rdOffset.
	return string(s.src[s.start:s.offset])
}
