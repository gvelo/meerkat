package executor

import (
	"math"
	"meerkat/internal/storage"
	"strings"
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
	Count() float64
	SumSq() float64
	Sum() float64
	StdDev() float64
	Mean() float64
	Min() float64
	Update(i float64)
	ValueOf(aggType AggType) float64
}

type counter struct {
	sum   float64
	sumSq float64
	count float64
	max   float64
	min   float64
	t     AggType
}

func NewCounter() Counter {
	return &counter{
		min: math.MaxFloat64,
		max: math.MaxFloat64 * -1.0,
	}
}

// Counter returns the number of values received.
func (c *counter) Count() float64 { return c.count }

// Sum returns the sum of Counter values.
func (c *counter) Sum() float64 { return c.sum }

// SumSq returns the squared sum of Counter values.
func (c *counter) SumSq() float64 { return c.sumSq }

// Mean returns the mean Counter value.
func (c *counter) Mean() float64 {
	if c.count == 0 {
		return 0
	}
	return float64(c.sum) / float64(c.count)
}

// StdDev returns the standard deviation counter value.
func (c *counter) StdDev() float64 {
	return stdev(c.count, float64(c.sumSq), float64(c.sum))
}

func stdev(count float64, sumSq, sum float64) float64 {
	div := count * (count - 1)
	if div == 0 {
		return 0.0
	}
	num := (count)*sumSq - sum*sum
	return math.Sqrt(num / div)
}

// Min returns the minimum Counter value.
func (c *counter) Min() float64 { return c.min }

// Max returns the maximum Counter value.
func (c *counter) Max() float64 { return c.max }

func (c *counter) Update(i float64) {

	c.count++
	c.sum = c.sum + i

	if i > c.max {
		c.max = i
	}

	if i < c.min {
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
		return c.Min()
	case Max:
		return c.Max()
	case Mean:
		return c.Mean()
	case Count:
		return c.Count()
	case Sum:
		return c.Sum()
	case SumSq:
		return c.SumSq()
	case StdDev:
		return c.StdDev()
	default:
		return 0
	}
}

func createSlice(idx []byte, ctx Context) interface{} {
	c := ctx.Segment().Col(idx)

	// Not found this should be an aggregated column.
	str := string(idx)
	if c == nil {
		if i := strings.IndexByte(str, '_'); i >= 0 {
			str = str[:i]
		}
	}

	// try to find the column
	c = ctx.Segment().Col([]byte(str))
	switch c.(type) {

	case storage.FloatColumn:
		return make([]float64, 0)
	case storage.IntColumn:
		return make([]int, 0)
	case storage.StringColumn:
		return make([][]byte, 0)
	case storage.TextColumn:
		return make([][]byte, 0)
	}
	panic(" No mapping Found.")
}

func createResultVector(n []storage.Vector, keyCols []int, aggCols []Aggregation, i int) []interface{} {
	// Create the result array
	c := make([]interface{}, 0)

	// append the keys
	for _, it := range keyCols {
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
	for range aggCols {
		c = append(c, NewCounter())
	}

	return c
}

type Aggregation struct {
	AggType AggType
	AggCol  int
}

func createNewKeyMap(ctx Context, keyCols []int, aggCols []Aggregation) map[int][]byte {
	// new index column map
	nkv := make(map[int][]byte)
	// old index column map
	if kv, ok := ctx.Get(ColumnIndexToColumnName); ok == true {

		okv := kv.(map[int][]byte)

		// for each key (group by) col we a create a result vector
		last := 0
		for i := 0; i < len(keyCols); i++ {
			nkv[i] = okv[keyCols[i]]
			last = i + 1
		}

		// for each aggregated column we create a result vector
		for x := 0; x < len(aggCols); x++ {

			agg := aggCols[x]

			resName := make([]byte, 0)
			resName = append(resName, okv[agg.AggCol]...)
			resName = append(resName, '_')
			resName = append(resName, []byte(agg.AggType.String())...)

			nkv[last+x] = resName
		}

	} else {
		panic("Error parsing column keys")
	}

	return nkv

}
