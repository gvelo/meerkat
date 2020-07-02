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
	"github.com/google/uuid"
	"meerkat/internal/schema"
	"meerkat/internal/storage/io"
	"time"
)

const (
	SegmentVersion1 = 1
)

type colData struct {
	colType schema.ColumnType
	id      string
	name    string
	bounds  io.Bounds
}

func ReadSegment(path string) (*Segment, error) {

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

func NewSegment(f *io.MMFile) *Segment {
	return &Segment{
		f:       f,
		columns: make(map[string]interface{}),
	}
}

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
	//columns   map[string]Column
	columns map[string]interface{}
}

// start = entry
func (s *Segment) read() error {

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

	for i := 0; i < s.numOfCol; i++ {
		c := colData{}
		c.id = br.ReadString()
		c.name = br.ReadString()
		c.colType = schema.ColumnType(br.ReadByte())
		c.bounds.End = br.ReadUVarint()
		if i == 0 {
			c.bounds.Start = s.start
		} else {
			c.bounds.Start = cd[i-1].bounds.End
		}
		cd[i] = c
	}

	s.readColumns(cd)

	// TODO(gvelo) recover from panic and return the err.

	return nil

}

func (s *Segment) readColumns(cd []colData) {

	for _, cData := range cd {

		//var col Column
		var col interface{}

		switch cData.colType {
		case schema.ColumnType_TIMESTAMP:
			col = NewIntColumn(s.f.Bytes, cData.bounds, s.numOfRows)
		case schema.ColumnType_LONG:
			col = NewIntColumn(s.f.Bytes, cData.bounds, s.numOfRows)
		case schema.ColumnType_STRING:
			col = NewBinaryColumn(s.f.Bytes, cData.bounds, s.numOfRows)
		default:
			panic("unknown column type")
		}

		s.columns[cData.id] = col

	}

}

func (s *Segment) IndexName() string {
	panic("implement me")
}

func (s *Segment) IndexID() string {
	panic("implement me")
}

func (s *Segment) From() time.Time {
	return time.Unix(0, int64(s.from))
}

func (s *Segment) To() time.Time {
	return time.Unix(0, int64(s.to))
}

func (s *Segment) Rows() int {
	return s.numOfRows
}

func (s *Segment) Col(id string) interface{} {
	return s.columns[id]
}
