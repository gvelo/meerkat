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

package io

import (
	"encoding/binary"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"meerkat/internal/util/testutil"
	"testing"
	"time"
)

const (
	numOfValues = 1024 * 200
	bufCap      = binary.MaxVarintLen64 * numOfValues
)

func TestBuffer(t *testing.T) {

	values := make([]interface{}, numOfValues)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	maxBufSize := 0

	//TODO(gvelo): to cover all varint cases generate values
	// uniformly distributed across the word size instead of
	// distributed across the range of integer values.

	for i := 0; i < numOfValues; i++ {
		n := rnd.Intn(4)
		switch n {
		case 0:
			values[i] = rnd.Int()
			maxBufSize += binary.MaxVarintLen64
		case 1:
			values[i] = -rnd.Int()
			maxBufSize += binary.MaxVarintLen64
		case 2:
			values[i] = uint(rnd.Int())
			maxBufSize += binary.MaxVarintLen64
		case 3:
			b := testutil.RandomBytes(255)
			values[i] = b
			maxBufSize += binary.MaxVarintLen64 + len(b)
		}
	}

	buf := Buffer{buf: make([]byte, maxBufSize)}

	for _, value := range values {
		switch v := value.(type) {
		case int:
			buf.WriteVarInt(v)
		case uint:
			buf.WriteUVarInt64(uint64(v))
		case []byte:
			buf.WriteBytes(v)
		default:
			panic("unknown type")
		}
	}

	buf.Reset()

	for _, value := range values {
		switch v := value.(type) {
		case int:
			r := buf.ReadVarInt()
			assert.Equal(t, v, r)
		case uint:
			r := buf.ReadUVarInt()
			assert.Equal(t, v, r)
		case []byte:
			r := buf.ReadBytes()
			assert.Equal(t, v, r)
		default:
			panic("unknown type")
		}
	}

}

func TestDecoderBuffer_Uvarint(t *testing.T) {

	values := RandomUVarInt(numOfValues)

	e := NewEncoderBuffer(bufCap)

	for _, v := range values {
		e.WriteUvarint(v)
	}

	d := NewDecoderBuffer()

	d.SetBytes(e.Bytes())

	for i := 0; i < numOfValues; i++ {
		n := d.ReadUvarint()

		assert.Equal(t, values[i], n, "wrong int values")

	}

}

func RandomUVarInt(size int) []int {

	var r []int

	for j := 0; j < size; j++ {
		n := uint64(0)
		s := rand.Intn(8) + 1
		for i := 0; i < s; i++ {
			n = n | uint64(byte(rand.Int()))<<uint(8*i)
		}
		r = append(r, int(n))
	}

	return r

}
