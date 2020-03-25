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
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSegmentWriterPool(t *testing.T) {

	indexInfo := createIndexInfo()
	buf := createBuffers(indexInfo)

	p := NewSegmentWriterPool(10, 2, "/Users/sebad/meerkat")
	p.Start()

	p.inChan <- buf

	// Flaky test.
	timer := time.AfterFunc(time.Millisecond*1000, func() {
		close(p.done)
	})
	defer timer.Stop()

	select {
	case <-p.done: // worker has received job
	}

	assert.Equal(t, 0, len(p.inChan), "This channel should be empty")

}
