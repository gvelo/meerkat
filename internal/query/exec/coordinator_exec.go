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
	"context"
	"encoding/gob"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/serf/serf"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"io"
	"meerkat/internal/cluster"
	"meerkat/internal/query/execbase"
	"meerkat/internal/query/execpb"
	"meerkat/internal/query/logical"
	"meerkat/internal/query/parser"
	"meerkat/internal/query/physical"
	"meerkat/internal/storage"
	"sync"
)

type nodeManager interface {
	sendCancel(err *execpb.ExecError)
	sendQueryFragments(queryId uuid.UUID, fragments []*logical.Fragment)
	Close()
}

func newNodeManager(
	queryId uuid.UUID,
	execCtx execbase.ExecutionContext,
	connRegistry cluster.ConnRegistry,
) nodeManager {

	manager := &defaultNodeManager{
		queryId: queryId,
	}

	connRegistry.Range(func(nodeName string, conn *grpc.ClientConn) bool {
		client := newNodeClient(execCtx, nodeName, conn)
		manager.nodeClients = append(manager.nodeClients, client)
		return true
	})

	return manager

}

type defaultNodeManager struct {
	queryId     uuid.UUID
	mu          sync.Mutex
	nodeClients []*nodeClient
}

func (d *defaultNodeManager) sendCancel(err *execpb.ExecError) {

	defer d.mu.Unlock()

	cmd := d.buildCancelCmd(err)

	d.mu.Lock()

	for _, client := range d.nodeClients {
		client.sendCancelCmd(cmd)
	}

}

func (d *defaultNodeManager) buildCancelCmd(err *execpb.ExecError) *execpb.ExecCmd {

	return &execpb.ExecCmd{
		Cmd: &execpb.ExecCmd_ExecCancel{
			ExecCancel: &execpb.ExecCancel{
				Error: err,
			},
		},
	}

}

func (d *defaultNodeManager) sendQueryFragments(queryId uuid.UUID, fragments []*logical.Fragment) {

	defer d.mu.Unlock()
	d.mu.Lock()

	msg := d.buildQueryMsg(queryId, fragments)
	fmt.Println("client len :", len(d.nodeClients))
	for _, client := range d.nodeClients {
		fmt.Println("sending query to :", client.nodeName)
		client.sendQueryCmd(msg)
	}

}

func (d *defaultNodeManager) buildQueryMsg(queryId uuid.UUID, fragments []*logical.Fragment) *execpb.ExecCmd {

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(fragments)

	if err != nil {
		panic(fmt.Sprintf("cannot encode query fragments: %v", err))
	}

	return &execpb.ExecCmd{
		Cmd: &execpb.ExecCmd_ExecQuery{
			ExecQuery: &execpb.ExecQuery{
				Id:   queryId[:],
				Plan: buf.Bytes(),
			},
		},
	}

}

func (d *defaultNodeManager) Close() {

	defer d.mu.Unlock()
	d.mu.Lock()

	for _, nodeClient := range d.nodeClients {
		nodeClient.close()
	}

}

func newNodeClient(execCtx execbase.ExecutionContext, nodeName string, conn *grpc.ClientConn) *nodeClient {

	return &nodeClient{
		nodeName:   nodeName,
		execClient: execpb.NewExecutorClient(conn),
		execCtx:    execCtx,
	}

}

type nodeClient struct {
	mu            sync.Mutex
	nodeName      string
	execClient    execpb.ExecutorClient
	controlClient execpb.Executor_ControlClient
	execCtx       execbase.ExecutionContext
}

func (n *nodeClient) sendCancelCmd(cmd *execpb.ExecCmd) {

	go func() {

		defer n.mu.Unlock()
		n.mu.Lock()

		// controlClient is nil if the client has not yet
		// sent the query or there was an error in sendQueryCmd
		// and the client has never been initialized.
		// TODO(gvelo): flag the client as canceled to avoid
		// send a query after a cancel cmd.
		if n.controlClient == nil {
			return
		}

		err := n.controlClient.Send(cmd)

		if err != nil {
			// doesn't make sense to call n.execCtx.Cancel(err) here
			// because the query is being cancel so just log the err
			// TODO(gvelo) log
		}

		err = n.controlClient.CloseSend()

	}()

}

func (n *nodeClient) sendQueryCmd(cmd *execpb.ExecCmd) {

	go func() {

		defer n.mu.Unlock()
		n.mu.Lock()

		controlClient, err := n.execClient.Control(context.TODO())

		if err != nil {

			n.execCtx.CancelWithPropagation(err, execbase.NewExecError(
				fmt.Sprintf("cannot open control stream to node %s [%s]", n.nodeName, err),
				"",
			))

			return
		}

		n.controlClient = controlClient

		go n.handleStream()

		err = n.controlClient.Send(cmd)

		if err != nil {

			n.execCtx.CancelWithPropagation(err, execbase.NewExecError(
				fmt.Sprintf("cannot send query to node %s [%s]", n.nodeName, err),
				"",
			))

			return
		}

	}()

}

func (n *nodeClient) handleStream() {

	for {

		execEvent, err := n.controlClient.Recv()

		if err != nil {

			if err == io.EOF {
				// this is a gracefully control stream shutdown
				return
			}

			execErr := execbase.ExtractExecError(err)

			if execErr == nil {
				execErr = execbase.NewExecError(
					fmt.Sprintf("error reading from control stream : %v", err),
					"", // TODO(gvelo): will be moved to execCtx
				)
			}

			n.execCtx.CancelWithExecError(execErr)

			return
		}

		switch event := execEvent.Event.(type) {
		case *execpb.ExecEvent_ExecOk:
			_ = event
			// do nothing for now, just log
		case *execpb.ExecEvent_ExecStats:
			_ = event
			// do nothing for now, just log
		}

	}

}

func (n *nodeClient) close() {
	go func() {
		_ = n.controlClient.CloseSend()
	}()
}

func NewCoordinatorExecutor(
	connReg cluster.ConnRegistry,
	segReg storage.SegmentRegistry,
	streamReg physical.StreamRegistry,
	cluster cluster.Cluster,
) *coordinatorExecutor {

	id := uuid.New()
	execCtx := execbase.NewExecutionContext()
	nodeManager := newNodeManager(id, execCtx, connReg)

	exec := &coordinatorExecutor{
		id:          id,
		segReg:      segReg,
		conReg:      connReg,
		nodeManager: nodeManager,
		streamReg:   streamReg,
		cluster:     cluster,
		execCtx:     execCtx,
		log: log.With().
			Str("component", "coordinatorExecutor").
			Str("queryId", id.String()).Logger(),
	}

	return exec

}

type coordinatorExecutor struct {
	id          uuid.UUID
	log         zerolog.Logger
	segReg      storage.SegmentRegistry
	conReg      cluster.ConnRegistry
	nodeManager nodeManager
	streamReg   physical.StreamRegistry
	cluster     cluster.Cluster
	execCtx     execbase.ExecutionContext
	outputOp    physical.RunnableOp
}

func (c *coordinatorExecutor) exec(query string, writer execbase.QueryOutputWriter) error {

	defer c.execCtx.Cancel()

	// Transform the string text into a abstract syntax tree
	ast, err := parser.Parse(query)

	if err != nil {
		return err
	}

	// perform the semantic validation
	ast, err = parser.Analyze(ast)

	if err != nil {
		return err
	}

	// transform the ast into a logical query plan
	logicalPlan := logical.ToLogical(ast)

	// optimize the plan
	optPlan := logical.Optimize(logicalPlan)

	// parallelize
	fragments := logical.Parallelize(
		optPlan,
		c.cluster.NodeName(),
		buildNodeNames(c.cluster.LiveMembers()),
	)

	go c.handleCtxCancel()

	go c.sendFragmentsToNodes(fragments.NodeFragments())

	dag, err := c.buildDAG(writer, fragments.AllFragments())

	if err != nil {
		c.Cancel(err)
		return err
	}

	dag.Run()

	return nil

}

func buildNodeNames(members []serf.Member) []string {
	nodeNames := make([]string, len(members))
	for i, member := range members {
		nodeNames[i] = member.Name
	}
	return nodeNames
}

func (c *coordinatorExecutor) buildDAG(
	writer execbase.QueryOutputWriter,
	fragments []*logical.Fragment,
) (physical.DAG, error) {

	dagBuilder := physical.NewDAGBuilder(
		c.conReg,
		c.segReg,
		c.streamReg,
		c.cluster.NodeName(),
	)

	dag, err := dagBuilder.BuildDAG(fragments, c.id, writer, c.execCtx)

	if err != nil {
		return nil, err
	}

	return dag, nil

}

func (c *coordinatorExecutor) sendFragmentsToNodes(fragments []*logical.Fragment) {
	c.nodeManager.sendQueryFragments(c.id, fragments)
}

func (c *coordinatorExecutor) handleCtxCancel() {

	<-c.execCtx.Done()

	if c.execCtx.Err() != nil {
		c.nodeManager.sendCancel(c.execCtx.Err())
	}

}

func (c *coordinatorExecutor) Cancel(err error) {
	execError := execbase.NewExecError(err.Error(), c.cluster.NodeName())
	c.execCtx.CancelWithExecError(execError)
}
