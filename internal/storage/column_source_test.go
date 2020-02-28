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

package storage

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"meerkat/internal/buffer"
	"meerkat/internal/utils"
	"testing"
)

func TestNullableIntColumnSource(t *testing.T) {

	bufSize := 50000
	pageSize := 1024

	b, p, numOfNulls := createIntBuffer(bufSize, true)

	src := NewIntColumnSource(b, pageSize, p)

	testIntColumnSource(t, src, b, bufSize, pageSize, numOfNulls, p)

}

func TestIntColumnSource(t *testing.T) {

	bufSize := 50000
	pageSize := 1024

	b, p, numOfNulls := createIntBuffer(bufSize, false)

	src := NewIntColumnSource(b, pageSize, p)

	testIntColumnSource(t, src, b, bufSize, pageSize, numOfNulls, p)

}

func TestTsColumnSource(t *testing.T) {

	bufSize := 50000
	pageSize := 1024

	b, p, numOfNulls := createIntBuffer(bufSize, false)

	src := NewTsColumnSource(b, pageSize)

	testIntColumnSource(t, src, b, bufSize, pageSize, numOfNulls, p)

}

func TestByteColumnSource(t *testing.T) {

	bufSize := 10000
	pageSize := 1024

	b, p, _ := createByteSliceBuffer(bufSize, false)

	src := NewByteSliceColumnSource(b, pageSize, p)

	testByteSliceColumnSource(t, src, b, p)

}

func TestNullableByteColumnSource(t *testing.T) {

	bufSize := 10000
	pageSize := 1024

	b, p, _ := createByteSliceBuffer(bufSize, true)

	src := NewByteSliceColumnSource(b, pageSize, p)

	testByteSliceColumnSource(t, src, b, p)

}

func testIntColumnSource(t *testing.T,
	src IntColumSource,
	b *buffer.IntBuffer,
	bufSize int,
	pageSize int,
	numOfNulls int,
	p []int) {

	numOfValues := bufSize - numOfNulls
	numOfPages := numOfValues / pageSize
	remainder := numOfValues % pageSize

	currPage := 0

	for src.HasNext() {

		v := src.Next()

		for i := 0; i < v.Len(); i++ {
			assert.Equal(t, v.Values()[i], b.Values()[p[v.Rid()[i]]])
		}

		if currPage < numOfPages {
			assert.Equal(t, pageSize, v.Len(), "wrong page length")
		} else {
			assert.Equal(t, remainder, v.Len(), "wrong page length")
		}

		currPage++

	}

}

func testByteSliceColumnSource(t *testing.T,
	src ByteSliceColumSource,
	b *buffer.ByteSliceBuffer,
	p []int) {

	for src.HasNext() {

		v := src.Next()

		for i := 0; i < v.Len(); i++ {
			assert.Equal(t, v.Get(i), b.Get(p[v.Rid()[i]]))
		}

	}

}

func createIntBuffer(bufSize int, nullable bool) (*buffer.IntBuffer, []int, int) {

	b := buffer.NewIntBuffer(nullable, bufSize)
	p := make([]int, bufSize)
	nulls := 0

	for i := 0; i < bufSize; i++ {
		p[i] = i
		if nullable && rand.Intn(2) == 1 {
			b.AppendNull()
			nulls++
			continue
		}
		b.AppendInt(i)
	}

	return b, p, nulls

}

func createByteSliceBuffer(bufSize int, nullable bool) (*buffer.ByteSliceBuffer, []int, int) {

	b := buffer.NewByteSliceBuffer(nullable, 0)
	p := make([]int, bufSize)
	nulls := 0

	for i := 0; i < bufSize; i++ {
		p[i] = i
		if nullable && rand.Intn(2) == 1 {
			b.AppendNull()
			nulls++
			continue
		}
		b.AppendSlice([]byte(utils.RandomString(50)))
	}

	return b, p, nulls

}
