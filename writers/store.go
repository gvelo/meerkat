package writers

import (
	"eventdb/io"
	"eventdb/segment"
)

const tsField = "ts"

func WriteEvents(name string, evts []segment.Event, infos []segment.FieldInfo) ([]int64,error) {

	bw, err := io.NewBinaryWriter(name)

	if err != nil {
		return nil,err
	}

	defer bw.Close()

	err = bw.WriteHeader(io.RowStoreV1)

	if err != nil {
		return nil,err
	}

	var max, min = evts[0][tsField].(uint64), evts[0][tsField].(uint64)

	offsets := make([]int64,len(evts))

	for i, e := range evts {

		if e[tsField].(uint64) > max {
			max = e[tsField].(uint64)
		}

		if e[tsField].(uint64) < min {
			min = e[tsField].(uint64)
		}

		offsets[i]=bw.Offset

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
		return nil,err
	}

	bw.WriteEncodedVarint(uint64(min))
	bw.WriteEncodedVarint(uint64(max))
	bw.WriteEncodedVarint(uint64(len(evts)))

	return offsets,nil
}

