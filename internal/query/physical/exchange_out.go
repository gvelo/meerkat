package physical

import (
	"context"
	"github.com/google/uuid"
	"meerkat/internal/query/exec"
	"meerkat/internal/storage/vector"
)

type ExchangeOut struct {
	input      BatchOperator
	execClient exec.ExecutorClient
	queryId    uuid.UUID
	streamId   int64
}

func NewExchangeOut(
	input BatchOperator,
	execClient exec.ExecutorClient,
	queryId uuid.UUID,
	streamId int64,
) *ExchangeOut {

	return &ExchangeOut{
		input:      input,
		execClient: execClient,
		queryId:    queryId,
		streamId:   streamId,
	}

}

func (e *ExchangeOut) Init()  { e.input.Init() }
func (e *ExchangeOut) Close() { e.input.Close() }

func (e *ExchangeOut) Run() {

	// TODO(gvelo) use a execCtx child here
	vectorExchangeClient, err := e.execClient.VectorExchange(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	headerMsg := &exec.VectorExchangeMsg{
		Msg: &exec.VectorExchangeMsg_Header{
			Header: &exec.StreamHeader{
				QueryId:  e.queryId[:],
				StreamId: e.streamId,
			},
		},
	}

	err = vectorExchangeClient.Send(headerMsg)

	if err != nil {
		panic(err)
	}

	for {

		// TODO(gvelo) if we get a panic signal we should send a
		// VectorExchangeMsg_Error

		v := e.input.Next()

		if v.Len == 0 {

			_, err := vectorExchangeClient.CloseAndRecv() // EOF

			if err != nil {
				panic(err) // TODO(gvelo) panic ?
			}

			return
		}

		vectorMsg := &exec.VectorExchangeMsg{
			Msg: &exec.VectorExchangeMsg_VectorBatch{
				Vector: &exec.VectorBatch{
					Len:     int64(v.Len),
					Columns: buildColumns(v),
				},
			},
		}

		err := vectorExchangeClient.Send(vectorMsg)

		if err != nil {
			panic(err)
		}

	}
}

func (e *ExchangeOut) Accept(v Visitor) {
	e.input = Walk(e.input, v).(BatchOperator)
}

func buildColumns(batch Batch) []*exec.Column {

	columns := make([]*exec.Column, len(batch.Columns))

	for name, col := range batch.Columns {

		colProto := &exec.Column{
			Name:     name,
			Group:    col.Group,
			Order:    col.Order,
			ColType:  col.ColumnType,
			Vector:   col.Vec.AsBytes(),
			Validity: col.Vec.ValidityAsBytes(),
			Offsets:  getOffset(col.Vec),
		}

		columns = append(columns, colProto)
	}

	return columns

}

func getOffset(vec vector.Vector) []byte {
	switch v := vec.(type) {
	case *vector.ByteSliceVector:
		return v.OffsetsAsBytes()
	default:
		return nil
	}
}
