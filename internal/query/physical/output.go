package physical

type OutputOp interface {
	// Run the graph and write the output the client.
	Run()
}
