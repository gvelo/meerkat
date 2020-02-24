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
	"meerkat/internal/buffer"
	"meerkat/internal/storage/io"
	"sort"
)

const (
	MagicNumber    = "MEERKAT"
	SegmentVersion = 1
	TSColID        = "_ts" // TODO(gvelo) change to []byte
)

func NewSegmentWriter(path string, id uuid.UUID, table *buffer.Table) *SegmentWriter {
	return &SegmentWriter{
		path:    path,
		table:   table,
		id:      id,
		offsets: make(map[string]int, len(table.Cols())),
	}
}

type SegmentWriter struct {
	path     string
	table    *buffer.Table
	bw       *io.BinaryWriter
	id       uuid.UUID
	offsets  map[string]int
	fromDate int
	toDate   int
}

func (sw *SegmentWriter) Write() error {

	var err error

	sw.bw, err = io.NewBinaryWriter(sw.path)

	if err != nil {
		return err
	}

	defer sw.bw.Close()

	err = sw.writeHeader()

	if err != nil {
		return err
	}

	perm, err := sw.writeTSColumn()

	if err != nil {
		return err
	}

	err = sw.writeColumns(perm)

	if err != nil {
		return err
	}

	err = sw.writeMetadata()

	if err != nil {
		return err
	}

	return nil

}

func (sw *SegmentWriter) writeHeader() error {

	_, err := sw.bw.Write([]byte(MagicNumber))

	if err != nil {
		return err
	}

	return nil

}

func (sw *SegmentWriter) writeTSColumn() ([]int, error) {

	c, ok := sw.table.Col(TSColID)

	if !ok {
		return nil, errors.New("missing TS column")
	}

	tsColumn, ok := c.(*buffer.IntBuffer)

	if !ok {
		return nil, errors.New("wrong TS column type")
	}

	perm := sortTSColumn(tsColumn.Values())

	// set the date range
	sw.fromDate = tsColumn.Values()[0]
	sw.toDate = tsColumn.Values()[tsColumn.Len()]

	cw := NewTSColumnWriter(tsColumn, perm, sw.bw)

	err := cw.Write()

	if err != nil {
		return nil, err
	}

	sw.offsets[TSColID] = sw.bw.Offset

	return perm, nil

}

func (sw *SegmentWriter) writeColumns(perm []int) error {

	for _, f := range sw.table.Index().Fields {

		// skip the timestamp column.
		if f.Id == TSColID {
			continue
		}

		b, ok := sw.table.Col(f.Id)

		if !ok {
			panic("error getting buffer for column")
		}

		w := NewColumWriter(f.FieldType, b, perm, sw.bw)

		err := w.Write()

		if err != nil {
			return err
		}

		sw.offsets[f.Id] = sw.bw.Offset

	}

	return nil

}

func (sw *SegmentWriter) writeMetadata() error {

	metadataStart := sw.bw.Offset

	err := sw.bw.WriteByte(byte(SegmentVersion))

	if err != nil {
		return err
	}

	_, err = sw.bw.Write(sw.id[:])

	if err != nil {
		return err
	}

	// TODO(gvelo) refactor to [16]byte
	err = sw.bw.WriteString(sw.table.Index().Id)

	if err != nil {
		return err
	}

	err = sw.bw.WriteString(sw.table.Index().Name)

	if err != nil {
		return err
	}

	err = sw.bw.WriteFixedInt(sw.fromDate)

	if err != nil {
		return err
	}

	err = sw.bw.WriteFixedInt(sw.toDate)

	if err != nil {
		return err
	}

	err = sw.bw.WriteUvarint(sw.table.Len())

	if err != nil {
		return err
	}

	err = sw.bw.WriteFixedInt(len(sw.table.Cols()))

	if err != nil {
		return err
	}

	for _, f := range sw.table.Index().Fields {

		err = sw.bw.WriteString(f.Id)

		if err != nil {
			return err
		}

		err = sw.bw.WriteString(f.Name)

		if err != nil {
			return err
		}

		err = sw.bw.WriteUvarint(sw.offsets[f.Id])

		if err != nil {
			return err
		}

	}

	err = sw.bw.WriteFixedInt(metadataStart)

	if err != nil {
		return err
	}

	return nil

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
