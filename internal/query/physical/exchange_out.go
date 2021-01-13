package physical

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"meerkat/internal/query/exec"
	"meerkat/internal/storage/vector"
)

type ExchangeOutOp struct {
	input         BatchOperator
	execClient    exec.ExecutorClient
	queryId       uuid.UUID
	streamId      int64
	localNodeName string
}

func NewExchangeOutOp(
	input BatchOperator,
	execClient exec.ExecutorClient,
	queryId uuid.UUID,
	streamId int64,
	localNodeName string,
) *ExchangeOutOp {

	return &ExchangeOutOp{
		input:         input,
		execClient:    execClient,
		queryId:       queryId,
		streamId:      streamId,
		localNodeName: localNodeName,
	}

}

func (e *ExchangeOutOp) Init()  { e.input.Init() }
func (e *ExchangeOutOp) Close() { e.input.Close() }

func (e *ExchangeOutOp) Run() {

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

		v, err := e.safeNext()

		if err != nil {

			// extract the execError to propagate it over the stream.
			execErr := exec.ExtractExecError(err)

			if execErr == nil {
				execErr = exec.NewExecError(
					fmt.Sprintf("error executing query : %v", err),
					e.localNodeName,
				)
			}

			vectorMsg := &exec.VectorExchangeMsg{
				Msg: &exec.VectorExchangeMsg_Error{
					Error: execErr,
				},
			}

			err := vectorExchangeClient.Send(vectorMsg)

			if err != nil {
				// TODO(gvelo): nothing we can do here , just log properly
				fmt.Println(err)
			}

			panic(execErr)

		}

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

func (e *ExchangeOutOp) safeNext() (batch Batch, err interface{}) {

	defer func() {
		if r := recover(); r != nil {
			err = r
		}
	}()

	batch = e.input.Next()
}

func (e *ExchangeOutOp) Accept(v Visitor) {
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
