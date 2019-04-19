package readers

import (
	"errors"
	"eventdb/io"
	"eventdb/segment"
	"fmt"
)

func ReadEvent(name string, id uint64) (*segment.Event, error) {

	offset, ok := findOffset(name, id)
	if ok {
		evt, _ := findEvent(name, offset)
		return evt, nil
	} else {
		return nil, errors.New(fmt.Sprintf(" %d not found ", id))
	}

}

func findOffset(name string, id uint64) (uint64, bool) {

	br, err := io.NewBinaryReader(name + ".idx")

	if err != nil {
		return 0, false
	}

	defer br.Close()

	fileType, err := br.ReadHeader()

	if fileType != io.RowStoreIDXV1 {
		panic("invalid file type")
	}

	br.Offset = int64(br.Size - 16)
	offset, _ := br.DecodeFixed64()
	lvl, _ := br.DecodeFixed64()

	r, ok := processLevel(br, offset, lvl, id, 100)

	return r, ok

}

func processLevel(br *io.BinaryReader, offset uint64, lvl uint64, id uint64, e uint64) (uint64, bool) {

	br.Offset = int64(offset)
	// total offset in this lvl
	t, err := br.DecodeVarint()
	if err != nil {
		panic(err)
	}

	// search this lvl
	for i := 1; i <= int(t); i++ {

		lvlOffset, _ := br.DecodeVarint()
		calcId := (i + 1) * int(lvl*e)

		if calcId > int(id) {
			return processLevel(br, lvlOffset, lvl-1, id, e)
		}

		if calcId == int(id) {
			return lvlOffset, true
		} else {
			return 0, false
		}

	}
	return 0, false
}

func findEvent(name string, offset uint64, infos []segment.FieldInfo) (segment.Event, error) {

	br, err := io.NewBinaryReader(name)

	if err != nil {
		return nil, err
	}

	defer br.Close()

	br.Offset = int64(offset)
	m := make(segment.Event)
	for _, info := range infos {
		br.DecodeVarint()
		value, err := br.ReadValue(info)
		if err != nil {
			return nil, err
		}
		m[info.FieldName] = value
	}
	return m, nil
}
