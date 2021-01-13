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
	"meerkat/internal/storage"
)

//go:generate protoc -I . -I ../../../build/proto/ -I ../../../internal/storage/ --plugin ../../../build/protoc-gen-gogofaster --gogofaster_out=plugins=grpc,paths=source_relative:.  ./exec.proto

type QueryOutputWriter interface {
	Write([]byte) (int, error)
	Flush()
	CloseNotify() <-chan bool
}

type Executor interface {
	ExecuteQuery(query string, writer QueryOutputWriter) error
	CancelQuery(queryId uuid.UUID)
	Stop()
}

func NewExecutor(
	connReg cluster.ConnRegistry,
	segReg storage.SegmentRegistry,
	streamReg StreamRegistry,
	cluster cluster.Cluster,
) Executor {

	return &executor{
		connReg:   connReg,
		segReg:    segReg,
		cluster:   cluster,
		streamReg: streamReg,
	}

}

type executor struct {
	connReg   cluster.ConnRegistry
	segReg    storage.SegmentRegistry
	streamReg StreamRegistry
	cluster   cluster.Cluster
}

func (e executor) ExecuteQuery(query string, writer QueryOutputWriter) error {

	coordinator := NewCoordinatorExecutor(e.connReg, e.segReg, e.streamReg, e.cluster)

	return coordinator.exec(query, writer)

}

func (e *executor) CancelQuery(queryId uuid.UUID) {
	panic("implement me")
}

func (e *executor) Stop() {
	panic("implement me")
}

func NewServer(connReg cluster.ConnRegistry, segReg storage.SegmentRegistry, streamReg StreamRegistry, localNodeName string) *Server {
	return &Server{
		streamReg:     streamReg,
		connReg:       connReg,
		segReg:        segReg,
		localNodeName: localNodeName,
	}
}

type Server struct {
	streamReg     StreamRegistry
	connReg       cluster.ConnRegistry
	segReg        storage.SegmentRegistry
	localNodeName string
}

func (s *Server) Control(controlSrv Executor_ControlServer) error {

	fmt.Println("new control stream")

	nodeExec := NewNodeExec(s.connReg, s.segReg, s.streamReg, controlSrv, s.localNodeName)
	nodeExec.Start()

	<-nodeExec.ExecutionContext().Done()

	execErr := nodeExec.ExecutionContext().Err()

	if execErr != nil {
		return execErr.Err()
	}

	return nil

}

func (s *Server) VectorExchange(vectorServer Executor_VectorExchangeServer) error {

	msg, err := vectorServer.Recv()

	if err != nil {
		return err
	}

	headerMsg, ok := msg.Msg.(*VectorExchangeMsg_Header)

	if !ok {
		return errors.New("invalid stream header")
	}

	queryId, err := uuid.FromBytes(headerMsg.Header.QueryId)

	if err != nil {
		return err
	}

	ch := make(chan *ExecError, 1)

	stream := NewVectorExchangeStream(vectorServer, ch)

	err = s.streamReg.RegisterStream(queryId, headerMsg.Header.StreamId, stream)

	if err != nil {
		return err
	}

	execErr := <-ch

	if execErr != nil {
		return execErr.Err()
	}

	return nil

}
