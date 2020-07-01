package ingestion

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/storage"
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
	running       bool
	mu            sync.Mutex
	segmentWriter storage.SegmentWriter
}

func NewSegmentBuffer(maxSize int,
	flushInterval time.Duration,
	chSize int,
	wg *sync.WaitGroup,
	tableName string,
	partitionID uint64,
	segmentWriter storage.SegmentWriter) *SegmentBuffer {

	return &SegmentBuffer{
		inCh:          make(chan *Partition, chSize),
		buf:           NewTableBuffer(tableName, partitionID),
		wg:            wg,
		flushInterval: flushInterval,
		maxSize:       maxSize,
		tableName:     tableName,
		partitionID:   partitionID,
		segmentWriter: segmentWriter,
		log: log.With().
			Str("src", "SegmentBuffer").
			Str("table", tableName).
			Uint64("partitionID", partitionID).
			Logger(),
	}
}

func (b *SegmentBuffer) Start() {

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		return
	}

	b.running = true
	b.log.Debug().Msg("start")
	b.timer = time.NewTimer(b.flushInterval)
	go b.run()

}

func (b *SegmentBuffer) Stop() {

	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return
	}

	b.log.Debug().Msg("stop")
	close(b.inCh)
}

func (b *SegmentBuffer) Append(partition *Partition) {
	b.inCh <- partition
}

func (b *SegmentBuffer) run() {

	defer func() {
		b.stopTimer()
		b.wg.Done()
	}()

	for {

		select {

		case partition, ok := <-b.inCh:

			if !ok {
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

func (b *SegmentBuffer) add(partition *Partition) {

	b.buf.Append(partition)

	// TODO(gvelo): check buffer size and flush if it
	//  is above the configured threshold

}

func (b *SegmentBuffer) flush() {

	b.log.Debug().Int("rows", b.buf.len).Msg("flush")

	segmentSource := NewSegmentSource(b.buf)

	b.segmentWriter.Write(segmentSource)

	b.buf = NewTableBuffer(b.tableName, b.partitionID)

	// reset the timer.
	b.stopTimer()
	b.timer.Reset(b.flushInterval)

}

func (b *SegmentBuffer) stopTimer() {
	if !b.timer.Stop() {
		// drain the timer channel if the timer is expired.
		select {
		case <-b.timer.C:
		default:
		}
	}
}
