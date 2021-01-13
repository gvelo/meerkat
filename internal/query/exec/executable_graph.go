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
	"errors"
	"fmt"
	"github.com/google/uuid"
	"meerkat/internal/cluster"
	"meerkat/internal/query/logical"
	"meerkat/internal/query/physical"
	"meerkat/internal/storage"
)

type DAGBuilder interface {
	BuildDAG(
		fragments []*logical.Fragment,
		queryId uuid.UUID,
		writer QueryOutputWriter,
		execCtx ExecutionContext,
	) (physical.DAG, error)
}

func NewDAGBuilder(
	connReg cluster.ConnRegistry,
	segReg storage.SegmentRegistry,
	streamReg StreamRegistry,
	localNodeName string,
) DAGBuilder {

	return &dagBuilder{
		connReg:       connReg,
		segReg:        segReg,
		streamReg:     streamReg,
		localNodeName: localNodeName,
	}

}

type dagBuilder struct {
	connReg       cluster.ConnRegistry
	segReg        storage.SegmentRegistry
	streamReg     StreamRegistry
	localNodeName string
}

func (e *dagBuilder) BuildDAG(
	fragments []*logical.Fragment,
	queryId uuid.UUID,
	writer QueryOutputWriter,
	execCtx ExecutionContext,
) (physical.DAG, error) {

	var segments []storage.Segment
	var roots []physical.RunnableOp
	var runnables []physical.RunnableOp
	localStreamMap := make(map[int64]physical.BatchOperator)

	for _, fragment := range fragments {

		builder := &dagBuilderVisitor{
			outputWriter:   writer,
			queryId:        queryId,
			localNodeName:  e.localNodeName,
			connReg:        e.connReg,
			streamReg:      e.streamReg,
			segReg:         e.segReg,
			localStreamMap: make(map[int64]physical.BatchOperator),
		}

		if len(fragment.Roots) > 1 {
			return nil, errors.New("fragments with multiple roots ( forkOp ) not supported yet")
		}

		// TODO(gvelo) catch panics, release adquired segments
		//  and return error
		logical.Walk(fragment.Roots[0], builder)

		segments = append(segments, builder.segments...)
		roots = append(roots, builder.roots...)
		runnables = append(runnables, builder.runnableOps...)

		for streamId, outputOp := range builder.localStreamMap {
			localStreamMap[streamId] = outputOp
		}

	}

	// connect local output streams to local input streams

	streamVisitor := newStreamRewriteVisitor(localStreamMap)

	for _, rootOp := range roots {
		physical.Walk(rootOp, streamVisitor)
	}

	dag := physical.NewDAG(
		execCtx,
		runnables,
		roots,
		queryId,
		segments,
		e.segReg,
		e.localNodeName,
	)

	return dag, nil

}

type dagBuilderVisitor struct {
	child         []physical.BatchOperator
	roots         []physical.RunnableOp
	outputWriter  QueryOutputWriter
	queryId       uuid.UUID
	localNodeName string
	connReg       cluster.ConnRegistry
	streamReg     StreamRegistry
	segReg        storage.SegmentRegistry
	// localStreamMap map local streams to the output operator. This operator
	// will be used as input operator for local streams instead of a ExchangeInOp
	localStreamMap map[int64]physical.BatchOperator
	segments       []storage.Segment
	runnableOps    []physical.RunnableOp
}

func (g *dagBuilderVisitor) VisitPre(n logical.Node) logical.Node { return n }

func (g *dagBuilderVisitor) VisitPost(n logical.Node) logical.Node {

	switch node := n.(type) {
	case *logical.OutputOp:
		g.assertSingleInput()
		jsonOutputOp := physical.NewJsonOutputOp(g.child[0], g.outputWriter)
		g.roots = append(g.roots, jsonOutputOp)
		g.runnableOps = append(g.runnableOps, jsonOutputOp)
	case *logical.SourceOp:

		// TODO(gvelo) call Segments() with real params
		g.segments = g.segReg.Segments(nil, "", "")

		var child []physical.BatchOperator

		for _, segment := range g.segments {
			op := buildBatchOp(segment) // TODO(gvelo) add columns and filter exp
			child = append(child, op)
		}

		g.child = child

	case *logical.NodeOutOp:

		var input physical.BatchOperator

		if len(g.child) > 1 {
			// TODO(gvelo) user merge sort op
			input = physical.NewMergeOp(g.child)

		} else {
			input = g.child[0]
		}

		streamId, found := node.StreamMap[g.localNodeName]

		if !found {
			panic(fmt.Sprintf("cannot found stream id for node %v", g.localNodeName))
		}

		// output to local node
		if node.Dst == g.localNodeName {

			// we add the input op to the output stream map
			// this operator will be used later as an input for a local
			// fragment.
			g.localStreamMap[streamId] = input

		} else {

			conn := g.connReg.ClientConn(node.Dst)

			if conn == nil {
				panic(fmt.Sprintf("cannot found grpc client conn for node %v", node.Dst))
			}

			client := NewExecutorClient(conn)

			exchangeOutOp := physical.NewExchangeOutOp(
				input,
				client,
				g.queryId,
				streamId,
				g.localNodeName,
			)

			g.roots = append(g.roots, exchangeOutOp)
			g.runnableOps = append(g.runnableOps, exchangeOutOp)

		}

	case *logical.MergeSortOp:

		var inputs []physical.BatchOperator

		for srcNodeName, streamId := range node.StreamMap {

			var op physical.BatchOperator

			if srcNodeName == g.localNodeName {
				op = physical.NewLocalExchangeInOp(streamId)
			} else {
				op = physical.NewExchangeInOp(g.streamReg, streamId, g.queryId)
			}

			inputs = append(inputs, op)

		}

		mergeOp := physical.NewMergeOp(inputs)

		g.child = []physical.BatchOperator{mergeOp}

	default:
		panic("unknown operator")

	}

	return n

}

func (g *dagBuilderVisitor) assertSingleInput() {
	if len(g.child) != 1 {
		panic("expected single input")
	}
}

func buildBatchOp(segment storage.Segment) physical.BatchOperator {

	info := segment.Info()

	var input []physical.ColumnOperator
	var colNames []string

	for _, columnInfo := range info.Columns {

		colNames = append(colNames, columnInfo.Name)
		col := segment.Column(columnInfo.Name)
		var iter storage.Iterator

		switch c := col.(type) {
		case storage.Int64Column:
			iter = c.Iterator()
		case storage.ByteSliceColumn:
			iter = c.Iterator()
		case storage.Float64Column:
			iter = c.Iterator()
		default:
			panic("unknown column type")

		}

		op := physical.NewColumnReaderOp(iter)
		input = append(input, op)
	}

	batchBuilder := physical.NewBatchBuilderOp(input, colNames)

	return batchBuilder

}

// streamRewriteVisitor rewrite the operator DAG replacing the LocalExchangeInOp
// by the local output operator.
type streamRewriteVisitor struct {
	localStreamMap [int64]physical.BatchOperator
}

func newStreamRewriteVisitor(localStreamMap [int64]physical.BatchOperator) *streamRewriteVisitor {
	return &streamRewriteVisitor{localStreamMap: localStreamMap}
}

func (s *streamRewriteVisitor) VisitPre(n physical.Operator) physical.Operator { return n }

func (s *streamRewriteVisitor) VisitPost(n physical.Operator) physical.Operator {

	switch op := n.(type) {
	case *physical.LocalExchangeInOp:

		localOutputOp := s.localStreamMap[op.streamId]

		if localOutputOp == nil {
			panic(fmt.Errorf("error, local stream not found : %v", op.streamId))
		}

		return localOutputOp

	}

	return n

}
