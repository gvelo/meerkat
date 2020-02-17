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
	"meerkat/internal/schema"
	"meerkat/internal/storage/io"
)

func NewByteSliceColumnWriter(fieldType schema.FieldType,
	src ByteSliceColumSource,
	encoder ByteSliceEncoder,
	index ByteSliceIndexWriter,
	pageIndex PageIndexWriter,
	validityIndex ValidityIndexWriter,
	bw *io.BinaryWriter) *ByteSliceColumnWriter {

	return &ByteSliceColumnWriter{
		fieldType: fieldType,
		src:       src,
		bw:        bw,
		encoder:   encoder,
		index:     index,
		pageIndex: pageIndex,
		validity:  validityIndex,
	}

}

type ByteSliceColumnWriter struct {
	fieldType      schema.FieldType
	bw             *io.BinaryWriter
	src            ByteSliceColumSource
	encoder        ByteSliceEncoder
	index          ByteSliceIndexWriter
	pageIndex      PageIndexWriter
	validity       ValidityIndexWriter
	numOfValues    int
	cardinality    int
	startOffset    int
	pageOffset     int
	pageIdxOffset  int
	encOffset      int
	indexOffset    int
	validityOffset int
}

func (w *ByteSliceColumnWriter) Write() error {

	w.startOffset = w.bw.Offset

	for w.src.HasNext() {

		vec := w.src.Next()

		w.numOfValues = w.numOfValues + vec.Len()

		err := w.encoder.Encode(vec)

		if err != nil {
			return err
		}

		if w.index != nil {
			w.index.Index(vec)
		}

		if w.src.HasNulls() {
			w.validity.Index(vec.Rid())
		}

	}

	err := w.encoder.FlushPages()

	if err != nil {
		return err
	}

	w.pageOffset = w.bw.Offset

	err = w.encoder.Flush()

	if err != nil {
		return err
	}

	w.encOffset = w.bw.Offset

	err = w.pageIndex.Flush()

	if err != nil {
		return err
	}

	w.pageIdxOffset = w.bw.Offset

	// TODO(gvelo) if the column is not indexed estimate
	// cardinality anyways using datasketches.
	if w.index != nil {

		err = w.index.Flush()

		if err != nil {
			return err
		}

		w.indexOffset = w.bw.Offset

		w.cardinality = w.index.Cardinality()

	}

	if w.src.HasNulls() {

		err = w.validity.Flush()

		if err != nil {
			return err
		}

		w.validityOffset = w.bw.Offset

	}

	err = w.WriteMetadata()

	if err != nil {
		return err
	}

	return nil

}

func (w *ByteSliceColumnWriter) WriteMetadata() error {

	metadata := []int{
		int(w.fieldType),
		int(w.encoder.Type()),
		w.startOffset,
		w.pageOffset,
		w.encOffset,
		w.pageIdxOffset,
		w.indexOffset,
		w.validityOffset,
		w.numOfValues,
		w.cardinality,
	}

	metadataStart := w.bw.Offset

	err := w.bw.WriteVarIntSlice(metadata)

	if err != nil {
		return err
	}

	err = w.bw.WriteFixedInt(metadataStart)

	if err != nil {
		return err
	}

	return nil

}
