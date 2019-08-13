package rel

type LiteralNode struct {
	value interface{}
}

type VarNode struct {
	varName string
}

type ColNode struct {
	colName string
}
