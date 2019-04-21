package writers

import (
	"eventdb/io"
	"eventdb/segment"
	"fmt"
	"log"
)

const tsField = "ts"

func WriteEvents(name string, evts []segment.Event, infos []segment.FieldInfo) ([]uint64, error) {

	bw, err := io.NewBinaryWriter(name)

	if err != nil {
		return nil, err
	}

	defer bw.Close()

	err = bw.WriteHeader(io.RowStoreV1)

	if err != nil {
		return nil, err
	}

	var max, min = evts[0][tsField].(uint64), evts[0][tsField].(uint64)

	offsets := make([]uint64, len(evts))

	for i, e := range evts {

		if e[tsField].(uint64) > max {
			max = e[tsField].(uint64)
		}

		if e[tsField].(uint64) < min {
			min = e[tsField].(uint64)
		}

		offsets[i] = uint64(bw.Offset)

		// write the field count.
		bw.WriteEncodedVarint(uint64(len(e)))

		// for now, lets take the idx as the FieldId
		for fieldId, info := range infos {

			n := info.FieldName
			v, ok := e[n]
			if ok { // got it
				bw.WriteEncodedVarint(uint64(fieldId))
				bw.WriteValue(v, info)
			}

		}

	}

	if err != nil {
		return nil, err
	}
	o := bw.Offset
	bw.WriteEncodedVarint(uint64(min))
	bw.WriteEncodedVarint(uint64(max))
	bw.WriteEncodedVarint(uint64(len(evts)))
	bw.WriteEncodedFixed64(uint64(o))

	WriteStoreIdx(name, offsets)

	return offsets, nil
}

func WriteStoreIdx(name string, offsets []uint64) error {

	bw, err := io.NewBinaryWriter(name + ".idx")

	if err != nil {
		return err
	}

	defer bw.Close()
	log.Println("Header")
	err = bw.WriteHeader(io.RowStoreIDXV1)

	log.Println(fmt.Sprintf("Offsets: %v", offsets))

	err, l, lvlOffset := processLevel(bw, offsets, 0, 5, int(bw.Offset))
	if err != nil {
		panic(err)
	}
	log.Println("Info")
	bw.WriteEncodedFixed64(uint64(lvlOffset))
	bw.WriteEncodedFixed64(uint64(l))

	return nil
}

func processLevel(bw *io.BinaryWriter, offsets []uint64, lvl int, ixl int, prevOffset int) (error, int, int) {

	offset := int(bw.Offset)
	log.Println(fmt.Sprintf("[W processLevel] lvl %d list size %d offset %d", lvl, len(offsets), offset))

	if len(offsets) <= 1 {
		log.Println(fmt.Sprintf("[W processLevel] returning lvl %d offset %d", lvl-1, prevOffset))
		return nil, lvl - 1, prevOffset
	}

	nl := make([]uint64, 0)

	// #items in this lvl
	if lvl > 0 {
		bw.WriteEncodedVarint(uint64(len(offsets)))
	}
	for i := 0; i < int(uint64(len(offsets))); i++ {
		if lvl > 0 {
			bw.WriteEncodedVarint(offsets[i])
		}
		if i%ixl == 0 {
			nl = append(nl, offsets[i])
		}
	}
	log.Println(fmt.Sprintf("[W processLevel] offsets lvl %d %v", lvl-1, nl))
	return processLevel(bw, nl, lvl+1, ixl, offset)
}
