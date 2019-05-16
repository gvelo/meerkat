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
)

type SegmentReader struct {
	log  zerolog.Logger
	path string
}

func (sr *SegmentReader) Read() (*ondsk.Segment, error) {

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

func (sr *SegmentReader) readIndexInfo() (*segment.IndexInfo, error) {

	file := filepath.Join(sr.path, "info")

	br, err := io.NewBinaryReader(file)

	defer br.Close()

	if err != nil {
		return nil, err
	}

	fType, err := br.ReadHeader()

	if fType != io.SegmentInfo {
		return nil, errors.New("invalid file type")
	}

	indexName, err := br.DecodeStringBytes()

	if err != nil {
		return nil, err
	}

	fieldCount, err := br.DecodeVarint()

	if err != nil {
		return nil, err
	}

	fieldInfos := make([]*segment.FieldInfo, fieldCount)

	for i := 0; i < int(fieldCount); i++ {

		name, err := br.DecodeStringBytes()

		if err != nil {
			return nil, err
		}

		fieldType, err := br.DecodeVarint()

		if err != nil {
			return nil, err
		}

		indexed, err := br.DecodeVarint()

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

func (sr *SegmentReader) readIdx(fields []*segment.FieldInfo) ([]interface{}, error) {

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
				//index, err = ReadTrie(fileName)
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

func NewSegmentReader(path string) *SegmentReader {

	sr := &SegmentReader{
		path: path,
	}

	sr.log = log.With().
		Str("component", "SegmentReader").
		Str("path", path).
		Logger()

	return sr

}
