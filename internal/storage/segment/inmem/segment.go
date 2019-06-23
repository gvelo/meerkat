//go:generate stringer -type=State

package inmem

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/storage/segment"
)

type State int

const (
	InMem State = iota
	Writing
	OnDisk
)

type Segment struct {
	IndexInfo  *segment.IndexInfo
	ID         string
	EventCount uint32
	State      State
	WriterChan chan *Segment
	Columns    []Column
	log        zerolog.Logger
}

func NewSegment(
	indexInfo *segment.IndexInfo,
	ID string,
	writerChan chan *Segment) *Segment {

	s := &Segment{
		IndexInfo:  indexInfo,
		ID:         ID,
		WriterChan: writerChan,
		State:      InMem,
		Columns:    make([]Column, len(indexInfo.Fields)),
	}

	s.log = log.With().
		Str("component", "inmem.Segment").
		Str("index", indexInfo.Name).Str("segmentID", ID).
		Logger()

	for i, fInfo := range indexInfo.Fields {
		c := NewColumnt(fInfo)
		s.Columns[i] = c
	}

	s.log.Debug().Msg("New Segment Created")

	return s

}

func (s *Segment) Add(event segment.Event) error {

	if s.State != InMem {
		log.Panic().
			Str("state", s.State.String()).
			Msg("trying to add event on invalid segment state")
	}

	for _, col := range s.Columns {

		if value, found := event[col.FieldInfo().Name]; found {
			col.Add(value)
		} else {
			// TODO: null values are not supported yet.
			return fmt.Errorf("missing value for field [%s]", col.FieldInfo().Name)
		}

	}

	s.EventCount++

	return nil
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
		Uint32("eventCount", s.EventCount).
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
		Uint32("eventCount", s.EventCount).
		Msg("segment successfully written to disk")

}
