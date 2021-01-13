package physical

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"meerkat/internal/query/exec"
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
	"meerkat/internal/util/sliceutil"
)

type ExchangeInOp struct {
	streamReg         exec.StreamRegistry
	streamId          int64
	queryId           uuid.UUID
	vectorExchangeSrv exec.Executor_VectorExchangeServer
	stream            exec.VectorExchangeStream
	hasStream         bool
	execCtx           exec.ExecutionContext
}

func NewExchangeInOp(
	streamReg exec.StreamRegistry,
	streamId int64,
	queryId uuid.UUID,
) *ExchangeInOp {

	return &ExchangeInOp{
		streamReg: streamReg,
		streamId:  streamId,
		queryId:   queryId,
	}

}

func (e *ExchangeInOp) Init()  {}
func (e *ExchangeInOp) Close() {}

func (e *ExchangeInOp) initStream() {

	stream, err := e.streamReg.GetStream(e.queryId, e.streamId)

	if err != nil {
		panic(fmt.Sprintf("cannot init stream : %v", err))
	}

	e.stream = stream
	e.vectorExchangeSrv = stream.Server()

	go e.handleStreamClose()

}

func (e *ExchangeInOp) handleStreamClose() {
	<-e.execCtx.Done()
	e.stream.Close(e.execCtx.Err())
}

func (e *ExchangeInOp) Next() Batch {

	if !e.hasStream {
		e.initStream()
		e.hasStream = true
	}

	msg, err := e.vectorExchangeSrv.Recv()

	if err == io.EOF {
		e.stream.Close(nil)
		return Batch{
			Len:     0,
			Columns: nil,
		}
	}

	if err != nil {
		panic(err)
	}

	switch m := msg.GetMsg().(type) {

	case *exec.VectorExchangeMsg_Error:
		panic(m.Error)
	case *exec.VectorExchangeMsg_VectorBatch:
		return buildBatch(m.VectorBatch)
	default:
		panic(fmt.Errorf("unexpected message type: %v", m))

	}

}

func (e *ExchangeInOp) Accept(v Visitor) {}

func buildBatch(batchMsg *exec.VectorBatch) Batch {

	columns := make(map[string]Col)

	for _, column := range batchMsg.Columns {
		columns[column.Name] = Col{
			Group:      column.Group,
			Order:      column.Order,
			Vec:        buildVector(column),
			ColumnType: column.ColType,
		}
	}

	batch := Batch{
		Len:     int(batchMsg.Len),
		Columns: columns,
	}

	return batch

}

func buildVector(column *exec.Column) vector.Vector {

	var validity []uint64

	if len(column.Validity) != 0 {
		validity = sliceutil.B2U64(column.Validity)
	}

	switch column.ColType {

	case storage.ColumnType_TIMESTAMP, storage.ColumnType_INT64, storage.ColumnType_DATETIME:
		data := sliceutil.B2I64(column.Vector)
		vec := vector.NewInt64Vector(data, validity)
		return &vec
	case storage.ColumnType_FLOAT64:
		data := sliceutil.B2F(column.Vector)
		vec := vector.NewFloat64Vector(data, validity)
		return &vec
	case storage.ColumnType_STRING:
		offset := sliceutil.B2I(column.Offsets)
		vec := vector.NewByteSliceVector(column.Vector, offset, validity)
		return &vec
	default:
		panic(fmt.Errorf("column type %v not supported yet", column))

	}

}

type LocalExchangeInOp struct {
	streamId int64
}

func (op *LocalExchangeInOp) Init()            {}
func (op *LocalExchangeInOp) Close()           {}
func (op *LocalExchangeInOp) Accept(v Visitor) {}

func (op *LocalExchangeInOp) Next() Batch {
	panic("Next() called on LocalExchangeInOp")
}

func NewLocalExchangeInOp(streamId int64) *LocalExchangeInOp {
	return &LocalExchangeInOp{streamId: streamId}
}
