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
	"encoding/base64"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"sync"
)

const segmentFolderName = "segments"

type SegmentWriter interface {
	Write(src SegmentSource)
}

func NewSegmentWriterPool(chanSize int, poolSize int, dbPath string) *SegmentWriterPool {
	return &SegmentWriterPool{
		poolSize: poolSize,
		inChan:   make(chan SegmentSource, chanSize),
		path:     path.Join(dbPath, segmentFolderName),
		log:      log.With().Str("src", "segmentWriterPool").Logger(),
	}
}

type SegmentWriterPool struct {
	poolSize int
	inChan   chan SegmentSource
	wg       sync.WaitGroup
	mu       sync.Mutex
	running  bool
	path     string
	log      zerolog.Logger
}

func (w *SegmentWriterPool) Start() error {

	w.mu.Lock()
	defer w.mu.Unlock()

	if w.running {
		return nil
	}

	w.log.Info().Msg("start")

	w.running = true

	err := os.MkdirAll(w.path, 0770)

	if err != nil {
		return fmt.Errorf("could not create dir %v , %v", w.path, err)
	}

	for i := 0; i < w.poolSize; i++ {

		worker := &segmentWriterWorker{
			id:     i,
			inChan: w.inChan,
			wg:     &w.wg,
			path:   w.path,
			log:    log.With().Str("src", "SegmentWriterWorker").Int("id", i).Logger(),
		}

		w.wg.Add(1)
		go func() { worker.start() }()

	}

	return nil

}

func (w *SegmentWriterPool) Stop() {

	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.running {
		return
	}

	w.log.Info().Msg("stopping segment writer pool")

	w.running = false

	close(w.inChan)
	w.log.Info().Msg("waiting for workers to finish")

	w.wg.Wait()

	w.log.Info().Msg("stopped")

}

func (w *SegmentWriterPool) Write(src SegmentSource) {
	w.inChan <- src
}

type segmentWriterWorker struct {
	id     int
	inChan chan SegmentSource
	wg     *sync.WaitGroup
	path   string
	log    zerolog.Logger
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

	info := src.Info()

	fileName := base64.RawURLEncoding.EncodeToString(info.Id)

	filePath := path.Join(w.path, fileName)

	logger := w.log.With().
		Str("sid", fileName).
		Str("table", info.TableName).
		Uint64("partition", info.PartitionId).Logger()

	logger.Debug().Msg("writing segment")

	err := WriteSegment(filePath, src)

	if err != nil {
		logger.Error().Err(err).Msg("error writing segment")
		return
	}

	logger.Debug().Msg("segment successfully written")

}
