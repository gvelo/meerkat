// Copyright 2020 The Meerkat Authors
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

package storage

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"meerkat/internal/schema"
	"meerkat/internal/storage/io"
)

const (
	SegmentVersion1 = 1
)

func ReadSegment(path string) (*segment, error) {

	f, err := io.MMap(path)

	if err != nil {
		return nil, err
	}

	br := f.NewBinaryReader()

	header := br.ReadSlice(0, len(MagicNumber))

	if string(header) != MagicNumber {
		panic("unknown file type")
	}

	segmentVersion := br.ReadByte()

	// we only have just one segment version

	switch segmentVersion {
	case SegmentVersion1:
		s := NewSegment(f)
		err := s.read()
		return s, err
	default:
		return nil, errors.New("unknown segment version")
	}

}

func NewSegment(f *io.MMFile) *segment {
	return &segment{
		f:       f,
		columns: make(map[string]interface{}),
	}
}

type segment struct {
	f         *io.MMFile
	id        uuid.UUID
	from      int
	to        int
	numOfRows int
	start     int
	tableId   string
	tableName string
	numOfCol  int
	//columns   map[string]Column
	columns map[string]interface{}
}

type colData struct {
	colType     schema.FieldType
	id          string
	name        string
	offsetStart int
	offsetEnd   int
}

// start = entry
func (s *segment) read() error {

	// magicNumber + version
	s.start = len(MagicNumber) + 1

	br := s.f.NewBinaryReader()

	br.Entry()

	br.ReadRaw(s.id[:])

	s.tableId = br.ReadString()

	s.tableName = br.ReadString()

	s.from = br.ReadFixed64()

	s.to = br.ReadFixed64()

	s.numOfRows = br.ReadUVarint()

	s.numOfCol = br.ReadUVarint()

	cd := make([]colData, s.numOfCol)

	fmt.Println("numofcol", s.numOfCol, s.numOfRows)

	for i := 0; i < s.numOfCol; i++ {
		c := colData{}
		c.id = br.ReadString()
		c.name = br.ReadString()
		c.colType = schema.FieldType(br.ReadByte())
		c.offsetEnd = br.ReadUVarint()
		if i == 0 {
			c.offsetStart = s.start
		} else {
			c.offsetStart = cd[i-1].offsetEnd
		}
		cd[i] = c
	}

	s.readColumns(cd)

	// TODO(gvelo) recover from panic and return the err.

	return nil

}

func (s *segment) readColumns(cd []colData) {

	for _, cData := range cd {

		//var colStart, colEnd int

		//if i == 0 {
		//	colStart = s.start
		//} else {
		//	colStart = colEnd
		//}

		fmt.Println("cdata ", cData.offsetStart, cData.offsetEnd)

		//colEnd = cData.offsetEnd

		fmt.Println("================", cData.offsetStart, cData.offsetEnd, cData.colType, cData.id)

		//var col Column
		var col interface{}

		switch cData.colType {
		case schema.FieldType_TIMESTAMP:
			col = NewIntColumn(s.f.Bytes, cData.offsetStart, cData.offsetEnd)
		case schema.FieldType_INT:
			col = NewIntColumn(s.f.Bytes, cData.offsetStart, cData.offsetEnd)
		case schema.FieldType_UINT:
		case schema.FieldType_FLOAT:
		case schema.FieldType_STRING:
		case schema.FieldType_TEXT:
		default:
			panic("unknown column type")
		}

		s.columns[cData.id] = col

	}

}
