package inmem

import (
	"fmt"
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
		bt.Add("test", 0)
	}

	assert.Equal(1, bt.Cardinality, "same term should not change index Cardinality")

}

func Test_btrie_multiple_cardinality(t *testing.T) {

	assert := assert.New(t)

	ps := NewPostingStore()

	bt := NewBtrie(ps)

	for i := 0; i < 10; i++ {
		bt.Add(fmt.Sprintf("test%v", i), 0)
	}

	assert.Equal(10, bt.Cardinality, "unique term should change index Cardinality")

}

func Test_btrie_bucketSize(t *testing.T) {

	assert := assert.New(t)

	ps := NewPostingStore()

	bt := &BTrie{
		BucketSize:   10,
		PostingStore: ps,
	}

	Root = bt.newNode()

	for i := 0; i < 10; i++ {
		bt.Add(fmt.Sprintf("test%v", i), 0)
	}

	assert.Equal(1, Size, "Size should be 1")

	bt.Add("Test burst", 0)

	assert.Equal(2, bt.Size, "Node burst should result in 2 nodes")

}

func Test_btrie_singlelookup(t *testing.T) {

	assert := assert.New(t)

	ps := NewPostingStore()

	bt := NewBtrie(ps)

	assert.Nil(Lookup("test"), "non empty index")

	bt.Add("test", 0)

	bitmap := bt.Lookup("test")

	assert.NotNil(bitmap, "result should not be nil")

	assert.Equal(uint64(1), bitmap.GetCardinality(), "wrong bitmap Cardinality")

}

func Test_btrie_lookup_cardinality(t *testing.T) {

	assert := assert.New(t)

	ps := NewPostingStore()

	bt := NewBtrie(ps)

	assert.Nil(Lookup("test"), "non empty index")

	for i := 0; i < 100; i++ {
		bt.Add("test", uint32(i))
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
		bt.Add("t", uint32(i))
	}

	bt.DumpTrie()

	bitmap := bt.Lookup("t")

	assert.NotNil(bitmap, "result should not be nil")

	assert.Equal(uint64(100), bitmap.GetCardinality(), "wrong bitmap Cardinality")
	assert.Equal(1, bt.Cardinality, "wrong cardinallity")
}
