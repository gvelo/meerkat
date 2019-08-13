package rel

import "meerkat/internal/tools"

type Builder interface {
	Scan(name string) Builder
	Filter(filter interface{}) Builder
	Project(e ...interface{}) Builder
	Aggregate(groupKey string, aggCall ...interface{}) Builder
	Distinct() Builder
	Sort(exp ...interface{}) Builder
	Limit(offset int) Builder
	SemiJoin(expr interface{}) Builder
	AntiJoin(expr interface{}) Builder
	Union(expr interface{}) Builder
	Intersect(expr interface{}) Builder
	Minus(expr interface{}) Builder
	Match(regex string) Builder
	And(a interface{}) Builder
	Or(o interface{}) Builder
	Build() *ParsedTree
}

func NewRelBuilder() Builder {
	b := new(relationalAlgBuilder)
	b.queue = make([]interface{}, 0)
	return b
}

type relationalAlgBuilder struct {
	queue []interface{}
}

func (r *relationalAlgBuilder) push(n interface{}) {
	r.queue = append(r.queue, n)
}

func (r *relationalAlgBuilder) peek() interface{} {
	return r.queue[0]
}

func (r *relationalAlgBuilder) pop() interface{} {
	e := r.queue[0]
	r.queue = r.queue[1:]
	return e
}

func (r *relationalAlgBuilder) Scan(name string) Builder {
	ts := NewTableScan(name)
	r.push(ts)
	return r
}

func (r *relationalAlgBuilder) And(a interface{}) Builder {
	return r
}

func (r *relationalAlgBuilder) Or(o interface{}) Builder {
	return r
}

func (r *relationalAlgBuilder) Filter(f interface{}) Builder {
	is := r.peek().(*IndexScan)
	is.SetFilter(f.(*Filter))
	return r
}

func (r *relationalAlgBuilder) Project(e ...interface{}) Builder {
	return r
}

func (r *relationalAlgBuilder) Aggregate(groupKey string, aggCall ...interface{}) Builder {
	return r
}

func (r *relationalAlgBuilder) Distinct() Builder {
	return r
}

func (r *relationalAlgBuilder) Sort(exp ...interface{}) Builder {
	return r
}

func (r *relationalAlgBuilder) Limit(offset int) Builder {
	return r
}

func (r *relationalAlgBuilder) SemiJoin(expr interface{}) Builder {
	tools.Log("implement me")
	return r
}

func (r *relationalAlgBuilder) AntiJoin(expr interface{}) Builder {
	tools.Log("implement me")
	return r
}

func (r *relationalAlgBuilder) Union(expr interface{}) Builder {
	tools.Log("implement me")
	return r
}

func (r *relationalAlgBuilder) Intersect(expr interface{}) Builder {
	tools.Log("implement me")
	return r
}

func (r *relationalAlgBuilder) Minus(expr interface{}) Builder {
	tools.Log("implement me")
	return r
}

func (r *relationalAlgBuilder) Match(regex string) Builder {
	tools.Log("implement me")
	return r
}

func (r *relationalAlgBuilder) Build() *ParsedTree {
	return &ParsedTree{IndexScan: r.peek().(*IndexScan)}
}
