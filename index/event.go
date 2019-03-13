package index

import (
	"strings"

	"github.com/RoaringBitmap/roaring"
	"github.com/derekparker/trie"
)

// FieldType represent the type of a field.
type FieldType int

const (
	// FieldTypeInt represent a numeric int64 field
	FieldTypeInt FieldType = iota
	// FieldTypeText represent a tokenizable text string.
	FieldTypeText
	// FieldTypeKeyword represent a single string field
	FieldTypeKeyword
)

// Event represent an indexable event.
type Event struct {
	Timestamp uint64
	Fields    map[string]interface{}
}

type PostingStore interface {
	Get(int) *roaring.Bitmap
	New() (int, *roaring.Bitmap)
}

type postingInfo struct {
	numOfRows uint32
	postingID int
	posting   *roaring.Bitmap
}

type EventStore interface {
	store(eventID uint32, event *Event)
	retrieve(eventID uint32) *Event
	retrieveFields(fieldNames []string, eventID uint32) map[string]interface{}
}

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

func newInMemEventStore() EventStore {
	return &InMemEventStore{
		eventStore: make([]*Event, 1024),
	}
}

type Tokenizer interface {
	tokenize(text string) []string
}

type NaiveTokenizer struct{}

func (tokenizer *NaiveTokenizer) tokenize(text string) []string {
	return strings.Fields(text)
}

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

func newInMemPostingStore() PostingStore {
	return &InMemPostingStore{
		postings:  make([]*roaring.Bitmap, 1024),
		postingID: -1,
	}
}

type Dictionary interface {
	addTerm(fieldName string, term string, eventID uint32)
	addTerms(fieldName string, terms []string, eventID uint32)
	lookupTerm(fieldName string, term string) *roaring.Bitmap
	lookupTermPrefix(fieldName string, termPrefix string) *roaring.Bitmap
}

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
	return nil
}
func (dict *InMemDict) lookupTermPrefix(fieldName string, termPrefix string) *roaring.Bitmap {
	return nil
}

type Index struct {
	fieldInfo    map[string]FieldType
	eventID      uint32
	tokenizer    Tokenizer
	dict         Dictionary
	store        EventStore
	postingStore PostingStore
}

func newIndex(name string, fieldsInfo map[string]FieldType) *Index {
	return &Index{
		fieldInfo:    fieldsInfo,
		eventID:      0,
		tokenizer:    &NaiveTokenizer{},
		dict:         newDict(),
		store:        newInMemEventStore(),
		postingStore: newInMemPostingStore(),
	}

}

func (index *Index) addEvent(event *Event) {

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
