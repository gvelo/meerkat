// Copyright 2020 The Meerkat Authors
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

package inmem

import (
	"fmt"
	"math/rand"
	"meerkat/internal/util/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_btrie_emptycardinality(t *testing.T) {

	assert := assert.New(t)

	ps := NewPostingStore()

	bt := NewBtrie(ps)

	assert.Equal(0, bt.Cardinality, "empty index should have cero Cardinality")

}

func Test_btrie_cardinality(t *testing.T) {

	assert := assert.New(t)

	ps := NewPostingStore()

	bt := NewBtrie(ps)

	for i := 0; i < 10; i++ {
		bt.Add([]byte("test"), 0)
	}

	assert.Equal(1, bt.Cardinality, "same term should not change index Cardinality")

}

func Test_btrie_multiple_cardinality(t *testing.T) {

	assert := assert.New(t)

	ps := NewPostingStore()

	bt := NewBtrie(ps)

	for i := 0; i < 10; i++ {
		bt.Add([]byte(fmt.Sprintf("test%v", i)), 0)
	}

	assert.Equal(10, bt.Cardinality, "unique term should change index Cardinality")

}

func Test_btrie_bucketSize(t *testing.T) {

	assert := assert.New(t)

	ps := NewPostingStore()

	bt := &BTrie{
		MaxBucketSize: 10,
		PostingStore:  ps,
	}

	bt.Root = bt.newNode()

	for i := 0; i < 10; i++ {
		bt.Add([]byte(fmt.Sprintf("test%v", i)), 0)
	}

	assert.Equal(1, bt.NumOfNodes, "NumOfNodes should be 1")

	bt.Add([]byte("Test burst"), 0)

	assert.Equal(2, bt.NumOfNodes, "Node burst should result in 2 nodes")

}

func Test_btrie_singlelookup(t *testing.T) {

	assert := assert.New(t)

	ps := NewPostingStore()

	bt := NewBtrie(ps)

	assert.Nil(bt.Lookup("test"), "non empty index")

	bt.Add([]byte("test"), 0)

	bitmap := bt.Lookup("test")

	assert.NotNil(bitmap, "result should not be nil")

	assert.Equal(uint64(1), bitmap.GetCardinality(), "wrong bitmap Cardinality")

}

func Test_btrie_lookup_cardinality(t *testing.T) {

	assert := assert.New(t)

	ps := NewPostingStore()

	bt := NewBtrie(ps)

	assert.Nil(bt.Lookup("test"), "non empty index")

	for i := 0; i < 100; i++ {
		bt.Add([]byte("test"), uint32(i))
	}

	bitmap := bt.Lookup("test")

	assert.NotNil(bitmap, "result should not be nil")

	assert.Equal(uint64(100), bitmap.GetCardinality(), "wrong bitmap Cardinality")
	assert.Equal(1, bt.Cardinality, "wrong cardinallity")
}

func Test_btrie_lookup_cardinality2(t *testing.T) {

	assert := assert.New(t)

	ps := NewPostingStore()

	bt := NewBtrie(ps)

	assert.Nil(bt.Lookup("t"), "non empty index")

	for i := 0; i < 100; i++ {
		bt.Add([]byte("t"), uint32(i))
	}

	bitmap := bt.Lookup("t")

	assert.NotNil(bitmap, "result should not be nil")

	assert.Equal(uint64(100), bitmap.GetCardinality(), "wrong bitmap Cardinality")
	assert.Equal(1, bt.Cardinality, "wrong cardinallity")
}

func Test_random_test(t *testing.T) {

	const numOfTokens = 1024 * 100

	tokenMap := make(map[string]bool, numOfTokens)
	tokenVec := make([]string, numOfTokens)
	freq := make([]int, numOfTokens)

	for i := 0; i < numOfTokens; i++ {

		t := testutil.RandomString(30)

		if tokenMap[t] {
			i--
			continue
		}

		tokenMap[t] = true

	}

	i := 0
	for t, _ := range tokenMap {

		tokenVec[i] = t
		i++

	}

	ps := NewPostingStore()
	bt := NewBtrie(ps)

	for i := 0; i < numOfTokens*10; i++ {

		idx := rand.Intn(numOfTokens)
		bt.Add([]byte(tokenVec[idx]), uint32(i))
		freq[idx]++
	}

	for i := 0; i < numOfTokens; i++ {
		posting := bt.Lookup(tokenVec[i])
		if posting == nil {
			assert.Zero(t, freq[i], "wrong cardinality")
			continue
		}
		assert.Equalf(t, freq[i], int(posting.GetCardinality()), "wrong cardinality for token %v", tokenVec[i])
	}

}
