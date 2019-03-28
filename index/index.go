package index

import (
	"github.com/RoaringBitmap/roaring"
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
	//
	FieldTypeTimestamp
)

// Field Info represents the info about a field.
type FieldInfo struct {
	fieldName string
	fieldType FieldType
}

// Event represent an indexable event. Event is the unit of indexing.
// An Event is a set of fields, Each Field has a name and a type that
// represent the kind of data that it holds.
type Event map[string]interface{}

// PostingStore implements a storage facility for posting lists
// a posting list for each term is stored in disk and indexed
// by the posting ID.
type PostingStore interface {
	Get(int) *roaring.Bitmap
	New() (int, *roaring.Bitmap)
}

// EventStore is a column oriented store for Events. Every field
// in the Event is store in their own specialized container
// ie. integers are packed using delta encoding.
type EventStore interface {
	store(eventID uint32, event Event)
	retrieve(eventID uint32) Event
	retrieveFields(fieldNames []string, eventID uint32) map[string]interface{}
}

// A Dictionary represents a map between terms and posting list holding
// the list of events id for that term.
type Dictionary interface {
	addTerm(fieldInfo FieldInfo, term string, eventID uint32)
	addNumber(fieldInfo FieldInfo, number uint64, eventID uint32)
	addTerms(fieldInfo FieldInfo, terms []string, eventID uint32)
	lookupTerm(fieldInfo FieldInfo, term string) *roaring.Bitmap
	lookupNumber(fieldInfo FieldInfo, term uint64) *roaring.Bitmap
	lookupTermPrefix(fieldInfo FieldInfo, termPrefix string) *roaring.Bitmap
}

// An Index is the basic engine. Is the event indexing entry point.
// It keeps together the term dictionary, the posting store and the
// event store providing a unified api to queries executors.
type Index interface {
	addEvent(event *Event)
	lookup(fieldName string, term string) *roaring.Bitmap
}
