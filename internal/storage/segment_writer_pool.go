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
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"meerkat/internal/buffer"
	"os"
	"path"
)

func NewSegmentWriterPool(chanSize int, poolSize int, dbPath string) *SegmentWriterPool {
	return &SegmentWriterPool{
		poolSize: poolSize,
		inChan:   make(chan *buffer.Table, chanSize),
		done:     make(chan struct{}),
		path:     dbPath,
		log:      log.With().Str("src", "SegmentWriterPool").Logger(),
	}
}

type SegmentWriterPool struct {
	poolSize int
	inChan   chan *buffer.Table
	done     chan struct{}
	path     string
	log      zerolog.Logger
}

func (s *SegmentWriterPool) Start() error {

	s.log.Info().Msg("start")

	s.path = path.Join(s.path, "segments")

	err := os.MkdirAll(s.path, 0770)

	if err != nil {
		return fmt.Errorf("could not create dir %v , %v", s.path, err)
	}

	for i := 0; i < s.poolSize; i++ {

		worker := &segmentWriterWorker{
			id:     i,
			inChan: s.inChan,
			done:   s.done,
			path:   s.path,
			log:    log.With().Str("src", "SegmentWriterWorker").Int("id", i).Logger(),
		}

		go func() { worker.Start() }()

	}

	return nil

}

func (s *SegmentWriterPool) InChan() chan *buffer.Table {
	return s.inChan
}

type segmentWriterWorker struct {
	id     int
	inChan chan *buffer.Table
	done   chan struct{}
	path   string
	log    zerolog.Logger
}

// start worker
func (w *segmentWriterWorker) Start() {

	for {
		select {
		case table, ok := <-w.inChan:
			if !ok {
				w.log.Info().Msg("inChan closed, quiting")
				return
			}
			w.writeTable(table)
		case <-w.done:
			w.log.Info().Msg("done")
			return
		}
	}
}

func (w *segmentWriterWorker) writeTable(t *buffer.Table) {

	//sid := uuid.New()
	//fileName := filepath.Join(w.path, sid.String())
	//
	//w.log.Debug().Str("sid", sid.String()).Msg("writing segment")
	//
	//err := WriteSegment(fileName, sid, t)
	//
	//if err != nil {
	//	// TODO: (sebad) what to do in this case?
	//	w.log.Error().Err(err).Str("sid", sid.String()).Msg("error writing segment")
	//}

}
