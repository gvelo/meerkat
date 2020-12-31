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
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/google/uuid"
	"io"
	"meerkat/internal/cluster"
	"meerkat/internal/query/logical"
	"meerkat/internal/query/physical"
	"meerkat/internal/storage"
	"sync"
)

type nodeExec struct {
	connReg    cluster.ConnRegistry
	segReg     storage.SegmentRegistry
	execCtx    ExecutionContext
	controlSrv Executor_ControlServer
	outputOps  []physical.OutputOp
	outputWg   *sync.WaitGroup
	streamReg  StreamRegistry
}

func NewNodeExec(
	connReg cluster.ConnRegistry,
	segReg storage.SegmentRegistry,
	streamReg StreamRegistry,
	controlSrv Executor_ControlServer,
) *nodeExec {
	return &nodeExec{
		connReg:    connReg,
		segReg:     segReg,
		controlSrv: controlSrv,
		execCtx:    NewExecutionContext(),
		streamReg:  streamReg,
	}
}

func (n *nodeExec) Start() {
	go n.handleControlStream()
}

func (n *nodeExec) handleControlStream() {

	fmt.Println("handleControlStream()")

	for {

		execCmd, err := n.controlSrv.Recv()

		fmt.Println("n.controlSrv.Recv()", execCmd, err)

		if err == io.EOF {
			return
		}

		if err != nil {
			n.execCtx.CancelWithPropagation(err,
				newExecError(fmt.Sprintf("error receiving on control stream: %v  ", err)))
			return
		}

		switch cmd := execCmd.Cmd.(type) {
		case *ExecCmd_ExecQuery:
			go n.execQuery(cmd.ExecQuery)
		case *ExecCmd_ExecCancel:
			n.execCtx.CancelWithExecError(cmd.ExecCancel.Error)
			return
		default:
			n.execCtx.CancelWithExecError(newExecError("unknown exec command"))
			return

		}
	}
}

func (n *nodeExec) ExecutionContext() ExecutionContext {
	return n.execCtx
}

func (n *nodeExec) execQuery(query *ExecQuery) {

	fmt.Println("execQuery(query *ExecQuery)", query)

	id, err := uuid.FromBytes(query.Id)

	fmt.Println("query.id = ", id, err)

	if err != nil {
		n.execCtx.CancelWithExecError(
			newExecError(fmt.Sprintf("cannot unmarshal query id: %v", err)))
		return
	}

	fragments, err := decodePlan(query)

	fmt.Println("decodePlan(query)", fragments, err)

	if err != nil {
		n.execCtx.CancelWithExecError(
			newExecError(fmt.Sprintf("cannot unmarshal query plan on query: %v : %v", id, err)))
		return
	}

	err = n.buildExecutableGraph(fragments)

	if err != nil {
		n.execCtx.CancelWithExecError(
			newExecError(fmt.Sprintf("cannot build ejecutable graph on query: %v : %v", id, err)))
		return
	}

	n.run()

	n.execCtx.Cancel()

}

func (n *nodeExec) run() {
	n.outputWg.Add(len(n.outputOps))
	for _, op := range n.outputOps {
		go op.Run()
	}
	n.outputWg.Wait()
}

func decodePlan(query *ExecQuery) ([]*logical.Fragment, error) {

	reader := bytes.NewReader(query.Plan)

	decoder := gob.NewDecoder(reader)

	var fragments []*logical.Fragment

	err := decoder.Decode(&fragments)

	return fragments, err

}

func (n *nodeExec) buildExecutableGraph(fragments []*logical.Fragment) error {

	builder := NewExecutableGraphBuilder(n.connReg, n.segReg, n.streamReg)

	var err error

	n.outputOps, err = builder.BuildNodeGraph(fragments, n.outputWg)

	if err != nil {
		return err
	}

	return nil

}
