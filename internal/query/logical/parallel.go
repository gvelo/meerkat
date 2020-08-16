package logical

func Parallelize(roots []Node) []*Fragment {

	parallelizer := NewNaiveParallelizer()

	if len(roots) > 1 {
		panic("multiple roots not supported yet")
	}

	output := &OutputOp{Child: roots[0]}

	_ = Walk(output, parallelizer)

	return parallelizer.fagments

}

func NewNaiveParallelizer() *NaiveParallelizer {
	return &NaiveParallelizer{isParallel: true}
}

type NaiveParallelizer struct {
	isParallel bool
	fagments   []*Fragment
}

func (p *NaiveParallelizer) VisitPre(n Node) Node { return n }

func (p *NaiveParallelizer) VisitPost(n Node) Node {

	switch n := n.(type) {
	case *SummarizeOp:
		if p.isParallel {
			p.isParallel = false
			return p.buildDistSummary(n)
		}
		return buildLocalSummary(n)
	case *OutputOp:
		return p.buildOutput(n)
	}

	return n

}

func (p *NaiveParallelizer) buildOutput(op *OutputOp) Node {

	if !p.isParallel {
		fragment := &Fragment{
			IsParallel: false,
			Roots:      []Node{op},
		}
		p.addFragment(fragment)
		return op
	}

	exchange := &ExchangeOutOp{Child: op.Child}
	fragment := &Fragment{
		IsParallel: true,
		Roots:      []Node{exchange},
	}
	p.addFragment(fragment)

	merge := &MergeOp{}
	op.Child = merge

	rootFragment := &Fragment{
		IsParallel: false,
		Roots:      []Node{op},
	}

	p.addFragment(rootFragment)

	return op

}

func (p *NaiveParallelizer) buildDistSummary(summaryOp *SummarizeOp) Node {
	distSummary := &DistSummaryOp{Child: summaryOp.Child}
	exchangeOut := &ExchangeOutOp{Child: distSummary}
	fragment := &Fragment{
		IsParallel: true,
		Roots:      []Node{exchangeOut},
	}
	p.addFragment(fragment)
	exchangeIn := &ExchangeInOp{}
	summCollector := &SummaryCollector{Child: exchangeIn}
	return summCollector
}

func buildLocalSummary(summaryOp *SummarizeOp) Node {
	localSummary := &LocalSummaryOp{Child: summaryOp.Child}
	return localSummary
}

func (p *NaiveParallelizer) addFragment(f *Fragment) {
	p.fagments = append(p.fagments, f)
}
