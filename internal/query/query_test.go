package query

import (
	"fmt"
	"meerkat/internal/query/logical"
	"meerkat/internal/query/parser"
	"testing"
)

func TestLogical(t *testing.T) {

	ast, err := parser.Parse("T|where a>10|limit 20")

	if err != nil {
		t.Error(err)
	}

	roots := logical.ToLogical(ast)

	fmt.Println(roots)

}
