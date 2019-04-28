package readers

import (
	"errors"
	"eventdb/io"
	"github.com/RoaringBitmap/roaring"
)

type PostingReader struct {
	*io.BinaryReader
}

func (pr PostingReader) Read(offset int64) (*roaring.Bitmap, error) {

	b := pr.SliceAt(offset)

	//TODO reuse Bitmaps.
	bitmap := roaring.NewBitmap()

	_, err := bitmap.FromBuffer(b)

	if err != nil {
		return nil, err
	}

	return bitmap, nil

}

func NewPostingReader(name string) (*PostingReader, error) {

	br, err := io.NewBinaryReader(name)

	if err != nil {
		return nil, err
	}

	fileType, err := br.ReadHeader()

	if fileType != io.PostingListV1 {
		return nil, errors.New("invalid file type")
	}

	pr := &PostingReader{
		br,
	}

	return pr, nil
}

func (pr *PostingReader) Close() error {
	return pr.Close()
}
