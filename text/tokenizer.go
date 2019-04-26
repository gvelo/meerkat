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
	tokenize(text string) []string
}

// NaiveTokenizer tokenize fields spliting their contents around
// each instance of one or more consecutive white space
// characters, as defined by unicode.IsSpace, .
type NaiveTokenizer struct{}

func (tokenizer *NaiveTokenizer) tokenize(text string) []string {
	return strings.Fields(text)
}

// Tokenizer parse segment a unicode byte string according to
// http://unicode.org/reports/tr29/#WB16.
// TODO: provide reuse to avoid allocation on each event.
type UnicodeTokenizer struct{}

func (tokenizer *UnicodeTokenizer) tokenize(text string) []string {

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
