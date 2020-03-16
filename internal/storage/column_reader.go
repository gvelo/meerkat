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
	"meerkat/internal/storage/index"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/vector"
)

type IntColumnIterator struct {
	dec    encoding.IntDecoder
	br     encoding.BlockReader
	pool   vector.Pool
	colLen int
	rid    int
}

func NewIntColumnIterator(
	dec encoding.IntDecoder,
	br encoding.BlockReader,
	pool vector.Pool,
	colLen int,
) *IntColumnIterator {

	return &IntColumnIterator{
		dec:    dec,
		br:     br,
		pool:   pool,
		colLen: colLen,
	}

}

func (i *IntColumnIterator) Next() vector.IntVector {

	if i.rid >= i.colLen {
		panic("column EOF")
	}

	v := i.pool.GetIntVector()
	l := 0

	for i.rid < i.colLen && v.RemainingLen() > 0 {
		b := i.br.Next()
		i.dec.Decode(b.Bytes(), v.Remaining())
		i.rid += b.Len()
		l += b.Len()
		v.SetLen(l)
	}

	return v

}

func (i *IntColumnIterator) HasNext() bool {
	return i.rid < i.colLen
}

type IntNullColumnIterator struct {
	decoder  encoding.IntDecoder
	br       encoding.BlockReader
	validity *roaring.Bitmap
	pool     vector.Pool
	buf      []int
	bufLen   int
	valid    []uint32
	valIter  roaring.ManyIntIterable
	pos      int
	rid      uint32
	colLen   int
	eof      bool
}

func NewIntNullColumnIterator(
	decoder encoding.IntDecoder,
	br encoding.BlockReader,
	validity *roaring.Bitmap,
	colLen int,
	pool vector.Pool,
) *IntNullColumnIterator {

	return &IntNullColumnIterator{
		decoder:  decoder,
		br:       br,
		validity: validity,
		pool:     pool,
		buf:      make([]int, blockLen),
		valid:    make([]uint32, blockLen),
		valIter:  validity.ManyIterator(),
		colLen:   colLen,
	}

}

func (i *IntNullColumnIterator) Next() vector.IntVector {

	if int(i.rid) >= i.colLen {
		panic("column EOF")
	}

	v := i.pool.GetIntVector()
	vbuf := v.Buf()

	r := 0

	for r < v.Cap() && int(i.rid) < i.colLen {

		if (i.pos == i.bufLen || i.bufLen == 0) && !i.eof {

			i.pos = 0

			if i.br.HasNext() {
				i.readBlock()
			} else {
				i.eof = true
			}

		}

		if !i.eof && (i.rid == i.valid[i.pos]) {
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
	i.bufLen = b.Len()
	i.decoder.Decode(b.Bytes(), i.buf)
	vn := i.valIter.NextMany(i.valid)

	if vn != i.bufLen {
		panic("validity block length doesn't match value block length")
	}

}

func (i *IntNullColumnIterator) HasNext() bool {
	return int(i.rid) < i.colLen
}

type intColumn struct {
	b              []byte
	valid          *roaring.Bitmap
	blkIdx         index.BlockIndexReader
	encoding       encoding.EncodingType
	colBounds      io.Bounds
	blkBounds      io.Bounds
	blkIdxBounds   io.Bounds
	encBounds      io.Bounds
	colIdxBounds   io.Bounds
	validityBounds io.Bounds
	numOfValues    int
	numOfRows      int
	cardinality    int
	blockLen       int
	vectorPool     vector.Pool
}

func (cr *intColumn) String() string {
	return fmt.Sprintf(`
	encoding       %v
	colBounds      %v
	blkBounds      %v
	blkIdxBounds   %v
	encBounds      %v
	colIdxBounds   %v
	validityBounds %v
	numOfValues    %v
    numOfRows      %v
	cardinality    %v
	blockLen       %v
`, cr.encoding,
		cr.colBounds,
		cr.blkBounds,
		cr.blkIdxBounds,
		cr.encBounds,
		cr.colIdxBounds,
		cr.validityBounds,
		cr.numOfValues,
		cr.numOfRows,
		cr.cardinality,
		cr.blockLen)

}

func NewIntColumn(b []byte, bounds io.Bounds, numOfRows int) *intColumn {

	cr := &intColumn{
		colBounds: bounds,
		b:         b,
		numOfRows: numOfRows,
	}

	br := io.NewBinaryReader(b)

	br.SetOffset(cr.colBounds.End - 8)
	br.SetOffset(br.ReadFixed64())

	e := 0

	cr.encoding = encoding.EncodingType(br.ReadUVarint())
	cr.blkBounds.Start = cr.colBounds.Start
	cr.blkBounds.End = br.ReadUVarint()

	e = cr.blkBounds.End

	cr.blkIdxBounds.End = br.ReadUVarint()

	if cr.blkIdxBounds.End != 0 {
		cr.blkIdxBounds.Start = e
		e = cr.blkIdxBounds.End
	}

	cr.encBounds.Start = e
	cr.encBounds.End = br.ReadUVarint()
	e = cr.encBounds.End

	cr.colIdxBounds.End = br.ReadUVarint()

	if cr.colIdxBounds.End != 0 {
		cr.colIdxBounds.Start = e
		e = cr.colIdxBounds.End
	}

	cr.validityBounds.End = br.ReadUVarint()

	if cr.validityBounds.End != 0 {
		cr.validityBounds.Start = e
		e = cr.validityBounds.End
	}

	cr.numOfValues = br.ReadUVarint()
	cr.cardinality = br.ReadUVarint()

	// TODO(gvelo) make column specific.
	cr.blockLen = blockLen
	cr.read()

	cr.vectorPool = vector.DefaultVectorPool()

	return cr

}

func (cr *intColumn) read() {

	//fmt.Println(cr)

	if cr.validityBounds.End != 0 {
		cr.valid = roaring.NewBitmap()
		_, err := cr.valid.FromBuffer(cr.b[cr.validityBounds.Start:cr.validityBounds.End])
		if err != nil {
			panic(err)
		}
	}

	cr.blkIdx = index.NewBlockIndexReader(cr.b, cr.blkIdxBounds)

}

func (cr *intColumn) Iterator() IntIterator {

	dec, blockReader := encoding.GetIntDecoder(
		cr.encoding,
		cr.b,
		cr.blkBounds,
		cr.blockLen,
	)

	if cr.valid != nil {

		return NewIntNullColumnIterator(
			dec,
			blockReader,
			cr.valid,
			cr.numOfRows,
			cr.vectorPool,
		)
	}

	return NewIntColumnIterator(
		dec,
		blockReader,
		cr.vectorPool,
		cr.numOfRows,
	)

}

func (cr *intColumn) Reader() IntColumnReader {

	dec, blockReader := encoding.GetIntDecoder(
		cr.encoding,
		cr.b,
		cr.blkBounds,
		cr.blockLen,
	)

	if cr.valid != nil {
		return NewIntNullColumnReader(blockReader,
			cr.blkIdx,
			dec,
			cr.valid,
			cr.numOfRows,
			cr.vectorPool,
		)
	}

	return NewIntColumnReader(blockReader,
		cr.blkIdx,
		dec,
		cr.vectorPool,
		cr.numOfRows,
		cr.blockLen,
	)
}

type intColumnReader struct {
	idx     index.BlockIndexReader
	br      encoding.BlockReader
	dec     encoding.IntDecoder
	pool    vector.Pool
	colLen  int
	buf     []int
	bufLen  int
	baseRID uint32
	nextRid uint32 // next block RID
}

func NewIntColumnReader(br encoding.BlockReader,
	idx index.BlockIndexReader,
	dec encoding.IntDecoder,
	pool vector.Pool,
	colLen int,
	blockLen int,
) IntColumnReader {

	return &intColumnReader{
		br:     br,
		idx:    idx,
		dec:    dec,
		pool:   pool,
		colLen: colLen,
		buf:    make([]int, blockLen),
	}
}

func (r *intColumnReader) Read(rids []uint32) vector.IntVector {

	v := r.pool.GetIntVector()
	vBuf := v.Buf()

	for i, rid := range rids {

		if r.bufLen == 0 || rid >= r.nextRid {
			r.ReadBlock(rid)
		}

		j := rid - r.baseRID

		vBuf[i] = r.buf[j]

	}

	v.SetLen(len(rids))

	return v

}

func (r *intColumnReader) ReadBlock(rid uint32) {
	var b encoding.Block
	b, r.baseRID, r.nextRid = r.FindBlock(rid)
	r.dec.Decode(b.Bytes(), r.buf)
	r.bufLen = b.Len()
}

func (r *intColumnReader) FindBlock(rid uint32) (encoding.Block, uint32, uint32) {

	blockRid, offset := r.idx.Lookup(rid)

	var nextRid uint32

	for {

		b := r.br.ReadBlock(offset)

		nextRid = blockRid + uint32(b.Len())

		if nextRid > rid {
			return b, blockRid, nextRid
		}

		blockRid = nextRid

	}

}

// null reader

type intNullColumnReader struct {
	idx      index.BlockIndexReader
	br       encoding.BlockReader
	dec      encoding.IntDecoder
	validity *roaring.Bitmap
	valIter  roaring.IntPeekable
	pool     vector.Pool
	buf      []int
	bRID     []uint32
	bufLen   int
	minRID   uint32
	maxRID   uint32
	colLen   int
	pos      int // the position into the decoded buffer
}

func NewIntNullColumnReader(br encoding.BlockReader,
	idx index.BlockIndexReader,
	dec encoding.IntDecoder,
	validity *roaring.Bitmap,
	colLen int,
	pool vector.Pool,
) IntColumnReader {

	return &intNullColumnReader{
		br:       br,
		idx:      idx,
		dec:      dec,
		validity: validity,
		pool:     pool,
		buf:      make([]int, blockLen), // TODO: pass by parameter
		bRID:     make([]uint32, blockLen),
		colLen:   colLen,
		valIter:  validity.Iterator(),
	}
}

func (r *intNullColumnReader) Read(rids []uint32) vector.IntVector {

	v := r.pool.GetIntVector()
	vBuf := v.Buf()

	for i, rid := range rids {

		if int(rid) < r.colLen && r.validity.Contains(rid) {

			if r.bufLen == 0 || rid > r.maxRID {
				r.ReadBlock(rid)
				r.pos = 0
			}

			for ; r.pos < r.bufLen; r.pos++ {

				if rid == r.bRID[r.pos] {
					vBuf[i] = r.buf[r.pos]
					v.SetValid(i)
					break
				}

			}

			// if the loop exited without a match
			// something is wrong
			if r.pos == r.bufLen {
				panic("cannot find RID in buffer")
			}

		} else {
			// TODO(gvelo) use setmem(0) in vec initialization
			//  to avoid this branch
			v.SetInvalid(i)
		}

	}

	v.SetLen(len(rids))

	return v

}

func (r *intNullColumnReader) ReadBlock(rid uint32) {
	var b encoding.Block
	b, r.minRID, r.maxRID = r.FindBlock(rid)
	r.dec.Decode(b.Bytes(), r.buf)
	r.bufLen = b.Len()
}

func (r *intNullColumnReader) FindBlock(rid uint32) (encoding.Block, uint32, uint32) {

	minRID, offset := r.idx.Lookup(rid)

	r.valIter.AdvanceIfNeeded(minRID)

	for {

		b := r.br.ReadBlock(offset)

		for i := 0; i < b.Len(); i++ {
			r.bRID[i] = r.valIter.Next()
		}

		maxRID := r.bRID[b.Len()-1]

		if maxRID >= rid {
			return b, minRID, maxRID
		}

	}

}
