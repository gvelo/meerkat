package index

import (
	"fmt"
	"github.com/deepfabric/bkdtree"
	"github.com/tinylib/msgp/msgp"
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/derekparker/trie"
)

// InMemEventStore is a naive implementation of an EventStore
// for testing pourposes. It holds the events list in an slice.
type InMemEventStore struct {
	eventStore []*Event
}

func (s InMemEventStore) store(eventID uint32, event *Event) {
	s.eventStore = append(s.eventStore, event)
}

func (s InMemEventStore) retrieve(eventID uint32) *Event {
	return s.eventStore[eventID]
}

func (s InMemEventStore) retrieveFields(fieldNames []string, eventID uint32) map[string]interface{} {
	//TODO
	return nil
}

// Create a new in memory EventStore.
func newInMemEventStore() EventStore {
	return &InMemEventStore{
		eventStore: make([]*Event, 1024),
	}
}

// NaiveTokenizer tokenize fields spliting their contents around
// each instance of one or more consecutive white space
// characters, as defined by unicode.IsSpace, .
type NaiveTokenizer struct{}

func (tokenizer *NaiveTokenizer) tokenize(text string) []string {
	return strings.Fields(text)
}

// InMemPostingStore implement a naive posting store using a slice.
type InMemPostingStore struct {
	postings  []*roaring.Bitmap
	postingID int
}

func (postingStore InMemPostingStore) New() (int, *roaring.Bitmap) {
	bitmap := roaring.New()
	postingStore.postings = append(postingStore.postings, bitmap)
	postingStore.postingID++
	return postingStore.postingID, bitmap
}

func (postingStore InMemPostingStore) Get(id int) *roaring.Bitmap {
	return postingStore.postings[id]
}

// Creates a new InMemPostingStore
func newInMemPostingStore() PostingStore {
	return &InMemPostingStore{
		postings:  make([]*roaring.Bitmap, 1024),
		postingID: -1,
	}
}

// InMemDict implements a naive in memory non persistent Dictionary
// using an in-memory prefix trie.
type InMemDict struct {
	trees        map[string]interface{}
	postingStore PostingStore
}

func newDict() *InMemDict {
	return &InMemDict{
		trees:        make(map[string]interface{}),
		postingStore: newInMemPostingStore(),
	}
}

// return the trie for the specified Field, if not trie was found
// create a new one.
func (dict *InMemDict) getTreeForField(fieldInfo FieldInfo) interface{} {

	fieldTree := dict.trees[fieldInfo.fieldName]
	if fieldTree == nil {
		if fieldInfo.fieldType == FieldTypeInt {
			// TODO: REVISAR!
			const MaxMBInMem = 100                  // para que no guarde a disco...
			const DefaultMaxPointsInLeafNode = 1024 // default en solr
			const IntraCap = 4                      // Ni idea, lo vi en los test de la lib
			const NumDimensions = 1                 // 1 numero solo
			const BytesPerDim = 8                   // bytes por dimension, como lo guarda.
			const Dir = "tmp"                       // directorio donde baja el indice.
			const Prefix = "bdk"                    // lo vi en los test de la lib

			// que hago si explota, ver como manejarlo.
			fieldTree, _ = bkdtree.NewBkdTree(MaxMBInMem, DefaultMaxPointsInLeafNode, IntraCap, NumDimensions, BytesPerDim, Dir, Prefix)

		} else {
			fieldTree = trie.New()
			dict.trees[fieldInfo.fieldName] = fieldTree
		}
	}
	return fieldTree
}

func (dict *InMemDict) addTerm(fieldInfo FieldInfo, term string, eventID uint32) {
	fieldTrie := dict.getTreeForField(fieldInfo)
	dict.addTermToTrie(fieldTrie.(*trie.Trie), term, eventID)
}

func (dict *InMemDict) addTermToTrie(trie *trie.Trie, term string, eventID uint32) {
	node, found := trie.Find(term)

	if found {
		pinfo := node.Meta().(*postingInfo)
		pinfo.numOfRows++
		pinfo.posting.Add(eventID)
		return
	}

	postingID, bitmap := dict.postingStore.New()
	bitmap.Add(eventID)
	pinfo := &postingInfo{
		numOfRows: 1,
		postingID: postingID,
		posting:   bitmap,
	}
	trie.Add(term, pinfo)

}

func (dict *InMemDict) addNumber(fieldName FieldInfo, number msgp.Number, eventID uint32) {
	fieldBkdTree := dict.getTreeForField(fieldName)
	dict.addNumberToBkdTree(fieldBkdTree.(*bkdtree.BkdTree), number, eventID)
}

func (dict *InMemDict) addNumberToBkdTree(bkdTree *bkdtree.BkdTree, number msgp.Number, eventID uint32) {
	/*
		numbers := make([]uint64 ,1)
		append(numbers, number.Int())
		node, found := bkdTree.Insert(bkdtree.Point{ Vals:  })

		if found {
			pinfo := node.Meta().(*postingInfo)
			pinfo.numOfRows++
			pinfo.posting.Add(eventID)
			return
		}

		postingID, bitmap := dict.postingStore.New()
		bitmap.Add(eventID)
		pinfo := &postingInfo{
			numOfRows: 1,
			postingID: postingID,
			posting:   bitmap,
		}
		trie.Add(term, pinfo)
	*/
}

func (dict *InMemDict) addTerms(fieldInfo FieldInfo, terms []string, eventID uint32) {
	fieldTrie := dict.getTreeForField(fieldInfo)
	for _, term := range terms {
		dict.addTermToTrie(fieldTrie.(*trie.Trie), term, eventID)
	}
}

func (dict *InMemDict) lookupTerm(fieldInfo FieldInfo, term string) *roaring.Bitmap {
	node, found := dict.getTreeForField(fieldInfo).(*trie.Trie).Find(term)
	if found {
		return node.Meta().(*postingInfo).posting
	}
	return nil
}

func (dict *InMemDict) lookupTermPrefix(fieldInfo FieldInfo, termPrefix string) *roaring.Bitmap {
	return nil
}

// InMemoryIndex implements a in memory Index using in-memory
// implementations of Dictionary, EventStore and PostingStore.
type InMemoryIndex struct {
	fieldInfo    []FieldInfo
	eventID      uint32
	tokenizer    Tokenizer
	dict         Dictionary
	store        EventStore
	postingStore PostingStore
}

// Create a new in-memory Index.
func newInMemoryIndex(fieldsInfo []FieldInfo) *InMemoryIndex {
	return &InMemoryIndex{
		fieldInfo:    fieldsInfo,
		eventID:      0,
		tokenizer:    &NaiveTokenizer{},
		dict:         newDict(),
		store:        newInMemEventStore(),
		postingStore: newInMemPostingStore(),
	}

}

func (index *InMemoryIndex) addEvent(event *Event) {

	index.eventID++

	for _, fieldInfo := range index.fieldInfo {
		if fieldValue, ok := event.Fields[fieldInfo.fieldName]; ok {
			switch fieldInfo.fieldType {
			case FieldTypeKeyword:
				index.dict.addTerm(fieldInfo, fieldValue.(string), index.eventID)
			case FieldTypeInt:
				index.dict.addNumber(fieldInfo, selectNumberType(fieldValue), index.eventID)
			case FieldTypeText:
				terms := index.tokenizer.tokenize(fieldValue.(string))
				index.dict.addTerms(fieldInfo, terms, index.eventID)
			}
			index.store.store(index.eventID, event)
		}
	}

}

func selectNumberType(number interface{}) msgp.Number {
	res := msgp.Number{}
	switch x := number.(type) {
	case uint8:
	case uint16:
	case uint32:
	case uint64:
		res.AsUint(x)
		return res
	case int8:
	case int16:
	case int32:
	case int64:
		res.AsInt(x)
		return res
	case float32:
		res.AsFloat32(x)
		return res
	case float64:
		res.AsFloat64(x)
		return res
	}
	panic(fmt.Sprint("Invalid input for a Number field ", number))

}

func (index *InMemoryIndex) lookup(fieldInfo FieldInfo, term string) *roaring.Bitmap {
	return index.dict.lookupTerm(fieldInfo, term)
}

// posingInfo holds information about the term posting list.
type postingInfo struct {
	// the cardinality of the term. not very usefull
	// in an in-memory implementation. In on-disk
	// implementation should save a seek.
	numOfRows uint32

	// Used as postig list index in on-disk implementation.
	// not used on in-memeroy implementation given that
	// we use a pointer to a in-memory posting list
	postingID int

	// the term's in-memory posting list.
	posting *roaring.Bitmap
}
