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

package encoding

import (
	"encoding/binary"
	"meerkat/internal/storage"
	"meerkat/internal/storage/io"
)

const (
	maxSlicesPerPage = 1024 * 2
)

type ByteSlicePlainEncoder struct {
	pw        storage.PageWriter
	buf       *io.EncoderBuffer
	offsetBuf []int
}

func NewByteSlicePlainEncodeer(pw storage.PageWriter) *ByteSlicePlainEncoder {
	return &ByteSlicePlainEncoder{
		pw: pw,
		//TODO(gvelo) use max page size for the enc.
		buf:       io.NewEncoderBuffer(32 * 1024),
		offsetBuf: make([]int, maxSlicesPerPage),
	}
}

func (e *ByteSlicePlainEncoder) Flush() error {
	return nil
}

func (e *ByteSlicePlainEncoder) FlushPages() error {
	return e.pw.Flush()
}

func (e *ByteSlicePlainEncoder) Type() storage.EncodingType {
	return storage.Plain
}

func (e *ByteSlicePlainEncoder) Encode(vec storage.ByteSliceVector) error {

	size := binary.MaxVarintLen64*(vec.Len()+2) + len(vec.Data())

	e.buf.Reset(size)

	DeltaEncode(vec.Offsets(), e.offsetBuf)

	// left enough room at the beginning of the page to write the
	// page length encoded as uvarint
	e.buf.WriteVarUintSliceAt(binary.MaxVarintLen64, e.offsetBuf[:vec.Len()])

	// write the vector data
	e.buf.WriteBytes(vec.Data())

	pageSize := e.buf.Len() - binary.MaxVarintLen64

	offset := binary.MaxVarintLen64 - io.SizeUVarint(uint64(pageSize))

	e.buf.WriteUvarintAt(offset, pageSize)

	page := e.buf.Bytes()[offset:]

	return e.pw.WritePage(page, vec.Rid()[len(vec.Rid())-1])

}
