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
	"meerkat/internal/query/execbase"
	"meerkat/internal/query/execpb"
	"meerkat/internal/query/logical"
	"meerkat/internal/query/physical"
	"meerkat/internal/storage"
)

type nodeExec struct {
	cluster     cluster.Cluster
	segReg      storage.SegmentRegistry
	execCtx     execbase.ExecutionContext
	controlSrv  execpb.Executor_ControlServer
	dag         physical.DAG
	streamReg   physical.StreamRegistry
	localNodeId string
}

func NewNodeExec(
	cluster cluster.Cluster,
	segReg storage.SegmentRegistry,
	streamReg physical.StreamRegistry,
	controlSrv execpb.Executor_ControlServer,
	localNodeId string,
) *nodeExec {
	return &nodeExec{
		cluster:     cluster,
		segReg:      segReg,
		controlSrv:  controlSrv,
		execCtx:     execbase.NewExecutionContext(),
		streamReg:   streamReg,
		localNodeId: localNodeId,
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
				execbase.NewExecError(
					fmt.Sprintf("error receiving on control stream: %v  ", err),
					"",
				))
			return
		}

		switch cmd := execCmd.Cmd.(type) {
		case *execpb.ExecCmd_ExecQuery:
			go n.execQuery(cmd.ExecQuery)
		case *execpb.ExecCmd_ExecCancel:
			n.execCtx.CancelWithExecError(cmd.ExecCancel.Error)
			return
		default:
			n.execCtx.CancelWithExecError(execbase.NewExecError("unknown exec command", ""))
			return

		}
	}
}

func (n *nodeExec) ExecutionContext() execbase.ExecutionContext {
	return n.execCtx
}

func (n *nodeExec) execQuery(query *execpb.ExecQuery) {

	fmt.Println("execQuery(query *ExecQuery)", query)

	id, err := uuid.FromBytes(query.Id)

	fmt.Println("query.id = ", id, err)

	if err != nil {
		n.execCtx.CancelWithExecError(
			execbase.NewExecError(fmt.Sprintf("cannot unmarshal query id: %v", err), n.localNodeId))
		return
	}

	fragments, err := decodePlan(query)

	fmt.Println("decodePlan(query)", fragments, err)

	if err != nil {
		n.execCtx.CancelWithExecError(
			execbase.NewExecError(fmt.Sprintf("cannot unmarshal query plan on query: %v : %v", id, err), n.localNodeId))
		return
	}

	err = n.buildDAG(fragments, query.Id)

	if err != nil {
		n.execCtx.CancelWithExecError(
			execbase.NewExecError(fmt.Sprintf("cannot build ejecutable graph on query: %v : %v", id, err), n.localNodeId))
		return
	}

	n.dag.Run()

	n.execCtx.Cancel()

}

func decodePlan(query *execpb.ExecQuery) ([]*logical.Fragment, error) {

	reader := bytes.NewReader(query.Plan)

	decoder := gob.NewDecoder(reader)

	var fragments []*logical.Fragment

	err := decoder.Decode(&fragments)

	return fragments, err

}

func (n *nodeExec) buildDAG(fragments []*logical.Fragment, queryId []byte) error {

	id, err := uuid.FromBytes(queryId)

	if err != nil {
		return err
	}

	builder := physical.NewDAGBuilder(n.cluster, n.segReg, n.streamReg, n.localNodeId)

	n.dag, err = builder.BuildDAG(fragments, id, nil, n.execCtx)

	if err != nil {
		return err
	}

	return nil

}
