package readers

import (
	"eventdb/io"
	"fmt"
	"log"
	"math"
)

func ReadSkip(ip string, id float64) (float64, uint64, error) {

	br, err := io.NewBinaryReader(ip)

	if err != nil {
		return 0, 0, err
	}

	defer br.Close()

	fileType, err := br.ReadHeader()

	if fileType != io.SkipListV1 {
		panic("invalid file type")
	}

	key, o, ok, err := findOffsetSkipList(ip, id)
	if ok {
		return key, o, nil
	}
	return 0, 0, err
}

func findOffsetSkipList(name string, id float64) (float64, uint64, bool, error) {

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

	return float64(r), start, true, err

}

func readSkipList(br *io.BinaryReader, offset uint64, lvl uint64, id float64) (float64, uint64, error) {

	br.Offset = int64(offset)

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
			br.Offset = int64(offset)
			bits, _ := br.DecodeFixed64()
			k := math.Float64frombits(bits)

			kOffset, _ := br.DecodeVarint()
			next := br.Offset
			kn, _ := br.DecodeFixed64()
			br.DecodeVarint()

			if k == float64(id) {
				log.Printf(fmt.Sprintf("Loading offset %d, lvl %d , id  %v ", kOffset, lvl-1, id))
				return readSkipList(br, kOffset, lvl-1, id)
			}

			if kn > uint64(id) {
				log.Printf(fmt.Sprintf("Loading offset %d, lvl %d , id  %v ", kOffset, lvl-1, id))
				// done, not found
				if lvl == 0 {
					return 0, 0, nil
				}
				return readSkipList(br, kOffset, lvl-1, id)
			}
			offset = uint64(next)
		}

	}
	return 0, 0, nil

}
