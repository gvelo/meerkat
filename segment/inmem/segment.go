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
	IndexInfo    *segment.IndexInfo
	ID           string
	EventID      uint32
	Idx          []interface{}
	FieldStorage []*segment.Event
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
	indexInfo *segment.IndexInfo,
	ID string,
	writerChan chan *Segment) *Segment {

	s := &Segment{
		IndexInfo:    indexInfo,
		ID:           ID,
		PostingStore: NewPostingStore(),
		Tokenizer:    text.NewTokenizer(),
		WriterChan:   writerChan,
		Monotonic:    false,
		State:        InMem,
		Idx:          make([]interface{}, len(indexInfo.Fields)),
	}

	s.log = log.With().
		Str("component", "inmem.Segment").
		Str("index", indexInfo.Name).Str("segmentID", ID).
		Logger()

	for i, fInfo := range indexInfo.Fields {
		switch fInfo.Type {
		case segment.FieldTypeTimestamp:
			// TODO add the proper index here
		case segment.FieldTypeInt:
			// TODO add the proper index here
		case segment.FieldTypeKeyword:
			s.Idx[i] = NewBtrie(s.PostingStore)
		case segment.FieldTypeText:
			s.Idx[i] = NewBtrie(s.PostingStore)
		default:
			log.Panic().Int("Type", int(fInfo.Type)).Msg("Invalid Type")
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

	s.EventID++

	// TODO compute min and max timestamp
	// TODO computa monotonic Flag

	for i, info := range s.IndexInfo.Fields {

		if info.Index {

			switch info.Type {

			case segment.FieldTypeInt:
				//TODO Add to the proper index.
			case segment.FieldTypeKeyword:
				idx := s.Idx[i].(*BTrie)
				eventValue := event[info.Name].(string)
				idx.Add(eventValue, s.EventID)
			case segment.FieldTypeText:
				idx := s.Idx[i].(*BTrie)
				eventValue := event[info.Name].(string)
				tokens := s.Tokenizer.Tokenize(eventValue)
				for _, token := range tokens {
					idx.Add(token, s.EventID)
				}
			case segment.FieldTypeTimestamp:
			//TODO Add to the proper index.
			default:
				log.Panic().Int("Type", int(info.Type)).Msg("Invalid Type")
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
		Uint32("eventCount", s.EventID).
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
		Uint32("eventCount", s.EventID).
		Msg("segment successfully written to disk")

}
