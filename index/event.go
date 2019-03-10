package index

import "strings"
import "github.com/derekparker/trie"

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

type EventID uint32

// Event represent an indexable event.
type Event struct {
	Timestamp uint64
	Fields    map[string]interface{}
}

func newIndex(name string, fieldsInfo map[string]FieldType) {

}

type Dictionary interface {
	addTerm(fieldName string, term string, eventID EventID)
	addTerms(fieldName string, terms []string, eventID EventID)
	lookupTerm(fieldName string, term string) *PostingList
	lookupTermPrefix(fieldName string, termPrefix string) *PostingList
}

type PostingList interface {
	add(eventID EventID)
}

type PostingStore interface {
	Get(id uint) PostingList
	New() (uint, PostingList)
}

type postingInfo struct {
	numOfRows EventID
	postingID uint
}

type EventStore interface {
	store(eventID EventID, event *Event)
	retrieve(eventID EventID) *Event
	retrieveFields(fieldNames []string, eventID EventID) map[string]interface{}
}

type InMemEventStore struct {
	eventStore []*Event
}

func (s InMemEventStore) store(eventID EventID, event *Event) {
	s.eventStore = append(s.eventStore, event)
}

func (s InMemEventStore) retrieve(eventID EventID) *Event {
	return s.eventStore[eventID]
}

func (s InMemEventStore) retrieveFields(fieldNames []string, eventID EventID) map[string]interface{} {
	//TODO
	return nil
}

type Tokenizer interface {
	tokenize(text string) []string
}

type Index struct {
	fieldInfo    map[string]FieldType
	eventID      EventID
	tokenizer    Tokenizer
	dict         Dictionary
	store        EventStore
	postingStore PostingStore
}

type NaiveTokenizer struct{}

func (tokenizer *NaiveTokenizer) tokenize(text string) []string {
	return strings.Fields(text)
}

type InMemPostingStore struct {
	postings []PostingList
}

func (postingStore InMemPostingStore) New() (uint, *InMemPostingList) {

	newPostingList := &InMemPostingList{
		list: make([]EventID, 1024),
	}

	postingStore.postings = append(postingStore.postings, newPostingList)

	return uint(len(postingStore.postings)), newPostingList

}

func (postingStore InMemPostingStore) get(id uint) PostingList {
	return postingStore.postings[id]
}

type InMemPostingList struct {
	list []EventID
}

func (postingList InMemPostingList) add(id EventID) {
	postingList.list = append(postingList.list, id)
}

type InMemDict struct {
	tries        map[string]*trie.Trie
	postingStore PostingStore
}

func newDict() *InMemDict {
	return &InMemDict{
		tries: make(map[string]*trie.Trie),
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
func (dict *InMemDict) addTerm(fieldName string, term string, eventID EventID) {

	fieldTrie := dict.getTrieForField(fieldName)

	node, found := fieldTrie.Find(term)

	if found {
		pinfo := node.Meta().(*postingInfo)
		pinfo.numOfRows++
		postingList := dict.postingStore.Get(pinfo.postingID)
		postingList.add(eventID)
		return
	}

	postingID, postingList := dict.postingStore.New()
	postingList.add(eventID)
	pinfo := &postingInfo{
		numOfRows: 1,
		postingID: postingID,
	}
	fieldTrie.Add(term, pinfo)

}

func (dict *InMemDict) addTerms(fieldName string, terms []string, eventID EventID) {

}
func (dict *InMemDict) lookupTerm(fieldName string, term string) *PostingList {

}
func (dict *InMemDict) lookupTermPrefix(fieldName string, termPrefix string) *PostingList {

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

//index
//posting
//column store
//dict por campo
