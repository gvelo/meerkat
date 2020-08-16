package logical

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"meerkat/internal/query/parser"
	"testing"
)

func TestParallelize(t *testing.T) {

	//ast, err := parser.Parse("T|where A>10|summarize avg=avg(C) by colA=bin(d)")
	ast, err := parser.Parse("T|where A>10|where B>5")

	fmt.Println("====== AST ===== ")
	spew.Dump(ast)
	fmt.Println()

	if err != nil {
		t.Fatal(err)
	}

	logical := ToLogical(ast)
	fmt.Println("====== LOGICAL ===== ")
	spew.Dump(logical)
	fmt.Println()

	fragments := Parallelize(logical)


	fmt.Println("====== PARALLEL ===== ")
	spew.Dump(fragments)

}
