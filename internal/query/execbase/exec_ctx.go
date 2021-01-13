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

package execbase

import (
	"meerkat/internal/query/execpb"
	"sync"
)

type ExecutionContext interface {
	Done() <-chan struct{}
	Cancel()
	CancelWithExecError(execError *execpb.ExecError)
	CancelWithPropagation(err error, execError *execpb.ExecError)
	Err() *execpb.ExecError
}

func NewExecutionContext() ExecutionContext {
	return &executionContext{
		done: make(chan struct{}),
	}
}

type executionContext struct {
	mu        sync.Mutex
	execError *execpb.ExecError
	done      chan struct{}
	canceled  bool
}

func (c *executionContext) CancelWithExecError(execError *execpb.ExecError) {

	defer c.mu.Unlock()

	c.mu.Lock()

	if c.canceled {
		return
	}

	c.canceled = true
	c.execError = execError

	close(c.done)

}

func (c *executionContext) CancelWithPropagation(err error, execError *execpb.ExecError) {

	e := ExtractExecError(err)

	if e == nil {
		e = execError
	}

	c.CancelWithExecError(e)

}

func (c *executionContext) Cancel()                { c.CancelWithExecError(nil) }
func (c *executionContext) Done() <-chan struct{}  { return c.done }
func (c *executionContext) Err() *execpb.ExecError { return c.execError }
