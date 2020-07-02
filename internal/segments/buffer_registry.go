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

const (
	eventHandler = "segment.buffer.registry.handler"
)

// TODO(gvelo) move to ingest package
//type SegmentBufferRegistry struct {
//	buffers              sync.Map
//	schemaChan           chan schema.IndexUpdateEvent
//	done                 chan struct{}
//	segmentBufferFactory SegmentBufferFactory
//	log                  zerolog.Logger
//}
//
//func NewSegmentBufferRegistry(sch schema.Schema, factory SegmentBufferFactory) *SegmentBufferRegistry {
//
//	reg := &SegmentBufferRegistry{
//		schemaChan:           make(chan schema.IndexUpdateEvent, 256),
//		segmentBufferFactory: factory,
//		log:                  log.With().Str("src", "SegmentBufferRegistry").Logger(),
//	}
//
//	sch.AddEventHandler(eventHandler, reg.schemaChan)
//
//	for _, i := range sch.AllIndex() {
//		reg.addIndex(i)
//	}
//
//	go reg.handleSchemaEvents()
//
//	return reg
//}
//
//func (sbr *SegmentBufferRegistry) handleSchemaEvents() {
//
//	for {
//		select {
//		case e := <-sbr.schemaChan:
//			switch e.OpType {
//			case schema.OpIndexUpdate:
//				// TODO(gvelo): implement schema mutation.
//				sbr.log.Debug().Str("indexID", e.IndexInfo.Id).Msg("updating index")
//			case schema.OpIndexCreate:
//				sbr.addIndex(e.IndexInfo)
//			case schema.OpIndexDelete:
//				// TODO(gvelo): should we shutdown the Buffer ?
//				sbr.buffers.Delete(e.IndexInfo.Id)
//			}
//		case <-sbr.done:
//			log.Info().Msg("shutting down")
//			return
//		}
//	}
//}
//
//func (sbr *SegmentBufferRegistry) addIndex(index schema.IndexInfo) {
//	sbr.log.Debug().Str("indexID", index.Id).Msg("adding new index")
//	sb := sbr.segmentBufferFactory(index)
//	sbr.buffers.Store(index.Id, sb)
//}
//
//func (sbr *SegmentBufferRegistry) Buffer(id string) *SegmentBuffer {
//
//	b, ok := sbr.buffers.Load(id)
//
//	if !ok {
//		return nil
//	}
//
//	return b.(*SegmentBuffer)
//
//}
