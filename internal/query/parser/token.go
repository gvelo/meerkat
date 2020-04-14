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

import "fmt"

//go:generate stringer -type=TokenType

type TokenType int

const (

	// Special tokens
	EOF TokenType = iota
	COMMENT

	literal_beg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT    // main
	INT      // 12345
	FLOAT    // 123.45
	STRING   // "abc"
	TIME     // time(15 seconds)
	DATETIME // datetime(2007-11-01)
	BOOL     // true false
	literal_end

	operator_beg
	// Operators and delimiters
	ADD               // +
	SUB               // -
	MUL               // *
	QUO               // /
	REM               // %
	ASSIGN            // =
	EQL               // ==
	EQL_CI            // =~
	NEQ               // !=
	NEQ_CI            // !~
	LSS               // <
	GTR               // >
	LEQ               // <=
	GEQ               // >=
	AND               // and
	OR                // or
	IN                // in
	NOT_IN            // !in
	IN_CI             // in~
	NOT_IN_CI         // !in~
	HAS               // has
	NOT_HAS           // !has
	HAS_CS            // has_cs
	NOT_HAS_CS        // !has_cs
	HASPREFIX         // hasprefix
	NOT_HASPREFIX     // !hasprefix
	HASPREFIX_CS      // hasprefix_cs
	NOT_HASPREFIX_CS  // !hasprefix_cs
	HASSUFFIX         // hassuffix
	NOT_HASSUFFIX     // !hassuffix
	HASSUFFIX_CS      // hassuffix_cs
	NOT_HASSUFFIX_CS  // !hassuffix_cs
	CONTAINS          // contains
	NOT_CONTAINS      // !contains
	CONTAINS_CS       // contains_cs
	NOT_CONTAINS_CS   // !contains_cs
	STARTSWITH        // startswith
	NOT_STARTSWITH    // !startswith
	STARTSWITH_CS     // startswith_cs
	NOT_STARTSWITH_CS // !startswith_cs
	ENDSWITH          // endswith
	NOT_ENDSWITH      // !endswith
	ENDSWITH_CS       // endswith_cs
	NOT_ENDSWITH_CS   // !endswith_cs
	MATCHES           // matches
	HAS_ANY           // has_any
	BETWEEN           // between
	NOT_BETWEEN       // !between
	RANGE             // ..
	operator_end

	LPAREN // (
	LBRACK // [
	LBRACE // {
	COMMA  // ,
	PERIOD // .

	RPAREN    // )
	RBRACK    // ]
	RBRACE    // }
	SEMICOLON // ;
	COLON     // :

	PIPE // |

)

var tokens = map[string]TokenType{
	"and":            AND,
	"or":             OR,
	"in":             IN,
	"!in":            NOT_IN,
	"in~":            IN_CI,
	"!in~":           NOT_IN_CI,
	"has":            HAS,
	"!has":           NOT_HAS,
	"has_cs":         HAS_CS,
	"!has_cs":        NOT_HAS_CS,
	"hasprefix":      HASPREFIX,
	"!hasprefix":     NOT_HASPREFIX,
	"hasprefix_cs":   HASPREFIX_CS,
	"!hasprefix_cs":  NOT_HASPREFIX_CS,
	"hassuffix":      HASSUFFIX,
	"!hassuffix":     NOT_HASSUFFIX,
	"hassuffix_cs":   HASSUFFIX_CS,
	"!hassuffix_cs":  NOT_HASSUFFIX_CS,
	"contains":       CONTAINS,
	"!contains":      NOT_CONTAINS,
	"contains_cs":    CONTAINS_CS,
	"!contains_cs":   NOT_CONTAINS_CS,
	"startswith":     STARTSWITH,
	"!startswith":    NOT_STARTSWITH,
	"startswith_cs":  STARTSWITH_CS,
	"!startswith_cs": NOT_STARTSWITH_CS,
	"endswith":       ENDSWITH,
	"!endswith":      NOT_ENDSWITH,
	"endswith_cs":    ENDSWITH_CS,
	"!endswith_cs":   NOT_ENDSWITH_CS,
	"matches":        MATCHES,
	"has_any":        HAS_ANY,
	"between":        BETWEEN,
	"!between":       NOT_BETWEEN,
	"time":           TIME,
	"datetime":       DATETIME,
}

var timeSuffix = map[string]bool{
	"ns":           true,
	"nanosecond":   true,
	"nanoseconds":  true,
	"ms":           true,
	"millisecond":  true,
	"milliseconds": true,
	"s":            true,
	"second":       true,
	"seconds":      true,
	"m":            true,
	"minute":       true,
	"minutes":      true,
	"h":            true,
	"hour":         true,
	"hours":        true,
	"d":            true,
	"day":          true,
	"days":         true,
	"month":        true,
	"months":       true,
	"year":         true,
	"years":        true,
}

// A set of constants for precedence-based expression parsing.
// Non-operators have lowest precedence, followed by operators
// starting with precedence 1 up to unary operators. The highest
// precedence serves as "catch-all" precedence for selector,
// indexing, and other operator and delimiter tokens.
//
const (
	LowestPrec  = 0 // non-operators
	UnaryPrec   = 6
	HighestPrec = 7
)

func getTokenType(lit string) TokenType {

	t, ok := tokens[lit]

	if ok {
		return t
	}

	return IDENT

}

type Token struct {
	Type    TokenType
	Literal string
	Offset  int
	Column  int
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("token %s (%v) offset:%v line:%v column:%v", t.Type.String(), t.Literal, t.Offset, t.Line, t.Column)
}

// Precedence returns the operator precedence of the binary
// operator op. If op is not a binary operator, the result
// is LowestPrecedence.
//
func (t Token) Precedence() int {
	switch t.Type {
	case OR:
		return 1
	case AND:
		return 2
	case EQL, NEQ, LSS, LEQ, GTR, GEQ:
		return 3
	case ADD, SUB:
		return 4
	case MUL, QUO, REM:
		return 5
	}
	return LowestPrec
}

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (t Token) IsLiteral() bool {
	return literal_beg < t.Type && t.Type < literal_end
}

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (t Token) IsOperator() bool {
	return operator_beg < t.Type && t.Type < operator_end
}
