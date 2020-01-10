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

package writers

import (
	"github.com/stretchr/testify/assert"
	"meerkat/internal/storage/readers"
	"meerkat/internal/storage/segment/inmem"

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

	err := WritePosting("/tmp/posting-test.bin", ps)

	if !assert.NoErrorf(err, "an error occurred while writing the posting list: %v", err) {
		return
	}

	t.Logf("posting list write took %v", time.Since(start))

	// write trie from disk.

	start = time.Now()

	err = WriteTrie("/tmp/trie.bin", trie)

	if !assert.NoErrorf(err, "an error occurred while writing trie to disk: %v", err) {
		return
	}

	t.Logf("writing trie to disk took %v", time.Since(start))

	// Read trie from disk

	start = time.Now()

	odt, err := readers.ReadTrie("/tmp/trie.bin")

	if !assert.NoErrorf(err, "an error occurred while creating trie reader: %v", err) {
		return
	}

	postingStore, err := readers.ReadPostingStore("/tmp/posting-test.bin")

	if !assert.NoErrorf(err, "an error occurred while creating posting reader: %v", err) {
		return
	}

	// compare the values on disk with the expected values.

	for term, count := range terms {

		offset, err := odt.Lookup(term)

		if !assert.NoErrorf(err, "an error occurred while searching on trie: %v", err) {
			return
		}

		bitmap, err := postingStore.Read(offset)

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
