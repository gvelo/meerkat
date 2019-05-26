package writers

import (
	"eventdb/readers"
	"eventdb/segment/inmem"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadWritePosting(t *testing.T) {

	assert := assert.New(t)

	postingStore := inmem.NewPostingStore()

	for i := 0; i < 1000; i++ {
		postingStore.NewPostingList(uint32(i))
	}

	file := "/tmp/posting.bin"

	err := WritePosting(file, postingStore)

	if !assert.NoErrorf(err, "an error occurred while writing the posting list: %v", err) {
		return
	}

	ps, err := readers.ReadPostingStore(file)

	if !assert.NoErrorf(err, "an error occurred while reading the posting list: %v", err) {
		return
	}

	for i, p := range postingStore.Store {

		b, err := ps.Read(p.Offset)

		if !assert.NoErrorf(err, "an error occurred while writing the posting list at offset %v: %v", p.Offset, err) {
			return
		}

		assert.Equal(uint64(1), b.GetCardinality(), "Wrong bitmap cardinality")
		assert.True(b.ContainsInt(i), "Bitmap doesn't contain expected value")

	}

}
