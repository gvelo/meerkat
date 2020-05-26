package ingestion

import (
	"encoding/binary"
	"meerkat/internal/storage/io"
)

type RowSetWriter struct {
	Buf *io.Buffer
}

func NewRowSetWriter(cap int) *RowSetWriter {
	return &RowSetWriter{
		Buf: io.NewBuffer(cap),
	}
}

func (rs *RowSetWriter) WriteString(colId int, str string) {

	rs.Reserve(len(str) + binary.MaxVarintLen64*2)

	rs.Buf.WriteIntAsUVarInt(colId)
	rs.Buf.WriteString(str)

}

func (rs *RowSetWriter) WriteInt(colId int, i int) {

	rs.Reserve(binary.MaxVarintLen64 * 2)

	rs.Buf.WriteIntAsUVarInt(colId)
	rs.Buf.WriteVarInt(i)

}

func (rs *RowSetWriter) WriteIntAsUVarInt(colId int, i int) {

	rs.Reserve(binary.MaxVarintLen64 * 2)

	rs.Buf.WriteIntAsUVarInt(colId)
	rs.Buf.WriteIntAsUVarInt(i)

}

func (rs *RowSetWriter) WriteFixedUInt64(colId int, i uint64) {

	rs.Reserve(binary.MaxVarintLen64 * 2)

	rs.Buf.WriteIntAsUVarInt(colId)
	rs.Buf.WriteFixedUInt64(i)

}

func (rs *RowSetWriter) WriteFixedInt64(colId int, i int64) {
	rs.WriteFixedUInt64(colId, uint64(i))
}

func (rs *RowSetWriter) WriteFloat(colId int, f float64) {

	rs.Reserve(binary.MaxVarintLen64 * 2)

	rs.Buf.WriteIntAsUVarInt(colId)
	rs.Buf.WriteFloat64(f)

}

func (rs *RowSetWriter) Reserve(size int) {
	if size > rs.Buf.Available() {
		rs.Buf.Grow((rs.Buf.Cap() + size) * 2)
	}
}

func (rs *RowSetWriter) grow(newSize int) {
	rs.Buf.Grow(newSize)
}
