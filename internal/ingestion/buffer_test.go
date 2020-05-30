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

package ingestion

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"meerkat/internal/schema"
	"meerkat/internal/util/testutil"
	"testing"
)

const (
	rows = 1024
)

func TestByteSliceBuffer(t *testing.T) {
	t.Run("test dense", func(t *testing.T) {
		testByteSliceBuffer(t, true)
	})
	t.Run("test sparse", func(t *testing.T) {
		testByteSliceBuffer(t, false)
	})
}

func TestFloatBuffer(t *testing.T) {
	t.Run("test dense", func(t *testing.T) {
		testFloat64Buffer(t, true)
	})
	t.Run("test sparse", func(t *testing.T) {
		testFloat64Buffer(t, false)
	})
}

func testByteSliceBuffer(t *testing.T, isDense bool) {

	bytes, size, validCount := generateRandomSlice([]byte{}, rows, isDense)

	bsb := NewByteSliceSparseBuffer(0, 0)

	bsb.Reserve(validCount, size)

	for i := uint32(0); i < rows; i++ {

		switch v := bytes[i].(type) {
		case nil:
			continue
		case []byte:
			bsb.Append(i, v)
		}

	}

	bsDense := bsb.ToDenseBuffer(rows)
	validateBSBuff(bsDense, bytes, t)

	bsb.Reserve(bsb.len*4, bsb.size*4)
	bsDense = bsb.ToDenseBuffer(rows)
	validateBSBuff(bsDense, bytes, t)

}

func validateBSBuff(bsDense *ByteSliceDenseBuffer, bytes []interface{}, t *testing.T) {

	valids := bsDense.Valids()

	for i := uint32(0); i < rows; i++ {

		switch v := bytes[i].(type) {
		case nil:
			assert.False(t, valids[i], "valid value should be false")
		case []byte:
			assert.True(t, valids[i], "valid value should be true")
			assert.Equal(t, v, bsDense.Value(i))
		}

	}

}

func testFloat64Buffer(t *testing.T, isDense bool) {

	values, _, validCount := generateRandomSlice(float64(0), rows, isDense)

	buffer := NewFloat64SparseBuffer(0)

	buffer.Reserve(validCount)

	for i := uint32(0); i < rows; i++ {

		switch v := values[i].(type) {
		case nil:
			continue
		case float64:
			buffer.Append(i, v)
		}

	}

	denseBuffer := buffer.ToDenseBuffer(rows)
	validateFloatBuff(denseBuffer, values, t)

	buffer.Reserve(buffer.len * 4)
	denseBuffer = buffer.ToDenseBuffer(rows)
	validateFloatBuff(denseBuffer, values, t)

}

func validateFloatBuff(denseBuffer *Float64DenseBuffer, values []interface{}, t *testing.T) {
	valids := denseBuffer.Valids()
	actualVal := denseBuffer.Values()

	for i := uint32(0); i < rows; i++ {

		switch expectedVal := values[i].(type) {
		case nil:
			assert.False(t, valids[i], "valid value should be false")
		case []byte:
			assert.True(t, valids[i], "valid value should be true")
			assert.Equal(t, expectedVal, actualVal)
		}

	}
}

func TestTSBuffer(t *testing.T) {

	values, _, validCount := generateRandomSlice(int64(0), rows, true)

	buffer := NewTSBuffer(0)

	buffer.Reserve(validCount)

	for i := uint32(0); i < rows; i++ {
		buffer.Append(values[i].(int64))
	}

	validateTSBuff(buffer, values, t)

	buffer.Reserve(buffer.len * 4)
	validateTSBuff(buffer, values, t)

}

func validateTSBuff(buffer *TSBuffer, values []interface{}, t *testing.T) {

	actuals := buffer.Values()

	for i := uint32(0); i < rows; i++ {
		expectedVal := values[i].(int64)
		assert.Equal(t, expectedVal, actuals[i])
	}

}

func generateRandomSlice(t interface{}, l int, dense bool) ([]interface{}, int, int) {

	var size int
	var valid int
	var values []interface{}

	for i := 0; i < l; i++ {

		r := rand.Intn(3)

		if r > 0 || dense {

			var value interface{}

			switch t.(type) {
			case []byte:
				b := testutil.RandomBytes(10)
				size = size + len(b)
				value = b
			case float64:
				value = rand.Float64()
			case int64:
				value = int64(rand.Int())
			default:
				panic("unknown type")

			}

			values = append(values, value)

			valid++

		} else {
			values = append(values, nil)
		}

	}

	return values, size, valid

}

func TestTableBuffer(t *testing.T) {

	const numBatch = 5
	const batchLen = 1024

	var tsColumnData []interface{}
	var testColumnData []interface{}
	var testColumnDataSize int

	tb := NewTableBuffer("testtable", 0)

	for i := 0; i < numBatch; i++ {

		tsDataBatch, _, _ := generateRandomSlice(int64(0), batchLen, true)
		testColDataBatch, testColDataBatchSize, _ := generateRandomSlice([]byte{}, batchLen, false)

		tsColumnData = append(tsColumnData, tsDataBatch...)
		testColumnData = append(testColumnData, testColDataBatch...)

		testColumnDataSize += testColDataBatchSize

		rsw := NewRowSetWriter(0)

		tsCol := &Column{
			Idx:     0,
			Name:    "_ts",
			ColSize: 0,
			Len:     0,
			Type:    schema.ColumnType_TIMESTAMP,
		}

		testCol := &Column{
			Idx:     1,
			Name:    "testCol",
			ColSize: 0,
			Len:     0,
			Type:    schema.ColumnType_STRING,
		}

		for i := 0; i < batchLen; i++ {

			rsw.WriteFixedInt64(0, tsDataBatch[i].(int64))
			tsCol.Len++

			if testColDataBatch[i] != nil {
				str := string(testColDataBatch[i].([]byte))
				rsw.WriteString(1, str)
				testCol.Len++
				testCol.ColSize += uint64(len(str))
			}

		}

		partition := &Partition{
			Id:      0,
			Columns: []*Column{tsCol, testCol},
			Data:    rsw.Buf.Data(),
		}

		tb.Append(partition)

	}

	tsBuf := tb.Columns()["_ts"].buff.(*TSBuffer)

	assert.Equal(t, len(tsColumnData), len(tsBuf.Values()))
	assert.Equal(t, numBatch*batchLen, len(tsColumnData))
	assert.Equal(t, numBatch*batchLen, int(tb.len))
	assert.Equal(t, tsBuf.len, numBatch*batchLen)

	for i, v := range tsBuf.Values() {
		assert.Equal(t, tsColumnData[i].(int64), v)
	}

	testColBuf := tb.Columns()["testCol"].buff.(*ByteSliceSparseBuffer).ToDenseBuffer(int(tb.len))

	assert.Equal(t, testColumnDataSize, testColBuf.size)
	assert.Equal(t, len(testColumnData), testColBuf.len)

	for i, v := range testColumnData {

		if v == nil {
			assert.False(t, testColBuf.Valids()[i])
			continue
		}

		assert.True(t, testColBuf.Valids()[i])
		assert.Equal(t, testColBuf.Value(uint32(i)), v)

	}

}
