package inmem

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
)

type NodeType int

const (
	Head NodeType = iota
	Tail
	Internal
)

type SLNode struct {
	key      interface{}
	UserData interface{}
	level    int
	forward  []*SLNode
	t        NodeType
}

func NewSLNode(k interface{}, v interface{}, l int, t NodeType) *SLNode {

	forward := make([]*SLNode, l)
	for i := 0; i <= l-1; i++ {
		forward[i] = nil
	}
	return &SLNode{key: k, UserData: v, level: l, forward: forward, t: t}
}

// Level returns the level of a node in the skiplist
func (n *SLNode) Level() int {
	return n.level
}

func (n *SLNode) String() string {
	return fmt.Sprintf("key: %d , level: %d , me: %p", n.key, n.level, &n)
}

type Comparator interface {
	Compare(a interface{}, b interface{}) int
	tailValue() interface{}
	headValue() interface{}
}

type Uint64Comparator struct{}

func (c Uint64Comparator) Compare(a interface{}, b interface{}) int {
	if a.(uint64) < b.(uint64) {
		return -1
	}
	if a.(uint64) > b.(uint64) {
		return 1
	}
	return 0
}
func (c Uint64Comparator) tailValue() interface{} {
	return ^uint64(0)
}
func (c Uint64Comparator) headValue() interface{} {
	return uint64(0)
}

type Float64Comparator struct{}

func (c Float64Comparator) Compare(a interface{}, b interface{}) int {
	if a.(float64) < b.(float64) {
		return -1
	}
	if a.(float64) > b.(float64) {
		return 1
	}
	return 0
}
func (c Float64Comparator) tailValue() interface{} {
	return math.MaxFloat64
}
func (c Float64Comparator) headValue() interface{} {
	return math.MaxFloat64 * -1
}

type SkipList struct {
	maxLevel       int     // In gral log 1/p ( N )
	p              float32 // 1/p
	head           *SLNode
	tail           *SLNode
	level          int
	updateCallback OnUpdate
	comparator     Comparator
	length         int
	postingStore   *PostingStore
}

func (s *SkipList) Level() int {
	return s.level
}

func NewSkipList(p *PostingStore, u OnUpdate, c Comparator) *SkipList {
	sl := NewSL(.6, 16, u, c)
	sl.postingStore = p
	return sl
}

func NewSL(p float32, maxLevel int, u OnUpdate, c Comparator) *SkipList {

	if c == nil {
		panic("you should provide a comparator")
	}

	maxUint64 := c.tailValue()
	minUint64 := c.headValue()

	var head = NewSLNode(minUint64, nil, maxLevel, Head)
	var tail = NewSLNode(maxUint64, nil, maxLevel, Tail)

	list := SkipList{maxLevel: maxLevel, p: p, head: head, tail: tail, level: 0, updateCallback: u, comparator: c, length: 0}

	for i := maxLevel - 1; i >= 0; i-- {
		head.forward[i] = tail
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

func (list *SkipList) Search(key interface{}) (node *SLNode, found bool) {
	x := list.head
	for i := list.maxLevel - 1; i >= 0; i-- {
		for list.comparator.Compare(x.forward[i].key, key) == -1 {
			x = x.forward[i]
		}

	}
	x = x.forward[0]
	found = false
	if list.comparator.Compare(x.key, key) == 0 {
		node = x
		found = true
	}
	return node, found
}

type OnUpdate func(n *SLNode, v interface{}) interface{}

func (list *SkipList) Add(key interface{}, v uint32) *SkipList {

	return nil
}

func (list *SkipList) InsertOrUpdate(key interface{}, v interface{}) *SkipList {
	var update = make([]*SLNode, list.maxLevel)
	x := list.head
	for i := list.maxLevel - 1; i >= 0; i-- {
		for list.comparator.Compare(x.forward[i].key, key) == -1 {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]
	if list.comparator.Compare(x.key, key) == 0 {
		if list.updateCallback != nil {
			x.UserData = list.updateCallback(x, v)
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
		x = NewSLNode(key, v, lvl, Internal)
		list.length++
		if list.updateCallback != nil {
			list.updateCallback(x, v)
		}

		for i := 0; i < lvl; i++ {
			x.forward[i] = update[i].forward[i]
			update[i].forward[i] = x
		}
	}
	return list
}

func (list *SkipList) Delete(key interface{}) {
	var update = make([]*SLNode, list.maxLevel)
	x := list.head
	for i := list.maxLevel - 1; i >= 0; i-- {
		for list.comparator.Compare(x.forward[i].key, key) == -1 {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]
	if list.comparator.Compare(x.key, key) == 0 {
		for i := 0; i < list.level; i++ {
			if update[i].forward[i] != x {
				break
			}
			update[i].forward[i] = x.forward[i]
		}
		x = nil
		list.length--
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
