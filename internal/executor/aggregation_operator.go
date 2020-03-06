package executor

import (
	"meerkat/internal/storage"
	"unsafe"
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

func (r *HashAggregateOperator) Next() []storage.Vector {

	mKey := make(map[string]int)

	result := make([][]interface{}, 0)

	l := len(r.keyCols)

	n := r.child.Next()

	// Iterate over all vectors...
	for ; n != nil; n = r.child.Next() {

		// iterate over all "rows"
		for i := 0; i < n[0].Len(); i++ {

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
				case storage.IntVector:
					counter.Update(float64(col.ValuesAsInt()[i]))
				case storage.FloatVector:
					counter.Update(col.ValuesAsFloat()[i])
				case storage.ByteSliceVector:
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

	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result[i]); j++ {
			// sacar este if de mierda
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

	resVec := make([]storage.Vector, 0)
	for _, it := range res {
		switch it.(type) {
		case []int:
			resVec = append(resVec, storage.NewIntVectorFromSlice(it.([]int)))
		case []float64:
			resVec = append(resVec, storage.NewFloatVectorFromSlice(it.([]float64)))
		case [][]byte:
			resVec = append(resVec, storage.NewByteSliceVectorSlice(it.([][]byte)))
		}
	}

	return resVec
}

func createKey(n []storage.Vector, keyCols []int, index int) []byte {
	k := make([]byte, 0)
	for _, it := range keyCols {
		switch t := n[it].(type) {
		case storage.ByteSliceVector:
			k = append(k, t.Get(index)...)
		case storage.IntVector:
			b := (*[8]byte)(unsafe.Pointer(&t.ValuesAsInt()[index]))[:]
			k = append(k, b...)
		case storage.FloatVector:
			b := (*[8]byte)(unsafe.Pointer(&t.ValuesAsFloat()[index]))[:]
			k = append(k, b...)
		}
	}

	return k
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
//
type SortedAggregateOperator struct {
	ctx       Context
	child     MultiVectorOperator
	aggCols   []Aggregation // Aggregation columns
	keyCols   []int         // Group by columns
	resultVec []storage.Vector
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

func (r *SortedAggregateOperator) Next() []storage.Vector {

	result := make([][]interface{}, 0)

	l := len(r.keyCols)

	n := r.child.Next()

	// Iterate over all vectors...
	for ; n != nil; n = r.child.Next() {

		// iterate over all "rows"
		for i := 0; i < n[0].Len(); i++ {

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
					case storage.IntVector:
						counter.Update(float64(col.ValuesAsInt()[i]))
					case storage.FloatVector:
						counter.Update(col.ValuesAsFloat()[i])
					case storage.ByteSliceVector:
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

	resVec := make([]storage.Vector, 0)
	for _, it := range res {
		switch it.(type) {
		case []int:
			resVec = append(resVec, storage.NewIntVectorFromSlice(it.([]int)))
		case []float64:
			resVec = append(resVec, storage.NewFloatVectorFromSlice(it.([]float64)))
		case [][]byte:
			resVec = append(resVec, storage.NewByteSliceVectorSlice(it.([][]byte)))
		}
	}

	return resVec
}
