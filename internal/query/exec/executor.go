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
	"github.com/pkg/errors"
	"meerkat/internal/cluster"
	"meerkat/internal/query/execbase"
	"meerkat/internal/query/execpb"
	"meerkat/internal/query/physical"
	"meerkat/internal/storage"
)

type Executor interface {
	ExecuteQuery(query string, writer execbase.QueryOutputWriter) error
	CancelQuery(queryId uuid.UUID)
	Stop()
}

func NewExecutor(
	segReg storage.SegmentRegistry,
	streamReg physical.StreamRegistry,
	cluster cluster.Cluster,
) Executor {

	return &executor{
		segReg:    segReg,
		cluster:   cluster,
		streamReg: streamReg,
	}

}

type executor struct {
	segReg    storage.SegmentRegistry
	streamReg physical.StreamRegistry
	cluster   cluster.Cluster
}

func (e executor) ExecuteQuery(query string, writer execbase.QueryOutputWriter) error {

	coordinator := NewCoordinatorExecutor(e.segReg, e.streamReg, e.cluster)

	return coordinator.exec(query, writer)

}

func (e *executor) CancelQuery(queryId uuid.UUID) {
	panic("implement me")
}

func (e *executor) Stop() {
	panic("implement me")
}

func NewServer(cluster cluster.Cluster, segReg storage.SegmentRegistry, streamReg physical.StreamRegistry, localNodeName string) *Server {
	return &Server{
		streamReg:   streamReg,
		cluster:     cluster,
		segReg:      segReg,
		localNodeId: localNodeName,
	}
}

type Server struct {
	streamReg   physical.StreamRegistry
	cluster     cluster.Cluster
	segReg      storage.SegmentRegistry
	localNodeId string
}

func (s *Server) Control(controlSrv execpb.Executor_ControlServer) error {

	fmt.Println("new control stream")

	nodeExec := NewNodeExec(s.cluster, s.segReg, s.streamReg, controlSrv, s.localNodeId)
	nodeExec.Start()

	<-nodeExec.ExecutionContext().Done()

	execErr := nodeExec.ExecutionContext().Err()

	if execErr != nil {
		return execbase.BuildGRPCError(execErr)
	}

	// TODO(gvelo): check if we need to cancel the execCtx before return.
	return nil

}

func (s *Server) VectorExchange(vectorServer execpb.Executor_VectorExchangeServer) error {

	msg, err := vectorServer.Recv()

	if err != nil {
		return err
	}

	headerMsg, ok := msg.Msg.(*execpb.VectorExchangeMsg_Header)

	if !ok {
		return errors.New("invalid stream header")
	}

	queryId, err := uuid.FromBytes(headerMsg.Header.QueryId)

	if err != nil {
		return err
	}

	ch := make(chan *execpb.ExecError, 1)

	stream := physical.NewVectorExchangeStream(vectorServer, ch)

	err = s.streamReg.RegisterStream(queryId, headerMsg.Header.StreamId, stream)

	if err != nil {
		return err
	}

	execErr := <-ch

	if execErr != nil {
		return execbase.BuildGRPCError(execErr)
	}

	return nil

}
