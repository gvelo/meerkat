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

type SkipListInterface interface {
	Compare(a interface{}, b interface{}) int
	TailValue() interface{}
	HeadValue() interface{}
}

type IntInterface struct{}

func (c IntInterface) Compare(a interface{}, b interface{}) int {
	if a.(int) < b.(int) {
		return -1
	}
	if a.(int) > b.(int) {
		return 1
	}
	return 0
}

func (c IntInterface) TailValue() interface{} {
	return math.MaxInt64
}

func (c IntInterface) HeadValue() interface{} {
	return math.MinInt64
}

type Uint64Interface struct{}

func (c Uint64Interface) Compare(a interface{}, b interface{}) int {
	if a.(uint64) < b.(uint64) {
		return -1
	}
	if a.(uint64) > b.(uint64) {
		return 1
	}
	return 0
}

func (c Uint64Interface) TailValue() interface{} {
	return ^uint64(0)
}

func (c Uint64Interface) HeadValue() interface{} {
	return uint64(0)
}

type Float64Interface struct{}

func (c Float64Interface) Compare(a interface{}, b interface{}) int {
	if a.(float64) < b.(float64) {
		return -1
	}
	if a.(float64) > b.(float64) {
		return 1
	}
	return 0
}
func (c Float64Interface) TailValue() interface{} {
	return math.MaxFloat64
}
func (c Float64Interface) HeadValue() interface{} {
	return math.MaxFloat64 * -1
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
	return fmt.Sprintf("key: %v , level: %d , me: %p", n.key, n.level, &n)
}

type SkipList struct {
	maxLevel       int     // In gral log 1/p ( N )
	p              float32 // 1/p
	head           *SLNode
	tail           *SLNode
	level          int
	updateCallback OnUpdate
	comparator     SkipListInterface
	length         int
	postingStore   *PostingStore
}

func (s *SkipList) Level() int {
	return s.level
}

func NewSkipList(p *PostingStore, c SkipListInterface) *SkipList {

	var u OnUpdate = func(n *SLNode, v interface{}) interface{} {
		if n.UserData == nil {
			n.UserData = p.NewPostingList(uint32(v.(int)))
		} else {
			n.UserData.(*PostingList).Bitmap.Add(uint32(v.(int)))
		}
		return n.UserData
	}

	sl := NewSL(.6, 16, u, c)
	sl.postingStore = p
	return sl
}

func NewSL(p float32, maxLevel int, u OnUpdate, c SkipListInterface) *SkipList {

	if c == nil {
		panic("you should provide a comparator")
	}

	max := c.TailValue()
	min := c.HeadValue()

	var head = NewSLNode(min, nil, maxLevel, Head)
	var tail = NewSLNode(max, nil, maxLevel, Tail)

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

func (list *SkipList) Add(key interface{}, v interface{}) *SkipList {
	list.InsertOrUpdate(key, v)
	return list
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
		x = NewSLNode(key, nil, lvl, Internal)
		list.length++
		if list.updateCallback != nil {
			list.updateCallback(x, v)
		} else {
			x.UserData = v
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

// Iterator is used for lookup and range operations on skiplist
type Iterator struct {
	s     *SkipList
	curr  *SLNode
	valid bool
	i     int
	lvl   int
}

// NewIterator creates an iterator for a skiplist lvl
func (s *SkipList) NewIterator(lvl int) *Iterator {
	it := &Iterator{s: s,
		curr:  s.head,
		i:     0,
		lvl:   lvl,
		valid: true,
	}
	return it
}

// Gets the idx
func (i *Iterator) Idx() int {
	return i.i
}

// Reset the iterator
func (i *Iterator) Reset() {
	i.curr = i.s.head
	i.i = 0
	i.valid = true
}

// Next Item in Level
func (i *Iterator) HasNext() bool {
	i.curr = i.curr.forward[i.lvl]
	i.i += 1
	return i.curr.t != Tail
}

// Get next item Item in Level
func (i *Iterator) Next() *SLNode {
	return i.curr
}

// Get next key Item in Level
func (i *Iterator) Key() interface{} {
	return i.curr.key
}

// Get next item Item in Level
func (i *Iterator) Error() error {
	if !i.valid {
		return errors.New("Invalid state")
	}
	return nil
}
