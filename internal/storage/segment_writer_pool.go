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
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"log"
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
	}
}

type SegmentWriterPool struct {
	poolSize int
	inChan   chan *buffer.Table
	done     chan struct{}
	path     string
}

func (s *SegmentWriterPool) Start() {

	s.path = path.Join(s.path, "segments")

	err := os.MkdirAll(s.path, 0770)

	if err != nil {
		panic(fmt.Sprintf("Could not create dir %v , %v", s.path, err))
	}

	for i := 0; i < s.poolSize; i++ {

		worker := segmentWriterWorker{
			id:     i,
			inChan: s.inChan,
			done:   s.done,
			path:   s.path,
		}
		go func() { worker.Start() }()
	}

}

func (s *SegmentWriterPool) InChan() chan *buffer.Table {
	return s.inChan
}

type segmentWriterWorker struct {
	id     int
	inChan chan *buffer.Table
	done   chan struct{}
	path   string
}

// start worker
func (w segmentWriterWorker) Start() {

	for {
		select {
		case table := <-w.inChan:
			w.writeTable(table)
		case <-w.done:
			return
		}
	}
}

func (w segmentWriterWorker) writeTable(t *buffer.Table) {
	// TODO: meter en la config.

	viper.Get("")

	sid := uuid.New()
	sw := NewSegmentWriter(w.path, sid, t)

	if err := sw.Write(); err != nil {
		// TODO: (sebad) what to do in this case?
		log.Printf("Error writing segment %v %v", sid, err)
	}

}
