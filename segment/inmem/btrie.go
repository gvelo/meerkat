package inmem

import (
	"fmt"
	"strings"

	"github.com/RoaringBitmap/roaring"
)

type node struct {
	children map[byte]*node
	bucket   []*record
	posting  *PostingList
	offset   int64
}

type record struct {
	value   string
	posting *PostingList
}

type btrie struct {
	root        *node
	bucketSize  int
	size        int
	cardinality int
}

func (bt *btrie) newNode() *node {
	bt.size++
	return &node{
		children: make(map[byte]*node),
		bucket:   make([]*record, bt.bucketSize)[:0],
	}
}

func (bt *btrie) add(str string, eventID uint32) {

	current := bt.root

	for i := 0; i < len(str); i++ {

		if child, found := current.children[str[i]]; found {
			current = child
			continue
		}

		for _, record := range current.bucket {
			if record.value == str[i:] {
				// key exist, return
				record.posting.add(eventID)
				return
			}
		}

		if len(current.bucket) == bt.bucketSize {
			// burst
			n := bt.newNode()
			newBucket := current.bucket[:0]
			for _, c := range current.bucket {

				if c.value[0] == str[i] {
					suffix := c.value[1:]
					if len(suffix) == 0 {
						n.posting = NewPostingList(eventID)
						bt.cardinality++
						continue
					}
					newRecord := &record{
						value:   c.value[1:],
						posting: c.posting,
					}
					n.bucket = append(n.bucket, newRecord)
				} else {
					newBucket = append(newBucket, c)
				}
			}
			current.bucket = newBucket
			current.children[str[i]] = n
			current = n
			continue
		}
		newRecord := &record{
			value:   str[i:],
			posting: NewPostingList(eventID),
		}
		current.bucket = append(current.bucket, newRecord)
		bt.cardinality++
		return
	}

	if current.posting == nil {
		current.posting = NewPostingList(eventID)
		bt.cardinality++
		return
	}

	current.posting.Bitmap.Add(eventID)
}

func (bt *btrie) dumpTrie() {
	bt.dumpNode("ROOT", bt.root, 0)
}

func (bt *btrie) dumpNode(value string, node *node, level int) {

	lPad := strings.Repeat(" ", level)
	fmt.Printf("%v[node]\n", lPad)
	fmt.Printf("%v %v \n", lPad, value)

	if len(node.bucket) > 0 {
		fmt.Printf("%v [bucket]\n", lPad)
		for _, r := range node.bucket {
			fmt.Printf("%v  %v\n", lPad, r.value)
		}
	}

	for c, n := range node.children {
		bt.dumpNode(string(c), n, level+1)
	}

}

func (bt *btrie) lookup(term string) *roaring.Bitmap {

	current := bt.root

	for i := 0; i < len(term); i++ {

		if child, found := current.children[term[i]]; found {
			current = child
			continue
		}

		for _, record := range current.bucket {
			if record.value == term[i:] {
				return record.posting.Bitmap
			}
		}

		return nil
	}

	if current.posting != nil {
		return current.posting.Bitmap
	}

	return nil

}

func newBtrie() *btrie {
	idx := &btrie{
		bucketSize: 64,
	}
	idx.root = idx.newNode()
	return idx
}
