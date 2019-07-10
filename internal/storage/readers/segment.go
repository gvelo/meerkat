// Copyright 2019 The Meerkat Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package readers

import (
	"errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/ondsk"
	"path/filepath"
)

const (
	idxPosExt = ".ipos" // index for value in posting
	idxPagExt = ".ipag" // index
	pagExt    = ".pag"  // pages encoded
	posExt    = ".pos"  // posting lists for values
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

	cols, err := sr.readColumns(indexInfo.Fields)

	if err != nil {
		return nil, err
	}

	return &ondsk.Segment{
		IndexInfo: indexInfo,
		Columns:   cols,
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

		fieldType, _ := br.ReadVarInt()

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

func (sr *segmentReader) readColumns(fields []*segment.FieldInfo) ([]ondsk.Column, error) {

	cols := make([]ondsk.Column, len(fields))

	for i, fieldInfo := range fields {

		var idxPag *ondsk.OnDiskColumnIdx
		var pag *ondsk.OnDiskColumn
		var err error

		idxPagFilename := filepath.Join(sr.path, fieldInfo.Name+idxPagExt)
		pagFilename := filepath.Join(sr.path, fieldInfo.Name+pagExt)

		idxPag, err = ReadColumnIdx(idxPagFilename, fieldInfo)
		if err != nil {
			sr.log.Error().Msgf("Error reading idxcolumn: %v,  %v", idxPagFilename, err)
		}
		pag, err = ReadColumn(pagFilename, fieldInfo)
		if err != nil {
			sr.log.Error().Msgf("Error reading column:  %v,  %v", idxPagFilename, err)
		}
		cols[i] = &ondsk.ColumnImpl{IdxPag: idxPag, Pag: pag}
	}

	return cols, nil

}

func (sr *segmentReader) readIdx(fields []*segment.FieldInfo) ([]interface{}, error) {

	idx := make([]interface{}, len(fields))

	for i, fieldInfo := range fields {

		if fieldInfo.Index {

			fileName := filepath.Join(sr.path, fieldInfo.Name+idxPosExt)

			var index interface{}
			var err error

			switch fieldInfo.Type {
			case segment.FieldTypeText, segment.FieldTypeKeyword:
				index, err = ReadTrie(fileName)
			case segment.FieldTypeInt, segment.FieldTypeTimestamp:
				index, err = ReadSkipList(fileName, ondsk.IntInterface{})
			case segment.FieldTypeFloat:
				index, err = ReadSkipList(fileName, ondsk.FloatInterface{})
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
