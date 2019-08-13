package rel

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"meerkat/internal/query/mql_parser"
)

func ProcessQuery(q string) *ParsedTree {

	fmt.Println("Query ", q)
	is := antlr.NewInputStream(q)

	// Create the Lexer
	lexer := mql_parser.NewMqlLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := mql_parser.NewMqlParser(stream)

	l := NewListener()
	// Finally parse the expression
	antlr.ParseTreeWalkerDefault.Walk(l, p.Start())

	return l.builder.Build()

}
