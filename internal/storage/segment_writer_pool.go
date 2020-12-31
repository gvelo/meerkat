// Copyright 2020 The Meerkat Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type SegmentWriter interface {
	Write(src SegmentSource)
}

type SegmentWriterPool struct {
	poolSize    int
	inChan      chan SegmentSource
	wg          sync.WaitGroup
	mu          sync.Mutex
	isRunning   bool
	segStorage  SegmentStorage
	segRegistry SegmentRegistry
	log         zerolog.Logger
}

func NewSegmentWriterPool(
	chanSize int,
	poolSize int,
	segStore SegmentStorage,
	segRegistry SegmentRegistry,
) *SegmentWriterPool {
	return &SegmentWriterPool{
		poolSize:    poolSize,
		inChan:      make(chan SegmentSource, chanSize),
		segStorage:  segStore,
		segRegistry: segRegistry,
		log:         log.With().Str("src", "segmentWriterPool").Logger(),
	}
}

func (w *SegmentWriterPool) Start() {

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isRunning {
		return
	}

	w.log.Info().Msg("start")

	w.isRunning = true

	for i := 0; i < w.poolSize; i++ {

		worker := &segmentWriterWorker{
			id:              i,
			inChan:          w.inChan,
			wg:              &w.wg,
			segStorage:      w.segStorage,
			segmentRegistry: w.segRegistry,
			log:             log.With().Str("src", "SegmentWriterWorker").Int("id", i).Logger(),
		}

		w.wg.Add(1)
		go worker.start()

	}

}

func (w *SegmentWriterPool) Stop() {

	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.isRunning {
		return
	}

	w.log.Info().Msg("stopping segment writer pool")

	w.isRunning = false

	close(w.inChan)
	w.log.Info().Msg("waiting for workers to finish")

	w.wg.Wait()

	w.log.Info().Msg("stopped")

}

func (w *SegmentWriterPool) Write(src SegmentSource) {
	w.inChan <- src
}

type segmentWriterWorker struct {
	id              int
	inChan          chan SegmentSource
	wg              *sync.WaitGroup
	segmentRegistry SegmentRegistry
	segStorage      SegmentStorage
	log             zerolog.Logger
}

// start worker
func (w *segmentWriterWorker) start() {

	defer w.wg.Done()

	w.log.Info().Msg("started")

	for {
		select {
		case src, ok := <-w.inChan:
			if !ok {
				w.log.Info().Msg("inChan closed, quiting")
				return
			}
			w.writeSegment(src)
		}
	}

}

func (w *segmentWriterWorker) writeSegment(src SegmentSource) {

	defer func() {
		if r := recover(); r != nil {
			w.log.
				Error().
				Interface("error", r).
				Str("sid", buildSegmentFileNameFromUUID(src.Info().Id)).
				Msg("error writing segment")
		}
	}()

	info := src.Info()

	logger := w.log.With().
		Str("sid", buildSegmentFileNameFromUUID(info.Id)).
		Str("table", info.TableName).
		Uint64("partition", info.PartitionId).Logger()

	logger.Debug().Msg("writing segment")

	segInfo := w.segStorage.WriteSegment(src)

	w.segmentRegistry.AddSegment(segInfo)

	logger.Debug().Msg("segment successfully written")

}
