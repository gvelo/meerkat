//go:generate stringer -type=State

package inmem

import (
	"eventdb/segment"
	"eventdb/text"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type State int

const (
	InMem State = iota
	Writing
	OnDisk
)

type Segment struct {
	IndexName    string
	ID           string
	FieldInfo    []segment.FieldInfo
	eventID      uint32
	Idx          []interface{}
	FieldStorage []segment.Event
	PostingStore *PostingStore
	MinTS        int64
	MaxTS        int64
	Monotonic    bool
	State        State
	Tokenizer    text.Tokenizer
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
		Monotonic:    false,
		State:        InMem,
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
			s.Idx[i] = NewBtrie(s.PostingStore)
		case segment.FieldTypeText:
			s.Idx[i] = NewBtrie(s.PostingStore)
		default:
			log.Panic().Int("FieldType", int(fInfo.FieldType)).Msg("Invalid FieldType")
		}
	}

	s.log.Debug().Msg("New Segment Created")

	return s

}

func (s *Segment) Add(event map[string]interface{}) {

	if s.State != InMem {
		log.Panic().
			Str("state", s.State.String()).
			Msg("trying to add event on invalid segment state")
	}

	s.eventID++

	// TODO compute min and max timestamp
	// TODO computa monotonic Flag

	for i, info := range s.FieldInfo {

		if info.Index {

			switch info.FieldType {

			case segment.FieldTypeInt:
				//TODO Add to the proper index.
			case segment.FieldTypeKeyword:
				idx := s.Idx[i].(BTrie)
				eventValue := event[info.FieldName].(string)
				idx.Add(eventValue, s.eventID)
			case segment.FieldTypeText:
				idx := s.Idx[i].(BTrie)
				eventValue := event[info.FieldName].(string)
				tokens := s.Tokenizer.Tokenize(eventValue)
				for _, token := range tokens {
					idx.Add(token, s.eventID)
				}
			case segment.FieldTypeTimestamp:
			//TODO Add to the proper index.
			default:
				log.Panic().Int("FieldType", int(info.FieldType)).Msg("Invalid FieldType")
			}

		}

	}

}

func (s *Segment) Write() {

	if s.State != InMem {
		log.Panic().
			Str("state", s.State.String()).
			Msg("trying to write a segment on invalid segment state")
	}

	s.State = Writing

	s.log.Debug().
		Str("state", s.State.String()).
		Uint32("eventCount", s.eventID).
		Msg("writing segment")

	s.WriterChan <- s

}

func (s *Segment) Close() {

	if s.State != Writing {
		log.Panic().
			Str("state", s.State.String()).
			Msg("error trying to close an inmem segment, invalid segment state")
	}

	s.State = OnDisk

	s.log.Debug().
		Str("status", s.State.String()).
		Uint32("eventCount", s.eventID).
		Msg("segment successfully written to disk")

}
