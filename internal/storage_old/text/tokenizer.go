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

package text

import (
	"bufio"
	"strings"

	"github.com/blevesearch/segment"
)

// Tokenizer parse textual fields and provide the stream of tokens.
// TODO: I guess we gonna have to optimize this interface for low
// allocation.
type Tokenizer interface {
	Tokenize(text string) []string
}

// NaiveTokenizer tokenize fields spliting their contents around
// each instance of one or more consecutive white space
// characters, as defined by unicode.IsSpace, .
type NaiveTokenizer struct{}

func (tokenizer *NaiveTokenizer) Tokenize(text string) []string {
	return strings.Fields(text)
}

// Tokenizer parse segment a unicode byte string according to
// http://unicode.org/reports/tr29/#WB16.
// TODO: provide reuse to avoid allocation on each event.
type UnicodeTokenizer struct{}

func (tokenizer *UnicodeTokenizer) Tokenize(text string) []string {

	tokens := make([]string, 128)
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(segment.SplitWords)
	for scanner.Scan() {
		tokens = append(tokens, scanner.Text())
	}
	return tokens
}

func NewTokenizer() *UnicodeTokenizer {
	return &UnicodeTokenizer{}
}
