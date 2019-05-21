package readers

import (
	"errors"
	"eventdb/io"
	"eventdb/segment"
	"math"
	"strings"
)

type OnDiskColumn struct {
	fi *segment.FieldInfo
	br *io.BinaryReader
}

type OnDiskStore struct {
	br         *io.BinaryReader
	rootOffset int
	ixl        int
	lvl        int
	indexInfo  *segment.IndexInfo
	columns    []*OnDiskColumn
}

func ReadStore(path string, ii *segment.IndexInfo, ixl int) (*OnDiskStore, error) {

	file, err := io.MMap(path)

	if err != nil {
		return nil, err
	}

	br := file.NewBinaryReader()

	fType, err := br.ReadHeader()

	if fType != io.RowStoreIDXV1 {
		return nil, errors.New("invalid file type")
	}

	br.Offset = br.Size - 16
	rootOffset, _ := br.ReadFixed64()
	lvl, _ := br.ReadFixed64()

	if err != nil {
		return nil, err
	}

	cols := make([]*OnDiskColumn, 0)

	for _, fi := range ii.Fields {
		col, _ := ReadColumn(path, fi)
		cols = append(cols, col)
	}

	return &OnDiskStore{
		br:         br,
		rootOffset: int(rootOffset),
		indexInfo:  ii,
		lvl:        int(lvl),
		ixl:        ixl,
		columns:    cols,
	}, nil
}

// TODO: Revisar interface.
func (sl *OnDiskStore) Lookup(id int) (segment.Event, error) {
	e := make(segment.Event)
	offset, start, _ := sl.findOffset(id)
	offsets, _ := sl.findIdxEvents(offset, start, id)

	for i, info := range sl.indexInfo.Fields {
		c := sl.columns[i]
		v, ok := c.findEvent(offsets[i])
		if !ok {
			e[info.Name] = nil
		}
		e[info.Name] = v
	}
	return e, nil
}

func (sl *OnDiskStore) findOffset(id int) (int, int, error) {

	r, start := sl.processLevel(int(sl.rootOffset), sl.lvl, id, 0)

	return r, start, nil

}

func (sl *OnDiskStore) processLevel(offset int, lvl int, id int, start int) (int, int) {
	// if it is the 1st lvl & the offsets are less than
	// the ixl then return the offset 0 to search from
	if lvl == 0 {
		return offset, start
	}
	sl.br.Offset = offset

	// search this lvl
	var lvlPtr = 0
	// TODO: FIX IT esto puede traer quilombos cuando el id sea mas grande que el ultimo del nivel.
	for i := 0; i < int(math.MaxUint32); i++ {

		lvlOffset, _ := sl.br.ReadVarint64()
		dlvlOffset, _ := sl.br.ReadVarint64()

		calcId := i * int(math.Pow(float64(sl.ixl), float64(lvl)))

		ptr := (i - 1) * int(math.Pow(float64(sl.ixl), float64(lvl)))

		if calcId > int(id) {
			return sl.processLevel(lvlPtr, lvl-1, id, ptr)
		}

		if calcId == int(id) {
			return int(lvlOffset), calcId
		}

		lvlPtr = int(dlvlOffset)
	}
	return 0, 0
}

func (sl *OnDiskStore) findIdxEvents(offset int, startFrom int, id int) ([]int, bool) {

	sl.br.Offset = offset

	for i := startFrom; i <= sl.br.Size; i++ {

		calcId, _ := sl.br.ReadVarint64()

		if int(calcId) < id {
			// TODO: put an offset and skip it by record
			sl.LoadIdx()
			continue
		}

		if int(calcId) == id {
			evt, err := sl.LoadIdx()
			return evt, err == nil
		} else {
			return nil, false
		}
	}
	return nil, false
}

func (sl *OnDiskStore) LoadIdx() ([]int, error) {
	offsets := make([]int, len(sl.indexInfo.Fields))
	for i := range sl.indexInfo.Fields {
		value, err := sl.br.ReadVarint64()
		if err != nil {
			return nil, err
		}
		offsets[i] = int(value)
	}

	return offsets, nil
}

func ReadColumn(path string, info *segment.FieldInfo) (*OnDiskColumn, error) {

	n := strings.Replace(path, idxExt, "."+info.Name+binExt, 1)

	file, err := io.MMap(n)

	if err != nil {
		return nil, err
	}

	br := file.NewBinaryReader()

	fileType, _ := br.ReadHeader()

	if fileType != io.RowStoreV1 {
		panic("invalid file type")
	}

	return &OnDiskColumn{br: br, fi: info}, nil

}

func (c *OnDiskColumn) findEvent(offset int) (interface{}, bool) {

	c.br.Offset = offset

	c.br.ReadVarint64() // ID
	value, err := c.br.ReadValue(c.fi)
	return value, err == nil

}
