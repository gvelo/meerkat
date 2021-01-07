package physical

// TODO(gvelo) this merge implementation in only for test purposes.
type MergeOp struct {
	input    []BatchOperator
	inputIdx int
}

func NewMergeOp(input []BatchOperator) *MergeOp {

	return &MergeOp{
		input:    input,
		inputIdx: 0,
	}

}

func (m *MergeOp) Init() {
	// TODO(gvelo) Init() it is ok to call init multiple times so we need to
	// track initialization state.
	for _, operator := range m.input {
		operator.Init()
	}
}

func (m *MergeOp) Close() {
	for _, operator := range m.input {
		operator.Close()
	}
}

func (m *MergeOp) Next() Batch {

	for len(m.input) != 0 {

		batch := m.input[m.inputIdx].Next()

		if batch.Len == 0 {
			m.input = append(m.input[:m.inputIdx], m.input[m.inputIdx+1:]...)
			continue
		}

		m.inputIdx++

		if m.inputIdx == len(m.input) {
			m.inputIdx = 0
		}

		return batch

	}

	return Batch{
		Len: 0,
	}

}

func (m *MergeOp) Accept(v Visitor) {
	for i, operator := range m.input {
		m.input[i] = Walk(operator, v).(BatchOperator)
	}
}
