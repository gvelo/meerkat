// Copyright 2019 The Meerkat Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this path except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package writers

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment"
	"meerkat/internal/storage/segment/inmem"
	"path/filepath"
)

const (
	idxPosExt = ".ipos" // index for value in posting
	idxPagExt = ".ipag" // index
	pagExt    = ".pag"  // pages encoded
	posExt    = ".pos"  // posting lists for values
)

type SegmentWriter struct {
	log     zerolog.Logger
	path    string
	segment *inmem.Segment
}

func NewSegmentWriter(path string, segment *inmem.Segment) *SegmentWriter {

	log := log.With().
		Str("component", "SegmentWriter").
		Str("segmentID", segment.ID).
		Str("indexName", segment.IndexInfo.Name).
		Logger()

	if segment.State != inmem.InMem {
		log.Panic().Msgf("invalid segment state [%s]", segment.State)
	}

	sw := &SegmentWriter{
		segment: segment,
		path:    path,
		log:     log,
	}

	return sw

}

func (sw *SegmentWriter) Write() error {

	sw.log.Info().Msg("Starting to write segment")

	// TS column must be be processed first because it could
	// be sorted and and in this case it will determine the order
	// of the rest of segment columns.

	tsColumn := sw.segment.Columns[0].(*inmem.ColumnTimeStamp)

	if !tsColumn.Sorted() {
		tsColumn.Sort()
	}

	for _, col := range sw.segment.Columns {
		col.SetSortMap(tsColumn.SortMap())
		sw.writeColumn(col)
	}

	err := sw.writeSegmentInfo()

	if err != nil {
		return err
	}

	return nil
}

func (sw *SegmentWriter) writeColumn(col inmem.Column) (pd []*inmem.Page, err error) {

	var chainSkip = []Middleware{
		BuildSkip,
		EncoderHandler,
	}

	var chainBTrie = []Middleware{
		BuildBTrie,
		EncoderHandler,
	}

	var chain []Middleware
	var mp *MiddlewarePayload

	switch col.FieldInfo().Type {
	case segment.FieldTypeInt, segment.FieldTypeTimestamp:
		chain = chainSkip
		mp = NewMiddlewarePayload(sw.path, col)
	case segment.FieldTypeKeyword, segment.FieldTypeText:
		chain = chainBTrie
		mp = NewMiddlewarePayload(sw.path, col)
	case segment.FieldTypeFloat:
		chain = chainSkip
		mp = NewMiddlewarePayload(sw.path, col)
	default:
		sw.log.Panic().Msgf("invalid column type [%v]", col.FieldInfo().Type)
	}

	executeChain := BuildChain(WriteToFile, chain...)
	err = executeChain(mp)
	if err != nil {
		pd = mp.Pages
	}

	return nil, nil
}

func (sw *SegmentWriter) writeSegmentInfo() error {

	file := filepath.Join(sw.path, "info")

	bw, err := io.NewBinaryWriter(file)

	if err != nil {
		return err
	}

	// Header
	err = bw.WriteHeader(io.SegmentInfo)

	if err != nil {
		return err
	}

	// Index name
	err = bw.WriteString(sw.segment.IndexInfo.Name)

	if err != nil {
		return err
	}

	// Field Count
	err = bw.WriteVarInt(len(sw.segment.IndexInfo.Fields))

	if err != nil {
		return err
	}

	// Field info
	for _, field := range sw.segment.IndexInfo.Fields {

		err = bw.WriteString(field.Name)
		if err != nil {
			return err
		}

		err = bw.WriteVarInt(int(field.Type))
		if err != nil {
			return err
		}

		i := byte(0)
		if field.Index {
			i = 1
		}

		err = bw.WriteByte(i)
		if err != nil {
			return err
		}

	}

	// segment stats.

	// Event count
	err = bw.WriteVarInt(int(sw.segment.EventCount))

	if err != nil {
		return err
	}

	err = bw.Close()

	if err != nil {
		return err
	}

	return nil

	// TODO add fields cardinality, max/min TS and SegmentID

}

func WriteSegment(path string, segment *inmem.Segment) error {
	sw := NewSegmentWriter(path, segment)
	return sw.Write()
}
