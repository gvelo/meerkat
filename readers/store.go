package readers

import (
	"errors"
	"eventdb/io"
	"eventdb/segment"
	"fmt"
	"log"
	"math"
)

func ReadEvent(name string, id uint64, infos []segment.FieldInfo) (segment.Event, error) {

	offset, start, found, _ := findOffset(name, id)
	evt, evtFound := findEvent(name, offset, infos, found, start, id)
	if !evtFound {
		return nil, errors.New(fmt.Sprintf(" %d not found ", id))
	}
	return evt, nil
}

func findOffset(name string, id uint64) (uint64, uint64, bool, error) {

	br, err := io.NewBinaryReader(name + ".idx")

	if err != nil {
		return 0, 0, false, err
	}

	defer br.Close()

	fileType, err := br.ReadHeader()

	if fileType != io.RowStoreIDXV1 {
		panic("invalid file type")
	}

	br.Offset = int64(br.Size - 16)
	offset, _ := br.DecodeFixed64()
	lvl, _ := br.DecodeFixed64()

	r, start, ok := processLevel(br, offset, lvl, id, 100)

	return r, start, ok, nil

}

func processLevel(br *io.BinaryReader, offset uint64, lvl uint64, id uint64, ixl uint64) (uint64, uint64, bool) {
	log.Println(fmt.Sprintf("Processing lvl %d , offset %d ", lvl, offset))
	// if it is the 1st lvl & the offsets are less than
	// the ixl then return the offset 0 to search from
	if lvl == 1 {
		return offset, 0, false
	}
	br.Offset = int64(offset)
	// total offsets in this lvl
	t, err := br.DecodeVarint()
	if err != nil {
		panic(err)
	}

	// search this lvl
	for i := 0; i < int(t); i++ {

		lvlOffset, _ := br.DecodeVarint()
		calcId := (i + 1) * int(math.Pow(float64(ixl), float64(lvl)))
		log.Println(fmt.Sprintf("Calc Id  %d in lvl %d i %d", calcId, lvl, i))
		if calcId > int(id) {
			log.Println(fmt.Sprintf("Process Next lvl Calc Id  %d , id %d ", calcId, id))
			return processLevel(br, lvlOffset, lvl-1, id, ixl)
		} else {
			continue
		}

		if calcId == int(id) {
			return lvlOffset, uint64(calcId), true
		} else {
			return 0, 0, false
		}

	}
	return 0, 0, false
}

func findEvent(name string, offset uint64, infos []segment.FieldInfo, found bool, startFrom uint64, id uint64) (segment.Event, bool) {

	br, err := io.NewBinaryReader(name)

	if err != nil {
		return nil, false
	}

	defer br.Close()

	fileType, _ := br.ReadHeader()
	if fileType != io.RowStoreV1 {
		panic("invalid file type")
	}

	br.Offset = int64(offset)

	if found {
		evt, err := loadEvent(br, infos)
		return evt, err == nil
	} else {
		for i := startFrom; i <= uint64(br.Size); i++ {

			calcId := i + 1
			//log.Println(fmt.Sprintf("Loading event", calcId))

			if calcId < id {
				loadEvent(br, infos)
				continue
			}

			if calcId == id {
				evt, err := loadEvent(br, infos)
				log.Println(fmt.Sprintf("evt %v , err %v", evt, err))
				return evt, err == nil
			} else {
				return nil, false
			}

		}
	}
	return nil, false
}

func loadEvent(br *io.BinaryReader, infos []segment.FieldInfo) (segment.Event, error) {

	m := make(segment.Event)
	br.DecodeVarint()
	for i, info := range infos {
		id, _ := br.DecodeVarint()
		// it is an empty field
		if int(id) > i {
			continue
		}
		value, err := br.ReadValue(info)
		if err != nil {
			return nil, err
		}
		m[info.FieldName] = value
	}
	return m, nil
}
