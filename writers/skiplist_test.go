package writers

import (
	"eventdb/collection"
	"eventdb/readers"
	"github.com/RoaringBitmap/roaring"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setUpValuesInsplitList(qty int) *collection.SkipList {
	skipList := collection.NewSL(.5, 16)

	for i := 0; i < qty; i++ {

		r := createRndRoaring(i)
		skipList.InsertOrUpdate(uint64(i), r, nil)

	}

	return skipList
}

func createRndRoaring(n int) *roaring.Bitmap {
	bitmap := roaring.New()
	for i := 0; i < 3; i++ {
		bitmap.AddInt(n + i)
	}
	return bitmap
}

func TestSkipList_Read10(t *testing.T) {
	a := assert.New(t)
	p := "/tmp/skip.bin"

	start := time.Now()
	sl := setUpValuesInsplitList(200)
	t.Logf("setup sl took %v ", time.Since(start))

	start = time.Now()
	err := WriteSkip(p, sl, 5)
	if err != nil {
		t.Fail()
	}
	t.Logf("write sl took %v ", time.Since(start))

	start = time.Now()
	k, bit, _ := readers.ReadSkip(p, 10, 15)
	a.NotNil(bit)
	a.Equal(k, 10)

	t.Logf("find sl took %v ", time.Since(start))

}

func TestSkipList_Read0(t *testing.T) {
	a := assert.New(t)
	p := "/tmp/skip.bin"

	start := time.Now()
	sl := setUpValuesInsplitList(200)
	t.Logf("setup sl took %v ", time.Since(start))

	start = time.Now()
	err := WriteSkip(p, sl, 5)
	if err != nil {
		t.Fail()
	}
	t.Logf("write sl took %v ", time.Since(start))

	start = time.Now()
	k, bit, _ := readers.ReadSkip(p, 0, 15)
	a.NotNil(bit)
	a.Equal(k, 0)

	t.Logf("find sl took %v ", time.Since(start))

}

func TestSkipList_Read200(t *testing.T) {
	a := assert.New(t)
	p := "/tmp/skip.bin"

	start := time.Now()
	sl := setUpValuesInsplitList(20)
	t.Logf("setup sl took %v ", time.Since(start))

	start = time.Now()
	err := WriteSkip(p, sl, 5)
	if err != nil {
		t.Fail()
	}
	t.Logf("write sl took %v ", time.Since(start))

	start = time.Now()
	k, bit, _ := readers.ReadSkip(p, 20, 15)
	a.NotNil(bit)
	a.Equal(k, 200)

	t.Logf("find sl took %v ", time.Since(start))

}
