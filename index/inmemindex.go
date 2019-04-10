package index

import (
	"unsafe"

	"bitbucket.org/sebad/skiplist/list"
	"github.com/RoaringBitmap/roaring"
	"github.com/derekparker/trie"
)

// InMemEventStore is a naive implementation of an EventStore
// for testing pourposes. It holds the events list in an slice.
type InMemEventStore struct {
	eventStore []Event
}

func (s InMemEventStore) store(eventID uint32, event Event) {
	s.eventStore = append(s.eventStore, event)
}

func (s InMemEventStore) retrieve(eventID uint32) Event {
	return s.eventStore[eventID]
}

func (s InMemEventStore) retrieveFields(fieldNames []string, eventID uint32) map[string]interface{} {
	//TODO
	return nil
}

// Create a new in memory EventStore.
func newInMemEventStore() EventStore {
	return &InMemEventStore{
		eventStore: make([]Event, 1024),
	}
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
			fieldTree = list.New(.5, 16)
		} else {
			fieldTree = trie.New()
		}
		dict.trees[fieldInfo.fieldName] = fieldTree
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
		pinfo.postingSize++
		pinfo.bitmap.Add(eventID)
		return
	}

	postingID, bitmap := dict.postingStore.New()
	bitmap.Add(eventID)
	pinfo := &postingInfo{
		numOfEvents: 1,
		posting:     bitmap,
	}
	trie.Add(term, pinfo)

}

func (dict *InMemDict) addNumber(fieldName FieldInfo, number uint64, eventID uint32) {
	fieldBkdTree := dict.getTreeForField(fieldName)
	dict.addNumberToBkdTree(fieldBkdTree.(*list.SkipList), number, eventID)
}

func (dict *InMemDict) addNumberToBkdTree(skipList *list.SkipList, number uint64, eventID uint32) {

	node, found := skipList.Search(number)
	if found {
		pinfo := (*postingInfo)(node.UserData)
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

	skipList.InsertOrUpdate(number, unsafe.Pointer(pinfo))

}

func (dict *InMemDict) addTerms(fieldInfo FieldInfo, terms []string, eventID uint32) {
	fieldTrie := dict.getTreeForField(fieldInfo)
	for _, term := range terms {
		dict.addTermToTrie(fieldTrie.(*trie.Trie), term, eventID)
	}
}

func (dict *InMemDict) lookupNumber(fieldInfo FieldInfo, number uint64) *roaring.Bitmap {
	skipList := dict.getTreeForField(fieldInfo).(*list.SkipList)
	node, found := skipList.Search(number)
	if found {
		return (*postingInfo)(node.UserData).posting
	}
	return nil
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
		tokenizer:    &UnicodeTokenizer{},
		dict:         newDict(),
		store:        newInMemEventStore(),
		postingStore: newInMemPostingStore(),
	}
}

func (index *InMemoryIndex) addEvent(event Event) {

	index.eventID++

	for _, fieldInfo := range index.fieldInfo {
		if fieldValue, ok := event[fieldInfo.fieldName]; ok {
			switch fieldInfo.fieldType {
			case FieldTypeKeyword:
				index.dict.addTerm(fieldInfo, fieldValue.(string), index.eventID)
			case FieldTypeInt:
				index.dict.addNumber(fieldInfo, fieldValue.(uint64), index.eventID)
			case FieldTypeText:
				terms := index.tokenizer.tokenize(fieldValue.(string))
				index.dict.addTerms(fieldInfo, terms, index.eventID)
			}
			index.store.store(index.eventID, event)
		}
	}

}

func (index *InMemoryIndex) lookup(fieldInfo FieldInfo, term interface{}) *roaring.Bitmap {
	if fieldInfo.fieldType == FieldTypeInt {
		return index.dict.lookupNumber(fieldInfo, uint64(term.(uint64)))
	} else {
		return index.dict.lookupTerm(fieldInfo, term.(string))
	}
	return nil
}

// postingList holds  the term posting list.
type postingList struct {
	// the posting list size
	size uint32

	// offset on disk
	offset int

	// the bitmap backing the list
	bitmap *roaring.Bitmap
}

func newPostingList(eventID uint32) *postingList {
	p := &postingInfo{
		numOfEvents: 1,
		bitmap:      roaring.New(),
	}
	p.bitmap.Add(eventID)
	return p
}

func (posting *postingList) add(eventID uint32) {
	posting.bitmap.Add(eventID)
	posting.size++
}
