package index

import (
	"fmt"
	"strings"

	"github.com/RoaringBitmap/roaring"
)

type node struct {
	children map[byte]*node
	bucket   []*record
	posting  *postingList
}

type record struct {
	value   string
	posting *postingList
}

type inMemTermIndex struct {
	root        *node
	bucketSize  int
	size        int
	cardinality int
}

func (index *inMemTermIndex) newNode() *node {
	index.size++
	return &node{
		children: make(map[byte]*node),
		bucket:   make([]*record, index.bucketSize)[:0],
	}
}

func (index *inMemTermIndex) addTerm(term string, eventID uint32) {

	current := index.root

	for i := 0; i < len(term); i++ {

		if child, found := current.children[term[i]]; found {
			current = child
			continue
		}

		for _, record := range current.bucket {
			if record.value == term[i:] {
				// key exist, return
				record.posting.add(eventID)
				return
			}
		}

		if len(current.bucket) == index.bucketSize {
			// burst it
			n := index.newNode()
			newBucket := current.bucket[:0]
			for _, c := range current.bucket {

				if c.value[0] == term[i] {
					suffix := c.value[1:]
					if len(suffix) == 0 {
						n.posting = newPostingList(eventID)
						index.cardinality++
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
			current.children[term[i]] = n
			current = n
			continue
		}
		newRecord := &record{
			value:   term[i:],
			posting: newPostingList(eventID),
		}
		current.bucket = append(current.bucket, newRecord)
		index.cardinality++
		return
	}

	if current.posting == nil {
		current.posting = newPostingList(eventID)
		index.cardinality++
		return
	}

	current.posting.bitmap.Add(eventID)
}

func (index *inMemTermIndex) dumpTrie() {
	index.dumpNode("ROOT", index.root, 0)
}

func (index *inMemTermIndex) dumpNode(value string, node *node, level int) {

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
		index.dumpNode(string(c), n, level+1)
	}

}

func (index *inMemTermIndex) lookup(term string) *roaring.Bitmap {

	current := index.root

	for i := 0; i < len(term); i++ {

		if child, found := current.children[term[i]]; found {
			current = child
			continue
		}

		for _, record := range current.bucket {
			if record.value == term[i:] {
				return record.posting.bitmap
			}
		}

		return nil
	}

	if current.posting != nil {
		return current.posting.bitmap
	}

	return nil

}

func newInMemTermIndex() *inMemTermIndex {
	idx := &inMemTermIndex{
		bucketSize: 64,
	}
	idx.root = idx.newNode()
	return idx
}
