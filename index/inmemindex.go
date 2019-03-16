package index

import (
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

// InMemDict implements a naive in memory non persisten Dictionary
// using an in-memory prefix trie.
type InMemDict struct {
	tries        map[string]*trie.Trie
	postingStore PostingStore
}

func newDict() *InMemDict {
	return &InMemDict{
		tries:        make(map[string]*trie.Trie),
		postingStore: newInMemPostingStore(),
	}
}

// return the trie for the specified Field, if not trie was found
// create a new one.
func (dict *InMemDict) getTrieForField(fieldName string) *trie.Trie {
	fieldTrie := dict.tries[fieldName]
	if fieldTrie == nil {
		fieldTrie = trie.New()
		dict.tries[fieldName] = fieldTrie
	}
	return fieldTrie
}
func (dict *InMemDict) addTerm(fieldName string, term string, eventID uint32) {
	fieldTrie := dict.getTrieForField(fieldName)
	dict.addTermToTrie(fieldTrie, term, eventID)
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

func (dict *InMemDict) addTerms(fieldName string, terms []string, eventID uint32) {
	fieldTrie := dict.getTrieForField(fieldName)
	for _, term := range terms {
		dict.addTermToTrie(fieldTrie, term, eventID)
	}
}

func (dict *InMemDict) lookupTerm(fieldName string, term string) *roaring.Bitmap {
	node, found := dict.getTrieForField(fieldName).Find(term)
	if found {
		return node.Meta().(*postingInfo).posting
	}
	return nil
}

func (dict *InMemDict) lookupTermPrefix(fieldName string, termPrefix string) *roaring.Bitmap {
	return nil
}

// InMemoryIndex implements a in memory Index using in-memory
// implementations of Dictionary, EventStore and PostingStore.
type InMemoryIndex struct {
	fieldInfo    map[string]FieldType
	eventID      uint32
	tokenizer    Tokenizer
	dict         Dictionary
	store        EventStore
	postingStore PostingStore
}

// Create a new in-memory Index.
func newInMemoryIndex(name string, fieldsInfo map[string]FieldType) *InMemoryIndex {
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

	for fieldName, fieldType := range index.fieldInfo {
		if fieldValue, ok := event.Fields[fieldName]; ok {
			if fieldType == FieldTypeText {
				terms := index.tokenizer.tokenize(fieldValue.(string))
				index.dict.addTerms(fieldName, terms, index.eventID)
			} else {
				index.dict.addTerm(fieldName, fieldValue.(string), index.eventID)
			}
			index.store.store(index.eventID, event)
		}
	}

}

func (index *InMemoryIndex) lookup(fieldName string, term string) *roaring.Bitmap {
	return index.dict.lookupTerm(fieldName, term)
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
