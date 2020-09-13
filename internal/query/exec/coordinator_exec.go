package exec

import (
	"context"
	"meerkat/internal/query/logical"
	"meerkat/internal/query/parser"
)

func NewCoordinatorExecutor() *coordinatorExecutor {
	return &coordinatorExecutor{
		execContext:NewExecutionContext()
	}
}

type coordinatorExecutor struct {
	// done channel
	// cancel function
	// port registry
	// out channel
	// cluster
	// SegmentRegistry.
	// Wg for drain
	// state

	execContext ExecutionContext



}

func (c *coordinatorExecutor) Cancel() {

	// once set error
	// cancel
	// state change

}

func (c *coordinatorExecutor) exec(query string ) error {

	// Transform the string text into a abstract sintax tree
	ast, err := parser.Parse(query)

	if err != nil {
		return err
	}

	// perform the semantic validation
	err = parser.Analyze(ast)

	if err != nil {
		return err
	}

	// transform the ast into a logical query plan
	logicalPlan := logical.ToLogical(ast)

	// optimize the plan
	optPlan := logical.Optimize(logicalPlan)

	// parallelize
	fragments := logical.Parallelize(optPlan)

	// run local ( build executablegraph )
	//    create the local graph.
	//    output to channel
	//    return query object
	// distribute
	//  broadcat to all the nodes. if somethign fails cancel

	// build local
	// distribute










}

func (c *coordinatorExecutor) execLocal(root *logical.Fragment, local []*logical.Fragment) {

	graph := buildPhysicalGraph(root, local)
	// run

}

func (c *coordinatorExecutor) broadcastFragments() {

}

func (c *coordinatorExecutor) buildPhysicalGraph(root *logical.Fragment, local []*logical.Fragment) {

	// merge exchange
	// group by exchange

}
