package writers

import (
	"eventdb/collection"
	"eventdb/io"
	"eventdb/segment"
)

const tsField = "ts"

func WriteEvents(name string, evts []segment.Event, infos []segment.FieldInfo) error {

	bw, err := io.NewBinaryWriter(name)

	if err != nil {
		return err
	}

	defer bw.Close()

	err = bw.WriteHeader(io.RowStoreV1)

	if err != nil {
		return err
	}

	var sl = collection.NewSL(0.25, 16)

	var max, min = evts[0][tsField].(uint64), evts[0][tsField].(uint64)

	for _, e := range evts {

		fm := make(map[string]int64)

		if e[tsField].(uint64) > max {
			max = e[tsField].(uint64)
		}

		if e[tsField].(uint64) < min {
			min = e[tsField].(uint64)
		}

		for _, info := range infos {
			n := info.FieldName
			v, ok := e[n]
			var offset int64 = -1
			if ok { // got it
				bw.StoreValue(v)
				offset = bw.Offset
			}
			fm[n] = offset
		}

		sl.InsertOrUpdate(e[tsField].(uint64), fm, nil)

	}

	bw.WriteEncodedVarint(uint64(min))
	bw.WriteEncodedVarint(uint64(max))
	bw.WriteEncodedVarint(uint64(len(evts)))

	return nil
}
