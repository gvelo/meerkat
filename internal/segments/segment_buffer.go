// Copyright 2019 The Meerkat Authors
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

package segments

// TODO(gvelo) move to ingest package

//type SegmentBufferFactory func(indexInfo schema.IndexInfo) *SegmentBuffer
//
//func NewSegmentBufferFactory(ingestChanSize int, t time.Duration, writerChan chan *buffer.Table) SegmentBufferFactory {
//	return func(indexInfo schema.IndexInfo) *SegmentBuffer {
//		return NewSegmentBuffer(indexInfo, ingestChanSize, t, writerChan)
//	}
//}
//
//type SegmentBuffer struct {
//	indexInfo  schema.IndexInfo
//	ingestChan chan *buffer.Table
//	timer      *time.Ticker
//	done       chan struct{}
//	table      *buffer.Table
//	log        zerolog.Logger
//	writerChan chan *buffer.Table
//}
//
//func NewSegmentBuffer(
//	indexInfo schema.IndexInfo,
//	ingestChanSize int,
//	t time.Duration,
//	writerChan chan *buffer.Table) *SegmentBuffer {
//
//	buf := &SegmentBuffer{
//		indexInfo:  indexInfo,
//		ingestChan: make(chan *buffer.Table, ingestChanSize),
//		timer:      time.NewTicker(t),
//		done:       make(chan struct{}),
//		table:      buffer.NewTable(indexInfo),
//		writerChan: writerChan,
//		log: log.With().
//			Str("src", "SegmentBuffer").
//			Str("index-id", indexInfo.Id).
//			Str("index-name", indexInfo.Name).
//			Logger(),
//	}
//
//	go buf.start()
//
//	return buf
//
//}
//
//func (b *SegmentBuffer) start() {
//
//	for {
//		select {
//		case table := <-b.ingestChan:
//			// TODO(gvelo): check for gracefully shutdown.
//			b.addBuffer(table)
//		case <-b.timer.C:
//			if b.table.Len() == 0 {
//				continue
//			}
//			b.writeBuffer()
//		case <-b.done:
//			// TODO(gvelo): check for pending data on gracefully shutdown.
//			return
//		}
//	}
//}
//
//// TODO(gvelo): move this method to the table struct.
//func (b *SegmentBuffer) addBuffer(t *buffer.Table) {
//
//	b.log.Debug().Msgf("add buffer with len %v", t.Len())
//
//	b.table.AppendTable(t)
//
//}
//
//func (b *SegmentBuffer) writeBuffer() {
//
//	select {
//	case b.writerChan <- b.table:
//		b.log.Debug().Msg("segment buffer queued for writing")
//	default:
//		b.writerChan <- b.table
//		b.log.Error().Msg("segment write channel blocked")
//	}
//
//	// TODO(gvelo): pool buffers to avoid allocs.
//	b.table = buffer.NewTable(b.indexInfo)
//
//}
//
//func (b *SegmentBuffer) IngestChan() chan *buffer.Table {
//	return b.ingestChan
//}
