package index

import (
	"fmt"
	"strings"
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

	//fmt..Printf("agregando string: %q\n", s)

	for i := 0; i < len(term); i++ {

		fmt.Printf("iterando charactor %v %q de %q \n", i, string(term[i]), term)
		if child, found := current.children[term[i]]; found {
			//fmt..Printf("navegando a child %q\n", s[i])
			current = child
			continue
		}

		//fmt..Printf("buscando en contenido\n")
		for _, record := range current.bucket {
			fmt.Printf("record=%v term=%q\n", record, term[i])
			if record.value == term[i:] {
				// key exist, return
				//fmt..Printf("encontrado en contenido , retornando\n")
				record.posting.add(eventID)
				return
			}
		}

		//fmt..Printf("content lenght %v\n", len(current.content))
		if len(current.bucket) == index.bucketSize {
			//fmt..Printf("BURST !!!!\n")
			//fmt..Printf("contenido %v\n", current.content)
			// burst it
			n := index.newNode()
			newBucket := current.bucket[:0]
			for _, c := range current.bucket {
				//fmt..Printf("c= %q", c)

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
		//fmt..Printf("agregando a contenido  %v\n", current.content)
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

func (index *inMemTermIndex) printTrie() {
	index.printNode(index.root, 0)
}

func (index *inMemTermIndex) printNode(node *node, level int) {

	leveStr := strings.Repeat("|", level)
	fmt.Printf("%v====content===\n", leveStr)
	for _, c := range node.bucket {
		fmt.Printf("%v%v %v\n", leveStr, c.value, c.posting.numOfEvents)
	}
	fmt.Printf("%v====childs===\n", leveStr)
	for b, n := range node.children {
		fmt.Printf("%v%v\n", leveStr, string(b))
		index.printNode(n, level+1)
	}

	fmt.Printf("%v=============\n", leveStr)
}

func newInMemTermIndex() *inMemTermIndex {
	idx := &inMemTermIndex{
		bucketSize: 64,
	}
	idx.root = idx.newNode()
	return idx
}
