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

package index

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"meerkat/internal/storage/io"
	"testing"
)

const (
	path       = "/tmp/blockindex" //TODO(gvelo) use a tmp file.
	testLength = 100000
)

func TestBlockIndex(t *testing.T) {

	var buff []byte

	bw, err := io.NewBinaryWriter(path)

	if err != nil {
		t.Error(err)
		return
	}

	pw := NewBlockIndexWriter(bw)

	maxBlockSize := 16 * 1024

	block := make([]byte, maxBlockSize)

	for i := 0; i < testLength; i++ {
		l := rand.Intn(maxBlockSize-8) + 8
		binary.LittleEndian.PutUint64(block, uint64(i))
		pw.IndexBlock(block[:l], uint32(i))
		buff = append(buff, block[:l]...)
	}

	e, err := pw.Flush()

	if err != nil {
		t.Error(err)
		return
	}

	err = bw.Close()

	if err != nil {
		t.Error(err)
	}

	// read

	f, err := io.MMap(path)

	if err != nil {
		t.Error(t)
		return
	}

	br := f.NewBinaryReader()

	br.Offset = e

	pr := NewBlockIndexReader(br)

	err = pr.read()

	if err != nil {
		t.Error(err)
		return
	}

	for i := 0; i < testLength; i++ {
		rid, offset := pr.Lookup(uint32(i))
		n := binary.LittleEndian.Uint64(buff[offset:])
		assert.Equal(t, uint64(rid), n, "block rid doesn't match")
	}

}
