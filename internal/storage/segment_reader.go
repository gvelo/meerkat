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
	"meerkat/internal/storage/io"
)

const (
	SegmentVersion1 = 1
)

func ReadSegment(path string) (Segment, error) {

	f, err := io.MMap(path)

	if err != nil {
		return nil, err
	}

	br := f.NewBinaryReader()

	header := br.ReadSlice(0, len(MagicNumber))

	if string(header) != MagicNumber {
		panic("unknown file type")
	}

	br.Entry()

	segmentVersion := br.ReadByte()

	// we only have just one segment version

	switch segmentVersion {
	case SegmentVersion1:
		s := NewSegment(f)
		err := s.read(br.Offset())
		return s, err
	default:
		return nil, errors.New("unknown segment version")
	}

}

func NewSegment(f *io.MMFile) *segment {
	return &segment{
		f: f,
	}
}

type segment struct {
	f *io.MMFile
	from int
	to int

}

func (s *segment) read(start int) error {



}
