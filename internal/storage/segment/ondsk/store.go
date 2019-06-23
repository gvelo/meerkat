package ondsk

import (
	"math"
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment"
)

type OnDiskColumn struct {
	Fi *segment.FieldInfo
	Br *io.BinaryReader
}

type OnDiskStore struct {
	Br         *io.BinaryReader
	RootOffset int
	Ixl        int
	Lvl        int
	IndexInfo  *segment.IndexInfo
	Columns    []*OnDiskColumn
}

// TODO: Revisar interface.
func (sl *OnDiskStore) Lookup(id int) (segment.Event, error) {
	e := make(segment.Event)
	offset, start, _ := sl.findOffset(id)
	offsets, _ := sl.findIdxEvents(offset, start, id)

	for i, info := range sl.IndexInfo.Fields {
		c := sl.Columns[i]
		v, ok := c.findEvent(offsets[i])
		if !ok {
			e[info.Name] = nil
		}
		e[info.Name] = v
	}
	return e, nil
}

func (sl *OnDiskStore) findOffset(id int) (int, int, error) {

	r, start := sl.processLevel(int(sl.RootOffset), sl.Lvl, id, 0)

	return r, start, nil

}

func (sl *OnDiskStore) processLevel(offset int, lvl int, id int, start int) (int, int) {
	// if it is the 1st lvl & the offsets are less than
	// the Ixl then return the offset 0 to search from
	if lvl == 0 {
		return offset, start
	}
	sl.Br.Offset = offset

	// search this Lvl
	var lvlPtr = 0
	// TODO: FIX IT esto puede traer quilombos cuando el id sea mas grande que el ultimo del nivel.
	for i := 0; i < int(math.MaxUint32); i++ {

		lvlOffset, _ := sl.Br.ReadVarint64()
		dlvlOffset, _ := sl.Br.ReadVarint64()

		calcId := i * int(math.Pow(float64(sl.Ixl), float64(lvl)))

		ptr := (i - 1) * int(math.Pow(float64(sl.Ixl), float64(lvl)))

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

	sl.Br.Offset = offset

	for i := startFrom; i <= sl.Br.Size; i++ {

		calcId, _ := sl.Br.ReadVarint64()

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
	offsets := make([]int, len(sl.IndexInfo.Fields))
	for i := range sl.IndexInfo.Fields {
		value, err := sl.Br.ReadVarint64()
		if err != nil {
			return nil, err
		}
		offsets[i] = int(value)
	}

	return offsets, nil
}

func (c *OnDiskColumn) findEvent(offset int) (interface{}, bool) {

	c.Br.Offset = offset

	c.Br.ReadVarint64() // ID
	value, err := c.Br.ReadValue(c.Fi)
	return value, err == nil

}
