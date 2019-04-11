package inmem

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"unsafe"
)

func setUpValuesInsplitList() *SkipList {
	skipList := New(.5, 16)

	var a CreateNodeValue = func() unsafe.Pointer {
		return nil
	}
	skipList.InsertOrUpdate(13, nil, a)
	skipList.InsertOrUpdate(1, nil, a)
	skipList.InsertOrUpdate(123, nil, a)
	skipList.InsertOrUpdate(555, nil, a)
	skipList.InsertOrUpdate(553, nil, a)
	skipList.InsertOrUpdate(554, nil, a)
	skipList.InsertOrUpdate(124, nil, a)
	skipList.InsertOrUpdate(125, nil, a)
	skipList.InsertOrUpdate(1222, nil, a)

	return skipList
}

func TestSkipList_Creation(t *testing.T) {
	ass := assert.New(t)
	skipList := setUpValuesInsplitList()

	ass.Equal(skipList.p, float32(0.5))
	ass.Equal(skipList.maxLevel, 16)

}

func TestSkipList_Insert(t *testing.T) {

	ass := assert.New(t)
	skipList := setUpValuesInsplitList()

	node := skipList.head.forward[0]
	for node.forward[0] != nil {
		ass.True(node.key < node.forward[0].key)
		node = node.forward[0]
	}
}

func TestSkipList_Search(t *testing.T) {

	ass := assert.New(t)
	skipList := setUpValuesInsplitList()

	res, found := skipList.Search(553)
	ass.True(res.key == 553)
	ass.True(found == true)

	res, found = skipList.Search(99999)
	ass.Nil(res)
	ass.False(found)

}

func TestSkipList_Delete(t *testing.T) {
	// flaky test, deberiamos chequear todos los niveles.
	ass := assert.New(t)
	skipList := setUpValuesInsplitList()

	skipList.Delete(553)
	res, found := skipList.Search(553)
	ass.Nil(res)
	ass.False(found)

}

func TestSkipList_InsertOrUpdate(t *testing.T) {
	ass := assert.New(t)
	skipList := New(.5, 16)

	type Holder struct {
		aInt int
	}

	skipList.InsertOrUpdate(13, nil, func() unsafe.Pointer {
		return unsafe.Pointer(&Holder{13})
	})
	skipList.InsertOrUpdate(1, nil, func() unsafe.Pointer {
		return unsafe.Pointer(&Holder{1})
	})
	skipList.InsertOrUpdate(123, nil, func() unsafe.Pointer {
		return unsafe.Pointer(&Holder{123})
	})

	var a OnUpdate = func(value unsafe.Pointer) unsafe.Pointer {
		h := (*Holder)(value)
		h.aInt = h.aInt + 1
		return unsafe.Pointer(h)
	}

	skipList.InsertOrUpdate(123, a, nil)

	res, found := skipList.Search(123)

	ass.True(found)
	ass.NotNil(res)
	i := (*Holder)(res.UserData).aInt
	ass.Equal(i, 124)

}

func BenchmarkSkipList_Insert(b *testing.B) {

	list := createRandomList(5000)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		insertItems(list)
	}

}

func BenchmarkSkipList_Search(b *testing.B) {

	list := createRandomList(5000)
	splitList := insertItems(list)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		splitList.Search(12)
	}

}

func BenchmarkSkipList_Delete(b *testing.B) {

	list := createRandomList(5000)
	splitList := insertItems(list)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		splitList.Delete(12)
	}

}

func createRandomList(qty int) []uint64 {
	list := make([]uint64, qty)
	for i := 0; i < len(list); i++ {
		list[i] = rand.Uint64()
	}
	return list
}

func insertItems(list []uint64) *SkipList {
	splitList := New(.5, 16)
	var a CreateNodeValue = func() unsafe.Pointer {
		return nil
	}
	for i := 0; i < len(list); i++ {
		splitList.InsertOrUpdate(list[i], nil, a)
	}
	return splitList
}
