package storage

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/storage/colval"
	"meerkat/internal/storage/encoding"
	"os"
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
	Id           uuid.UUID
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
	Encoding  encoding.Type
	Nullable  bool
	Len       uint32
}

type SegmentSource interface {
	Info() SegmentSourceInfo
	ColumnSource(colName string, blockSize int) ColumnSource
}

type Segment interface {
	Info() *SegmentInfo
	Column(name string) Column
	Close()
}

type SegmentStorage interface {
	WriteSegment(src SegmentSource) *SegmentInfo
	// TODO(gvelo) currently we are using *SegmentInfo as a param just
	// to add a reference to SegmentInfo in the segment. We should serialize
	// and write the SegmentInfo into the segment and replaze this param by the
	// the segment id.
	OpenSegment(info *SegmentInfo) Segment
	DeleteSegment(id uuid.UUID)
}

type defaultStorage struct {
	dbPath string
	log    zerolog.Logger
}

func NewStorage(dbPath string) SegmentStorage {

	segmentFolderPath := path.Join(dbPath, segmentFolderName)

	err := os.MkdirAll(segmentFolderPath, 0770)

	if err != nil {
		panic(fmt.Errorf("could not create dir %v , %v", segmentFolderPath, err))
	}

	return &defaultStorage{
		dbPath: dbPath,
		log:    log.With().Str("component", "Storage").Logger(),
	}
}

func (d *defaultStorage) WriteSegment(src SegmentSource) *SegmentInfo {

	srcInfo := src.Info()

	filePath := d.buildSegmentPathFromUUID(srcInfo.Id)

	logger := d.log.With().
		Str("sid", filePath).
		Str("table", srcInfo.TableName).
		Uint64("partition", srcInfo.PartitionId).Logger()

	logger.Debug().Msg("writing segment")

	WriteSegment(filePath, src)

	logger.Debug().Msg("segment successfully written")

	// TODO(gvelo) build SegmentInfo

	var colInfos []*ColumnInfo

	for _, srcCol := range srcInfo.Columns {

		colInfo := &ColumnInfo{
			Name:           srcCol.Name,
			ColumnType:     srcCol.ColumnType,
			IndexType:      srcCol.IndexType,
			Encoding:       Encoding_PLAIN,
			Nullable:       srcCol.Nullable,
			Len:            srcCol.Len,
			Cardinality:    0,
			SizeOnDisk:     0,
			CompressedSize: 0,
			NullCount:      0,
		}

		colInfos = append(colInfos, colInfo)
	}

	segInfo := &SegmentInfo{
		Id:           srcInfo.Id[:],
		DatabaseName: srcInfo.DatabaseName,
		TableName:    srcInfo.TableName,
		PartitionId:  srcInfo.PartitionId,
		Len:          srcInfo.Len,
		Interval: &Interval{
			From: srcInfo.Interval.From,
			To:   srcInfo.Interval.To,
		},
		Columns: colInfos,
	}

	return segInfo

}

func (d *defaultStorage) OpenSegment(info *SegmentInfo) Segment {

	filePath := d.buildSegmentPath(info.Id)

	logger := d.log.With().
		Str("sid", filePath).
		Str("table", info.TableName).
		Uint64("partition", info.PartitionId).Logger()

	logger.Debug().Msg("opening segment")

	segment := ReadSegment(filePath)

	// TODO(segmentInfo should be written into the segment file)
	segment.segmentInfo = info

	return segment

}

func (d *defaultStorage) DeleteSegment(id uuid.UUID) {
	filepath := d.buildSegmentPathFromUUID(id)
	err := os.Remove(filepath)
	if err != nil {
		panic(err)
	}
}

func (d *defaultStorage) buildSegmentPath(id []byte) string {
	fileName := buildSegmentFileName(id)
	filePath := path.Join(d.dbPath, segmentFolderName, fileName)
	return filePath
}

func (d *defaultStorage) buildSegmentPathFromUUID(id uuid.UUID) string {
	return d.buildSegmentPath(id[:])
}

func buildSegmentFileName(id []byte) string {
	return base64.RawURLEncoding.EncodeToString(id)
}

func buildSegmentFileNameFromUUID(id uuid.UUID) string {
	return buildSegmentFileName(id[:])
}
