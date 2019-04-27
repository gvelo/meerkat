package writers

import (
	"eventdb/collection"
	"eventdb/readers"
	"eventdb/segment/inmem"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setUpPostingListStore(qty int) (*collection.SkipList, *inmem.PostingStore) {
	sl := collection.NewSL(.5, 16)
	s := inmem.NewPostingStore()
	for i := 1; i <= qty; i++ {
		p := createRndPostingList(s, i)
		sl.InsertOrUpdate(uint64(i), p, nil)
	}
	WritePosting("/tmp/skiplistposting-test.bin", s.Store)
	return sl, s
}

func createRndPostingList(s *inmem.PostingStore, e int) *inmem.PostingList {
	b := s.NewPostingList(uint32(e))
	for i := e; i < 3; i++ {
		b.Bitmap.AddInt(i)
	}
	return b
}

func TestSkipList_Read10(t *testing.T) {
	a := assert.New(t)
	ip := "/tmp/skip.bin"
	sp := "/tmp/posting-test.bin"

	start := time.Now()
	sl, st := setUpPostingListStore(200)
	err := WritePosting(sp, st.Store)
	t.Logf("setup sl & saving store took %v ", time.Since(start))

	start = time.Now()
	err = WriteSkip(ip, sl, 5)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	start = time.Now()
	k, _, _ := readers.ReadSkip(ip, sp, 10, 15)
	a.Equal(uint64(10), k)

	t.Logf("find sl in dist took %v ", time.Since(start))

}

func TestSkipList_Read6(t *testing.T) {
	a := assert.New(t)
	ip := "/tmp/skip.bin"
	sp := "/tmp/posting-test.bin"

	start := time.Now()
	sl, st := setUpPostingListStore(200)
	err := WritePosting(sp, st.Store)
	t.Logf("setup sl & saving store took %v ", time.Since(start))

	start = time.Now()
	err = WriteSkip(ip, sl, 5)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	start = time.Now()
	k, o, _ := readers.ReadSkip(ip, sp, 6, 15)
	a.Equal(uint64(6), k)
	t.Logf("find sl in dist took %v ", time.Since(start))

	pr, _ := readers.NewPostingReader(sp)
	b, _ := pr.Read(int64(o))

	print(b.GetCardinality())

}
