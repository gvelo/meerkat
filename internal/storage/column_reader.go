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
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/storage/encoding"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/vector"
)

type colBound struct {
	start int
	end   int
}

type IntColumnIterator struct {
	dec      encoding.IntDecoder
	br       encoding.BlockReader
	pool     vector.Pool
	blockLen int
}

func NewIntColumnIterator(
	dec encoding.IntDecoder,
	br encoding.BlockReader,
	pool vector.Pool,
	blockLen int,
) *IntColumnIterator {

	return &IntColumnIterator{
		dec:      dec,
		br:       br,
		pool:     pool,
		blockLen: blockLen,
	}

}

func (i *IntColumnIterator) Next() vector.IntVector {

	v := i.pool.GetIntVector()

	for l := 0; i.br.HasNext() && v.RemainingLen() > 0; l++ {

		b := i.br.Next()

		i.dec.Decode(b.Bytes(), v.Remaining())

		l += i.blockLen

		v.SetLen(l)

	}

	return v

}

func (i *IntColumnIterator) HasNext() bool {
	return i.br.HasNext()
}

type IntNullColumnIterator struct {
	decoder   encoding.IntDecoder
	br        encoding.BlockReader
	validity  *roaring.Bitmap
	pool      vector.Pool
	blockSize int
	buf       []int
	valid     []uint32
	valIter   roaring.ManyIntIterable
	pos       int
	rid       uint32
}

func NewIntNullColumnIterator(
	decoder encoding.IntDecoder,
	br encoding.BlockReader,
	validity *roaring.Bitmap,
	pool vector.Pool,
	blockLen int,
) *IntNullColumnIterator {

	return &IntNullColumnIterator{
		decoder:   decoder,
		br:        br,
		validity:  validity,
		pool:      pool,
		blockSize: blockLen,
	}

}

func (i *IntNullColumnIterator) Next() vector.IntVector {

	v := i.pool.GetIntVector()
	i.readBlock()

	vbuf := v.Buf()

	r := 0

	for r < v.Cap() {

		if i.pos == i.blockSize {

			if !i.br.HasNext() {
				break
			}

			i.readBlock()
			i.pos = 0

		}

		if i.rid == i.valid[i.pos] {
			vbuf[r] = i.buf[i.pos]
			v.SetValid(r)
			i.pos++
		} else {
			v.SetInvalid(r)
		}

		i.rid++
		r++

	}

	v.SetLen(r)

	return v
}

func (i *IntNullColumnIterator) readBlock() {
	b := i.br.Next()
	i.decoder.Decode(b.Bytes(), i.buf)
	i.valIter.NextMany(i.valid)
}

func (i *IntNullColumnIterator) HasNext() bool {

	return i.br.HasNext() || i.pos != 0
}

type intColumn struct {
	start            int
	end              int
	b                []byte
	valid            *roaring.Bitmap
	encoding         encoding.EncodingType
	blkStart         int
	blkEnd           int
	blkIdxStart      int
	blkIdxEnd        int
	encoderStart     int
	encoderEnd       int
	colIdxStart      int
	colIdxEnd        int
	validityIdxStart int
	validityIdxEnd   int
	numOfValues      int
	cardinality      int
	blockLen         int
	vectorPool       vector.Pool
}

func (cr *intColumn) debug() {

	fmt.Println("start ", cr.start)
	fmt.Println("end ", cr.end)
	fmt.Println("encoding ", cr.encoding)
	fmt.Println("blkStart ", cr.blkStart)
	fmt.Println("blkEnd ", cr.blkEnd)
	fmt.Println("blkIdxStart ", cr.blkIdxStart)
	fmt.Println("blkIdxEnd ", cr.blkIdxEnd)
	fmt.Println("encoderStart ", cr.encoderStart)
	fmt.Println("encoderEnd ", cr.encoderEnd)
	fmt.Println("colIdxStart ", cr.colIdxStart)
	fmt.Println("colIdxEnd ", cr.colIdxEnd)
	fmt.Println("validityIdxStart ", cr.validityIdxStart)
	fmt.Println("validityIdxEnd ", cr.validityIdxEnd)
	fmt.Println("numOfValues ", cr.numOfValues)
	fmt.Println("cardinality ", cr.cardinality)
	fmt.Println("blockLen ", cr.blockLen)

}

func NewIntColumn(b []byte, start int, end int) *intColumn {

	cr := &intColumn{
		start: start,
		end:   end,
		b:     b,
	}

	br := io.NewBinaryReader(b)

	br.SetOffset(cr.end - 8)
	br.SetOffset(br.ReadFixed64())

	e := 0

	cr.encoding = encoding.EncodingType(br.ReadUVarint())
	cr.blkStart = start
	cr.blkEnd = br.ReadUVarint()
	e = cr.blkEnd

	cr.blkIdxEnd = br.ReadUVarint()

	if cr.blkIdxEnd != 0 {
		cr.blkIdxStart = e
		e = cr.blkIdxEnd
	}

	cr.encoderStart = e
	cr.encoderEnd = br.ReadUVarint()
	e = cr.encoderEnd

	cr.colIdxEnd = br.ReadUVarint()

	if cr.colIdxEnd != 0 {
		cr.colIdxStart = e
		e = cr.colIdxEnd
	}

	cr.validityIdxEnd = br.ReadUVarint()

	if cr.validityIdxEnd != 0 {
		cr.validityIdxStart = e
		e = cr.validityIdxEnd
	}

	cr.numOfValues = br.ReadUVarint()
	cr.cardinality = br.ReadUVarint()

	// TODO(gvelo) make column specific.
	cr.blockLen = blockLen
	cr.read()

	return cr

}

func (cr *intColumn) read() {

	cr.debug()

	if cr.validityIdxEnd != 0 {
		cr.valid = roaring.NewBitmap()
		_, err := cr.valid.FromBuffer(cr.b[cr.validityIdxStart:cr.validityIdxEnd])
		if err != nil {
			panic(err)
		}
	}

}

func (cr *intColumn) Iterator() IntIterator {

	dec, blockReader := encoding.GetIntDecoder(
		cr.encoding,
		cr.b,
		cr.encoderStart,
		cr.encoderEnd,
		cr.blockLen,
	)

	if cr.valid != nil {
		return NewIntNullColumnIterator(
			dec,
			blockReader,
			cr.valid,
			cr.vectorPool,
			cr.blockLen,
		)
	}

	return NewIntColumnIterator(
		dec,
		blockReader,
		cr.vectorPool,
		cr.blockLen,
	)

}
