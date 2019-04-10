package index

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_inMemTermIndex_emptycardinality(t *testing.T) {

	assert := assert.New(t)

	idx := newInMemTermIndex()

	assert.Equal(0, idx.cardinality, "empty index should have cero cardinality")

}

func Test_inMemTermIndex_cardinality(t *testing.T) {

	assert := assert.New(t)

	idx := newInMemTermIndex()
	for i := 0; i < 10; i++ {
		idx.addTerm("test", 0)
	}

	assert.Equal(1, idx.cardinality, "same term should not change index cardinality")

}

func Test_inMemTermIndex_multiple_cardinality(t *testing.T) {

	assert := assert.New(t)

	idx := newInMemTermIndex()
	for i := 0; i < 10; i++ {
		idx.addTerm(fmt.Sprintf("test%v", i), 0)
	}

	assert.Equal(10, idx.cardinality, "unique term should change index cardinality")

}

func Test_inMemTermIndex_bucketSize(t *testing.T) {

	assert := assert.New(t)

	idx := &inMemTermIndex{
		bucketSize: 10,
	}
	idx.root = idx.newNode()

	for i := 0; i < 10; i++ {
		idx.addTerm(fmt.Sprintf("test%v", i), 0)
	}

	assert.Equal(1, idx.size, "size should be 1")

	idx.addTerm("Test burst", 0)

	assert.Equal(2, idx.size, "node burst should result in 2 nodes")

}

func Test_inMemTermIndex_singlelookup(t *testing.T) {

	assert := assert.New(t)

	idx := newInMemTermIndex()

	assert.Nil(idx.lookup("test"), "non empty index")

	idx.addTerm("test", 0)

	bitmap := idx.lookup("test")

	assert.NotNil(bitmap, "result should not be nil")

	assert.Equal(uint64(1), bitmap.GetCardinality(), "wrong bitmap cardinality")

}

func Test_inMemTermIndex_lookup_cardinality(t *testing.T) {

	assert := assert.New(t)

	idx := newInMemTermIndex()

	assert.Nil(idx.lookup("test"), "non empty index")

	for i := 0; i < 100; i++ {
		idx.addTerm("test", uint32(i))
	}

	bitmap := idx.lookup("test")

	assert.NotNil(bitmap, "result should not be nil")

	assert.Equal(uint64(100), bitmap.GetCardinality(), "wrong bitmap cardinality")
	assert.Equal(1, idx.cardinality, "wrong cardinallity")
}
