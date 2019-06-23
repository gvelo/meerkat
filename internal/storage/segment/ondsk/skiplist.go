package ondsk

import (
	"math"
	"meerkat/internal/storage/io"
)

// Devuelve el ofset de la pagina.
type OnDiskSkipList struct {
	Br         *io.BinaryReader
	RootOffset int
	Lvl        int
}

func (sl *OnDiskSkipList) Lookup(id float64) (uint64, error) {
	_, o, ok, _ := sl.findOffsetSkipList(id)
	if ok {
		return o, nil
	}
	return 0, nil
}

func (sl *OnDiskSkipList) findOffsetSkipList(id float64) (float64, uint64, bool, error) {

	r, start, err := sl.readSkipList(int(sl.RootOffset), sl.Lvl, id)
	return float64(r), start, true, err

}

func (sl *OnDiskSkipList) readSkipList(offset int, lvl int, id float64) (float64, uint64, error) {

	sl.Br.Offset = offset

	// search this lvl
	for i := 0; i < int(math.MaxUint32); i++ {

		if lvl == 0 {
			bits, _ := sl.Br.ReadFixed64()
			k := math.Float64frombits(bits)
			kOffset, _ := sl.Br.ReadVarint64()
			if k == float64(id) {
				return k, kOffset, nil
			}
			if k > float64(id) {
				// not found
				return k, kOffset, nil
			}
		} else {
			sl.Br.Offset = offset
			bits, _ := sl.Br.ReadFixed64()
			k := math.Float64frombits(bits)

			kOffset, _ := sl.Br.ReadVarint64()
			next := sl.Br.Offset
			bitsn, _ := sl.Br.ReadFixed64()
			kn := math.Float64frombits(bitsn)
			// kn , _:= Br.ReadFixed64()
			sl.Br.ReadVarint64()

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
