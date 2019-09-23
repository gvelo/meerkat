package logical

// logical aggregation be implemented by all functions
type Aggregation struct {
	limit    int // limited top x
	function string
}
