package storage

import (
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/storage/colval"
	"path"
)

//go:generate protoc -I . -I ../../build/proto/ --plugin ../../build/protoc-gen-gogofaster --gogofaster_out=plugins=grpc,paths=source_relative:.  ./storage.proto

const (
	TSColumnName      = "_ts"
	segmentFolderName = "segments"
)

type ColumnSource interface {
	HasNext() bool
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
	// TODO(gvelo): indexing and encoding concerns should be in another component.
	IndexType IndexType
	Encoding  Encoding
	Nullable  bool
	Len       uint32
}

type SegmentSource interface {
	Info() SegmentSourceInfo
	ColumnSource(colName string, blockSize int) ColumnSource
}

type SegmentStorage interface {
	WriteSegment(src SegmentSource) *SegmentInfo
	OpenSegment(info *SegmentInfo) *Segment
	DeleteSegment(id uuid.UUID)
}

type defaultStorage struct {
	dbPath string
	log    zerolog.Logger
}

func NewStorage(dbPath string) SegmentStorage {
	return &defaultStorage{
		dbPath: dbPath,
		log:    log.With().Str("component", "Storage").Logger(),
	}
}

func (d defaultStorage) WriteSegment(src SegmentSource) *SegmentInfo {

	info := src.Info()

	fileName := buildSegmentFileName(info.Id)

	filePath := path.Join(d.dbPath, segmentFolderName, fileName)

	logger := d.log.With().
		Str("sid", fileName).
		Str("table", info.TableName).
		Uint64("partition", info.PartitionId).Logger()

	logger.Debug().Msg("writing segment")

	WriteSegment(filePath, src)

	logger.Debug().Msg("segment successfully written")

	// TODO(gvelo) build SegmentInfo

	return &SegmentInfo{}

}

func (d defaultStorage) OpenSegment(info *SegmentInfo) *Segment {

	fileName := buildSegmentFileName(info.Id)

	filePath := path.Join(d.dbPath, segmentFolderName, fileName)

	logger := d.log.With().
		Str("sid", fileName).
		Str("table", info.TableName).
		Uint64("partition", info.PartitionId).Logger()

	logger.Debug().Msg("opening segment")

	segment := ReadSegment(filePath)

	return segment

}

func (d defaultStorage) DeleteSegment(id uuid.UUID) {
	panic("implement me")
}

func buildSegmentFileName(id []byte) string {
	return base64.RawURLEncoding.EncodeToString(id)
}
