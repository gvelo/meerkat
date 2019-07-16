// Copyright 2019 The Meerkat Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package writers

import (
	"github.com/stretchr/testify/assert"
	"meerkat/internal/config"
	"meerkat/internal/storage/readers"
	"meerkat/internal/storage/segment/inmem"
	"meerkat/internal/storage/segment/ondsk"
	"testing"
	"time"
)

func setUpPostingListStore(qty int) (*inmem.SkipList, *inmem.PostingStore) {
	sl := inmem.NewSL(.5, 16, nil, inmem.Float64Interface{})
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
	config.SkipLevelSize = 5
	err = WriteSkipList(ip, sl)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	idx, _ := readers.ReadSkipList(ip, ondsk.FloatInterface{})
	start = time.Now()
	o, _ := idx.Lookup(11.22220)
	a.Equal(0, o, "Error Element Found ")

	o, _ = idx.Lookup(0.1)

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
	config.SkipLevelSize = 5
	err = WriteSkipList(ip, sl)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	start = time.Now()
	idx, _ := readers.ReadSkipList(ip, ondsk.FloatInterface{})
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
	config.SkipLevelSize = 5
	err = WriteSkipList(ip, sl)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	start = time.Now()
	idx, _ := readers.ReadSkipList(ip, ondsk.FloatInterface{})
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
	config.SkipLevelSize = 20
	err = WriteSkipList(ip, sl)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	start = time.Now()
	idx, _ := readers.ReadSkipList(ip, ondsk.FloatInterface{})
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
	config.SkipLevelSize = 200
	err = WriteSkipList(ip, sl)
	if err != nil {
		t.Fail()
	}
	t.Logf("create sl took %v ", time.Since(start))

	start = time.Now()
	idx, _ := readers.ReadSkipList(ip, ondsk.FloatInterface{})
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
