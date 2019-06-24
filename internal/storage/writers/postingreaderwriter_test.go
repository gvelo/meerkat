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
