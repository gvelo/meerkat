package executor

import (
	"github.com/rs/zerolog/log"
	"math"
	"meerkat/internal/schema"
	"meerkat/internal/storage"
	"meerkat/internal/storage/vector"
	"strings"
	"time"
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

func createSlice(name string, ctx Context) interface{} {
	c := ctx.Segment().Col(name)

	// Not found this should be an aggregated column.
	str := name
	if c == nil {
		if i := strings.IndexByte(str, '_'); i >= 0 {
			str = str[:i]
		}
	}

	// try to find the column
	c = ctx.Segment().Col(str)
	switch c.(type) {

	case storage.FloatColumn:
		return make([]float64, 0)
	case storage.IntColumn:
		return make([]int, 0)
	case storage.StringColumn:
		return make([][]byte, 0)
	case storage.TextColumn:
		return make([][]byte, 0)
	default:
		panic(" No mapping Found.")
	}

}

func createResultVector(n []interface{}, keyCols []int, aggCols []Aggregation, i int) (rKey []interface{}, rAgg []interface{}) {
	// Create the result array
	rKey = make([]interface{}, 0)

	// append the keys
	for _, it := range keyCols {
		switch col := n[it].(type) {
		case vector.IntVector:
			rKey = append(rKey, col.Values()[i])
		case vector.FloatVector:
			rKey = append(rKey, col.Values()[i])
		case vector.ByteSliceVector:
			rKey = append(rKey, col.Get(i))
		default:
			log.Error().Msg("Error creating result vector")
		}

	}

	rAgg = make([]interface{}, 0)
	// append the counters.
	for range aggCols {
		rAgg = append(rAgg, NewCounter())
	}

	return
}

type Aggregation struct {
	AggType AggType
	AggCol  int
}

func updateNewKeyMap(ctx Context, keyCols []int, aggCols []Aggregation) {
	// new index column map
	nkv := make([]schema.Field, 0)

	for i := 0; i < (len(keyCols) + len(aggCols)); i++ {
		nkv = append(nkv, schema.Field{})
	}

	last := 0
	for i := 0; i < len(keyCols); i++ {
		nkv[i] = ctx.GetFieldProcessed().Fields[keyCols[i]]
		last = i + 1
	}

	// for each aggregated column we create a result vector
	for x := 0; x < len(aggCols); x++ {
		agg := aggCols[x]
		f := ctx.GetFieldProcessed().Fields[agg.AggCol]
		name := f.Name + "_" + agg.AggType.String()
		nf := schema.Field{
			Id:        name,
			Name:      name,
			Desc:      name,
			IndexId:   f.IndexId,
			FieldType: f.FieldType,
			Nullable:  false,
			Created:   time.Time{},
			Updated:   time.Time{},
		}

		nkv[last+x] = nf
	}

	ctx.SetFieldProcessed(nkv)

}
