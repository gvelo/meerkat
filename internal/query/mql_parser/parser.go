package mql_parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"meerkat/internal/query/logical"
	"meerkat/internal/schema"
)

type ParseError struct {
	Err string
}

func (pe *ParseError) Error() string {
	return pe.Err
}

//TODO:(sebad) revisar multithreading safe, mover a parser
func Parse(s schema.Schema, q string) ([]logical.Node, error) {

	fmt.Println("Query ", q)

	is := antlr.NewInputStream(q)

	// Create the Lexer
	lexer := NewMqlLexer(is)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := NewMqlParser(stream)

	b := NewRelBuilder(s)

	l := newListener(b, lexer)

	// Finally parse the expression
	antlr.ParseTreeWalkerDefault.Walk(l, p.Start())

	return l.GetTree()

}
