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

	err, l, _, lvlOffset := processLevel(bw, offsets, nil, 0, 5, uint64(bw.Offset), 0, 0)
	if err != nil {
		panic(err)
	}
	log.Println("Info")
	// cantidad falta
	bw.WriteEncodedFixed64(uint64(lvlOffset))
	bw.WriteEncodedFixed64(uint64(l))

	return nil
}

func processLevel(bw *io.BinaryWriter, offsets []uint64, idxOffsets []uint64, lvl int, ixl int, ts uint64, tb int, lastOffset int) (error, int, int, int) {

	offset := int(bw.Offset)
	log.Println(fmt.Sprintf("[W processLevel] lvl %d list size %d ", lvl, len(offsets)))
	tbb := len(offsets)
	if len(offsets) <= 1 {
		log.Println(fmt.Sprintf("[W processLevel] returning lvl %d ", lvl-1))
		return nil, lvl - 1, tb, lastOffset
	}

	nl := make([]uint64, 0)
	ns := make([]uint64, 0)

	// #items in this lvl

	if lvl > 0 {
		//bw.WriteEncodedVarint(uint64(len(offsets)))
		//log.Println(fmt.Sprintf("[W processLevel] len %d ", uint64(len(offsets))))
	}
	for i := 0; i < int(uint64(len(offsets))); i++ {
		if lvl > 0 {
			coso := bw.Offset
			bw.WriteEncodedVarint(offsets[i])
			bw.WriteEncodedVarint(idxOffsets[i])
			if i%ixl == 0 {
				ns = append(ns, uint64(coso))
			}
		}
		if i%ixl == 0 {
			nl = append(nl, offsets[i])
			if lvl == 0 {
				ns = append(ns, offsets[i])
			}
		}
	}
	log.Println(fmt.Sprintf("[W processLevel] offsets lvl %d %v", lvl+1, nl))
	log.Println(fmt.Sprintf("[W processLevel] offsets internos lvl %d %v", lvl+1, ns))

	lastOffset = offset
	return processLevel(bw, nl, ns, lvl+1, ixl, ts, tbb, offset)
}
