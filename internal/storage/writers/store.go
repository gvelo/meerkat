package writers

import (
	"meerkat/internal/storage/io"
	"meerkat/internal/storage/segment/inmem"
)

func WriteStoreIdx(name string, offsets []*inmem.PageDescriptor, ixl int) error {

	bw, err := io.NewBinaryWriter(name)

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

	bw.WriteFixedUint64(uint64(lvlOffset))
	bw.WriteFixedUint64(uint64(lvl))

	return nil
}

func writeLevel0(bw *io.BinaryWriter, offsets []*inmem.PageDescriptor) ([]uint64, error) {
	o := make([]uint64, 0)

	max := len(offsets)

	for i := 0; i < max; i++ {
		o = append(o, uint64(bw.Offset))
		//bw.WriteVarUint64(uint64(i)) // # columns
		//for _, colOffsets := range offsets {
		bw.WriteVarUint64(uint64(offsets[i].StartID))
		bw.WriteVarUint64(uint64(offsets[i].Offset))
		//}
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
			bw.WriteVarUint64(offsets[i])
			bw.WriteVarUint64(idxOffsets[i])
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
