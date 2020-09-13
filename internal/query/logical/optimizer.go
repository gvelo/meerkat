package logical

// Optimize rewrite and optimize the logical tree.
// Currently a very small set of heuristic rules
// is supported.
//
// - fold constant
// - evaluate functions with constant arguments ie. ago(1h)
// - capture/freeze timestamps
// - pushdown predicates or filters.
// = collapse filter
func Optimize(logicalTree []Node) []Node {

	// do nothing for now
	return logicalTree
}
