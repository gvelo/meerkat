package parser

// Analyze runs all the semantic validations rules.
//
// - type checking
// - type inference
// - column name resolution ( when possible )
// - function name resolution
//
// TODO(gvelo) TabularStmt will be replaced with a new ast root
// that include suport for other statements types.
func Analyze(ast *TabularStmt) (*TabularStmt, error) {

	// do nothing for now.
	return ast, nil

}
