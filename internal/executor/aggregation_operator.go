package executor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math"
	"meerkat/internal/storage/vector"
)

// HashAggregateOperator
func NewHashAggregateOperator(ctx Context, child MultiVectorOperator, aggCols []Aggregation, keyCols []int) MultiVectorOperator {
	return &HashAggregateOperator{
		ctx:     ctx,
		child:   child,
		aggCols: aggCols,
		keyCols: keyCols,
		log:     log.With().Str("src", "HashAggregateOperator").Logger(),
	}
}

// AggregateOperator
//
type HashAggregateOperator struct {
	ctx     Context
	child   MultiVectorOperator
	aggCols []Aggregation // Aggregation columns
	keyCols []int         // Group by columns
	sz      int
	log     zerolog.Logger
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

/*
  TODO(sebad): We need to spill out to disk we we do not have memory
*/
func (r *HashAggregateOperator) Next() []interface{} {

	mKey := make(map[string]int)

	rKey := make([][]interface{}, 0, r.sz)
	rAgg := make([][]interface{}, 0, r.sz)

	n := r.child.Next()

	// Iterate over all vectors...
	for ; n != nil; n = r.child.Next() {

		l1 := n[0].(vector.Vector).Len()
		// iterate over all "rows"
		for i := 0; i < l1; i++ {

			// create the key
			k := createKey(n, r.keyCols, i)

			// check the key
			if _, ok := mKey[string(k)]; ok != true {
				key, agg := createResultVector(n, r.keyCols, r.aggCols, i)
				// sets the key's position in the resulting slice
				mKey[string(k)] = len(rKey)
				rKey = append(rKey, key)
				rAgg = append(rAgg, agg)
			}

			// update values.
			updateCounters(n, rAgg[mKey[string(k)]], r.aggCols, i)
		}

	}

	return pivotAndBuildVectors(r.ctx, rKey, rAgg, r.aggCols)
}

// SortedAggregateOperator
func NewSortedAggregateOperator(ctx Context, child MultiVectorOperator, aggCols []Aggregation, keyCols []int) MultiVectorOperator {
	return &SortedAggregateOperator{
		ctx:     ctx,
		child:   child,
		aggCols: aggCols,
		keyCols: keyCols,
		log:     log.With().Str("src", "SortedAggregateOperator").Logger(),
	}
}

// AggregateOperator
type SortedAggregateOperator struct {
	ctx       Context
	child     MultiVectorOperator
	aggCols   []Aggregation // Aggregation columns
	keyCols   []int         // Group by columns
	resultVec []vector.Vector
	log       zerolog.Logger
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

func (r *SortedAggregateOperator) Next() []interface{} {

	mKey := make(map[string]int)

	rKey := make([][]interface{}, 0)
	rAgg := make([][]interface{}, 0)

	n := r.child.Next()

	// Iterate over all vectors...
	for ; n != nil; n = r.child.Next() {

		l1 := n[0].(vector.Vector).Len()
		var keyAnt []byte
		// iterate over all "rows"
		for i := 0; i < l1; i++ {

			// create the key
			k := createKey(n, r.keyCols, i)

			// check the key
			if !bytes.Equal(keyAnt, k) {
				key, agg := createResultVector(n, r.keyCols, r.aggCols, i)
				// sets the key's position in the resulting slice
				mKey[string(k)] = len(rKey)
				rKey = append(rKey, key)
				rAgg = append(rAgg, agg)
				keyAnt = k
			}

			// update values.
			updateCounters(n, rAgg[mKey[string(k)]], r.aggCols, i)
		}

	}

	return pivotAndBuildVectors(r.ctx, rKey, rAgg, r.aggCols)
}

func updateCounters(n []interface{}, rAgg []interface{}, aggCols []Aggregation, i int) {
	// TODO: Usando contadores se repiten por ahora...
	for x, it := range aggCols {

		counter := rAgg[x].(Counter)
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

func createKey(n []interface{}, keyCols []int, index int) []byte {
	buf := new(bytes.Buffer)
	var err error
	for _, it := range keyCols {
		switch t := n[it].(type) {
		case *vector.ByteSliceVector:
			err = binary.Write(buf, binary.LittleEndian, t.Get(index))
		case *vector.IntVector:
			err = binary.Write(buf, binary.LittleEndian, int64(t.Values()[index]))
		case *vector.FloatVector:
			err = binary.Write(buf, binary.LittleEndian, math.Float64bits(t.Values()[index]))
		}
		if err != nil {
			panic(fmt.Sprint("binary.Write failed:", err))
		}
	}
	return buf.Bytes()
}

func pivotAndBuildVectors(ctx Context, rKey [][]interface{}, rAgg [][]interface{}, aggCols []Aggregation) []interface{} {
	// Create columns
	resKey := make([]interface{}, 0)
	resAgg := make([]interface{}, 0)

	if kv, ok := ctx.Get(ColumnIndexToColumnName); ok == true {
		okv := kv.(map[int]string)
		// create key slices
		for i, _ := range rKey[0] {
			resKey = append(resKey, createSlice(okv[i], ctx))
		}

		// create aggregation slices
		for i, _ := range rAgg[0] {
			resAgg = append(resAgg, createSlice(okv[i+len(resKey)], ctx))

		}
	} else {
		panic("Error")
	}

	for j := 0; j < len(resKey); j++ { // columns
		switch resKey[j].(type) {
		case []int:
			for i := 0; i < len(rKey); i++ { // rows
				resKey[j] = append(resKey[j].([]int), rKey[i][j].(int))
			}
		case []float64:
			for i := 0; i < len(rKey); i++ { // rows
				resKey[j] = append(resKey[j].([]float64), rKey[i][j].(float64))
			}
		case [][]byte:
			for i := 0; i < len(rKey); i++ { // rows
				resKey[j] = append(resKey[j].([][]byte), rKey[i][j].([]byte))
			}
		}
	}

	for j := 0; j < len(resAgg); j++ { // columns
		switch resAgg[j].(type) {
		case []int:
			for i := 0; i < len(rAgg); i++ { // rows
				resAgg[j] = append(resAgg[j].([]int), int(rAgg[i][j].(Counter).ValueOf(aggCols[j].AggType)))
			}
		case []float64:
			for i := 0; i < len(rAgg); i++ { // rows
				resAgg[j] = append(resAgg[j].([]float64), rAgg[i][j].(Counter).ValueOf(aggCols[j].AggType))
			}
		case [][]byte:
			panic("?????")
		}
	}

	resVec := make([]interface{}, 0)
	res := append(resKey, resAgg...)
	for _, it := range res {
		switch it.(type) {
		case []int:
			v := vector.NewIntVector(it.([]int), []uint64{})
			v.SetLen(len(it.([]int)))
			resVec = append(resVec, &v)
		case []float64:
			v := vector.NewFloatVector(it.([]float64), []uint64{})
			v.SetLen(len(it.([]float64)))
			resVec = append(resVec, &v)
		case [][]byte:
			s := it.([][]byte)
			offsets := make([]int, len(s))
			sum := 0
			for i, it := range s {
				sum = sum + len(it)
				offsets[i] = sum
			}
			v := vector.NewByteSliceVector(bytes.Join(s, nil), offsets, []uint64{})
			resVec = append(resVec, &v)
		}
	}

	return resVec
}
