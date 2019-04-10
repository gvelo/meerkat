package list

import (
	"fmt"
	"unsafe"
)

type Node struct {
	key      uint64
	UserData unsafe.Pointer
	level    int
	forward  []*Node
}

func NewNode(key uint64, value unsafe.Pointer, level int) Node {

	forward := make([]*Node, level)
	for i := 0; i <= level-1; i++ {
		forward[i] = nil
	}
	return Node{key: key, UserData: value, level: level, forward: forward}
}

// Level returns the level of a node in the skiplist
func (n Node) Level() int {
	return n.level
}

func (n Node) String() string {
	return fmt.Sprintf("key: %d , level: %d , me: %p", n.key, n.level, &n)
}
