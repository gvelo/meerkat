package executor

import "math"

type AggregationOperator struct {
}

// Type defines an aggregation function.
type Type int

// Supported aggregation types.
const (
	Min Type = iota
	Max
	Mean
	Count
	Sum
	SumSq
	StdDev
)

type Counter interface {
	Count() int64
	SumSq() int64
	Sum() int64
	Stdev() float64
	Mean() float64
	Min() int64
	Update(i int64)
	ValueOf(aggType Type) float64
}

type counter struct {
	sum   int64
	sumSq int64
	count int64
	max   int64
	min   int64
	t     Type
}

func NewCounter(t Type) Counter {
	return &counter{t: t, max: math.MinInt64, min: math.MaxInt64}
}

// Counter returns the number of values received.
func (c *counter) Count() int64 { return c.count }

// Sum returns the sum of Counter values.
func (c *counter) Sum() int64 { return c.sum }

// SumSq returns the squared sum of Counter values.
func (c *counter) SumSq() int64 { return c.sumSq }

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

func stdev(count int64, sumSq, sum float64) float64 {
	div := count * (count - 1)
	if div == 0 {
		return 0.0
	}
	num := float64(count)*sumSq - sum*sum
	return math.Sqrt(num / float64(div))
}

// Min returns the minimum Counter value.
func (c *counter) Min() int64 { return c.min }

// Max returns the maximum Counter value.
func (c *counter) Max() int64 { return c.min }

func (c *counter) Update(i int64) {

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
func (c *counter) ValueOf(aggType Type) float64 {
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
func histogram(ts []int64) map[int64]int64 {

	m := make(map[int64]int64)
	for i := 0; i < len(ts); i++ {
		var s int64 = 1
		ant := ts[i]
		for ; i < len(ts) && ant == ts[i]; i++ {
			s++
		}
		m[ant] = s
	}

	return m
}
