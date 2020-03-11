package executor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"meerkat/internal/storage/vector"
)

// HashAggregateOperator
func NewHashAggregateOperator(ctx Context, child MultiVectorOperator, aggCols []Aggregation, keyCols []int) MultiVectorOperator {
	return &HashAggregateOperator{
		ctx:     ctx,
		child:   child,
		aggCols: aggCols,
		keyCols: keyCols,
	}
}

// AggregateOperator
//
type HashAggregateOperator struct {
	ctx     Context
	child   MultiVectorOperator
	aggCols []Aggregation // Aggregation columns
	keyCols []int         // Group by columns
}

func (r *HashAggregateOperator) Init() {
	r.child.Init()
	nkv := createNewKeyMap(r.ctx, r.keyCols, r.aggCols)
	// Update the key values
	r.ctx.Value(ColumnIndexToColumnName, nkv)
}

func (r *HashAggregateOperator) Destroy() {
	r.child.Destroy()
}

func (r *HashAggregateOperator) Next() []vector.Vector {

	mKey := make(map[string]int)

	result := make([][]interface{}, 0)

	l := len(r.keyCols)

	n := r.child.Next()

	// Iterate over all vectors...
	for ; n != nil; n = r.child.Next() {

		l1 := n[0].Len()
		// iterate over all "rows"
		for i := 0; i < l1; i++ {

			// create the key
			k := createKey(n, r.keyCols, i)

			// check the key
			if _, ok := mKey[string(k)]; ok != true {

				c := createResultVector(n, r.keyCols, r.aggCols, i)
				// sets the key's position in the resulting slice
				mKey[string(k)] = len(result)
				result = append(result, c)
			}

			// update values.
			// TODO: Usando contadores se repiten por ahora...
			for x, it := range r.aggCols {

				counter := result[mKey[string(k)]][l+x].(Counter)
				// aca tenemos que separar los operadores.
				switch col := n[it.AggCol].(type) {
				case *vector.IntVector:
					counter.Update(float64(col.Values()[i]))
				case *vector.FloatVector:
					counter.Update(col.Values()[i])
				case *vector.ByteSliceVector:
					counter.Update(1)
				}

			}

		}

	}

	// Create columns
	res := make([]interface{}, 0)

	if kv, ok := r.ctx.Get(ColumnIndexToColumnName); ok == true {
		okv := kv.(map[int][]byte)
		for i, _ := range result[0] {
			res = append(res, createSlice(okv[i], r.ctx))
		}
	}

	for i := 0; i < len(result); i++ { // rows
		for j := 0; j < len(result[i]); j++ { // columns
			// sacar este if de mierda con 2 vectores.
			if j < len(r.keyCols) {
				switch res[j].(type) {
				case []int:
					res[j] = append(res[j].([]int), result[i][j].(int))
				case []float64:
					res[j] = append(res[j].([]float64), result[i][j].(float64))
				case [][]byte:
					res[j] = append(res[j].([][]byte), result[i][j].([]byte))
				}
			} else {
				switch res[j].(type) {
				case []int:
					res[j] = append(res[j].([]int), int(result[i][j].(Counter).ValueOf(r.aggCols[j-len(r.keyCols)].AggType)))
				case []float64:
					res[j] = append(res[j].([]float64), result[i][j].(Counter).ValueOf(r.aggCols[j-len(r.keyCols)].AggType))
				case [][]byte:
					panic("Va a pasar... y fix")
				}
			}
		}
	}

	resVec := make([]vector.Vector, 0)
	for _, it := range res {
		switch it.(type) {
		case []int:
			v := vector.NewIntVector(it.([]int), []uint64{})
			resVec = append(resVec, &v)
		case []float64:
			v := vector.NewFloatVector(it.([]float64), []uint64{})
			resVec = append(resVec, &v)
		case [][]byte:
			s := it.([][]byte)
			offsets := make([]int, len(s))
			sum := 0
			for i, it := range s {
				sum = sum + len(it)
				offsets[i] = sum
			}
			v := vector.NewByteSliceVector(bytes.Join(s, nil), []uint64{}, offsets)
			resVec = append(resVec, &v)
		}
	}

	return resVec
}

func createKey(n []vector.Vector, keyCols []int, index int) []byte {
	buf := new(bytes.Buffer)
	var err error
	for _, it := range keyCols {
		switch t := n[it].(type) {
		case *vector.ByteSliceVector:
			err = binary.Write(buf, binary.LittleEndian, t.Get(index))
		case *vector.IntVector:
		case *vector.FloatVector:
			err = binary.Write(buf, binary.LittleEndian, &t.Values()[index])
		}
		if err != nil {
			panic(fmt.Sprint("binary.Write failed:", err))
		}
	}
	return buf.Bytes()
}

// SortedAggregateOperator
func NewSortedAggregateOperator(ctx Context, child MultiVectorOperator, aggCols []Aggregation, keyCols []int) MultiVectorOperator {
	return &SortedAggregateOperator{
		ctx:     ctx,
		child:   child,
		aggCols: aggCols,
		keyCols: keyCols,
	}
}

// AggregateOperator
type SortedAggregateOperator struct {
	ctx       Context
	child     MultiVectorOperator
	aggCols   []Aggregation // Aggregation columns
	keyCols   []int         // Group by columns
	resultVec []vector.Vector
}

func (r *SortedAggregateOperator) Init() {
	r.child.Init()
	nkv := createNewKeyMap(r.ctx, r.keyCols, r.aggCols)
	// Update the key values
	r.ctx.Value(ColumnIndexToColumnName, nkv)
}

func (r *SortedAggregateOperator) Destroy() {
	r.child.Destroy()
}

func (r *SortedAggregateOperator) Next() []vector.Vector {

	result := make([][]interface{}, 0)

	l := len(r.keyCols)

	n := r.child.Next()

	// Iterate over all vectors...
	for ; n != nil; n = r.child.Next() {

		// iterate over all "rows"
		for i := 0; i < n[0].(vector.Vector).Len(); i++ {

			// create the key
			k := createKey(n, r.keyCols, i)
			ant := k

			c := createResultVector(n, r.keyCols, r.aggCols, i)
			result = append(result, c)

			for string(ant) == string(k) {

				// update values.
				// TODO: Usando contadores se repiten por ahora...
				for x, it := range r.aggCols {

					counter := result[len(result)-1][l+x].(Counter)

					switch col := n[it.AggCol].(type) {
					case *vector.IntVector:
						counter.Update(float64(col.Values()[i]))
					case *vector.FloatVector:
						counter.Update(col.Values()[i])
					case *vector.ByteSliceVector:
						counter.Update(1)
					}

				}

				ant = k
				// create the key
				k = createKey(n, r.keyCols, i)

			}

		}

	}

	// Create columns
	res := make([]interface{}, 0)

	if kv, ok := r.ctx.Get(ColumnIndexToColumnName); ok == true {
		okv := kv.(map[int][]byte)
		for i, _ := range result[0] {
			res = append(res, createSlice(okv[i], r.ctx))
		}
	}

	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result[i]); j++ {

			if j < len(r.keyCols) {
				switch res[j].(type) {
				case []int:
					res[j] = append(res[j].([]int), result[i][j].(int))
				case []float64:
					res[j] = append(res[j].([]float64), result[i][j].(float64))
				case [][]byte:
					res[j] = append(res[j].([][]byte), result[i][j].([]byte))
				}
			} else {
				switch res[j].(type) {
				case []int:
					res[j] = append(res[j].([]int), int(result[i][j].(Counter).ValueOf(r.aggCols[j-len(r.keyCols)].AggType)))
				case []float64:
					res[j] = append(res[j].([]float64), result[i][j].(Counter).ValueOf(r.aggCols[j-len(r.keyCols)].AggType))
				case [][]byte:
					panic("Va a pasar... y fix")
				}
			}
		}
	}

	resVec := make([]vector.Vector, 0)
	for _, it := range res {
		switch it.(type) {
		case []int:
			v := vector.NewIntVector(it.([]int), []uint64{})
			resVec = append(resVec, &v)
		case []float64:
			v := vector.NewFloatVector(it.([]float64), []uint64{})
			resVec = append(resVec, &v)
		case [][]byte:
			v := vector.NewByteSliceVector(it.([]byte), []uint64{}, []int{})
			resVec = append(resVec, &v)
		}
	}

	return resVec
}
