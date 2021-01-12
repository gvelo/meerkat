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

package exec

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

const (
	streamReadTimeout = 2 * time.Second
)

type VectorExchangeStream interface {
	Server() Executor_VectorExchangeServer
	Close(err *ExecError)
}

type vectorExchangeStream struct {
	server Executor_VectorExchangeServer
	done   chan *ExecError
	closed bool
	mu     sync.Mutex
}

func NewVectorExchangeStream(server Executor_VectorExchangeServer, done chan *ExecError) VectorExchangeStream {
	return &vectorExchangeStream{
		server: server,
		done:   done,
	}
}

func (v *vectorExchangeStream) Server() Executor_VectorExchangeServer {
	return v.server
}

func (v *vectorExchangeStream) Close(err *ExecError) {

	v.mu.Lock()
	defer v.mu.Unlock()

	if v.closed {
		return
	}

	v.closed = true
	v.done <- err
	close(v.done)
}

type StreamRegistry interface {
	GetStream(queryId uuid.UUID, streamId int64) (*vectorExchangeStream, error)
	RegisterStream(queryId uuid.UUID, streamId int64, stream *vectorExchangeStream) error
}

type streamKey struct {
	queryId  uuid.UUID
	streamId int64
}

func (s *streamKey) String() string {
	return fmt.Sprintf("queryId: %v streamId: %v", s.queryId, s.streamId)
}

type streamRegistry struct {
	registry map[streamKey]chan *vectorExchangeStream
	mu       sync.Mutex
}

func NewStreamRegistry() StreamRegistry {
	return &streamRegistry{
		registry: make(map[streamKey]chan *vectorExchangeStream),
	}
}

func (s *streamRegistry) GetStream(queryId uuid.UUID, streamId int64) (*vectorExchangeStream, error) {

	key := streamKey{
		queryId:  queryId,
		streamId: streamId,
	}

	s.mu.Lock()

	ch, found := s.registry[key]

	if !found {
		ch = make(chan *vectorExchangeStream)
		s.registry[key] = ch
	}

	s.mu.Unlock()

	timeOutTimer := time.NewTimer(streamReadTimeout)

	select {
	case stream := <-ch:
		if !timeOutTimer.Stop() {
			<-timeOutTimer.C
		}
		delete(s.registry, key)
		return stream, nil
	case <-timeOutTimer.C:
		delete(s.registry, key)
		return nil, fmt.Errorf("timeout waiting for stream %v ", key)
	}

}

func (s *streamRegistry) RegisterStream(queryId uuid.UUID, streamId int64, stream *vectorExchangeStream) error {

	key := streamKey{
		queryId:  queryId,
		streamId: streamId,
	}

	s.mu.Lock()

	ch, found := s.registry[key]

	if !found {
		ch = make(chan *vectorExchangeStream)
		s.registry[key] = ch
	}

	s.mu.Unlock()

	timeOutTimer := time.NewTimer(streamReadTimeout)

	select {
	case ch <- stream:
		if !timeOutTimer.Stop() {
			<-timeOutTimer.C
		}
		return nil
	case <-timeOutTimer.C:
		delete(s.registry, key)
		return fmt.Errorf("timeout waiting for query execution %v", key)
	}

}
