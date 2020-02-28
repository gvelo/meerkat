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
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/utils"
	"strings"
)

// Represents a node of the trie
type Node struct {
	Children  map[byte]*Node
	Container map[string]*PostingList
	Posting   *PostingList
	Offset    int
}

// Btrie is a naive implementation of an in memory Burst-Trie
type BTrie struct {
	PostingStore  *PostingStore
	Root          *Node
	MaxBucketSize int
	NumOfNodes    int
	Cardinality   int
}

func (bt *BTrie) newNode() *Node {
	bt.NumOfNodes++
	return &Node{
		Children:  make(map[byte]*Node),
		Container: make(map[string]*PostingList, bt.MaxBucketSize),
	}
}

func (bt *BTrie) Add(key []byte, rid uint32) {

	str := utils.B2S(key)

	currentNode := bt.Root

	for i := 0; i < len(str); i++ {

		if child, found := currentNode.Children[str[i]]; found {
			currentNode = child
			continue
		}

		posting := currentNode.Container[str[i:]]

		if posting != nil {
			// key exist, add to the posting and return
			posting.Add(rid)
			return
		}

		// if the bucket reach the threshold burst it
		if len(currentNode.Container) == bt.MaxBucketSize {

			newNode := bt.newNode()

			for s, p := range currentNode.Container {

				if s[0] == str[i] {

					suffix := s[1:]

					if len(suffix) == 0 {
						newNode.Posting = p
						continue
					}

					newNode.Container[suffix] = p
					delete(currentNode.Container, s)

				}

			}

			currentNode.Children[str[i]] = newNode
			currentNode = newNode

			continue

		}

		currentNode.Container[str[i:]] = bt.PostingStore.NewPostingList(rid)
		bt.Cardinality++
		return
	}

	if currentNode.Posting == nil {
		currentNode.Posting = bt.PostingStore.NewPostingList(rid)
		bt.Cardinality++
		return
	}

	currentNode.Posting.Bitmap.Add(rid)

}

func (bt *BTrie) DumpTrie() {
	bt.dumpNode("ROOT", bt.Root, 0)
}

func (bt *BTrie) dumpNode(value string, node *Node, level int) {

	lPad := strings.Repeat(" ", level)
	fmt.Printf("%v[Node] %v ", lPad, value)
	if node.Posting != nil {
		fmt.Printf(" (%v) \n", node.Posting.Bitmap.GetCardinality())
	} else {
		fmt.Println()
	}

	if len(node.Container) > 0 {
		fmt.Printf("%v [Container]\n", lPad)
		for k, v := range node.Container {
			fmt.Printf("%v  %v (%v) \n", lPad, k, v.Bitmap.GetCardinality())
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

		posting := current.Container[term[i:]]

		if posting != nil {
			return posting.Bitmap
		}

		return nil

	}

	if current.Posting != nil {
		return current.Posting.Bitmap
	}

	return nil

}

func NewBtrie(postingStore *PostingStore) *BTrie {
	idx := &BTrie{
		MaxBucketSize: 64,
		PostingStore:  postingStore,
	}
	idx.Root = idx.newNode()
	return idx
}
