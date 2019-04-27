package readers

import (
	"eventdb/io"
	"fmt"
	"log"
	"math"
)

func ReadSkip(ip string, sp string, min int, max int) (uint64, uint64, error) {

	br, err := io.NewBinaryReader(ip)

	if err != nil {
		return 0, 0, err
	}

	defer br.Close()

	fileType, err := br.ReadHeader()

	if fileType != io.SkipListV1 {
		panic("invalid file type")
	}

	key, o, ok, err := findOffsetSkipList(ip, uint64(min))
	if ok {
		return key, o, nil
	}
	return 0, 0, err
}

func findOffsetSkipList(name string, id uint64) (uint64, uint64, bool, error) {

	br, err := io.NewBinaryReader(name)

	if err != nil {
		return 0, 0, false, err
	}

	defer br.Close()

	fileType, err := br.ReadHeader()

	if fileType != io.SkipListV1 {
		panic("invalid file type")
	}

	br.Offset = int64(br.Size - 16)
	offset, _ := br.DecodeFixed64()
	lvl, _ := br.DecodeFixed64()

	r, start, err := readSkipList(br, offset, lvl, id)

	return uint64(r), start, true, err

}

func readSkipList(br *io.BinaryReader, offset uint64, lvl uint64, min uint64) (int, uint64, error) {

	//br.Offset = int64(offset)

	// search this lvl
	for i := 0; i < int(math.MaxUint32); i++ {

		br.Offset = int64(offset)

		if lvl == 0 {
			k, _ := br.DecodeVarint()
			kOffset, _ := br.DecodeVarint()
			if k == uint64(min) {
				return int(k), kOffset, nil
			}
			if k > uint64(min) {
				// not found
				return int(k), kOffset, nil
			}
		} else {

			k, _ := br.DecodeVarint()
			kOffset, _ := br.DecodeVarint()
			next := br.Offset
			kn, _ := br.DecodeVarint()
			br.DecodeVarint()

			// TODO falta que pasa cuando esta en la ultima mitad.
			if k == uint64(min) {
				log.Printf(fmt.Sprintf("Loading offset %d, lvl %d , min  %d ", kOffset, lvl-1, min))
				return readSkipList(br, kOffset, lvl-1, min)
			}

			if kn > uint64(min) {
				log.Printf(fmt.Sprintf("bajando offset %d, lvl %d , min  %d ", kOffset, lvl-1, min))
				// done, not found
				if lvl == 0 {
					return 0, 0, nil
				}
				return readSkipList(br, kOffset, lvl-1, min)
			}
			offset = uint64(next)
		}

	}
	return 0, 0, nil

}
