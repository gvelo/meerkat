package ondsk

import (
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/storage/io"
)

type PostingStore struct {
	file *io.MMFile
}

func (ps PostingStore) Read(offset int) (*roaring.Bitmap, error) {

	b := ps.file.NewBinaryReader().SliceAt(offset)

	//TODO reuse Bitmaps.
	bitmap := roaring.NewBitmap()

	_, err := bitmap.FromBuffer(b)

	if err != nil {
		return nil, err
	}

	return bitmap, nil

}

func (ps *PostingStore) Close() error {
	return ps.file.UnMap()
}

func NewPostingStore(file *io.MMFile) *PostingStore {
	return &PostingStore{
		file: file,
	}
}
