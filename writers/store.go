package writers

import (
	"eventdb/io"
	"eventdb/segment"
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
				bw.StoreValue(v)
			}

		}

	}

	if err != nil {
		return nil, err
	}

	WriteStoreIdx(offsets)

	bw.WriteEncodedVarint(uint64(min))
	bw.WriteEncodedVarint(uint64(max))
	bw.WriteEncodedVarint(uint64(len(evts)))

	return offsets, nil
}

func WriteStoreIdx(name string, offsets []uint64) error {

	bw, err := io.NewBinaryWriter(name + ".idx")

	if err != nil {
		return err
	}

	defer bw.Close()

	err = bw.WriteHeader(io.RowStoreIDXV1)

	processLevel(bw, offsets, 1, 100)

	bw.WriteEncodedFixed64(uint64(bw.Offset - 8))

	return nil
}

func processLevel(bw *io.BinaryWriter, offsets []uint64, lvl int, ixl int) error {

	if len(offsets) <= 2 {
		return nil
	}

	nl := make([]uint64, len(offsets)/ixl^lvl)

	bw.WriteEncodedVarint(uint64(len(offsets)))

	for i := 1; i < len(offsets); i++ {
		bw.WriteEncodedVarint(offsets[i])
		if i%(ixl^lvl) == 0 {
			nl = append(nl, offsets[i])
		}
	}
	return processLevel(bw, nl, lvl+1, ixl)
}
