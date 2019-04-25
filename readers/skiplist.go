package readers

import (
	"eventdb/io"
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
			b, _ := br.DecodeRawBytes(false)
			if k == uint64(min) {
				bitmap := roaring.New()
				_, err := bitmap.FromBuffer(b)
				if err != nil {
					return 0, nil, err
				}
				return int(k), bitmap, nil
			}

		} else {
			br.DecodeVarint()
			kOffset, _ := br.DecodeVarint()

			kn, _ := br.DecodeVarint()
			br.DecodeVarint()

			if kn > uint64(min) {
				log.Printf("bajando")
				return readSkipList(br, kOffset, lvl-1, min, max)
			}

		}

	}
	return 0, nil, nil

}
