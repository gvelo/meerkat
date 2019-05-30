package ondsk

import (
	"eventdb/io"
	"github.com/RoaringBitmap/roaring"
)

type ColumnInt struct {
	file       io.MMFile
	eventCount int
	root       int
	end        int
}

func (c *ColumnInt) Scan() *CIntScanIter {
	return newCIntScannIter()
}

func (c *ColumnInt) Close() error {
	return c.file.UnMap()
}

type CIntScanIter struct {
	reader     *io.BinaryReader
	eventCount int
	idx        int
}

func (iter *CIntScanIter) next(buf []int) error {
	for i, j := iter.idx, 0; i < iter.eventCount && j < len(buf); i, j = i+1, j+1 {
		v, err := iter.reader.ReadVarInt()
		if err != nil {
			return err
		}
		buf[j] = v
	}
	iter.idx = i
	return nil
}

func newCIntScannIter(eventCount int, r *io.BinaryReader) *CIntScanIter {
	return &CIntScanIter{
		reader:     r,
		eventCount: eventCount,
	}
}

type ColumnInt struct {
	file       io.MMFile
	eventCount int
	root       int
	end        int
}

func (c *ColumnInt) Scan() *CIntScanIter {
	return newCIntScannIter()
}

func (c *ColumnInt) Close() error {
	return c.file.UnMap()
}

type CIntScanIter struct {
	reader     *io.BinaryReader
	eventCount int
	idx        int
}

func (iter *CIntScanIter) next(buf []int) error {
	for i, j := iter.idx, 0; i < iter.eventCount && j < len(buf); i, j = i+1, j+1 {
		v, err := iter.reader.ReadVarInt()
		if err != nil {
			return err
		}
		buf[j] = v
	}
	iter.idx = i
	return nil
}

func newCIntScannIter(eventCount int, r *io.BinaryReader) *CIntScanIter {
	return &CIntScanIter{
		reader:     r,
		eventCount: eventCount,
	}
}

type ColumnStr struct {
	file       io.MMFile
	eventCount int
	root       int
	end        int
}

func (c *ColumnStr) Scan() *CStrScanIter {
	return newCStrScannIter()
}

func (c *ColumnStr) Close() error {
	return c.file.UnMap()
}

type CStrScanIter struct {
	reader     *io.BinaryReader
	eventCount int
	idx        int
}

func (iter *CStrScanIter) next(buf []string) error {
	for i, j := iter.idx, 0; i < iter.eventCount && j < len(buf); i, j = i+1, j+1 {
		v, err := iter.reader.ReadString()
		if err != nil {
			return err
		}
		buf[j] = v
	}
	iter.idx = i
	return nil
}

func newCStrScannIter(eventCount int, r *io.BinaryReader) *CIntScanIter {
	return &CIntScanIter{
		reader:     r,
		eventCount: eventCount,
	}
}

type CStrFilterIter struct {
	reader     *io.BinaryReader
	eventCount int
	idx        int
	bitmap     roaring.Bitmap
}

func (iter *CStrFilterIter) next(buf []string) error {
	for i, j := iter.idx, 0; i < iter.eventCount && j < len(buf); i, j = i+1, j+1 {
		v, err := iter.reader.ReadString()
		if err != nil {
			return err
		}
		buf[j] = v
	}
	iter.idx = i
	return nil
}

func newCStrFilterIter(eventCount int, bitmap roaring.Bitmap, r *io.BinaryReader) *CStrFilterIter {
	return &CStrFilterIter{
		reader:     r,
		eventCount: eventCount,
		bitmap:     bitmap,
	}
}
