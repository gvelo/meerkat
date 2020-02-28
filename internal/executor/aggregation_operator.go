package executor

import (
	"fmt"
	"math"
	"meerkat/internal/storage"
	"unsafe"
)

type AggregationOperator struct {
}

// AggType defines an aggregation function.
type AggType int

// Supported aggregation types.
const (
	Min AggType = iota
	Max
	Mean
	Count
	Sum
	SumSq
	StdDev
)

func (d AggType) String() string {
	return [...]string{"Min", "Max", "Mean", "Count", "Sum", "SumSq", "StdDev"}[d]
}

// Should we use it??? change for functions.
type Counter interface {
	Count() int
	SumSq() int
	Sum() int
	Stdev() float64
	Mean() float64
	Min() int
	Update(i int)
	ValueOf(aggType AggType) float64
}

type counter struct {
	sum   int
	sumSq int
	count int
	max   int
	min   int
	t     AggType
}

func NewCounter() Counter {
	return &counter{
		max: math.MinInt64,
		min: math.MaxInt64,
	}
}

// Counter returns the number of values received.
func (c *counter) Count() int { return c.count }

// Sum returns the sum of Counter values.
func (c *counter) Sum() int { return c.sum }

// SumSq returns the squared sum of Counter values.
func (c *counter) SumSq() int { return c.sumSq }

// Mean returns the mean Counter value.
func (c *counter) Mean() float64 {
	if c.count == 0 {
		return 0
	}
	return float64(c.sum) / float64(c.count)
}

// StdDev returns the standard deviation counter value.
func (c *counter) Stdev() float64 {
	return stdev(c.count, float64(c.sumSq), float64(c.sum))
}

func stdev(count int, sumSq, sum float64) float64 {
	div := count * (count - 1)
	if div == 0 {
		return 0.0
	}
	num := float64(count)*sumSq - sum*sum
	return math.Sqrt(num / float64(div))
}

// Min returns the minimum Counter value.
func (c *counter) Min() int { return c.min }

// Max returns the maximum Counter value.
func (c *counter) Max() int { return c.min }

func (c *counter) Update(i int) {

	c.count++
	c.sum = c.sum + i

	if c.max < i {
		c.max = i
	}

	if c.min > i {
		c.min = i
	}

	if c.t == SumSq || c.t == StdDev {
		c.sumSq += i * i
	}

}

// ValueOf returns the value for the aggregation type.
func (c *counter) ValueOf(aggType AggType) float64 {
	switch aggType {
	case Min:
		return float64(c.Min())
	case Max:
		return float64(c.Max())
	case Mean:
		return c.Mean()
	case Count:
		return float64(c.Count())
	case Sum:
		return float64(c.Sum())
	case SumSq:
		return float64(c.SumSq())
	case StdDev:
		return c.Stdev()
	default:
		return 0
	}
}

// Histogram operator takes a list of ordered items and returns
// a map with the input as key and the count as value.
//
// Example:
//  Duration = 1 min
//
//  input      output
// 12:00:00    12:00:00, 1
// 12:01:00    12:01:00, 3
// 12:01:00    12:03:00, 1
// 12:01:00    12:04:00, 3
// 12:03:00    12:07:00, 3
// 12:04:00
// 12:04:00
// 12:04:00
// 12:07:00
// 12:07:00
// 12:07:00
//
//
func histogram(ts []int) map[int]int {

	m := make(map[int]int)
	for i := 0; i < len(ts); i++ {
		var s int = 1
		ant := ts[i]
		for ; i < len(ts) && ant == ts[i]; i++ {
			s++
		}
		m[ant] = s
	}

	return m
}

func NewHashAggregateOperator(ctx Context, child MultiVectorOperator, aggCols []Aggregation, keyCols []int) MultiVectorOperator {
	return &HashAggregateOperator{
		ctx:     ctx,
		child:   child,
		aggCols: aggCols,
		keyCols: keyCols,
	}
}

type Aggregation struct {
	AggType AggType
	AggCol  int
}

// AggregateOperator
//
type HashAggregateOperator struct {
	ctx       Context
	child     MultiVectorOperator
	aggCols   []Aggregation // Aggregation columns
	keyCols   []int         // Group by columns
	resultVec []storage.Vector
}

func createVector(idx []byte, ctx Context) storage.Vector {
	c := ctx.Segment().Col(idx)
	switch c.(type) {
	case storage.FloatColumn:
		return storage.NewFloatVector()
	case storage.IntColumn:
		return storage.NewIntVector()
	case storage.StringColumn:
		return storage.NewByteSliceVector()
	case storage.TextColumn:
		return storage.NewByteSliceVector()

	}
	panic(" No mapping Found.")
}

func (r *HashAggregateOperator) Init() {
	r.child.Init()
	r.resultVec = make([]storage.Vector, 0)
	// new index column map
	nkv := make(map[int][]byte)
	// old index column map
	// TODO: tengo que saber cual es la entrada en Vectores, para saber que salen, no estoy seguro que sea lo merjor.
	// aca lo mejor seria usar la salida del counter.
	if kv, ok := r.ctx.Get(ColumnIndexKeysKey); ok == true {

		okv := kv.(map[int][]byte)

		// for each key (group by) col we a create a result vector
		last := 0
		for i := 0; i < len(r.keyCols); i++ {
			colId := okv[r.keyCols[i]]
			r.resultVec = append(r.resultVec, createVector(colId, r.ctx))
			nkv[i] = okv[r.keyCols[i]]
			last = i + 1
		}

		// for each aggregated column we create a result vector
		for x := 0; x < len(r.aggCols); x++ {
			colId := okv[r.aggCols[x].AggCol]
			r.resultVec = append(r.resultVec, createVector(colId, r.ctx))

			resName := make([]byte, 0)
			resName = append(resName, okv[r.aggCols[x].AggCol]...)
			resName = append(resName, '_')
			resName = append(resName, []byte(r.aggCols[x].AggType.String())...)

			nkv[last+x] = resName
		}

		// Update the key values
		r.ctx.Value(ColumnIndexKeysKey, nkv)

	} else {
		panic("Error parsing column keys")
	}

}

func (r *HashAggregateOperator) Destroy() {
	r.child.Destroy()
}

func (r *HashAggregateOperator) Next() []storage.Vector {

	mKey := make(map[string]int)

	// the capacity should be the product of the cardinality of the n keys.
	result := make([][]interface{}, 0, 100)

	l := len(r.keyCols)

	// Iterate over all vectors...
	for n := r.child.Next(); n != nil; n = r.child.Next() {

		// iterate over all "rows"
		// For each item in first column deberia ser 1000? o algo por el estilo.?
		for i := 0; i < n[0].Len(); i++ {

			// create the key
			k := createKey(n, r.keyCols, i)

			// check the key
			if _, ok := mKey[string(k)]; ok != true {

				// Create the result array
				c := make([]interface{}, 0)

				// append the keys
				for _, it := range r.keyCols {
					switch col := n[it].(type) {
					case storage.IntVector:
						c = append(c, col.ValuesAsInt()[i])
					case storage.FloatVector:
						c = append(c, col.ValuesAsFloat()[i])
					case storage.ByteSliceVector:
						c = append(c, col.Get(i))
					}
				}

				// append the counters.
				for range r.aggCols {
					c = append(c, NewCounter())
				}

				// sets the key's position in the resulting slice
				mKey[string(k)] = len(result)
				result = append(result, c)
			}

			// update values.
			// TODO: Usando contadores se repiten por ahora...
			for x, it := range r.aggCols {
				idx := mKey[string(k)] // get the result array index
				// TODO: ver los demas tipos de valores.
				switch col := n[it.AggCol].(type) {
				case storage.IntVector:
					result[idx][l+x].(Counter).Update(col.ValuesAsInt()[i])
				case storage.FloatVector:
					result[idx][l+x].(Counter).Update(int(col.ValuesAsFloat()[i]))
				}

			}

		}

	}
	// Sin contadores esto seria mejor....
	for i := 0; i < len(result); i++ {
		for j := 0; j < len(result[i]); j++ {
			fmt.Printf("1 %p , %v\n", r.resultVec[j], r.resultVec[j])
			switch v := r.resultVec[j].(type) {
			//Check counters...
			case storage.IntVector:

				if j >= len(r.keyCols) {
					v.AppendValue(int(result[i][j].(Counter).ValueOf(r.aggCols[j-len(r.keyCols)].AggType)))
				} else {
					v.AppendValue(result[i][j].(int))
				}
				fmt.Printf("3 %p , %v\n", v, v)
			case storage.FloatVector:

				if j >= len(r.keyCols) {
					v.AppendValue(result[i][j].(Counter).ValueOf(r.aggCols[j-len(r.keyCols)].AggType))
				} else {
					v.AppendValue(result[i][j].(float64))
				}
				fmt.Printf("3 %p , %v\n ", v, v)
			case storage.ByteSliceVector:
				v.AppendValue(result[i][j].([]byte))
				fmt.Printf("3 %p , %v\n ", v, v)
			}
		}
	}

	return r.resultVec
	// habria que ir agrupando los vectores y claves.... [c1][c2][c3]  para despues hacer -> (sum1) (sum3)
	// pasar to do a vectores
	// Aca deberia devolver los vectores.
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
}

func (r *SortedAggregateOperator) Destroy() {
	r.child.Destroy()
}

func (r *SortedAggregateOperator) Next() []storage.Vector {
	r.child.Next()
	return nil
}
