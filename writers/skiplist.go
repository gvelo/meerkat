package writers

import (
	"eventdb/collection"
	"eventdb/io"
	"eventdb/segment/inmem"
	"math"
)

func WriteSkip(name string, sl *collection.SkipList, ixl int) error {

	bw, err := io.NewBinaryWriter(name)

	if err != nil {
		return err
	}

	defer bw.Close()

	err = bw.WriteHeader(io.SkipListV1)

	if err != nil {
		return err
	}

	offsets := make([]uint64, 0)
	keys := make([]uint64, 0)

	it := sl.NewIterator(0)
	for it.Next() {
		// TODO: FIX copio las claves, puede ocupar mucho...
		keys = append(keys, uint64(it.Key()))
		offsets = append(offsets, uint64(it.Get().UserData.(*inmem.PostingList).Offset))
	}

	writeSkipIdx(bw, keys, offsets, ixl)

	return nil
}

func writeSkipIdx(bw *io.BinaryWriter, keys []uint64, offsets []uint64, ixl int) error {

	err, l, lvlOffset := processSkip(bw, keys, offsets, 0, ixl, int(bw.Offset))
	if err != nil {
		panic(err)
	}
	bw.WriteEncodedFixed64(uint64(lvlOffset))
	bw.WriteEncodedFixed64(uint64(l))

	return nil
}

func processSkip(bw *io.BinaryWriter, keys []uint64, offsets []uint64, lvl int, ixl int, lastOffset int) (error, int, int) {

	offset := int(bw.Offset)
	if len(offsets) <= 2 {
		return nil, lvl - 1, lastOffset
	}

	nl := make([]uint64, 0)
	nk := make([]uint64, 0)

	for i := 0; i < int(uint64(len(offsets))); i++ {
		o := bw.Offset
		bw.WriteEncodedVarint(keys[i])
		bw.WriteEncodedVarint(offsets[i])
		if i%ixl == 0 {
			nk = append(nk, keys[i])
			nl = append(nl, uint64(o))
		}

	}

	bw.WriteEncodedVarint(math.MaxUint64)

	return processSkip(bw, nk, nl, lvl+1, ixl, offset)
}
