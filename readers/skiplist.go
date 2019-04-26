package readers

import (
	"eventdb/io"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"log"
	"math"
)

func ReadSkip(name string, min int, max int) (int, *roaring.Bitmap, error) {

	br, err := io.NewBinaryReader(name)

	if err != nil {
		return 0, nil, err
	}

	defer br.Close()

	fileType, err := br.ReadHeader()

	if fileType != io.SkipListV1 {
		panic("invalid file type")
	}

	br.Offset = int64(br.Size - 16)
	offset, _ := br.DecodeFixed64()
	lvl, _ := br.DecodeFixed64()

	k, roaring, error := readSkipList(br, offset, lvl, min, max)
	if error != nil {
		return 0, nil, error
	}
	return k, roaring, nil
}

func readSkipList(br *io.BinaryReader, offset uint64, lvl uint64, min int, max int) (int, *roaring.Bitmap, error) {

	br.Offset = int64(offset)
	// search this lvl
	for i := 0; i < int(math.MaxUint32); i++ {

		if lvl == 0 {
			k, _ := br.DecodeVarint()
			b, _ := br.DecodeRawBytes(true)
			if k == uint64(min) {
				bitmap := roaring.NewBitmap()
				_, err := bitmap.FromBuffer(b)
				if err != nil {
					return 0, nil, err
				}
				return int(k), bitmap, nil
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
				br.Offset = int64(kOffset)
				k, _ := br.DecodeVarint()
				b, _ := br.DecodeRawBytes(false)
				bitmap := roaring.NewBitmap()
				_, err := bitmap.FromBuffer(b)
				if err != nil {
					return int(k), nil, err
				}
				return int(k), bitmap, nil
			}

			if kn > uint64(min) {
				log.Printf(fmt.Sprintf("bajando offset %d, lvl %d , min  %d ", kOffset, lvl-1, min))
				// done, not found
				if lvl == 0 {
					return 0, nil, nil
				}
				return readSkipList(br, kOffset, lvl-1, min, max)
			}
			br.Offset = next
		}

	}
	return 0, nil, nil

}
