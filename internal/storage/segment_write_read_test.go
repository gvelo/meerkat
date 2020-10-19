package storage

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"meerkat/internal/storage/colval"
	"meerkat/internal/util/testutil"
	"os"
	"path"
	"testing"
)

func generateRandomColumnData(colType ColumnType, nullable bool, maxColLen uint32) ([]interface{}, uint32) {

	colData := make([]interface{}, 0, maxColLen)
	var colLen uint32

	for rid := uint32(0); rid < maxColLen; rid++ {

		if nullable {

			if rand.Intn(3) == 2 {
				colData = append(colData, nil)
				continue
			}

		}

		switch colType {
		case ColumnType_INT64:
			colData = append(colData, rand.Int63())
		case ColumnType_INT32:
			colData = append(colData, rand.Int31())
		case ColumnType_FLOAT64:
			colData = append(colData, rand.Float64())
		case ColumnType_STRING:
			colData = append(colData, testutil.RandomBytes(10))
		case ColumnType_TIMESTAMP:
			colData = append(colData, rand.Int())
		default:
			panic("column type not implemented yet")
		}

		colLen++

	}

	return colData, colLen

}

type Int64TestColumnSrc struct {
	data      []interface{}
	pos       int
	blockSize int
}

func NewInt64TestColumnSrc(data []interface{}, blockSize int) *Int64TestColumnSrc {
	return &Int64TestColumnSrc{
		data:      data,
		pos:       0,
		blockSize: blockSize,
	}
}

func (s *Int64TestColumnSrc) HasNext() bool {
	return s.pos < len(s.data)
}

func (s *Int64TestColumnSrc) HasNulls() bool {
	panic("implement me")
}

func (s *Int64TestColumnSrc) Next() colval.Int64ColValues {

	var rids []uint32
	var values []int64

	for ; s.pos < len(s.data) && len(values) < s.blockSize; s.pos++ {

		if s.data[s.pos] == nil {
			continue
		}

		value := s.data[s.pos].(int64)

		values = append(values, value)
		rids = append(rids, uint32(s.pos))

	}

	return colval.NewInt64ColValues(values, rids)

}

func TestInt64TestColumnSrc(t *testing.T) {

	const blockSize = 123
	const segmentLen = uint32(2000)

	colInfo := ColumnSourceInfo{
		Name:       "testColumn",
		ColumnType: ColumnType_INT64,
		IndexType:  IndexType_NONE,
		Encoding:   Encoding_PLAIN,
		Nullable:   true,
	}

	colData, colLen := generateRandomColumnData(ColumnType_INT64, colInfo.Nullable, segmentLen)
	colInfo.Len = colLen

	src := NewInt64TestColumnSrc(colData, blockSize)

	var values []interface{}
	var rids []uint32

	for src.HasNext() {

		v := src.Next()

		for _, value := range v.Values() {
			values = append(values, value)
		}

		rids = append(rids, v.Rid()...)

	}

	nullCount := countNulls(colData)

	assert.Equal(t, nullCount, int(segmentLen)-(len(values)))

	assertEqual(t, colData, values, rids)

}

func countNulls(values []interface{}) int {
	i := 0
	for _, value := range values {
		if value == nil {
			i++
		}
	}
	return i
}

type ByteSliceTestColumnSrc struct {
	data      []interface{}
	pos       int
	blockSize int
}

func NewByteSliceTestColumnSrc(data []interface{}, blockSize int) *ByteSliceTestColumnSrc {
	return &ByteSliceTestColumnSrc{
		data:      data,
		pos:       0,
		blockSize: blockSize,
	}
}

func (s *ByteSliceTestColumnSrc) HasNext() bool {
	return s.pos < len(s.data)
}

func (s *ByteSliceTestColumnSrc) HasNulls() bool {
	panic("implement me")
}

func (s *ByteSliceTestColumnSrc) Next() colval.ByteSliceColValues {

	var rids []uint32
	var data []byte
	var offsets []int

	for ; s.pos < len(s.data) && len(data) < s.blockSize; s.pos++ {

		if s.data[s.pos] == nil {
			continue
		}

		value := s.data[s.pos].([]byte)

		data = append(data, value...)
		offsets = append(offsets, len(data))
		rids = append(rids, uint32(s.pos))

	}

	return colval.NewByteSliceColValues(data, rids, offsets)

}

func TestByteSliceTestColumnSrc(t *testing.T) {

	const blockSize = 123
	const segmentLen = uint32(2000)

	colInfo := ColumnSourceInfo{
		Name:       "testColumn",
		ColumnType: ColumnType_STRING,
		IndexType:  IndexType_NONE,
		Encoding:   Encoding_PLAIN,
		Nullable:   true,
	}

	colData, colLen := generateRandomColumnData(ColumnType_STRING, colInfo.Nullable, segmentLen)
	colInfo.Len = colLen

	src := NewByteSliceTestColumnSrc(colData, blockSize)

	var values []interface{}
	var rids []uint32

	for src.HasNext() {

		v := src.Next()

		for i := 0; i < v.Len(); i++ {
			values = append(values, v.Get(i))
		}

		rids = append(rids, v.Rid()...)

	}

	nullCount := countNulls(colData)

	assert.Equal(t, nullCount, int(segmentLen)-len(values))

	assertEqual(t, colData, values, rids)

}

func assertEqual(t *testing.T, expected []interface{}, values []interface{}, rids []uint32) {

	actual := make([]interface{}, len(expected))

	for i, rid := range rids {
		actual[rid] = values[i]
	}

	assert.Equal(t, expected, actual)

}

type TestSegmentSource struct {
	columns map[string][]interface{}
	colInfo map[string]ColumnSourceInfo
	info    SegmentSourceInfo
}

func NewTestSegmentSource(info SegmentSourceInfo) *TestSegmentSource {

	src := &TestSegmentSource{
		columns: make(map[string][]interface{}),
		colInfo: make(map[string]ColumnSourceInfo),
		info:    info,
	}

	for _, columnSourceInfo := range info.Columns {
		src.colInfo[columnSourceInfo.Name] = columnSourceInfo
	}

	return src
}

func (t *TestSegmentSource) Info() SegmentSourceInfo {
	return t.info
}

func (t *TestSegmentSource) ColumnSource(colName string, blockSize int) ColumnSource {

	if colInfo, found := t.colInfo[colName]; found {

		switch colInfo.ColumnType {
		case ColumnType_INT64:
			return NewInt64TestColumnSrc(t.columns[colInfo.Name], blockSize)
		case ColumnType_INT32:
			panic("not implemented yet")
		case ColumnType_FLOAT64:
			panic("not implemented yet")
		case ColumnType_TIMESTAMP:
			panic("not implemented yet")
		case ColumnType_STRING:
			return NewByteSliceTestColumnSrc(t.columns[colInfo.Name], blockSize)
		default:
			panic("not implemented yet")
		}

	}

	panic("column not found")

}

func generateRandomSegmentSrc(info SegmentSourceInfo) *TestSegmentSource {

	src := NewTestSegmentSource(info)

	for i, columnSourceInfo := range info.Columns {
		colData, colLen := generateRandomColumnData(columnSourceInfo.ColumnType, columnSourceInfo.Nullable, info.Len)
		info.Columns[i].Len = colLen
		src.columns[columnSourceInfo.Name] = colData
	}

	return src

}

func TestReadWriteSegment(t *testing.T) {

	segmentLen := 1024 * 1024 * 2

	id := [16]byte(uuid.New())
	segmentSrc := generateRandomSegmentSrc(SegmentSourceInfo{
		Id:           id[:],
		DatabaseId:   id[:],
		DatabaseName: "test-db",
		TableName:    "test-table",
		PartitionId:  0,
		Len:          uint32(segmentLen),
		Interval:     Interval{},
		Columns: []ColumnSourceInfo{
			{
				Name:       "column-test-int64",
				ColumnType: ColumnType_INT64,
				IndexType:  IndexType_NONE,
				Encoding:   Encoding_PLAIN,
				Nullable:   false,
			},
			{
				Name:       "column-test-int64-nullable",
				ColumnType: ColumnType_INT64,
				IndexType:  IndexType_NONE,
				Encoding:   Encoding_PLAIN,
				Nullable:   true,
			},
			{
				Name:       "column-test-string",
				ColumnType: ColumnType_STRING,
				IndexType:  IndexType_NONE,
				Encoding:   Encoding_PLAIN,
				Nullable:   false,
			},
			{
				Name:       "column-test-string-nullable",
				ColumnType: ColumnType_STRING,
				IndexType:  IndexType_NONE,
				Encoding:   Encoding_PLAIN,
				Nullable:   true,
			},
		},
	})

	segmentFile := path.Join(os.TempDir(), "test-segment")

	WriteSegment(segmentFile, segmentSrc)

	segment := ReadSegment(segmentFile)

	checkColumnsIterators(t, segment, segmentSrc)
	checkColumnsReads(t, segment, segmentSrc)

}

func checkColumnsIterators(t *testing.T, segment *Segment, src *TestSegmentSource) {

	for _, info := range src.colInfo {

		t.Logf("testing iterator on column %v", info.Name)

		if col, found := segment.columns[info.Name]; found {
			actual := readColumnIter(col, info, src.Info().Len)
			expected := src.columns[info.Name]
			assert.Equal(t, expected, actual)
			continue
		}

		panic("column not found")

	}

}

func readColumnIter(column interface{}, info ColumnSourceInfo, srcLen uint32) []interface{} {

	values := make([]interface{}, 0, srcLen)

	switch c := column.(type) {

	case *int64Column:
		iter := c.Iterator()
		for iter.HasNext() {
			vec := iter.Next()
			for i, v := range vec.Values() {
				if info.Nullable && !vec.IsValid(i) {
					values = append(values, nil)
				} else {
					values = append(values, v)
				}
			}
		}

	case *binaryColumn:
		iter := c.Iterator()
		for iter.HasNext() {
			vec := iter.Next()
			for i := 0; i < vec.Len(); i++ {
				if info.Nullable && !vec.IsValid(i) {
					values = append(values, nil)
				} else {
					values = append(values, vec.Get(i))
				}
			}
		}

	default:
		panic("unknown column type")
	}

	return values

}

func checkColumnsReads(t *testing.T, segment *Segment, src *TestSegmentSource) {

	rids := []uint32{0, 123, 124, 1023, 2048, 8191, 8192, 8193, 10212, 20000, 24575}

	for _, colInfo := range src.Info().Columns {

		t.Logf("testing reads on column %v", colInfo.Name)

		if column, found := segment.columns[colInfo.Name]; found {

			var expected []interface{}

			for _, rid := range rids {
				v := src.columns[colInfo.Name][rid]
				expected = append(expected, v)
			}

			actual := readColumn(column, colInfo, rids)

			assert.Equal(t, expected, actual)

			continue

		}

		panic("column not found")

	}

}

func readColumn(column interface{}, colInfo ColumnSourceInfo, rids []uint32) []interface{} {

	var values []interface{}

	switch c := column.(type) {
	case *int64Column:
		vec := c.Reader().Read(rids)
		for i, v := range vec.Values() {
			if colInfo.Nullable && !vec.IsValid(i) {
				values = append(values, nil)
			} else {
				values = append(values, v)
			}
		}

	case *binaryColumn:
		vec := c.Reader().Read(rids)
		for i := 0; i < vec.Len(); i++ {
			if colInfo.Nullable && !vec.IsValid(i) {
				values = append(values, nil)
			} else {
				v := vec.Get(i)
				values = append(values, v)
			}
		}

	default:
		panic("unknown column type")
	}

	return values

}
