package physical

import (
	"encoding/json"
	"meerkat/internal/query/exec"
	"meerkat/internal/storage/vector"
)

type JsonOutputOp struct {
	input   BatchOperator
	writer  exec.QueryOutputWriter
	encoder *json.Encoder
}

func NewJsonOutputOp(input BatchOperator, writer exec.QueryOutputWriter) *JsonOutputOp {
	return &JsonOutputOp{
		input:   input,
		writer:  writer,
		encoder: json.NewEncoder(writer),
	}
}

func (o *JsonOutputOp) Init() { o.input.Init() }

func (o *JsonOutputOp) Close() { o.input.Close() }

func (o *JsonOutputOp) Run() {

	for {

		// TODO(gvelo) handle panic
		batch := o.input.Next()

		if batch.Len == 0 {
			o.writer.Flush()
			return
		}

		o.writeBatch(batch)

	}
}

func (o *JsonOutputOp) Accept(v Visitor) {
	o.input = Walk(o.input, v).(BatchOperator)
}

func (o *JsonOutputOp) writeBatch(batch Batch) {

	// TODO(gvelo): sort columns by group and order.

	m := map[string]interface{}{
		"type":    "column_batch",
		"columns": buildJsonColumns(batch.Columns),
	}

	err := o.encoder.Encode(m)

	if err != nil {
		panic(err)
	}

}

func buildJsonColumns(colMap map[string]Col) []map[string]interface{} {

	var jsonColumns []map[string]interface{}

	for name, column := range colMap {
		c := buildJsonColumn(name, column)
		jsonColumns = append(jsonColumns, c)
	}

	return jsonColumns

}

func buildJsonColumn(name string, column Col) map[string]interface{} {

	jsonCol := map[string]interface{}{
		"name":   name,
		"type":   column.ColumnType.String(),
		"values": buildJsonVectorValues(column.Vec),
	}

	return jsonCol
}

func buildJsonVectorValues(vec vector.Vector) interface{} {

	switch v := vec.(type) {
	case *vector.Int64Vector:
		return v.Values()
	case *vector.ByteSliceVector:
		s := make([]string, vec.Len())
		for i := 0; i < v.Len(); i++ {
			s[i] = string(v.Get(i))
		}
		return s
	default:
		panic("cannot convert vector to json")

	}

}
