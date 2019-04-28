package collection

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func setUpValuesInsplitList() *SkipList {
	skipList := NewSL(.5, 16)

	skipList.InsertOrUpdate(13, nil, nil)
	skipList.InsertOrUpdate(1, nil, nil)
	skipList.InsertOrUpdate(123, nil, nil)
	skipList.InsertOrUpdate(555, nil, nil)
	skipList.InsertOrUpdate(553, nil, nil)
	skipList.InsertOrUpdate(554, nil, nil)
	skipList.InsertOrUpdate(124, nil, nil)
	skipList.InsertOrUpdate(125, nil, nil)
	skipList.InsertOrUpdate(1222, nil, nil)

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

	a := assert.New(t)
	skipList := setUpValuesInsplitList()

	res, found := skipList.Search(553)
	a.True(res.key == 553)
	a.True(found == true)

	res, found = skipList.Search(99999)
	a.Nil(res)
	a.False(found)

}

func TestSkipList_Delete(t *testing.T) {
	// flaky test, deberiamos chequear todos los niveles.
	a := assert.New(t)
	skipList := setUpValuesInsplitList()

	skipList.Delete(553)
	res, found := skipList.Search(553)
	a.Nil(res)
	a.False(found)

}

func TestSkipList_InsertOrUpdate(t *testing.T) {
	a := assert.New(t)
	skipList := NewSL(.5, 16)

	type Holder struct {
		aInt int
	}

	skipList.InsertOrUpdate(13, Holder{13}, nil)
	skipList.InsertOrUpdate(1, Holder{1}, nil)
	skipList.InsertOrUpdate(123, Holder{123}, nil)

	var u OnUpdate = func(value interface{}) interface{} {
		h := value.(Holder)
		h.aInt = h.aInt + 1
		return h
	}

	skipList.InsertOrUpdate(123, nil, u)

	res, found := skipList.Search(123)

	a.True(found)
	a.NotNil(res)
	i := res.UserData.(Holder).aInt
	a.Equal(i, 124)

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
	splitList := NewSL(.5, 16)
	for i := 0; i < len(list); i++ {
		splitList.InsertOrUpdate(list[i], nil, nil)
	}
	return splitList
}
