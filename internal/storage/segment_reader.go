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
	"fmt"
	"github.com/google/uuid"
	"meerkat/internal/storage/io"
)

const (
	SegmentVersion1 = 1
)

type colData struct {
	colType ColumnType
	name    string
	bounds  io.Bounds
}

func ReadSegment(path string) *Segment {

	f, err := io.MMap(path)

	if err != nil {
		panic(err)
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
		s.read()
		return s
	default:
		panic("unknown segment version")
	}

}

func NewSegment(f *io.MMFile) *Segment {
	return &Segment{
		f:       f,
		columns: make(map[string]Column),
	}
}

// TODO(gvelo) extract interface for segment.
type Segment struct {
	f         *io.MMFile
	id        uuid.UUID
	from      int
	to        int
	numOfRows int
	start     int
	tableId   string
	tableName string
	numOfCol  int
	// columns   map[string]Column
	columns     map[string]Column
	segmentInfo *SegmentInfo
}

func (s *Segment) read() {

	// magicNumber + version
	s.start = len(MagicNumber) + 1

	br := s.f.NewBinaryReader()

	br.Entry()

	s.numOfRows = br.ReadUVarint()

	s.numOfCol = br.ReadUVarint()

	cd := make([]colData, s.numOfCol)

	for i := 0; i < s.numOfCol; i++ {
		c := colData{}
		c.name = br.ReadString()
		c.colType = ColumnType(br.ReadByte())
		c.bounds.End = br.ReadUVarint()
		if i == 0 {
			c.bounds.Start = s.start
		} else {
			c.bounds.Start = cd[i-1].bounds.End
		}
		cd[i] = c
	}

	s.readColumns(cd)

}

func (s *Segment) readColumns(cd []colData) {

	for _, cData := range cd {

		var col Column

		switch cData.colType {
		case ColumnType_TIMESTAMP:
			col = NewInt64Column(s.f.Bytes, cData.bounds, s.numOfRows)
		case ColumnType_INT64:
			col = NewInt64Column(s.f.Bytes, cData.bounds, s.numOfRows)
		case ColumnType_FLOAT64:
			col = NewFloat64Column(s.f.Bytes, cData.bounds, s.numOfRows)
		case ColumnType_STRING:
			col = NewBinaryColumn(s.f.Bytes, cData.bounds, s.numOfRows)
		default:
			panic(fmt.Sprintf("unknown column type %v", cData.colType))
		}

		s.columns[cData.name] = col

	}

}

func (s *Segment) Info() *SegmentInfo {
	return s.segmentInfo
}

func (s *Segment) Column(name string) Column {
	return s.columns[name]
}

func (s *Segment) Close() {
	err := s.f.UnMap()
	if err != nil {
		panic(err)
	}
}
