package writers

import (
	"eventdb/collection"
	"eventdb/io"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"log"
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

		offsets = append(offsets, uint64(bw.Offset))
		// TODO: FIX copio las claves, puede ocupar mucho 8 bytes x numero.
		keys = append(keys, uint64(it.Key()))
		// write the key
		bw.WriteEncodedVarint(it.Key())
		bin, _ := it.Get().UserData.(*roaring.Bitmap).MarshalBinary()
		bw.WriteEncodedRawBytes(bin)
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
		bw.WriteEncodedVarint(keys[i])
		bw.WriteEncodedVarint(offsets[i])
		if i%ixl == 0 {
			nk = append(nk, keys[i])
			nl = append(nl, offsets[i])
		}
	}

	bw.WriteEncodedVarint(math.MaxUint64 - 1)
	log.Printf("Writing MAX ")
	log.Printf(" 0")

	log.Printf(fmt.Sprintf("keys %v, lvl %d ", keys, lvl-1))
	log.Printf(fmt.Sprintf("offsets %v, lvl %d", offsets, lvl-1))
	return processSkip(bw, nk, nl, lvl+1, ixl, offset)
}
