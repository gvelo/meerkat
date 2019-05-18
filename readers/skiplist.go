package readers

import (
	"errors"
	"eventdb/io"
	"math"
)

type OnDiskSkipList struct {
	br         *io.BinaryReader
	rootOffset int
	lvl        int
}

func (sl *OnDiskSkipList) Lookup(id float64) (uint64, error) {
	_, o, ok, _ := sl.findOffsetSkipList(id)
	if ok {
		return o, nil
	}
	return 0, nil
}

func ReadSkipList(file string) (*OnDiskSkipList, error) {

	br, err := io.NewBinaryReader(file)

	if err != nil {
		return nil, err
	}

	fType, err := br.ReadHeader()

	if fType != io.SkipListV1 {
		return nil, errors.New("invalid file type")
	}

	br.Offset = br.Size - 16
	rootOffset, _ := br.DecodeFixed64()
	lvl, _ := br.DecodeFixed64()

	if err != nil {
		return nil, err
	}

	return &OnDiskSkipList{
		br:         br,
		rootOffset: int(rootOffset),
		lvl:        int(lvl),
	}, nil
}

func (sl *OnDiskSkipList) findOffsetSkipList(id float64) (float64, uint64, bool, error) {

	r, start, err := sl.readSkipList(int(sl.rootOffset), sl.lvl, id)
	return float64(r), start, true, err

}

func (sl *OnDiskSkipList) readSkipList(offset int, lvl int, id float64) (float64, uint64, error) {

	sl.br.Offset = offset

	// search this lvl
	for i := 0; i < int(math.MaxUint32); i++ {

		if lvl == 0 {
			bits, _ := sl.br.DecodeFixed64()
			k := math.Float64frombits(bits)
			kOffset, _ := sl.br.DecodeVarint()
			if k == float64(id) {
				return k, kOffset, nil
			}
			if k > float64(id) {
				// not found
				return k, kOffset, nil
			}
		} else {
			sl.br.Offset = offset
			bits, _ := sl.br.DecodeFixed64()
			k := math.Float64frombits(bits)

			kOffset, _ := sl.br.DecodeVarint()
			next := sl.br.Offset
			bitsn, _ := sl.br.DecodeFixed64()
			kn := math.Float64frombits(bitsn)
			// kn , _:= br.DecodeFixed64()
			sl.br.DecodeVarint()

			if k == float64(id) {
				return sl.readSkipList(int(kOffset), lvl-1, id)
			}

			if kn > float64(id) {
				// done, not found
				if lvl == 0 {
					return 0, 0, nil
				}
				return sl.readSkipList(int(kOffset), lvl-1, id)
			}
			offset = next
		}

	}
	return 0, 0, nil

}
