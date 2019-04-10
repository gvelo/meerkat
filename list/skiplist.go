package list

import (
	"errors"
	"fmt"
	"math/rand"
	"unsafe"
)

type SkipList struct {
	maxLevel int     // In gral log 1/p ( N )
	p        float32 // 1/p
	head     *Node
	tail     *Node
	level    int
}

func New(p float32, maxLevel int) *SkipList {

	maxUint64 := ^uint64(0)
	minUint64 := uint64(0)

	var head = NewNode(minUint64, nil, maxLevel)
	var tail = NewNode(maxUint64, nil, maxLevel)

	list := SkipList{maxLevel: maxLevel, p: p, head: &head, tail: &tail, level: 0}

	for i := maxLevel - 1; i >= 0; i-- {
		head.forward[i] = &tail
	}

	return &list
}

func (list *SkipList) randomLevel() int {
	lvl := 1
	for rand.Float32() < list.p && lvl < list.maxLevel {
		lvl = lvl + 1
	}
	return lvl
}

func (list *SkipList) Search(key uint64) (node *Node, found bool) {
	x := list.head
	for i := list.maxLevel - 1; i >= 0; i-- {
		for x.forward[i].key < key {
			x = x.forward[i]
		}

	}
	x = x.forward[0]
	found = false
	if x.key == key {
		node = x
		found = true
	}
	return node, found
}

type OnUpdate func(value unsafe.Pointer) unsafe.Pointer

type CreateNodeValue func() unsafe.Pointer

func (list *SkipList) InsertOrUpdate(key uint64, updateCallback OnUpdate, insertFactoryMethod CreateNodeValue) *SkipList {
	var update = make([]*Node, list.maxLevel)
	x := list.head
	for i := list.maxLevel - 1; i >= 0; i-- {
		for x.forward[i].key < key {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]
	if x.key == key {
		if updateCallback != nil {
			x.UserData = updateCallback(x.UserData)
		} else {
			if insertFactoryMethod == nil {
				panic("You should implement insertFactoryMethod function.")
			}
			x.UserData = insertFactoryMethod()
		}
	} else {
		lvl := list.randomLevel()
		if lvl > list.level {
			for i := list.level; i >= lvl; i-- {
				update[i] = list.head
			}
			list.level = lvl
		}
		if insertFactoryMethod == nil {
			panic("You should implement insertFactoryMethod function.")
		}
		var node = NewNode(key, insertFactoryMethod(), lvl)
		x = &node
		for i := 0; i < lvl; i++ {
			x.forward[i] = update[i].forward[i]
			update[i].forward[i] = x
		}
	}
	return list
}

func (list *SkipList) Delete(key uint64) {
	var update = make([]*Node, list.maxLevel)
	x := list.head
	for i := list.maxLevel - 1; i >= 0; i-- {
		for x.forward[i].key < key {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]
	if x.key == key {
		for i := 0; i < list.level; i++ {
			if update[i].forward[i] != x {
				break
			}
			update[i].forward[i] = x.forward[i]
		}
		x = nil
		for i := list.level; list.level > 0 && list.head.forward[list.level] == nil; i-- {
			list.level = list.level - 1
		}
	}

}

func (list SkipList) String() string {

	result := fmt.Sprintf("Head levels %v\n", list.head)
	result = result + fmt.Sprintf("Tail levels %v\n", list.tail)

	result = result + fmt.Sprint("========  Head  ==========\n")

	for i := 0; i < list.maxLevel; i++ {
		result = result + fmt.Sprintf("level %d %v\n", i, list.head.forward[i])
	}

	result = result + fmt.Sprint("======== Tail ========== \n")

	for x := 0; x < list.maxLevel; x++ {
		result = result + fmt.Sprintf("level %d %v\n", x, list.tail.forward[x])
	}

	result = result + fmt.Sprint("========= Totals ========= \n")

	for z := 0; z < list.maxLevel; z++ {
		ptr := list.head
		lvl := 0
		// el ultimo es un tail
		for ptr.forward[z].forward[z] != nil {
			ptr = ptr.forward[z]
			lvl++
		}
		result = result + fmt.Sprintf("Total level %d %d\n", z, lvl)
	}

	return result
}

func (list *SkipList) DebugItems(lvl int) (string, error) {
	result := ""
	if lvl > list.level {
		return "", errors.New("list > lvl")
	}
	for x := 0; x < lvl; x++ {
		node := list.head.forward[x]
		result = result + fmt.Sprintf("========= Lvl %d ========= \n", x)
		for node.forward[x] != nil {
			result = result + fmt.Sprintf(" { k: %d ,l: %d }", node.key, node.level)
			node = node.forward[x]
		}
		result = result + "\n"
	}

	return result, nil
}
