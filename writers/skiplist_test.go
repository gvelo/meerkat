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
	for i := e; i < e+3; i++ {
		b.Add(uint32(i))
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
	k, _, _ := readers.ReadSkip(ip, 10)
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
	k, o, _ := readers.ReadSkip(ip, 6)
	a.Equal(uint64(6), k)
	t.Logf("find sl in dist took %v ", time.Since(start))

	pr, _ := readers.NewPostingReader(sp)
	b, _ := pr.Read(int64(o))

	a.True(b.Contains(6))
	a.True(b.Contains(7))
	a.True(b.Contains(8))
	a.Equal(b.GetCardinality(), uint64(3))

}

func TestSkipList_Read200(t *testing.T) {
	a := assert.New(t)
	ip := "/tmp/skip.bin"
	sp := "/tmp/posting-test.bin"

	start := time.Now()
	sl, st := setUpPostingListStore(200)
	err := WritePosting(sp, st.Store)
	t.Logf("setup sl & saving store took %v ", time.Since(start))

	start = time.Now()
	err = WriteSkip(ip, sl, 20)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	start = time.Now()
	k, o, _ := readers.ReadSkip(ip, 200)
	a.Equal(uint64(200), k)
	t.Logf("find sl in dist took %v ", time.Since(start))

	pr, _ := readers.NewPostingReader(sp)
	b, _ := pr.Read(int64(o))

	a.True(b.Contains(200))
	a.True(b.Contains(201))
	a.True(b.Contains(202))
	a.Equal(b.GetCardinality(), uint64(3))

}

func TestSkipList_Read1M(t *testing.T) {
	a := assert.New(t)
	ip := "/tmp/skip.bin"
	sp := "/tmp/posting-test.bin"

	start := time.Now()
	sl, st := setUpPostingListStore(1000000)
	err := WritePosting(sp, st.Store)
	t.Logf("setup sl & saving store took %v ", time.Since(start))

	start = time.Now()
	err = WriteSkip(ip, sl, 200)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	start = time.Now()
	k, o, _ := readers.ReadSkip(ip, 1000000)
	a.Equal(uint64(1000000), k)
	t.Logf("find sl in dist took %v ", time.Since(start))

	pr, _ := readers.NewPostingReader(sp)
	b, _ := pr.Read(int64(o))

	a.True(b.Contains(1000000))
	a.True(b.Contains(1000001))
	a.True(b.Contains(1000002))
	a.Equal(b.GetCardinality(), uint64(3))

}
