package readers

import (
	"errors"
	"eventdb/io"
	"eventdb/segment"
	"eventdb/segment/ondsk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"path/filepath"
)

const (
	idxExt = ".idx"
	binExt = ".bin"
)

type segmentReader struct {
	log  zerolog.Logger
	path string
}

func (sr *segmentReader) read() (*ondsk.Segment, error) {

	indexInfo, err := sr.readIndexInfo()

	if err != nil {
		return nil, err
	}

	idx, err := sr.readIdx(indexInfo.Fields)

	if err != nil {
		return nil, err
	}

	return &ondsk.Segment{
		IndexInfo: indexInfo,
		Idx:       idx,
	}, nil

}

func (sr *segmentReader) readIndexInfo() (*segment.IndexInfo, error) {

	path := filepath.Join(sr.path, "info")

	file, err := io.MMap(path)

	if err != nil {
		return nil, err
	}

	defer file.UnMap()

	br := file.NewBinaryReader()

	fType, err := br.ReadHeader()

	if fType != io.SegmentInfo {
		return nil, errors.New("invalid file type")
	}

	indexName, err := br.ReadString()

	if err != nil {
		return nil, err
	}

	fieldCount, err := br.ReadVarInt()

	if err != nil {
		return nil, err
	}

	fieldInfos := make([]*segment.FieldInfo, fieldCount)

	for i := 0; i < int(fieldCount); i++ {

		name, err := br.ReadString()

		if err != nil {
			return nil, err
		}

		fieldType, err := br.ReadVarInt()

		if err != nil {
			return nil, err
		}

		indexed, err := br.ReadVarInt()

		if err != nil {
			return nil, err
		}

		fieldInfo := &segment.FieldInfo{
			ID:   i,
			Name: name,
			Type: segment.FieldType(fieldType),
		}

		if indexed == 1 {
			fieldInfo.Index = true
		}

		fieldInfos[i] = fieldInfo
	}

	return &segment.IndexInfo{
		Name:   indexName,
		Fields: fieldInfos,
	}, nil

}

func (sr *segmentReader) readIdx(fields []*segment.FieldInfo) ([]interface{}, error) {

	idx := make([]interface{}, len(fields))

	for i, fieldInfo := range fields {

		if fieldInfo.Index {

			fileName := filepath.Join(sr.path, fieldInfo.Name+idxExt)

			var index interface{}
			var err error

			switch fieldInfo.Type {
			case segment.FieldTypeText:
				index, err = ReadTrie(fileName)
			case segment.FieldTypeKeyword:
				index, err = ReadTrie(fileName)
			case segment.FieldTypeInt:
				index, err = ReadSkipList(fileName)
			case segment.FieldTypeFloat:
				index, err = ReadSkipList(fileName)
			case segment.FieldTypeTimestamp:
				index, err = ReadSkipList(fileName)
			default:
				sr.log.Panic().
					Str("field", fieldInfo.Name).
					Msgf("Unknown index type: %T", fieldInfo.Type)
			}

			if err != nil {
				return nil, err
			}

			idx[i] = index

		}
	}

	return idx, nil
}

func newSegmentReader(path string) *segmentReader {

	sr := &segmentReader{
		path: path,
	}

	sr.log = log.With().
		Str("component", "segmentReader").
		Str("path", path).
		Logger()

	return sr

}

func ReadSegment(path string) (*ondsk.Segment, error) {

	reader := newSegmentReader(path)
	s, err := reader.read()
	if err != nil {
		return nil, err
	}
	return s, nil

}
