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
	"github.com/stretchr/testify/assert"
	"testing"
)

type Test struct {
	name    string
	input   string
	tokens  []Token
	isError bool
}

// TODO(gvelo): add more test cases.

var tests = []Test{
	{
		name:  "IDENT",
		input: "foo(10)",
		tokens: []Token{
			{
				Type:    IDENT,
				Literal: "foo",
			},
			{
				Type:    LPAREN,
				Literal: "(",
			},
			{
				Type:    INT,
				Literal: "10",
			},
			{
				Type:    RPAREN,
				Literal: ")",
			},
		},
		isError: false,
	},
	{
		name:  "DATETIME",
		input: "foo(datetime(2007-11-01))",
		tokens: []Token{
			{
				Type:    IDENT,
				Literal: "foo",
			},
			{
				Type:    LPAREN,
				Literal: "(",
			},
			{
				Type:    DATETIME,
				Literal: "2007-11-01",
			},
			{
				Type:    RPAREN,
				Literal: ")",
			},
			{
				Type:    RPAREN,
				Literal: ")",
			},
		},
		isError: false,
	},
	{
		name:    "string not terminated",
		input:   "foo(\"not terminated)",
		isError: true,
	},
}

func TestScanner(t *testing.T) {

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			tokens, err := scan(test.input)

			if test.isError {
				t.Log(err)
				assert.Error(t, err, "ParseError expected")
			}

			assertTokens(t, tokens, test.tokens)

		})

	}

}

func assertTokens(t *testing.T, actual []Token, expected []Token) {

	for i, token := range actual {

		e := expected[i]

		assert.Equal(t, e.Type, token.Type, "token Type doesn't match")
		assert.Equal(t, e.Literal, token.Literal, "token literal doesn't match")

	}

}

func scan(s string) ([]Token, error) {

	scanner := NewScanner(s)
	var tokens []Token

	for {
		t, err := safeScan(scanner)

		if err != nil {
			return nil, err
		}

		if t.Type == EOF {
			break
		}

		tokens = append(tokens, t)

	}

	return tokens, nil

}

func safeScan(s *Scanner) (t Token, err error) {

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

	t = s.Scan()

	return

}
