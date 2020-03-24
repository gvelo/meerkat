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

import "meerkat/internal/buffer"

func NewSegmentWriterPool(chanSize int, poolSize int) *SegmentWriterPool {
	return &SegmentWriterPool{
		poolSize: poolSize,
		inChan:   make(chan *buffer.Table, chanSize),
		done:     make(chan struct{}),
	}
}

type SegmentWriterPool struct {
	poolSize int
	inChan   chan *buffer.Table
	done     chan struct{}
}

func (s *SegmentWriterPool) Start() {

}

func (s *SegmentWriterPool) Stop() {

}

func (s *SegmentWriterPool) InChan() chan *buffer.Table {
	return s.inChan
}
