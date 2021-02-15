package physical

type BatchBuilderOp struct {
	colNames []string
	input    []ColumnOperator
}

func NewBatchBuilderOp(input []ColumnOperator, colNames []string) *BatchBuilderOp {
	return &BatchBuilderOp{
		colNames: colNames,
		input:    input,
	}
}

func (b *BatchBuilderOp) Init() {
	for _, operator := range b.input {
		operator.Init()
	}
}

func (b *BatchBuilderOp) Close() {
	for _, operator := range b.input {
		operator.Close()
	}
}

func (b *BatchBuilderOp) Next() Batch {

	batch := NewBatch()
	lastVectorLen := 0

	for i, name := range b.colNames {

		v := b.input[i].Next()

		batch.Columns[name] = Col{
			Group: 0,
			Order: int64(i),
			Vec:   v,
		}

		if i != 0 && lastVectorLen != v.Len() {
			panic("vector length doesn't match")
		}

		lastVectorLen = v.Len()

	}

	batch.Len = lastVectorLen

	return batch

}

func (b *BatchBuilderOp) Accept(v Visitor) {
	for i, operator := range b.input {
		b.input[i] = Walk(operator, v).(ColumnOperator)
	}
}
