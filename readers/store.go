package readers

import (
	"eventdb/io"
	"eventdb/segment"
	"math"
)

func ReadEvent(name string, id uint64, infos []*segment.FieldInfo, ixl int) (segment.Event, error) {
	e := make(segment.Event)
	offset, start, _ := findOffset(name, id, ixl)
	offsets, _ := findIdxEvents(name, offset, infos, start, id)
	for i, info := range infos {
		idxName := name + "." + info.Name
		v, ok := findEvent(idxName, offsets[i], info)
		if !ok {
			e[info.Name] = nil
		}
		e[info.Name] = v
	}
	return e, nil
}

func findOffset(name string, id uint64, ixl int) (int, uint64, error) {

	br, err := io.NewBinaryReader(name + ".idx")

	if err != nil {
		return 0, 0, err
	}

	defer br.Close()

	fileType, err := br.ReadHeader()

	if fileType != io.RowStoreIDXV1 {
		panic("invalid file type")
	}

	br.Offset = br.Size - 16
	offset, _ := br.DecodeFixed64()
	lvl, _ := br.DecodeFixed64()

	r, start := processLevel(br, int(offset), lvl, id, ixl, 0)

	return r, start, nil

}

func processLevel(br *io.BinaryReader, offset int, lvl uint64, id uint64, ixl int, start uint64) (int, uint64) {
	// if it is the 1st lvl & the offsets are less than
	// the ixl then return the offset 0 to search from
	if lvl == 0 {
		return offset, start
	}
	br.Offset = offset

	// search this lvl
	var lvlPtr = 0
	// TODO: FIX IT esto puede traer quilombos cuando el id sea mas grande que el ultimo del nivel.
	for i := 0; i < int(math.MaxUint32); i++ {

		lvlOffset, _ := br.DecodeVarint()
		dlvlOffset, _ := br.DecodeVarint()

		calcId := i * int(math.Pow(float64(ixl), float64(lvl)))

		ptr := (i - 1) * int(math.Pow(float64(ixl), float64(lvl)))

		if calcId > int(id) {
			return processLevel(br, lvlPtr, lvl-1, id, ixl, uint64(ptr))
		}

		if calcId == int(id) {
			return int(lvlOffset), uint64(calcId)
		}

		lvlPtr = int(dlvlOffset)
	}
	return 0, 0
}

func findEvent(name string, offset int, infos *segment.FieldInfo) (interface{}, bool) {

	br, err := io.NewBinaryReader(name + ".bin")

	if err != nil {
		return nil, false
	}

	defer br.Close()

	fileType, _ := br.ReadHeader()
	if fileType != io.RowStoreV1 {
		panic("invalid file type")
	}

	br.Offset = offset

	br.DecodeVarint() // ID

	evt, err := LoadEvent(br, infos)
	return evt, err == nil

}

func LoadEvent(br *io.BinaryReader, info *segment.FieldInfo) (interface{}, error) {
	value, err := br.ReadValue(info)
	if err != nil {
		return nil, err
	}
	return value, err
}

func findIdxEvents(name string, offset int, infos []*segment.FieldInfo, startFrom uint64, id uint64) ([]int, bool) {

	br, err := io.NewBinaryReader(name + ".idx")

	if err != nil {
		return nil, false
	}

	defer br.Close()

	fileType, _ := br.ReadHeader()
	if fileType != io.RowStoreIDXV1 {
		panic("invalid file type")
	}

	br.Offset = offset

	for i := startFrom; i <= uint64(br.Size); i++ {

		calcId, _ := br.DecodeVarint()

		if calcId < id {
			// TODO: put an offset and skip it by record
			LoadIdx(br, infos)
			continue
		}

		if calcId == id {
			evt, err := LoadIdx(br, infos)
			return evt, err == nil
		} else {
			return nil, false
		}
	}
	return nil, false
}

func LoadIdx(br *io.BinaryReader, infos []*segment.FieldInfo) ([]int, error) {
	offsets := make([]int, len(infos))
	for i := range infos {
		value, err := br.DecodeVarint()
		if err != nil {
			return nil, err
		}
		offsets[i] = int(value)
	}

	return offsets, nil
}
