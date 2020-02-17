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

package encoders

import (
	"encoding/binary"
	"meerkat/internal/storage"
	"meerkat/internal/storage/io"
)

type ByteSlicePlain struct {
	pw  storage.PageWriter
	buf *io.EncoderBuffer
}

func NewByteSlicePlain(pw storage.PageWriter) *ByteSlicePlain {
	return &ByteSlicePlain{
		pw:  pw,
		buf: io.NewEncoderBuffer(32 * 1024),
	}
}

func (e *ByteSlicePlain) Flush() error {
	return nil
}

func (e *ByteSlicePlain) FlushPages() error {
	return e.pw.Flush()
}

func (e *ByteSlicePlain) Type() storage.EncodingType {
	return storage.Plain
}

func (e *ByteSlicePlain) Encode(vec storage.ByteSliceVector) error {

	size := binary.MaxVarintLen64*(vec.Len()+2) + len(vec.Data())

	e.buf.Reset(size)

	// left enough room at the beginning of the page to write the
	// page length encoded as uvarint
	e.buf.WriteVarUintSliceAt(binary.MaxVarintLen64, vec.Offsets())

	// write the vector data
	e.buf.WriteBytes(vec.Data())

	pageSize := e.buf.Len() - binary.MaxVarintLen64

	offset := binary.MaxVarintLen64 - io.SizeUVarint(uint64(pageSize))

	e.buf.WriteUvarintAt(offset, pageSize)

	page := e.buf.Bytes()[offset:]

	return e.pw.WritePage(page, vec.Rid()[len(vec.Rid())])

}
