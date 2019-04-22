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

	r, start, ok := processLevel(br, offset, lvl, id, 5, 0)

	return r, start, ok, nil

}

func processLevel(br *io.BinaryReader, offset uint64, lvl uint64, id uint64, ixl uint64, start uint64) (uint64, uint64, bool) {
	log.Println(fmt.Sprintf("[processLevel] Processing lvl %d , offset %d ", lvl, offset))
	// if it is the 1st lvl & the offsets are less than
	// the ixl then return the offset 0 to search from
	if lvl == 0 {
		log.Println(fmt.Sprintf("[processLevel] returning offset %d start %d in lvl %d", offset, start, lvl))
		return offset, start, start == id
	}
	br.Offset = int64(offset)
	// total offsets in this lvl
	//t, err := br.DecodeVarint()
	//log.Println(fmt.Sprintf("[processLevel] t  %d ", t))
	//if err != nil {
	//	panic(err)
	//}

	// search this lvl
	var lvlPtr uint64 = 0
	for i := 0; i < int(1000); i++ {

		lvlOffset, _ := br.DecodeVarint()
		dlvlOffset, _ := br.DecodeVarint()

		calcId := i * int(math.Pow(float64(ixl), float64(lvl)))

		ptr := (i - 1) * int(math.Pow(float64(ixl), float64(lvl)))
		log.Println(fmt.Sprintf("[processLevel] Calc Id  %d in lvl %d i %d , local offset %d , downoffset %d ", calcId, lvl, i, lvlOffset, dlvlOffset))

		if calcId > int(id) {
			log.Println(fmt.Sprintf("[processLevel] Process Next lvl Calc Id  %d , id %d ", calcId, id))
			return processLevel(br, lvlPtr, lvl-1, id, ixl, uint64(ptr))
		}

		if calcId == int(id) {
			log.Println(fmt.Sprintf("[processLevel] returning offset %d start %d in lvl %d", lvlOffset, calcId, lvl))
			return lvlOffset, uint64(calcId), true
		}

		lvlPtr = dlvlOffset
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
		evt, err := LoadEvent(br, infos)
		return evt, err == nil
	} else {
		for i := startFrom; i <= uint64(br.Size); i++ {

			calcId := i
			//log.Println(fmt.Sprintf("Loading event", calcId))

			if calcId < id {
				o := br.Offset
				evt, err := LoadEvent(br, infos)
				log.Println(fmt.Sprintf("[findEvent] evt %v , err %v , i %d , offset: %d ", evt, err, i, o))
				continue
			}

			if calcId == id {
				o := br.Offset
				evt, err := LoadEvent(br, infos)
				log.Println(fmt.Sprintf("[findEvent] evt %v , err %v , i %d , offset: %d ", evt, err, i, o))
				return evt, err == nil
			} else {
				return nil, false
			}

		}
	}
	return nil, false
}

func LoadEvent(br *io.BinaryReader, infos []segment.FieldInfo) (segment.Event, error) {

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
