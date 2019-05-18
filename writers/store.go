package writers

import (
	"eventdb/io"
	"eventdb/segment"
)

const tsField = "ts"

func WriteEvents(name string, evts []segment.Event, ii *segment.IndexInfo, ixl int) ([][]uint64, error) {

	iOffsets := make([][]uint64, len(ii.Fields))

	// for now, lets take the idx as the FieldId
	for idx, info := range ii.Fields {

		bw, err := io.NewBinaryWriter(name + "." + info.Name + ".bin")

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
			// it doesn't have any value
			if i > 1 && offsets[i-1] == offsets[i] {
				// sets 0 offset
				offsets[i] = 0
			}

			v, ok := e[info.Name]
			if ok { // got it
				bw.WriteEncodedVarint(uint64(i)) // # id
				bw.WriteValue(v, info)
			}

		}

		o := bw.Offset
		bw.WriteEncodedVarint(uint64(min))
		bw.WriteEncodedVarint(uint64(max))
		bw.WriteEncodedVarint(uint64(len(evts)))
		bw.WriteEncodedFixed64(uint64(o))

		iOffsets[idx] = offsets
	}

	WriteStoreIdx(name, iOffsets, ixl)

	return iOffsets, nil
}

func WriteStoreIdx(name string, offsets [][]uint64, ixl int) error {

	bw, err := io.NewBinaryWriter(name + ".idx")

	if err != nil {
		return err
	}

	defer bw.Close()
	err = bw.WriteHeader(io.RowStoreIDXV1)

	err, lvl, lvlOffset := nil, 0, bw.Offset
	lvl0Offsets, _ := writeLevel0(bw, offsets)
	if ixl < len(lvl0Offsets) {
		err, lvl, lvlOffset = processLevel(bw, lvl0Offsets, nil, lvl, ixl, 0, 0)
		if err != nil {
			panic(err)
		}
	}

	bw.WriteEncodedFixed64(uint64(lvlOffset))
	bw.WriteEncodedFixed64(uint64(lvl))

	return nil
}

func writeLevel0(bw *io.BinaryWriter, offsets [][]uint64) ([]uint64, error) {
	o := make([]uint64, 0)

	max := len(offsets[0])
	for i := 1; i < len(offsets); i++ {
		if max < len(offsets[i]) {
			max = len(offsets[i])
		}
	}

	for i := 0; i < max; i++ {
		o = append(o, uint64(bw.Offset))
		bw.WriteEncodedVarint(uint64(i)) // # columns
		for _, colOffsets := range offsets {
			bw.WriteEncodedVarint(colOffsets[i])
		}
	}
	return o, nil
}

func processLevel(bw *io.BinaryWriter, offsets []uint64, idxOffsets []uint64, lvl int, ixl int, ts uint64, lastOffset int) (error, int, int) {

	offset := int(bw.Offset)
	if len(offsets) <= 1 {
		return nil, lvl - 1, lastOffset
	}

	nl := make([]uint64, 0) // offsets storeFile
	ns := make([]uint64, 0) // offsets idxile

	for i := 0; i < int(uint64(len(offsets))); i++ {
		if lvl > 0 {
			o := bw.Offset
			bw.WriteEncodedVarint(offsets[i])
			bw.WriteEncodedVarint(idxOffsets[i])
			if i%ixl == 0 {
				ns = append(ns, uint64(o))
			}
		}
		if i%ixl == 0 {
			nl = append(nl, offsets[i])
			if lvl == 0 {
				ns = append(ns, offsets[i])
			}
		}
	}
	lastOffset = offset
	return processLevel(bw, nl, ns, lvl+1, ixl, ts, offset)
}
