package storage

import "meerkat/internal/storage/colval"

//go:generate protoc -I . -I ../../build/proto/ --plugin ../../build/protoc-gen-gogofaster --gogofaster_out=plugins=grpc,paths=source_relative:.  ./storage.proto

const (
	TSColumnName = "_ts"
)

type ColumnSource interface {
	HasNext() bool
	// TODO: remove, this flag is now on the columninfo
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

type SegmentSourceInfo struct {
	Id           []byte
	DatabaseId   []byte
	DatabaseName string
	TableName    string
	PartitionId  uint64
	Len          uint32
	Interval     Interval
	Columns      []ColumnSourceInfo
}

type ColumnSourceInfo struct {
	Name       string
	ColumnType ColumnType
	IndexType  IndexType
	Encoding   Encoding
	Nullable   bool
	Len        uint32
}

type SegmentSource interface {
	Info() SegmentSourceInfo
	// TODO(gvelo): add blockSize and blockLen.
	ColumnSource(colName string, blockSize int) ColumnSource
}
