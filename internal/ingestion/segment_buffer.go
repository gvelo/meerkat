package ingestion

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type SegmentBuffer struct {
	timer         *time.Timer
	inCh          chan *Partition
	buf           *TableBuffer
	wg            *sync.WaitGroup
	flushInterval time.Duration
	maxSize       int
	tableName     string
	partitionID   uint64
	log           zerolog.Logger
	// start stop mutex state
	// segment writer
}

func NewSegmentBuffer(maxSize int,
	flushInterval time.Duration,
	chSize int,
	wg *sync.WaitGroup,
	tableName string,
	partitionID uint64) *SegmentBuffer {

	return &SegmentBuffer{
		inCh:          make(chan *Partition, chSize),
		buf:           NewTableBuffer(tableName, partitionID),
		wg:            wg,
		flushInterval: flushInterval,
		maxSize:       maxSize,
		tableName:     tableName,
		partitionID:   partitionID,
		log: log.With().
			Str("src", "SegmentBuffer").
			Str("table", tableName).
			Uint64("partitionID", partitionID).
			Logger(),
	}
}

func (b *SegmentBuffer) Start() {
	b.log.Debug().Msg("start")
	b.timer = time.NewTimer(b.flushInterval)
	go b.run()
}

func (b *SegmentBuffer) Stop() {
	b.log.Debug().Msg("stop")
	close(b.inCh)
}

func (b *SegmentBuffer) Add(partition *Partition) {
	b.inCh <- partition
}

func (b *SegmentBuffer) run() {

	defer func() {
		// TODO(gvelo): handle panic
		b.timer.Stop()
		b.wg.Done()
	}()

	for {

		select {

		case partition, chClosed := <-b.inCh:

			if chClosed {
				b.flush()
				return
			}

			b.add(partition)

		case <-b.timer.C:

			if b.buf.len > 0 {
				b.flush()
			}

		}

	}
}

func (b SegmentBuffer) add(partition *Partition) {

	b.log.Debug().Int("size", len(partition.Data)).Msg("add partition")

	b.buf.Append(partition)

	//TODO(gvelo): check buffer size and flush if it
	// is above the configured threshold

}

func (b SegmentBuffer) flush() {

	b.log.Debug().Msg("flush")
	// flush here

	// segmentWriter.writer( b.buf )

	b.buf = NewTableBuffer(b.tableName, b.partitionID)

	// reset the timer.
	if ! b.timer.Stop() {
		<-b.timer.C
	}

	b.timer.Reset(b.flushInterval)

}
