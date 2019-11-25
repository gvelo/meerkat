package storage

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/buffer"
	"sync"
)

type SegmentWriter struct {
	inChan chan *buffer.Table
	wg     *sync.WaitGroup
	done   chan struct{}
	log    zerolog.Logger
}

func (sw SegmentWriter) start() {

	defer sw.wg.Done()

	sw.log.Debug().Msg("start")

	for {
		select {
		case t, ok := <-sw.inChan:
			if !ok {
				sw.log.Debug().Msg("stop")
				return
			}
			sw.log.Debug().Msgf("writing segment for index %v with %v rows", t.Index().Name, t.Len())
		case <-sw.done:
			return
		}
	}

}

func NewSegmentWriter(name string, inChan chan *buffer.Table, done chan struct{}, wg *sync.WaitGroup) *SegmentWriter {
	return &SegmentWriter{
		inChan: inChan,
		wg:     wg,
		done:   done,
		log:    log.With().Str("src", "SegmentWriter").Str("name", name).Logger(),
	}
}

type SegmentWriterPool struct {
	wg     *sync.WaitGroup
	inChan chan *buffer.Table
	done   chan struct{}
}

func (sp *SegmentWriterPool) InChan() chan *buffer.Table {
	return sp.inChan
}

func (sp *SegmentWriterPool) Wait() {
	sp.wg.Wait()
}

func NewSegmentWriterPool(chanSize int, workersCount int) *SegmentWriterPool {

	sp := &SegmentWriterPool{
		inChan: make(chan *buffer.Table, chanSize),
		done:   make(chan struct{}),
		wg:     &sync.WaitGroup{},
	}

	for i := 0; i < workersCount; i++ {
		sw := NewSegmentWriter(fmt.Sprintf("SegmentWriter-%v", i), sp.inChan, sp.done, sp.wg)
		sp.wg.Add(1)
		go sw.start()
	}

	return sp

}
