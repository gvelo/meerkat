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
}

func NewExchangeInOp(
	streamReg exec.StreamRegistry,
	streamId int64,
	queryId uuid.UUID,
	vectorExchangeSrv exec.Executor_VectorExchangeServer,
) *ExchangeInOp {

	return &ExchangeInOp{
		streamReg:         streamReg,
		streamId:          streamId,
		queryId:           queryId,
		vectorExchangeSrv: vectorExchangeSrv,
	}

}

func (e *ExchangeInOp) Init() {

	stream, err := e.streamReg.GetStream(e.queryId, e.streamId)

	if err != nil {
		panic(err)
	}

	e.vectorExchangeSrv = stream.Server()

}

func (e *ExchangeInOp) Close() {}

func (e *ExchangeInOp) Next() Batch {

	msg, err := e.vectorExchangeSrv.Recv()

	if err == io.EOF {
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
