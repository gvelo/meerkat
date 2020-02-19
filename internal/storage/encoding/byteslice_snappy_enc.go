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
	"github.com/golang/snappy"
	"meerkat/internal/storage"
	"meerkat/internal/storage/io"
)

type ByteSliceSnappyEncoder struct {
	pw        storage.PageWriter
	buf       *io.EncoderBuffer
	offsetBuf []int
}

func NewByteSliceSnappyEncodeer(pw storage.PageWriter) *ByteSliceSnappyEncoder {
	return &ByteSliceSnappyEncoder{
		pw:        pw,
		buf:       io.NewEncoderBuffer(64 * 1024),
		offsetBuf: make([]int, maxSlicesPerPage),
	}
}

func (e *ByteSliceSnappyEncoder) Flush() error {
	return nil
}

func (e *ByteSliceSnappyEncoder) FlushPages() error {
	// TODO(gvelo): check if we need a flusher interface
	//  on the page writer.
	return e.pw.Flush()
}

func (e *ByteSliceSnappyEncoder) Type() storage.EncodingType {
	return storage.Snappy
}

func (e *ByteSliceSnappyEncoder) Encode(vec storage.ByteSliceVector) error {

	// make sure that the buffer has enough space to accommodate
	// the offsets slice plus the encoded data. We need to avoid
	// allocation inside the snappy encoder.
	size := binary.MaxVarintLen64*(vec.Len()+2) + snappy.MaxEncodedLen(len(vec.Data()))

	e.buf.Reset(size)

	DeltaEncode(vec.Offsets(), e.offsetBuf)

	// left enough room at the beginning of the page to write the
	// page length encoded as uvarint
	e.buf.WriteVarUintSliceAt(binary.MaxVarintLen64, e.offsetBuf[:vec.Len()])

	// as we have enough room in the dst buffer the encoder will not
	// allocate a new slice. See MaxEncodedLen(srcLen int) int
	r := snappy.Encode(e.buf.Free(), vec.Data())

	e.buf.SetLen(e.buf.Len() + len(r))

	pageSize := e.buf.Len() - binary.MaxVarintLen64

	offset := binary.MaxVarintLen64 - io.SizeUVarint(uint64(pageSize))

	e.buf.WriteUvarintAt(offset, pageSize)

	page := e.buf.Bytes()[offset:]

	return e.pw.WritePage(page, vec.Rid()[len(vec.Rid())-1])

}
