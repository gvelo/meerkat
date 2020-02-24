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

func NewIntColumnWriter(fieldType schema.FieldType,
	src IntColumSource,
	encoder IntEncoder,
	colIndex IntIndexWriter,
	blockIndex BlockIndexWriter,
	validityIndex ValidityIndexWriter,
	bw *io.BinaryWriter) *IntColumnWriter {

	return &IntColumnWriter{
		fieldType:  fieldType,
		src:        src,
		bw:         bw,
		encoder:    encoder,
		colIndex:   colIndex,
		blockIndex: blockIndex,
		validity:   validityIndex,
	}

}

type IntColumnWriter struct {
	fieldType      schema.FieldType
	bw             *io.BinaryWriter
	src            IntColumSource
	encoder        IntEncoder
	colIndex       IntIndexWriter
	blockIndex     BlockIndexWriter
	validity       ValidityIndexWriter
	numOfValues    int
	cardinality    int
	startOffset    int
	blocksOffset   int
	blockIdxOffset int
	encOffset      int
	colIndexOffset int
	validityOffset int
}

func (w *IntColumnWriter) Write() error {

	w.startOffset = w.bw.Offset

	for w.src.HasNext() {

		vec := w.src.Next()

		w.numOfValues = w.numOfValues + vec.Len()

		err := w.encoder.Encode(vec)

		if err != nil {
			return err
		}

		if w.colIndex != nil {
			w.colIndex.Index(vec)
		}

		if w.src.HasNulls() {
			w.validity.Index(vec.Rid())
		}

	}

	err := w.encoder.FlushBlocks()

	if err != nil {
		return err
	}

	w.blocksOffset = w.bw.Offset

	err = w.encoder.Flush()

	if err != nil {
		return err
	}

	w.encOffset = w.bw.Offset

	err = w.blockIndex.Flush()

	if err != nil {
		return err
	}

	w.blockIdxOffset = w.bw.Offset

	// TODO(gvelo) if the column is not indexed estimate
	// cardinality anyways using datasketches.
	if w.colIndex != nil {

		err = w.colIndex.Flush()

		if err != nil {
			return err
		}

		w.colIndexOffset = w.bw.Offset

		w.cardinality = w.colIndex.Cardinality()

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

func (w *IntColumnWriter) WriteMetadata() error {

	metadata := []int{
		int(w.fieldType),
		int(w.encoder.Type()),
		w.startOffset,
		w.blocksOffset,
		w.encOffset,
		w.blockIdxOffset,
		w.colIndexOffset,
		w.validityOffset,
		w.numOfValues,
		w.cardinality,
	}

	metadataStart := w.bw.Offset

	err := w.bw.WriteUVarIntSlice(metadata)

	if err != nil {
		return err
	}

	err = w.bw.WriteFixedInt(metadataStart)

	if err != nil {
		return err
	}

	return nil

}
