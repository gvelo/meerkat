package inmem

import (
	"eventdb/segment"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"strings"
)

type Node struct {
	Children map[byte]*Node
	Bucket   []*Record
	Posting  *segment.PostingList
	Offset   int64
}

type Record struct {
	Value   string
	Posting *segment.PostingList
}

type BTrie struct {
	Root        *Node
	BucketSize  int
	Size        int
	Cardinality int
}

func (bt *BTrie) newNode() *Node {
	bt.Size++
	return &Node{
		Children: make(map[byte]*Node),
		Bucket:   make([]*Record, bt.BucketSize)[:0],
	}
}

func (bt *BTrie) Add(str string, eventID uint32) {

	current := bt.Root

	for i := 0; i < len(str); i++ {

		if child, found := current.Children[str[i]]; found {
			current = child
			continue
		}

		for _, record := range current.Bucket {
			if record.Value == str[i:] {
				// key exist, return
				record.Posting.Add(eventID)
				return
			}
		}

		if len(current.Bucket) == bt.BucketSize {
			// burst
			n := bt.newNode()
			//TODO Clear the garbage at the end of the slice.
			newBucket := current.Bucket[:0]
			for _, c := range current.Bucket {

				if c.Value[0] == str[i] {
					suffix := c.Value[1:]
					if len(suffix) == 0 {
						n.Posting = segment.NewPostingList(eventID)
						bt.Cardinality++
						continue
					}
					newRecord := &Record{
						Value:   c.Value[1:],
						Posting: c.Posting,
					}
					n.Bucket = append(n.Bucket, newRecord)
				} else {
					newBucket = append(newBucket, c)
				}
			}
			current.Bucket = newBucket
			current.Children[str[i]] = n
			current = n
			continue
		}
		newRecord := &Record{
			Value:   str[i:],
			Posting: segment.NewPostingList(eventID),
		}
		current.Bucket = append(current.Bucket, newRecord)
		bt.Cardinality++
		return
	}

	if current.Posting == nil {
		current.Posting = segment.NewPostingList(eventID)
		bt.Cardinality++
		return
	}

	current.Posting.Bitmap.Add(eventID)
}

func (bt *BTrie) DumpTrie() {
	bt.dumpNode("ROOT", bt.Root, 0)
}

func (bt *BTrie) dumpNode(value string, node *Node, level int) {

	lPad := strings.Repeat(" ", level)
	fmt.Printf("%v[Node]\n", lPad)
	fmt.Printf("%v %v \n", lPad, value)

	if len(node.Bucket) > 0 {
		fmt.Printf("%v [Bucket]\n", lPad)
		for _, r := range node.Bucket {
			fmt.Printf("%v  %v\n", lPad, r.Value)
		}
	}

	for c, n := range node.Children {
		bt.dumpNode(string(c), n, level+1)
	}

}

func (bt *BTrie) Lookup(term string) *roaring.Bitmap {

	current := bt.Root

	for i := 0; i < len(term); i++ {

		if child, found := current.Children[term[i]]; found {
			current = child
			continue
		}

		for _, record := range current.Bucket {
			if record.Value == term[i:] {
				return record.Posting.Bitmap
			}
		}

		return nil
	}

	if current.Posting != nil {
		return current.Posting.Bitmap
	}

	return nil

}

func NewBtrie() *BTrie {
	idx := &BTrie{
		BucketSize: 64,
	}
	idx.Root = idx.newNode()
	return idx
}
