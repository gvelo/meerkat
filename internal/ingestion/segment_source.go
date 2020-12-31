package ingestion

import (
	"github.com/google/uuid"
	"meerkat/internal/storage"
	"meerkat/internal/storage/colval"
	"meerkat/internal/storage/encoding"
	"meerkat/internal/util/sliceutil"
	"sort"
	"strconv"
	"time"
)

type SegmentSource struct {
	info  storage.SegmentSourceInfo
	perm  []int
	table *TableBuffer
}

func NewSegmentSource(table *TableBuffer) *SegmentSource {

	ts, found := table.Columns()[storage.TSColumnName].buf.(*TSBuffer)

	if !found {
		panic("cannot found TS column")
	}

	perm := sortTS(ts.Values())

	return &SegmentSource{
		info:  buildSrcInfo(table),
		table: table,
		perm:  perm,
	}

}

func buildSrcInfo(table *TableBuffer) storage.SegmentSourceInfo {

	id := uuid.New()

	tsCol := table.Columns()[storage.TSColumnName].buf.(*TSBuffer)

	tsFrom := tsCol.Values()[0]
	tsTo := tsCol.Values()[tsCol.Len()-1]

	info := storage.SegmentSourceInfo{
		Id:           id,
		DatabaseName: "",
		TableName:    table.tableName,
		PartitionId:  table.partitionID,
		Len:          uint32(table.len),
		Interval: storage.Interval{
			From: time.Unix(0, tsFrom),
			To:   time.Unix(0, tsTo),
		},
	}

	var columnsInfo []storage.ColumnSourceInfo

	for name, column := range table.Columns() {

		nullable := column.buf.Len() != table.Len()

		colInfo := storage.ColumnSourceInfo{
			Name:       name,
			ColumnType: column.colType,
			IndexType:  storage.IndexType_NONE,
			Encoding:   encoding.Plain,
			Nullable:   nullable,
			Len:        uint32(column.buf.Len()),
		}

		columnsInfo = append(columnsInfo, colInfo)

	}

	info.Columns = columnsInfo

	return info

}

func (s *SegmentSource) Info() storage.SegmentSourceInfo {
	return s.info
}

func (s *SegmentSource) ColumnSource(colName string, blockSize int) storage.ColumnSource {

	if c, found := s.table.columns[colName]; found {

		switch buf := c.buf.(type) {

		case *TSBuffer:
			return NewTsSource(buf, blockSize)
		case *ByteSliceSparseBuffer:
			return s.createSrc(buf, c.colType, blockSize)
		default:
			panic("unknown column type")
		}

	}

	panic("column not found")

}

func (s *SegmentSource) createSrc(buf *ByteSliceSparseBuffer, columnType storage.ColumnType, blockSize int) storage.ColumnSource {

	// TODO(gvelo): reuse the dense buffers.

	denseBuf := buf.ToDenseBuffer(s.table.len)

	switch columnType {
	case storage.ColumnType_INT64:
		return NewInt64DynamicSrc(denseBuf, s.perm, blockSize)
	case storage.ColumnType_STRING:
		return NewByteSliceDynamicSrc(denseBuf, blockSize, s.perm)
	default:
		panic("unknown column type")
	}

}

func sortTS(ts []int64) []int {

	perm := make([]int, len(ts))

	for i := 0; i < len(perm); i++ {
		perm[i] = i
	}

	tsSlice := &TSSlice{
		ts:   ts,
		perm: perm,
	}

	sort.Stable(tsSlice)

	return perm

}

type TSSlice struct {
	ts   []int64
	perm []int
}

func (t *TSSlice) Len() int {
	return len(t.ts)
}

func (t *TSSlice) Less(i, j int) bool {
	return t.ts[i] < t.ts[j]
}

func (t *TSSlice) Swap(i, j int) {
	t.ts[i], t.ts[j] = t.ts[j], t.ts[i]
	t.perm[i], t.perm[j] = t.perm[j], t.perm[i]
}

func NewTsSource(buf *TSBuffer, blockSize int) *TSSource {
	return &TSSource{
		dstSize: blockSize,
		srcBuf:  buf.Values(),
		rid:     make([]uint32, blockSize),
	}
}

type TSSource struct {
	srcBuf  []int64
	rid     []uint32
	start   int
	end     int
	dstSize int
	pos     int
}

func (cs *TSSource) HasNext() bool {
	return cs.end < len(cs.srcBuf)
}

func (cs *TSSource) Next() colval.Int64ColValues {

	cs.start = cs.end
	cs.end = cs.start + cs.dstSize
	dstLen := cs.dstSize

	if cs.end > len(cs.srcBuf) {
		cs.end = len(cs.srcBuf)
		dstLen = cs.end - cs.start
	}

	for i := 0; i < dstLen; i++ {
		cs.rid[i] = uint32(cs.pos)
		cs.pos++
	}

	return colval.NewInt64ColValues(cs.srcBuf[cs.start:cs.end], cs.rid[0:dstLen])

}

type Int64DynamicSrc struct {
	srcBuf *ByteSliceDenseBuffer
	dstBuf []int64
	valids []bool
	perm   []int
	pos    uint32
	rids   []uint32
}

func NewInt64DynamicSrc(buf *ByteSliceDenseBuffer, perm []int, blockSize int) *Int64DynamicSrc {
	return &Int64DynamicSrc{
		srcBuf: buf,
		dstBuf: make([]int64, blockSize),
		rids:   make([]uint32, blockSize),
		valids: buf.Valids(),
		perm:   perm,
	}
}

func (s *Int64DynamicSrc) HasNext() bool {
	return int(s.pos) < s.srcBuf.Len()
}

func (s *Int64DynamicSrc) Next() colval.Int64ColValues {

	i := 0

	for i < len(s.dstBuf) && int(s.pos) < s.srcBuf.Len() {

		p := uint32(s.perm[s.pos])

		if s.valids[p] {

			bytes := s.srcBuf.Value(p)

			intVal, err := strconv.ParseInt(sliceutil.B2S(bytes), 10, 64)

			if err != nil {
				panic(err)
			}

			s.dstBuf[i] = intVal
			s.rids[i] = s.pos

			i++
			s.pos++
			continue
		}

		s.pos++

	}

	return colval.NewInt64ColValues(s.dstBuf[:i], s.rids[:i])

}

type ByteSliceDynamicSrc struct {
	maxSize    int
	srcBuf     *ByteSliceDenseBuffer
	dstBuf     []byte
	dstOffsets []int
	valid      []bool
	rids       []uint32
	perm       []int
	rid        uint32
}

func NewByteSliceDynamicSrc(src *ByteSliceDenseBuffer, maxSize int, perm []int) *ByteSliceDynamicSrc {
	return &ByteSliceDynamicSrc{
		srcBuf:     src,
		maxSize:    maxSize,
		dstBuf:     make([]byte, maxSize),
		dstOffsets: make([]int, 0, 1024*8),
		valid:      src.Valids(),
		rids:       make([]uint32, 0, 1024*8),
		perm:       perm,
	}
}

func (s *ByteSliceDynamicSrc) HasNext() bool {
	return int(s.rid) < s.srcBuf.Len()
}

func (s *ByteSliceDynamicSrc) Next() colval.ByteSliceColValues {

	s.dstBuf = s.dstBuf[:0]
	s.dstOffsets = s.dstOffsets[:0]
	s.rids = s.rids[:0]

	for ; int(s.rid) < s.srcBuf.Len(); s.rid++ {

		pos := uint32(s.perm[s.rid])

		if s.valid[pos] {

			value := s.srcBuf.Value(pos)

			available := s.maxSize - len(s.dstBuf)

			if len(value) > available && len(s.dstBuf) > 0 {
				break
			}

			s.dstBuf = append(s.dstBuf, value...)
			s.dstOffsets = append(s.dstOffsets, len(s.dstBuf))
			s.rids = append(s.rids, s.rid)

		}

	}

	return colval.NewByteSliceColValues(s.dstBuf, s.rids, s.dstOffsets)

}
