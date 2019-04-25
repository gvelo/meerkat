package collection

import (
	"errors"
	"fmt"
	"log"
	"math"
)

// Iterator is used for lookup and range operations on skiplist
type Iterator struct {
	s     *SkipList
	curr  *Node
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
	log.Printf(fmt.Sprintf("Creating It { %v }", it))
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
func (i *Iterator) Next() bool {
	i.curr = i.curr.forward[i.lvl]
	i.i += 1
	return i.curr.key != math.MaxUint64
}

// Get next item Item in Level
func (i *Iterator) Get() *Node {
	return i.curr
}

// Get next item Item in Level
func (i *Iterator) Key() uint64 {
	return i.curr.key
}

// Get next item Item in Level
func (i *Iterator) Error() error {
	if !i.valid {
		return errors.New("Invalid state")
	}
	return nil
}