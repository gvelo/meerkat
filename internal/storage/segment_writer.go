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

//go:generate protoc -I . -I ../../build/proto/ --plugin ../../build/protoc-gen-gogofaster --gogofaster_out=plugins=grpc,paths=source_relative:.  ./storage.proto

import (
	"github.com/google/uuid"
	"meerkat/internal/buffer"
	"meerkat/internal/storage/io"
	"sort"
)

const (
	MagicNumber    = "MEERKAT"
	SegmentVersion = 1
	TSColID        = "_ts" // TODO(gvelo) change to []byte
)

func NewSegmentWriter(path string, src SegmentSource) *SegmentWriter {
	return &SegmentWriter{
		path:    path,
		table:   table,
		id:      id,
		offsets: make(map[string]int, len(table.Cols())),
	}
}

type SegmentWriter struct {
	path    string
	src     SegmentSource
	bw      *io.BinaryWriter
	offsets map[string]int
}

func (sw *SegmentWriter) Write() (err error) {

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

	//perm := sw.writeTSColumn()
	//
	//sw.writeColumns(perm)

	sw.writeFooter()

	return

}

func (sw *SegmentWriter) writeHeader() {

	sw.bw.WriteRaw([]byte(MagicNumber))
	sw.bw.WriteByte(byte(SegmentVersion))

}

func (sw *SegmentWriter) writeColumns() {
	for _, colInfo := range sw.src.Columns() {
		columnWriter := NewColumnWriter1(colInfo, sw.src, sw.bw)
		columnWriter.Write()
		sw.offsets[colInfo.Name] = sw.bw.Offset()
	}
}

//func (sw *SegmentWriter) writeTSColumn() []int {
//
//	c, ok := sw.table.Col(TSColID)
//
//	if !ok {
//		panic("missing TS column")
//	}
//
//	tsColumn, ok := c.(*buffer.IntBuffer)
//
//	if !ok {
//		panic("wrong TS column type")
//	}
//
//	perm := sortTSColumn(tsColumn.Values())
//
//	// set the date range
//	sw.fromDate = tsColumn.Values()[0]
//	sw.toDate = tsColumn.Values()[tsColumn.Len()-1]
//
//	cw := NewTSColumnWriter(tsColumn, sw.bw)
//
//	cw.Write()
//
//	sw.offsets[TSColID] = sw.bw.Offset()
//
//	return perm
//
//}

//func (sw *SegmentWriter) writeColumns(perm []int) {

	//for _, f := range sw.table.Index().Fields {
	//
	//	// skip the timestamp column.
	//	if f.Id == TSColID {
	//		continue
	//	}
	//
	//	b, ok := sw.table.Col(f.Id)
	//
	//	if !ok {
	//		panic("error getting buffer for column")
	//	}
	//
	//	w := NewColumWriter(f.FieldType, b, perm, sw.bw)
	//
	//	w.Write()
	//
	//	sw.offsets[f.Id] = sw.bw.Offset()
	//
	//}

//}

func (sw *SegmentWriter) writeFooter() {

	entry := sw.bw.Offset()

	sw.bw.WriteUvarint(int(sw.src.SegmentInfo().Len))

	sw.bw.WriteUvarint(len(sw.src.Columns()))

	for _, columnInfo := range sw.src.Columns() {

		sw.bw.WriteString(columnInfo.Name)

		sw.bw.WriteByte(byte(columnInfo.ColumnType)

		sw.bw.WriteUvarint(sw.offsets[columnInfo.Name])

	}

	sw.bw.WriteFixedInt(entry)

}

func sortTSColumn(values []int) []int {

	perm := make([]int, len(values))

	for i := 0; i < len(perm); i++ {
		perm[i] = i
	}

	tsSlice := &TSSlice{
		ts:   values,
		perm: perm,
	}

	sort.Stable(tsSlice)

	return perm

}

type TSSlice struct {
	ts   []int
	perm []int
}

func (t *TSSlice) Len() int {
	return len(t.ts)
}

func (t *TSSlice) Less(i, j int) bool {
	return t.ts[i] < t.ts[j]
}

func (t *TSSlice) Swap(i, j int) {
	t.ts[i], t.ts[j] = t.ts[j], t.ts[i]
	t.perm[i], t.perm[j] = t.perm[j], t.perm[i]
}

func WriteSegment(path string, id uuid.UUID, table *buffer.Table) error {
	w := NewSegmentWriter(path, id, table)
	return w.Write()
}
