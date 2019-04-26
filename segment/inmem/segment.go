package inmem

import (
	"eventdb/segment"
	"eventdb/text"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Status int

const (
	InMem Status = iota
	Writing
	OnDisk
)

type Segment struct {
	IndexName    string
	ID           string
	FieldInfo    *[]segment.FieldInfo
	eventID      int
	Idx          []interface{}
	FieldStorage *[]segment.Event
	PostingStore *PostingStore
	MinTS        int64
	MaxTS        int64
	Status       int
	Tokenizer    *text.Tokenizer
	WriterChan   chan *Segment
	log          zerolog.Logger
}

func NewSegment(
	indexName string,
	ID string,
	fieldInfo []segment.FieldInfo,
	writerChan chan *Segment) *Segment {

	s := &Segment{
		IndexName:    indexName,
		ID:           ID,
		FieldInfo:    fieldInfo,
		PostingStore: NewPostingStore(),
		Tokenizer:    text.NewTokenizer(),
		WriterChan:   writerChan,
		Status:       InMem,
		Idx:          make([]interface{}, len(fieldInfo)),
	}

	s.log = log.With().
		Str("component", "inmem.Segment").
		Str("index", indexName).Str("segmentID", ID).
		Logger()

	for i, fInfo := range fieldInfo {
		switch fInfo.FieldType {
		case segment.FieldTypeTimestamp:
			// TODO add the proper index here
		case segment.FieldTypeInt:
			// TODO add the proper index here
		case segment.FieldTypeKeyword:
			Idx[i] = NewBtrie(s.PostingStore)
		case segment.FieldTypeText:
			Idx[i] = NewBtrie(s.PostingStore)
		default:
			panic("Unknown field type")
		}
	}

	s.log.Debug().Msg("New Segment Created")

}

func (s *Segment) Add(event *segment.Event) {

	for _, info := range s.FieldInfo {

		if info.Index {

			switch info.FieldType {

			case segment.FieldTypeInt:
			case segment.FieldTypeKeyword
			case segment.FieldTypeText
			case segment.FieldTypeTimestamp

			}

		}

	}

}

func (s *Segment) Write() {

}
