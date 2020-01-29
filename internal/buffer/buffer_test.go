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

package buffer

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewIntBuffer(t *testing.T) {
	assert := assert.New(t)
	b := NewIntBuffer(false, 0)
	assert.False(b.Nullable(), "wrong nullable return value")
	assert.Zero(b.Len(), "wrong len return value")
}

func TestNewNullableIntBuffer(t *testing.T) {
	assert := assert.New(t)
	b := NewIntBuffer(true, 0)
	assert.True(b.Nullable(), "wrong nullable return value")
	assert.Zero(b.Len(), "wrong len return value")
}

func TestNewIntBufferWithCap(t *testing.T) {
	assert := assert.New(t)
	b := NewIntBuffer(true, 1024)
	assert.True(b.Nullable(), "wrong nullable return value")
	assert.Zero(b.Len(), "wrong len return value")
}

func TestIntBuffer_AddInt(t *testing.T) {
	assert := assert.New(t)
	b := NewIntBuffer(true, 0)
	b.AppendInt(1)
	assert.True(b.Nullable(), "wrong nullable return value")
	assert.Equal(b.Len(), 1, "wrong len return value")
}

func TestIntBuffer_AddNull(t *testing.T) {
	assert := assert.New(t)
	b := NewIntBuffer(true, 0)
	b.AppendNull()
	assert.True(b.Nullable(), "wrong nullable return value")
	assert.Equal(b.Len(), 1, "wrong len return value")
	assert.Equal(len(b.Nulls()), 1, "wrong null len return value")
}

func TestIntBuffer_AddBuffer(t *testing.T) {

	assert := assert.New(t)

	srcBuff := NewIntBuffer(true, 0)
	dstBuff := NewIntBuffer(true, 0)

	for i := 0; i < 100; i++ {
		srcBuff.AppendInt(i)
		dstBuff.AppendInt(i)
	}

	dstBuff.AppendIntBuffer(srcBuff)

	assert.Equal(dstBuff.Len(), 200, "wrong len return value")
	assert.Equal(len(dstBuff.nulls), 200, "wrong nulls len return value")
	for _, value := range dstBuff.Nulls() {
		assert.True(value, "wrong null value")
	}

}

func TestIntBuffer_AddBufferNull(t *testing.T) {

	assert := assert.New(t)

	srcBuff := NewIntBuffer(true, 0)
	dstBuff := NewIntBuffer(true, 0)

	for i := 0; i < 100; i++ {
		dstBuff.AppendNull()
		srcBuff.AppendInt(i)
	}

	dstBuff.AppendIntBuffer(srcBuff)

	assert.Equal(dstBuff.Len(), 200, "wrong len return value")
	assert.Len(dstBuff.Nulls(), 200, "wrong nulls len return value")

	for i := 0; i < 100; i++ {
		assert.Equal(i, dstBuff.buf[i+100], "wrong array content")
		assert.True(dstBuff.Nulls()[i])
	}

}

func TestSliceBuffer_AppendString(t *testing.T) {

	assert := assert.New(t)

	ss := []string{"one", "two", "three"}

	buf := NewSliceBuffer(false, 0)

	bufSize := 0
	for _, s := range ss {
		buf.AppendString(s)
		bufSize = bufSize + len(s)
	}

	assert.Equal(len(ss), len(buf.offsets))

	var sr []string

	buf.Each(func(i int, bytes []byte) bool {
		sr = append(sr, string(bytes))
		return true
	})

	assert.Equal(ss, sr, "wrong slice value")

}

func TestSliceBuffer_AppendSliceBuffer(t *testing.T) {

	assert := assert.New(t)

	sb1 := []string{"one", "two", "three"}
	sb2 := []string{"four", "five"}

	buf1 := NewSliceBuffer(false, 0)
	buf2 := NewSliceBuffer(false, 0)

	for _, s := range sb1 {
		buf1.AppendString(s)
	}

	for _, s := range sb2 {
		buf2.AppendString(s)
	}

	buf1.AppendSliceBuffer(buf2)
	sbr := append(sb1, sb2...)

	var bufr []string

	buf1.Each(func(i int, bytes []byte) bool {
		bufr = append(bufr, string(bytes))
		return true
	})

	assert.Equal(sbr, bufr, "wrong slice value")

}

func TestUUIDBuffer_Append(t *testing.T) {

	assert := assert.New(t)

	buf := NewUUIDBuffer(false, 0)

	uid := make([]uuid.UUID, 256)

	for i := 0; i < len(uid); i++ {
		uid[i] = uuid.New()
		buf.Append(uid[i])
	}

	assert.Equal(256, buf.Len())

	buf.Each(func(i int, actual uuid.UUID) bool {
		expected := uid[i]
		r := bytes.Equal(expected[:], actual[:])
		assert.True(r, "buffers are not equals")
		return true
	})

}
