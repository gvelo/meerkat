package readers

import (
	"eventdb/io"
	"eventdb/segment"
)

type StoreReader struct {
	fw  *io.BinaryReader
	rsh rowSubHeader
}

type rowSubHeader struct {
	tsMax uint64
	tsMin uint64
	total uint64
}

type StoreRow struct {
	offsets []uint16
	fields  []interface{}
}

func newStoreReader(name string) (*StoreReader, error) {
	fw, err := io.NewBinaryReader(name)
	if err != nil {
		return nil, err
	}
	sw := &StoreReader{
		fw:  fw,
		rsh: rowSubHeader{tsMax: 0, tsMin: 0, total: 0},
	}
	return sw, nil
}

func (sw *StoreReader) read() ([]segment.Event, []segment.FieldInfo) {
	return nil, nil
}
