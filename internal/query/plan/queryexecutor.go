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

package plan

import (
	"github.com/RoaringBitmap/roaring"
	"meerkat/internal/query/logical"
	"meerkat/internal/tools"
	"sync"
	"time"
)

type Executor interface {
	ExecuteQuery(t *logical.Projection) *ResultSet
}

type MeerkatExecutor struct {
}

func NewMeerkatExecutor() *MeerkatExecutor {
	return &MeerkatExecutor{}
}

func (e *MeerkatExecutor) ExecuteQuery(t *logical.Projection) *ResultSet {

	//e.exe(t.IndexScan.GetFilter())

	return nil
}

func (e *MeerkatExecutor) exe(f interface{}) *roaring.Bitmap {

	var jobQueue = make(chan interface{})

	var numWorkers = 10
	// the waitgroup will allow us to wait for all the goroutines to finish at the end
	var wg = new(sync.WaitGroup)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, wg, jobQueue)
	}

	e.executeFilters(jobQueue, f)

	// closing jobQueue will cause all goroutines to exit the loop on the channel.
	close(jobQueue)
	// Wait for all the goroutines to finish
	wg.Wait()

	return nil
}

func worker(i int, wg *sync.WaitGroup, in <-chan interface{}) {

	for exp := range in {

		time.Sleep(time.Millisecond * 500)

		switch v := exp.(type) {
		case *logical.Filter:
			tools.Logf(" %s ", v)
			// its a comp, need to execute

		case *logical.Exp:
			// tools.Logf("Exp %s ", v.Value)
		}

		//tools.Logf("WNode %d ", exp.(rel.Node).String())

		if exp.(logical.Node).GetParent() != nil {
			tools.Logf("GetParent %d ", exp.(logical.Node).GetParent().String())
		} else {
			tools.Log("GetParent nil ")
		}

	}
	wg.Done()
}

func (e *MeerkatExecutor) executeFilters(jobQueue chan interface{}, f interface{}) *roaring.Bitmap {
	if f == nil {
		return nil
	}

	switch v := f.(type) {
	case *logical.Filter:
		// its a comp, need to execute
		e.executeFilters(jobQueue, f.(*logical.Filter).Left)
		e.executeFilters(jobQueue, f.(*logical.Filter).Right)

		f.(*logical.Filter).Left.(logical.Node).SetParent(f.(logical.Node))
		f.(*logical.Filter).Right.(logical.Node).SetParent(f.(logical.Node))

		jobQueue <- v

	case *logical.Exp:
		// tools.Logf("Exp %s ", v.Value)
	}
	return nil
}
