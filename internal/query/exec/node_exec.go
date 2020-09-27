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
	"meerkat/internal/cluster"
	"meerkat/internal/storage"
)

type nodeExec struct {
	connReg cluster.ConnRegistry
	segReg  storage.SegmentRegistry
	execCtx ExecutionContext
}

func NewNodeExec(connReg cluster.ConnRegistry, segReg storage.SegmentRegistry, controlSrv Executor_ControlServer) *nodeExec {
	return &nodeExec{
		connReg: connReg,
		segReg:  segReg,
		execCtx: NewExecutionContext(),
	}
}

func (n *nodeExec) Run() {

}

func (n *nodeExec) ExecutionContext() ExecutionContext {
	return n.execCtx
}
