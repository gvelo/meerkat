package collection

import (
	"errors"
	"fmt"
	"math/rand"
)

type Node struct {
	key      uint64
	UserData interface{}
	level    int
	forward  []*Node
}

func NewSLNode(k uint64, v interface{}, l int) Node {

	forward := make([]*Node, l)
	for i := 0; i <= l-1; i++ {
		forward[i] = nil
	}
	return Node{key: k, UserData: v, level: l, forward: forward}
}

// Level returns the level of a node in the skiplist
func (n Node) Level() int {
	return n.level
}

func (n Node) String() string {
	return fmt.Sprintf("key: %d , level: %d , me: %p", n.key, n.level, &n)
}

type SkipList struct {
	maxLevel int     // In gral log 1/p ( N )
	p        float32 // 1/p
	head     *Node
	tail     *Node
	level    int
}

func NewSL(p float32, maxLevel int) *SkipList {

	maxUint64 := ^uint64(0)
	minUint64 := uint64(0)

	var head = NewSLNode(minUint64, nil, maxLevel)
	var tail = NewSLNode(maxUint64, nil, maxLevel)

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

type OnUpdate func(value interface{}) interface{}

func (list *SkipList) InsertOrUpdate(key uint64, v interface{}, updateCallback OnUpdate) *SkipList {
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
			x.UserData = v
		}
	} else {
		lvl := list.randomLevel()
		if lvl > list.level {
			for i := list.level; i >= lvl; i-- {
				update[i] = list.head
			}
			list.level = lvl
		}
		var node = NewSLNode(key, v, lvl)
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

	r := fmt.Sprintf("Head levels %v\n", list.head)
	r = r + fmt.Sprintf("Tail levels %v\n", list.tail)

	r = r + fmt.Sprint("========  Head  ==========\n")

	for i := 0; i < list.maxLevel; i++ {
		r = r + fmt.Sprintf("level %d %v\n", i, list.head.forward[i])
	}

	r = r + fmt.Sprint("======== Tail ========== \n")

	for x := 0; x < list.maxLevel; x++ {
		r = r + fmt.Sprintf("level %d %v\n", x, list.tail.forward[x])
	}

	r = r + fmt.Sprint("========= Totals ========= \n")

	for z := 0; z < list.maxLevel; z++ {
		ptr := list.head
		lvl := 0
		// el ultimo es un tail
		for ptr.forward[z].forward[z] != nil {
			ptr = ptr.forward[z]
			lvl++
		}
		r = r + fmt.Sprintf("Total level %d %d\n", z, lvl)
	}

	return r
}

func (list *SkipList) DebugItems(lvl int) (string, error) {
	r := ""
	if lvl > list.level {
		return "", errors.New("list > lvl")
	}
	for x := 0; x < lvl; x++ {
		node := list.head.forward[x]
		r = r + fmt.Sprintf("========= Lvl %d ========= \n", x)
		for node.forward[x] != nil {
			r = r + fmt.Sprintf(" { k: %d ,l: %d }", node.key, node.level)
			node = node.forward[x]
		}
		r = r + "\n"
	}

	return r, nil
}
