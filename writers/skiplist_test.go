package writers

import (
	"eventdb/readers"
	"eventdb/segment/inmem"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setUpPostingListStore(qty int) (*inmem.SkipList, *inmem.PostingStore) {
	sl := inmem.NewSL(.5, 16, nil, inmem.Float64Comparator{})
	s := inmem.NewPostingStore()
	c := 0.1
	for i := 1; i <= qty; i++ {
		p := createRndPostingList(s, i)
		sl.InsertOrUpdate(c, p)
		c += 1.0
	}
	WritePosting("/tmp/skiplistposting-test.bin", s)
	return sl, s
}

func createRndPostingList(s *inmem.PostingStore, e int) *inmem.PostingList {
	b := s.NewPostingList(uint32(e))
	for i := e; i < e+3; i++ {
		b.Add(uint32(i))
	}
	return b
}

func TestSkipList_Read1(t *testing.T) {
	a := assert.New(t)
	ip := "/tmp/skip.bin"
	sp := "/tmp/posting-test.bin"

	start := time.Now()
	sl, st := setUpPostingListStore(200)
	err := WritePosting(sp, st)
	t.Logf("setup sl & saving store took %v ", time.Since(start))

	start = time.Now()
	err = WriteSkip(ip, sl, 5)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	idx, _ := readers.ReadSkipList(ip)
	start = time.Now()
	o, _ := idx.Lookup(0.1)

	t.Logf("find sl in dist took %v ", time.Since(start))

	start = time.Now()
	pr, _ := readers.ReadPostingStore(sp)
	b, _ := pr.Read(int(o))

	t.Logf("read posting took %v ", time.Since(start))

	a.True(b.Contains(1))
	a.True(b.Contains(2))
	a.True(b.Contains(3))
	a.Equal(b.GetCardinality(), uint64(3))
}

func TestSkipList_Read6(t *testing.T) {
	a := assert.New(t)
	ip := "/tmp/skip.bin"
	sp := "/tmp/posting-test.bin"

	start := time.Now()
	sl, st := setUpPostingListStore(200)
	err := WritePosting(sp, st)
	t.Logf("setup sl & saving store took %v ", time.Since(start))

	start = time.Now()
	err = WriteSkip(ip, sl, 5)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	start = time.Now()
	idx, _ := readers.ReadSkipList(ip)
	o, _ := idx.Lookup(6.1)
	t.Logf("find sl in dist took %v ", time.Since(start))

	start = time.Now()
	pr, _ := readers.ReadPostingStore(sp)
	b, _ := pr.Read(int(o))
	t.Logf("read posting took %v ", time.Since(start))

	a.True(b.Contains(7))
	a.True(b.Contains(8))
	a.True(b.Contains(9))
	a.Equal(b.GetCardinality(), uint64(3))

}

func TestSkipList_Read100(t *testing.T) {
	a := assert.New(t)
	ip := "/tmp/skip.bin"
	sp := "/tmp/posting-test.bin"

	start := time.Now()
	sl, st := setUpPostingListStore(100)
	err := WritePosting(sp, st)
	t.Logf("setup sl & saving store took %v ", time.Since(start))

	start = time.Now()
	err = WriteSkip(ip, sl, 5)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	start = time.Now()
	idx, _ := readers.ReadSkipList(ip)
	o, _ := idx.Lookup(99.1)
	t.Logf("find sl in dist took %v ", time.Since(start))

	start = time.Now()
	pr, _ := readers.ReadPostingStore(sp)
	b, _ := pr.Read(int(o))
	t.Logf("read posting took %v ", time.Since(start))

	a.True(b.Contains(100))
	a.True(b.Contains(101))
	a.True(b.Contains(102))
	a.Equal(b.GetCardinality(), uint64(3))

}

func TestSkipList_Read200(t *testing.T) {
	a := assert.New(t)
	ip := "/tmp/skip.bin"
	sp := "/tmp/posting-test.bin"

	start := time.Now()
	sl, st := setUpPostingListStore(200)
	err := WritePosting(sp, st)
	t.Logf("setup sl & saving store took %v ", time.Since(start))

	start = time.Now()
	err = WriteSkip(ip, sl, 20)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	start = time.Now()
	idx, _ := readers.ReadSkipList(ip)
	o, _ := idx.Lookup(199.1)
	t.Logf("find sl in dist took %v ", time.Since(start))

	start = time.Now()
	pr, _ := readers.ReadPostingStore(sp)
	b, _ := pr.Read(int(o))
	t.Logf("read posting took %v ", time.Since(start))

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
	err := WritePosting(sp, st)
	t.Logf("setup sl & saving store took %v ", time.Since(start))

	start = time.Now()
	err = WriteSkip(ip, sl, 200)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	start = time.Now()
	idx, _ := readers.ReadSkipList(ip)
	o, _ := idx.Lookup(999999.1)
	t.Logf("find sl in dist took %v ", time.Since(start))

	start = time.Now()
	pr, _ := readers.ReadPostingStore(sp)
	b, _ := pr.Read(int(o))
	t.Logf("read posting took %v ", time.Since(start))

	a.True(b.Contains(1000000))
	a.True(b.Contains(1000001))
	a.True(b.Contains(1000002))
	a.Equal(b.GetCardinality(), uint64(3))

}
