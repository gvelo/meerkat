package index

import (
	"fmt"
	"strings"

	"github.com/RoaringBitmap/roaring"
)

type node struct {
	children  map[byte]*node
	container []*record
	posting   *postingInfo
}

type record struct {
	value   string
	posting *postingInfo
}

type inMemTermIndex struct {
	root        *node
	contentSize int
}

func (index *inMemTermIndex) newNode() *node {
	return &node{
		children:  make(map[byte]*node),
		container: make([]*record, index.contentSize)[:0],
	}
}

func (index *inMemTermIndex) addTerm(term string, eventID uint32) {

	current := index.root

	//fmt..Printf("agregando string: %q\n", s)

	for i := 0; i < len(term); i++ {

		fmt.Printf("iterando charactor %v %q de %q \n", i, string(term[i]), term)
		if n, found := current.children[term[i]]; found {
			//fmt..Printf("navegando a child %q\n", s[i])
			current = n
			continue
		}

		//fmt..Printf("buscando en contenido\n")
		for _, record := range current.container {
			fmt.Printf("record=%v term=%q\n", record, term[i])
			if record.value == term[i:] {
				// key exist, return
				//fmt..Printf("encontrado en contenido , retornando\n")
				record.posting.numOfRows++
				record.posting.posting.Add(eventID)
				return
			}
		}

		//fmt..Printf("content lenght %v\n", len(current.content))
		if len(current.container) == index.contentSize {
			//fmt..Printf("BURST !!!!\n")
			//fmt..Printf("contenido %v\n", current.content)
			// burst it
			n := index.newNode()
			newContainer := make([]*record, 0)
			for _, c := range current.container {
				//fmt..Printf("c= %q", c)

				if c.value[0] == term[i] {
					suffix := c.value[1:]
					if len(suffix) == 0 {
						n.posting = &postingInfo{
							numOfRows: 1,
							posting:   roaring.New(),
						}
						n.posting.posting.Add(eventID)
						continue
					}
					newRecord := &record{
						value:   c.value[1:],
						posting: c.posting,
					}
					n.container = append(n.container, newRecord)
				} else {
					newContainer = append(newContainer, c)
				}
			}
			current.container = newContainer
			current.children[term[i]] = n
			current = n
			continue
		}
		//fmt..Printf("agregando a contenido  %v\n", current.content)
		newRecord := &record{
			value: term[i:],
			posting: &postingInfo{
				numOfRows: 1,
				posting:   roaring.New(),
			},
		}
		newRecord.posting.posting.Add(eventID)
		current.container = append(current.container, newRecord)
		return
	}
	if current.posting == nil {
		current.posting = &postingInfo{
			numOfRows: 1,
			posting:   roaring.New(),
		}
	}
	current.posting.posting.Add(eventID)
}

func (index *inMemTermIndex) printTrie() {
	index.printNode(index.root, 0)
}

func (index *inMemTermIndex) printNode(node *node, level int) {

	leveStr := strings.Repeat("|", level)
	fmt.Printf("%v====content===\n", leveStr)
	for _, c := range node.container {
		fmt.Printf("%v%v %v\n", leveStr, c.value, c.posting.numOfRows)
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
		contentSize: 64,
	}
	idx.root = idx.newNode()
	return idx
}
