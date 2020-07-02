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
	"meerkat/internal/storage/io"
)

const (
	MagicNumber    = "MEERKAT"
	SegmentVersion = 1
)

func newSegmentWriter(path string, src SegmentSource) *segmentWriter {
	return &segmentWriter{
		path:    path,
		src:     src,
		offsets: make(map[string]int),
	}
}

type segmentWriter struct {
	path    string
	src     SegmentSource
	bw      *io.BinaryWriter
	offsets map[string]int
}

func (sw *segmentWriter) Write() (err error) {

	defer func() {
		if r := recover(); r != nil {
			panic(r)
			//e, ok := r.(error)
			//if ok {
			//	err = e
			//	return
			//} else {
			//	panic(r)
			//}
		}
	}()

	sw.bw, err = io.NewBinaryWriter(sw.path)

	if err != nil {
		return
	}

	defer sw.bw.Close()

	sw.writeHeader()

	sw.writeColumns()

	sw.writeFooter()

	return

}

func (sw *segmentWriter) writeHeader() {

	sw.bw.WriteRaw([]byte(MagicNumber))
	sw.bw.WriteByte(byte(SegmentVersion))

}

func (sw *segmentWriter) writeColumns() {

	for _, colInfo := range sw.src.Info().Columns {
		columnWriter := NewColumnWriter(colInfo, sw.src, sw.bw)
		columnWriter.Write()
		sw.offsets[colInfo.Name] = sw.bw.Offset()
	}

}

func (sw *segmentWriter) writeFooter() {

	entry := sw.bw.Offset()

	sw.bw.WriteUvarint(int(sw.src.Info().Len))

	sw.bw.WriteUvarint(len(sw.src.Info().Columns))

	for _, columnInfo := range sw.src.Info().Columns {

		sw.bw.WriteString(columnInfo.Name)

		sw.bw.WriteByte(byte(columnInfo.ColumnType))

		sw.bw.WriteUvarint(sw.offsets[columnInfo.Name])

	}

	sw.bw.WriteFixedInt(entry)

}

func WriteSegment(path string, src SegmentSource) error {
	w := newSegmentWriter(path, src)
	return w.Write()
}
