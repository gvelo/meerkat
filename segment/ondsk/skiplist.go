package ondsk

import (
	"eventdb/io"
	"math"
)

type OnDiskSkipList struct {
	p string
}

func (t OnDiskSkipList) Lookup(id float64) (uint64, error) {
	_, o, ok, _ := t.findOffsetSkipList(t.p, id)
	if ok {
		return o, nil
	}
	return 0, nil
}

func ReadSkipList(ip string) (OnDiskSkipList, error) {
	return OnDiskSkipList{p: ip}, nil
}

func (t OnDiskSkipList) findOffsetSkipList(name string, id float64) (float64, uint64, bool, error) {

	br, err := io.NewBinaryReader(name)

	if err != nil {
		return 0, 0, false, err
	}

	defer br.Close()

	fileType, err := br.ReadHeader()

	if fileType != io.SkipListV1 {
		panic("invalid file type")
	}

	br.Offset = br.Size - 16
	offset, _ := br.DecodeFixed64()
	lvl, _ := br.DecodeFixed64()

	r, start, err := readSkipList(br, int(offset), lvl, id)

	return float64(r), start, true, err

}

func readSkipList(br *io.BinaryReader, offset int, lvl uint64, id float64) (float64, uint64, error) {

	br.Offset = offset

	// search this lvl
	for i := 0; i < int(math.MaxUint32); i++ {

		if lvl == 0 {
			bits, _ := br.DecodeFixed64()
			k := math.Float64frombits(bits)
			kOffset, _ := br.DecodeVarint()
			if k == float64(id) {
				return k, kOffset, nil
			}
			if k > float64(id) {
				// not found
				return k, kOffset, nil
			}
		} else {
			br.Offset = offset
			bits, _ := br.DecodeFixed64()
			k := math.Float64frombits(bits)

			kOffset, _ := br.DecodeVarint()
			next := br.Offset
			bitsn, _ := br.DecodeFixed64()
			kn := math.Float64frombits(bitsn)
			// kn , _:= br.DecodeFixed64()
			br.DecodeVarint()

			if k == float64(id) {
				return readSkipList(br, int(kOffset), lvl-1, id)
			}

			if kn > float64(id) {
				// done, not found
				if lvl == 0 {
					return 0, 0, nil
				}
				return readSkipList(br, int(kOffset), lvl-1, id)
			}
			offset = next
		}

	}
	return 0, 0, nil

}
