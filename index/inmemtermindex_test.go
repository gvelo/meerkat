package index

import (
	"fmt"
	"testing"
)

func Test_inMemTermIndex_addTerm(t *testing.T) {

	fmt.Println("testing")

	idx := newInMemTermIndex()

	for i := 0; i < 65; i++ {
		idx.addTerm(fmt.Sprintf("%v%v", "A", i), uint32(i))
	}

	idx.printTrie()

}
