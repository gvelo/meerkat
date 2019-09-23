package mql_parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"meerkat/internal/query/logical"
)

//TODO:(sebad) revisar multithreading safe
func Parse(q string) *logical.Projection {

	fmt.Println("Query ", q)

	is := antlr.NewInputStream(q)

	// Create the Lexer
	lexer := NewMqlLexer(is)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := NewMqlParser(stream)

	l := newListener(lexer)

	// Finally parse the expression
	antlr.ParseTreeWalkerDefault.Walk(l, p.Start())

	return l.GetTree()

}
