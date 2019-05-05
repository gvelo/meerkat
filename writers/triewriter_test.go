package writers

import (
	"eventdb/readers"
	"eventdb/segment/inmem"
	"github.com/stretchr/testify/assert"

	//"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

var chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
	"abcdefghijklmnopqrstuvwxyzåäö" +
	"0123456789")


func TestBTrieReadWriter(t *testing.T) {

	assert := assert.New(t)

	rand.Seed(time.Now().UnixNano())

	terms := make(map[string]int)

	for i := 0; i < 1000; i++ {
		term := rndStr()
		count := rand.Intn(1000) + 1
		c := terms[term]
		terms[term] = count + c
	}

	ps := inmem.NewPostingStore()

	trie := inmem.NewBtrie(ps)

	var eventID uint32

	start := time.Now()

	for term, count := range terms {
		for c := 0; c < count; c++ {
			trie.Add(term, eventID)
			eventID++
		}
	}

	t.Logf("creating in-mem trie took %v\n", time.Since(start))

	// Validate the expected terms and event count
	// on the on memory trie.

	start = time.Now()

	for term, count := range terms {

		bitmap := trie.Lookup(term)

		if !assert.NotNil(bitmap, "term not found") {
			return
		}

		if !assert.Equal(uint64(count), bitmap.GetCardinality(), "wrong cardinality") {
			return
		}

	}

	t.Logf("in mem trie validation took %v", time.Since(start))

	// Write posting list to disk.

	start = time.Now()

	err := WritePosting("/tmp/posting-test.bin", ps.Store)

	if !assert.NoErrorf(err, "an error occurred while writing the posting list: %v", err) {
		return
	}

	t.Logf("posting list write took %v", time.Since(start))

	// write trie from disk.

	start = time.Now()

	writer, err := NewTrieWriter("/tmp/trie.bin")

	if !assert.NoErrorf(err, "an error occurred while creating trie writer: %v", err) {
		return
	}

	err = writer.Write(trie)

	if !assert.NoErrorf(err, "an error occurred while writing trie to disk: %v", err) {
		return
	}

	writer.Close()

	t.Logf("writing trie to disk took %v", time.Since(start))

	// Read trie from disk

	start = time.Now()

	odt, err := readers.ReadTrie("/tmp/trie.bin")

	if !assert.NoErrorf(err, "an error occurred while creating trie reader: %v", err) {
		return
	}

	pReader, err := readers.NewPostingReader("/tmp/posting-test.bin")

	if !assert.NoErrorf(err, "an error occurred while creating posting reader: %v", err) {
		return
	}

	// compare the values on disk with the expected values.

	for term, count := range terms {

		offset, err := odt.Lookup(term)

		if !assert.NoErrorf(err, "an error occurred while searching on trie: %v", err) {
			return
		}

		bitmap, err := pReader.Read(int64(offset))

		if !assert.NoErrorf(err, "an error occurred while fetching posting posting from disk: %v", err) {
			return
		}

		if !assert.Equal(uint64(count), bitmap.GetCardinality(), "wrong cardinality on disk for term %v", term) {
			return
		}

	}

	t.Logf("ondisk trie validation took %v ", time.Since(start))

}

func rndStr() string {
	length := rand.Intn(15) + 1
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
