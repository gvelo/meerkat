package logical

type Fragment struct {
	IsParallel bool
	Roots      []Node
}

// Fragments is the set of parallel execution fragments created by the
// parallelizer. Nodes execute a subset of these fragments whereas the
// coordinator executes the whole set ( output fragments + node fragments )
type Fragments struct {
	fragments []*Fragment
}

// AllFragments return the whole fragment set.
func (f *Fragments) AllFragments() []*Fragment {
	return f.fragments
}

// NodeFragments return the subset of fragments that should
// be executed on nodes.
func (f *Fragments) NodeFragments() []*Fragment {
	return f.fragments[:len(f.fragments)-1]
}

func (f *Fragments) append(fragment *Fragment) {
	f.fragments = append(f.fragments, fragment)
}

// Parallelize build the plan fragments or subtrees ( actually subgraphs )
// that will be executed in the cluster nodes. Each fragment is a
// subgraph of the logical plan with the addition of the exchange operators
// needed to transfer the vectors between nodes. Fragment will be evaluated
// on the cluster nodes and transformed into physical operators.
// Currently the strategy used to parallelize the graph is quite naive, we
// plan to evolved it over time. The graph is traversed in post-order and a
// fragment is created every time a distributed operator is found ( ie.
// Summary or  Join ).
// The last fragment in the slice  always executed in the query coordinator
// node.
func Parallelize(rootNodes []Node, localNodeName string, nodeNames []string) *Fragments {

	parallelizer := NewNaiveParallelizer(localNodeName, nodeNames)

	if len(rootNodes) > 1 {
		panic("multiple rootNodes not supported yet")
	}

	// add a logical output to the plan
	output := &OutputOp{Child: rootNodes[0]}

	_ = Walk(output, parallelizer)

	return parallelizer.fragments

}

func NewNaiveParallelizer(localNodeName string, nodeNames []string) *NaiveParallelizer {
	return &NaiveParallelizer{
		localNodeName:  localNodeName,
		nodeNames:      nodeNames,
		fragments:      &Fragments{},
		inParallelFlow: true,
	}
}

type NaiveParallelizer struct {
	localNodeName  string
	nodeNames      []string
	inParallelFlow bool
	fragments      *Fragments
	streamId       int64
}

func (p *NaiveParallelizer) VisitPre(n Node) Node { return n }

func (p *NaiveParallelizer) VisitPost(n Node) Node {

	switch n := n.(type) {
	case *SummarizeOp:
		if p.inParallelFlow {
			p.inParallelFlow = false
			return p.buildDistSummary(n)
		}
		return buildLocalSummary(n)
	case *OutputOp:
		return p.buildOutput(n)
	}

	return n

}

func (p *NaiveParallelizer) buildOutput(op *OutputOp) Node {

	if !p.inParallelFlow {

		fragment := &Fragment{
			IsParallel: false,
			Roots:      []Node{op},
		}

		p.fragments.append(fragment)

		return op

	}

	streamMap := p.buildStreamMap()

	nodeOutput := &NodeOutOp{
		Dst:       p.localNodeName,
		StreamMap: streamMap,
		Child:     op.Child,
	}

	fragment := &Fragment{
		IsParallel: true,
		Roots:      []Node{nodeOutput},
	}

	p.fragments.append(fragment)

	merge := &MergeSortOp{StreamMap: streamMap}

	op.Child = merge

	outputFragment := &Fragment{
		IsParallel: false,
		Roots:      []Node{op},
	}

	p.fragments.append(outputFragment)

	return op

}

func (p *NaiveParallelizer) buildDistSummary(summaryOp *SummarizeOp) Node {
	distSummary := &DistSummaryOp{Child: summaryOp.Child}
	exchangeOut := &NodeOutOp{Child: distSummary}
	fragment := &Fragment{
		IsParallel: true,
		Roots:      []Node{exchangeOut},
	}
	p.fragments.append(fragment)
	exchangeIn := &MergeSortOp{} // TODO: refactor to summarycollector ?
	summCollector := &SummaryCollector{Child: exchangeIn}
	return summCollector
}

func (p *NaiveParallelizer) newStreamId() int64 {
	p.streamId++
	return p.streamId
}

func (p *NaiveParallelizer) buildStreamMap() map[string]int64 {

	streamMap := make(map[string]int64, len(p.nodeNames)+1)

	for _, name := range p.nodeNames {
		streamMap[name] = p.newStreamId()
	}

	streamMap[p.localNodeName] = p.newStreamId()

	return streamMap

}

func buildLocalSummary(summaryOp *SummarizeOp) Node {
	localSummary := &LocalSummaryOp{Child: summaryOp.Child}
	return localSummary
}
