package executor

import (
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

func NewHashAggregateOperator(ctx Context, child MultiVectorOperator, aCols []Aggregation, gByCols []int) MultiVectorOperator {
	return &HashAggregateOperator{
		ctx:     ctx,
		child:   child,
		aCols:   aCols,
		keyCols: gByCols,
	}
}

type Aggregation struct {
	AggType AggType
	AgCol   int
}

// AggregateOperator
//
type HashAggregateOperator struct {
	ctx       Context
	child     MultiVectorOperator
	aCols     []Aggregation // Aggregation columns
	keyCols   []int         // Group by columns
	resultVec []storage.Vector
}

func (r *HashAggregateOperator) Init() {
	r.child.Init()
	r.resultVec = make([]storage.Vector, 0)
	// for each key col a create a result
	for i := 0; i < len(r.keyCols); i++ {
		// r.resultVec = append(r.resultVec, createVector() )
	}

}

func (r *HashAggregateOperator) Destroy() {
	r.child.Destroy()
}

func (r *HashAggregateOperator) Next() []storage.Vector {

	n := r.child.Next() // Iterate over all vectors...

	if n != nil && len(n) > 0 {

		mKey := make(map[string][]Counter)

		for i := 0; i < n[0].Len(); i++ {
			// For each item in first column
			// TODO(sebad): check if this is ok.

			// create the key
			k := createKey(n, r.keyCols, i)

			// check the key
			if _, ok := mKey[string(k)]; ok != true {
				//new key append counters.
				c := make([]Counter, 0)
				for range r.aCols {
					c = append(c, NewCounter())
				}
				mKey[string(k)] = c
			}

			// update values.
			for x, it := range r.aCols {
				c := mKey[string(k)] // get the counters array
				// TODO: ver los demas tipos de valores.
				switch col := n[it.AgCol].(type) {
				case storage.IntVector:
					c[x].Update(col.ValuesAsInt()[i])
				case storage.FloatVector:
					c[x].Update(int(col.ValuesAsFloat()[i]))
				}

			}

		}

		// habria que ir agrupando los vectores y claves.... [c1][c2][c3]  para despues hacer -> (sum1) (sum3)
		// o devolverlos a Nodo agrupador.
		// paralelizar?
		// Analizar como procesar esto de manera mas simple.
	}
	return nil
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

func NewSortedAggregateOperator(ctx Context, child MultiVectorOperator, aCols []AggType, gByCol int) MultiVectorOperator {
	return &SortedAggregateOperator{
		ctx:    ctx,
		child:  child,
		aCols:  aCols,
		gByCol: gByCol,
	}
}

// AggregateOperator
//
type SortedAggregateOperator struct {
	ctx    Context
	child  MultiVectorOperator
	aCols  []AggType
	gByCol int
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
