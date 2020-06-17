package storage

import "meerkat/internal/storage/colval"

const (
	TSColumnName = "_ts"
)

type ColumnSource interface {
	HasNext() bool
	HasNulls() bool
}

type Int64ColumnSource interface {
	ColumnSource
	Next() colval.Int64ColValues
}

type Int32ColumnSource interface {
	ColumnSource
	Next() colval.Int32ColValues
}

type Float64ColumnSource interface {
	ColumnSource
	Next() colval.Float64ColValues
}

type BoolColumnSource interface {
	ColumnSource
	Next() colval.BoolColValues
}

type ByteSliceColumnSource interface {
	ColumnSource
	Next() colval.ByteSliceColValues
}

type SegmentSource interface {
	SegmentInfo() *SegmentInfo
	Columns() []*ColumnInfo
	ColumnSource(colName string, blockSize int) ColumnSource
}
