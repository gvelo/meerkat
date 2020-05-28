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

package ingestion

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

// BufferRegistry is a registry of ingestion buffers.
type BufferRegistry interface {
	AddToBuffer(table *Table)
}

func NewBufferRegistry(inChSize int,
	evictIdleTime time.Duration,
	segmentMaxSize int,
	segmentFlushInterval time.Duration,
	segmentChSize int) BufferRegistry {
	return &bufferRegistry{
		buffers:              make(map[bufferKey]bufferEntry),
		inCh:                 make(chan *Table, inChSize),
		evictIdleTime:        evictIdleTime,
		segmentMaxSize:       segmentMaxSize,
		segmentFlushInterval: segmentFlushInterval,
		segmentChSize:        segmentChSize,
		log:                  log.With().Str("src", "BufferRegistry").Logger(),
	}
}

type bufferKey struct {
	tableName   string
	partitionID uint64
}

type bufferEntry struct {
	buffer     *SegmentBuffer
	lastAccess time.Time
}

type bufferRegistry struct {
	buffers              map[bufferKey]bufferEntry
	inCh                 chan *Table
	wg                   *sync.WaitGroup
	mu                   sync.Mutex
	running              bool
	log                  zerolog.Logger
	evictTicker          *time.Ticker
	evictIdleTime        time.Duration
	segmentMaxSize       int
	segmentFlushInterval time.Duration
	segmentChSize        int
}

func (b *bufferRegistry) Start() {

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		return
	}

	b.log.Info().Msg("start")
	b.running = true
	b.evictTicker = time.NewTicker(10 * time.Second)
	go b.run()

}

func (b *bufferRegistry) Stop() {

	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.running {
		return
	}

	b.log.Info().Msg("stop")

	b.running = false
	b.evictTicker.Stop()
	close(b.inCh)
	b.wg.Wait()

	b.log.Info().Msg("stopped")

}

func (b *bufferRegistry) run() {

	for {

		select {

		case table, done := <-b.inCh:

			if done {
				b.closeBuffers()
				return
			}

			b.add(table)

		case <-b.evictTicker.C:
			b.evictIdle()
		}

	}

}

func (b *bufferRegistry) add(table *Table) {

	for _, partition := range table.Partitions {

		fmt.Println("=========================")
		fmt.Println(table.Name)
		fmt.Println(table.Partitions)
		fmt.Println("=========================")

		key := bufferKey{
			tableName:   table.Name,
			partitionID: partition.Id,
		}

		segmentBuffer := b.getEntry(key).buffer

		segmentBuffer.add(partition)

	}

}

func (b *bufferRegistry) closeBuffers() {

	for key, entry := range b.buffers {
		delete(b.buffers, key)
		entry.buffer.Stop()
	}

}

func (b *bufferRegistry) evictIdle() {

	now := time.Now()

	for key, entry := range b.buffers {
		if now.Sub(entry.lastAccess) > b.evictIdleTime {
			delete(b.buffers, key)
			entry.buffer.Stop()
		}
	}

}

func (b bufferRegistry) AddToBuffer(table *Table) {

	b.inCh <- table

}

func (b *bufferRegistry) getEntry(bufKey bufferKey) bufferEntry {

	if entry, found := b.buffers[bufKey]; found {
		entry.lastAccess = time.Now()
		return entry
	}

	entry := bufferEntry{
		buffer: NewSegmentBuffer(b.segmentMaxSize,
			b.segmentFlushInterval,
			b.segmentChSize,
			b.wg,
			bufKey.tableName,
			bufKey.partitionID),
		lastAccess: time.Now(),
	}

	b.buffers[bufKey] = entry

	return entry
}
