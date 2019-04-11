package index

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_btrie_emptycardinality(t *testing.T) {

	assert := assert.New(t)

	bt := newBtrie()

	assert.Equal(0, bt.cardinality, "empty index should have cero cardinality")

}

func Test_btrie_cardinality(t *testing.T) {

	assert := assert.New(t)

	bt := newBtrie()
	for i := 0; i < 10; i++ {
		bt.add("test", 0)
	}

	assert.Equal(1, bt.cardinality, "same term should not change index cardinality")

}

func Test_btrie_multiple_cardinality(t *testing.T) {

	assert := assert.New(t)

	bt := newBtrie()
	for i := 0; i < 10; i++ {
		bt.add(fmt.Sprintf("test%v", i), 0)
	}

	assert.Equal(10, bt.cardinality, "unique term should change index cardinality")

}

func Test_btrie_bucketSize(t *testing.T) {

	assert := assert.New(t)

	bt := &btrie{
		bucketSize: 10,
	}
	bt.root = bt.newNode()

	for i := 0; i < 10; i++ {
		bt.add(fmt.Sprintf("test%v", i), 0)
	}

	assert.Equal(1, bt.size, "size should be 1")

	bt.add("Test burst", 0)

	assert.Equal(2, bt.size, "node burst should result in 2 nodes")

}

func Test_btrie_singlelookup(t *testing.T) {

	assert := assert.New(t)

	bt := newBtrie()

	assert.Nil(bt.lookup("test"), "non empty index")

	bt.add("test", 0)

	bitmap := bt.lookup("test")

	assert.NotNil(bitmap, "result should not be nil")

	assert.Equal(uint64(1), bitmap.GetCardinality(), "wrong bitmap cardinality")

}

func Test_btrie_lookup_cardinality(t *testing.T) {

	assert := assert.New(t)

	bt := newBtrie()

	assert.Nil(bt.lookup("test"), "non empty index")

	for i := 0; i < 100; i++ {
		bt.add("test", uint32(i))
	}

	bitmap := bt.lookup("test")

	assert.NotNil(bitmap, "result should not be nil")

	assert.Equal(uint64(100), bitmap.GetCardinality(), "wrong bitmap cardinality")
	assert.Equal(1, bt.cardinality, "wrong cardinallity")
}
