package inmem

import "github.com/RoaringBitmap/roaring"

// postingList holds  the term posting list.
type PostingList struct {
	// offset on disk
	Offset int64

	// the bitmap backing the list
	Bitmap *roaring.Bitmap
}

func (posting *PostingList) Add(eventID uint32) {
	posting.Bitmap.Add(eventID)
}

type PostingStore struct {
	Store []*PostingList
}

func NewPostingStore() *PostingStore {
	return &PostingStore{
		Store: make([]*PostingList, 0),
	}
}

func (s *PostingStore) NewPostingList(eventID uint32) *PostingList {
	p := &PostingList{
		Bitmap: roaring.New(),
	}
	p.Bitmap.Add(eventID)
	s.Store = append(s.Store, p)
	return p
}
