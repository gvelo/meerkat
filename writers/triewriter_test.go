package writers

import (
	"eventdb/segment/inmem"
	"fmt"
	"testing"
)

func TestBTrieWriter(t *testing.T) {

	trie := inmem.NewBtrie()

	for i := 0; i < 1000; i++ {
		s := fmt.Sprintf("Test number %v", i)
		trie.Add(s, uint32(i))
	}

	writer, err := NewTrieWriter("/tmp/trie.bin")

	if err != nil {
		t.Error(err)
		return
	}

	_,err = writer.Write(trie)

	if err != nil {
		t.Error(err)
	}

	writer.Close()

}
