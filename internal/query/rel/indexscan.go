package rel

type IndexScan struct {
	indexName string
	filter    *Filter
}

func NewTableScan(name string) *IndexScan {
	s := new(IndexScan)
	s.indexName = name
	return s
}

func (i *IndexScan) SetFilter(f *Filter) {
	i.filter = f
}
