package inmem

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

type Holder struct {
	aInt int
}

var u OnUpdate = func(n *SLNode, v interface{}) interface{} {
	if UserData == nil {
		UserData = Holder{aInt: 0}
	}
	h := UserData.(Holder)
	h.aInt = h.aInt + 1
	return h
}

func setUpValuesInSplitList() *SkipList {

	skipList := NewSL(.5, 16, u, Uint64Comparator{})

	InsertOrUpdate(uint64(13), nil)
	InsertOrUpdate(uint64(1), nil)
	InsertOrUpdate(uint64(123), nil)
	InsertOrUpdate(uint64(555), nil)
	InsertOrUpdate(uint64(553), nil)
	InsertOrUpdate(uint64(554), nil)
	InsertOrUpdate(uint64(124), nil)
	InsertOrUpdate(uint64(125), nil)
	InsertOrUpdate(uint64(1222), nil)

	return skipList
}

func TestSkipList_Creation(t *testing.T) {
	a := assert.New(t)
	skipList := setUpValuesInSplitList()
	a.Equal(length, 9)

	a.Equal(p, float32(0.5))
	a.Equal(maxLevel, 16)

}

func TestSkipList_Insert(t *testing.T) {

	a := assert.New(t)
	skipList := setUpValuesInSplitList()
	a.Equal(length, 9)

	node := forward[0]
	for forward[0] != nil {
		a.True(Compare(key, key) == -1)
		node = forward[0]
	}
}

func TestSkipList_Search(t *testing.T) {

	a := assert.New(t)
	skipList := setUpValuesInSplitList()

	res, found := Search(uint64(553))
	a.Equal(uint64(553), key)
	a.True(found)

	res, found = Search(uint64(99999))
	a.Nil(res)
	a.False(found)

}

func TestSkipList_Delete(t *testing.T) {
	// flaky test, deberiamos chequear todos los niveles.
	a := assert.New(t)
	skipList := setUpValuesInSplitList()
	a.Equal(length, 9)

	Delete(uint64(553))
	a.Equal(length, 8)
	res, found := Search(uint64(553))
	a.Nil(res)
	a.False(found)

}

func TestSkipList_InsertOrUpdate(t *testing.T) {
	a := assert.New(t)
	skipList := NewSL(.5, 16, u, Uint64Comparator{})

	InsertOrUpdate(uint64(13), Holder{13})
	InsertOrUpdate(uint64(1), Holder{1})
	InsertOrUpdate(uint64(123), Holder{123})
	a.Equal(length, 3)

	InsertOrUpdate(uint64(123), nil)
	a.Equal(length, 3)

	res, found := Search(uint64(123))

	a.True(found)
	a.NotNil(res)
	i := res.UserData.(Holder).aInt
	a.Equal(i, 1)

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
		Search(12)
	}

}

func BenchmarkSkipList_Delete(b *testing.B) {

	list := createRandomList(5000)
	splitList := insertItems(list)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		Delete(12)
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
	splitList := NewSL(.5, 16, u, Uint64Comparator{})
	for i := 0; i < len(list); i++ {
		InsertOrUpdate(list[i], nil)
	}
	return splitList
}
