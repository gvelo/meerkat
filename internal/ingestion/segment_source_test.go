package ingestion

import (
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"meerkat/internal/storage"
	"meerkat/internal/util/testutil"
	"strconv"
	"testing"
	"time"
)

const (
	testLen   = 1024*10 + 10
	maxStrLen = 10
)

var int64StrLen = len(strconv.Itoa(math.MaxInt64)) + 1

type TestData struct {
	values []interface{}
	rids   []uint32
}

func TestSegmentSource(t *testing.T) {

	testTableData := map[string]TestData{
		storage.TSColumnName: generateData(storage.ColumnType_TIMESTAMP, testLen, false),
		"bs-col":             generateData(storage.ColumnType_STRING, testLen, false),
		"bs-nul-col":         generateData(storage.ColumnType_STRING, testLen, true),
		"int64-dyn-col":      generateData(storage.ColumnType_INT64, testLen, false),
		"int64-null-dyn-col": generateData(storage.ColumnType_INT64, testLen, true),
	}

	table := createTestTable(testTableData)

	segmentSource := NewSegmentSource(table)

	for _, column := range segmentSource.Info().Columns {

		blockSize := 0

		if column.ColumnType == storage.ColumnType_STRING {
			blockSize = 16 * 1024
		} else {
			blockSize = 1024
		}

		columnSource := segmentSource.ColumnSource(column.Name, blockSize)

		actual := readSrc(columnSource)
		expected := testTableData[column.Name]

		t.Run(column.Name, func(t *testing.T) {
			assert.Equal(t, expected, actual)
		})

	}

}

func readSrc(src storage.ColumnSource) TestData {

	testData := TestData{}

	for src.HasNext() {

		switch s := src.(type) {

		case storage.ByteSliceColumnSource:
			values := s.Next()
			testData.rids = append(testData.rids, values.Rid()...)
			for i := 0; i < values.Len(); i++ {
				v := values.Get(i)
				// we need to copy the values given that colvalues
				// are only valid until the next call to next()
				c := make([]byte, len(v))
				copy(c, v)
				testData.values = append(testData.values, c)
			}

		case storage.Int64ColumnSource:
			values := s.Next()
			testData.rids = append(testData.rids, values.Rid()...)
			for i := 0; i < values.Len(); i++ {
				testData.values = append(testData.values, values.Values()[i])
			}

		default:
			panic("unknown source type")
		}

	}

	return testData

}

func createTestTable(testData map[string]TestData) *TableBuffer {

	table := &TableBuffer{
		partitionID: 0,
		tableName:   "test-table",
		len:         testLen,
		columns: map[string]*tableColumn{
			storage.TSColumnName: {
				colType: storage.ColumnType_TIMESTAMP,
				buf:     createTSBuff(testData[storage.TSColumnName]),
			},
			"bs-col": {
				colType: storage.ColumnType_STRING,
				buf:     createByteSliceBuf(testData["bs-col"]),
			},
			"bs-nul-col": {
				colType: storage.ColumnType_STRING,
				buf:     createByteSliceBuf(testData["bs-nul-col"]),
			},
			"int64-dyn-col": {
				colType: storage.ColumnType_INT64,
				buf:     createInt64Buf(testData["int64-dyn-col"]),
			},
			"int64-null-dyn-col": {
				colType: storage.ColumnType_INT64,
				buf:     createInt64Buf(testData["int64-null-dyn-col"]),
			},
		},
	}

	return table

}

func createTSBuff(testData TestData) *TSBuffer {
	ts := NewTSBuffer(testLen)
	for _, value := range testData.values {
		ts.Append(value.(int64))
	}
	return ts
}

func createByteSliceBuf(testData TestData) *ByteSliceSparseBuffer {
	buf := NewByteSliceSparseBuffer(testLen, testLen*maxStrLen)
	for i, rid := range testData.rids {
		buf.Append(rid, testData.values[i].([]byte))
	}
	return buf
}

func createInt64Buf(testData TestData) *ByteSliceSparseBuffer {
	buf := NewByteSliceSparseBuffer(testLen, testLen*int64StrLen)
	for i, rid := range testData.rids {
		str := strconv.FormatInt(testData.values[i].(int64), 10)
		buf.Append(rid, []byte(str))
	}
	return buf
}

func generateData(colType storage.ColumnType, l uint32, nullable bool) TestData {

	tc := TestData{}
	n := time.Now().UnixNano()

	for rid := uint32(0); rid < l; rid++ {

		if nullable && rand.Intn(3) == 2 {
			continue
		}

		var value interface{}

		switch colType {
		case storage.ColumnType_STRING:
			value = testutil.RandomBytes(maxStrLen)
		case storage.ColumnType_INT64:
			value = rand.Int63()
		case storage.ColumnType_TIMESTAMP:
			//value = time.Now().UnixNano()
			value = n + int64(rid)
		}

		tc.rids = append(tc.rids, rid)
		tc.values = append(tc.values, value)

	}

	return tc
}
