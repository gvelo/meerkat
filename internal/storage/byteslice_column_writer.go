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
	colIndex ByteSliceIndexWriter,
	blockIndex BlockIndexWriter,
	validityIndex ValidityIndexWriter,
	bw *io.BinaryWriter) *ByteSliceColumnWriter {

	return &ByteSliceColumnWriter{
		fieldType:  fieldType,
		src:        src,
		bw:         bw,
		encoder:    encoder,
		colIndex:   colIndex,
		blockIndex: blockIndex,
		validity:   validityIndex,
	}

}

type ByteSliceColumnWriter struct {
	fieldType   schema.FieldType
	bw          *io.BinaryWriter
	src         ByteSliceColumSource
	encoder     ByteSliceEncoder
	colIndex    ByteSliceIndexWriter
	blockIndex  BlockIndexWriter
	validity    ValidityIndexWriter
	numOfValues int
	cardinality int

	blkEnd         int
	blkIdxEnd      int
	encoderEnd     int
	colIdxEnd      int
	validityIdxEnd int
}

func (w *ByteSliceColumnWriter) Write() {

	for w.src.HasNext() {

		vec := w.src.Next()

		w.numOfValues += vec.Len()

		w.encoder.Encode(vec)

		if w.colIndex != nil {
			w.colIndex.Index(vec)
		}

		if w.src.HasNulls() {
			w.validity.Index(vec.Rid())
		}

	}

	w.encoder.FlushBlocks()

	w.blkEnd = w.bw.Offset()

	w.blockIndex.Flush()

	w.blkIdxEnd = w.bw.Offset()

	w.encoder.Flush()

	w.encoderEnd = w.bw.Offset()

	// TODO(gvelo) if the column is not indexed estimate
	// cardinality anyways using datasketches.
	if w.colIndex != nil {

		w.colIndex.Flush()

		w.colIdxEnd = w.bw.Offset()

		w.cardinality = w.colIndex.Cardinality()

	}

	if w.src.HasNulls() {

		w.validity.Flush()

		w.validityIdxEnd = w.bw.Offset()

	}

	w.writeFooter()

}

func (w *ByteSliceColumnWriter) writeFooter() {

	entry := w.bw.Offset()

	w.bw.WriteUvarint(int(w.fieldType))
	w.bw.WriteUvarint(int(w.encoder.Type()))
	w.bw.WriteUvarint(w.blkEnd)
	w.bw.WriteUvarint(w.blkIdxEnd)
	w.bw.WriteUvarint(w.encoderEnd)
	w.bw.WriteUvarint(w.colIdxEnd)
	w.bw.WriteUvarint(w.validityIdxEnd)
	w.bw.WriteUvarint(w.numOfValues)
	w.bw.WriteUvarint(w.cardinality)

	w.bw.WriteFixedInt(entry)

}

func NewUUIDColumnWriter(fieldType schema.FieldType,
	src UUIDColumSource,
	encoder ByteSliceEncoder,
	colIndex ByteSliceIndexWriter,
	blockIndex BlockIndexWriter,
	validityIndex ValidityIndexWriter,
	bw *io.BinaryWriter) *UUIDColumnWriter {

	return &UUIDColumnWriter{
		fieldType:  fieldType,
		src:        src,
		bw:         bw,
		encoder:    encoder,
		colIndex:   colIndex,
		blockIndex: blockIndex,
		validity:   validityIndex,
	}

}

type UUIDColumnWriter struct {
	fieldType   schema.FieldType
	bw          *io.BinaryWriter
	src         UUIDColumSource
	encoder     ByteSliceEncoder
	colIndex    ByteSliceIndexWriter
	blockIndex  BlockIndexWriter
	validity    ValidityIndexWriter
	numOfValues int
	cardinality int

	blkEnd         int
	blkIdxEnd      int
	encoderEnd     int
	colIdxEnd      int
	validityIdxEnd int
}

func (w *UUIDColumnWriter) Write() {

	for w.src.HasNext() {

		vec := w.src.Next()

		w.numOfValues += vec.Len()

		w.encoder.Encode(vec)

		if w.colIndex != nil {
			w.colIndex.Index(vec)
		}

		if w.src.HasNulls() {
			w.validity.Index(vec.Rid())
		}

	}

	w.encoder.FlushBlocks()

	w.blkEnd = w.bw.Offset()

	w.blockIndex.Flush()

	w.blkIdxEnd = w.bw.Offset()

	w.encoder.Flush()

	w.encoderEnd = w.bw.Offset()

	// TODO(gvelo) if the column is not indexed estimate
	// cardinality anyways using datasketches.
	if w.colIndex != nil {

		w.colIndex.Flush()

		w.colIdxEnd = w.bw.Offset()

		w.cardinality = w.colIndex.Cardinality()

	}

	if w.src.HasNulls() {

		w.validity.Flush()

		w.validityIdxEnd = w.bw.Offset()

	}

	w.writeFooter()

}

func (w *UUIDColumnWriter) writeFooter() {

	entry := w.bw.Offset()

	w.bw.WriteUvarint(int(w.fieldType))
	w.bw.WriteUvarint(int(w.encoder.Type()))
	w.bw.WriteUvarint(w.blkEnd)
	w.bw.WriteUvarint(w.blkIdxEnd)
	w.bw.WriteUvarint(w.encoderEnd)
	w.bw.WriteUvarint(w.colIdxEnd)
	w.bw.WriteUvarint(w.validityIdxEnd)
	w.bw.WriteUvarint(w.numOfValues)
	w.bw.WriteUvarint(w.cardinality)

	w.bw.WriteFixedInt(entry)

}
