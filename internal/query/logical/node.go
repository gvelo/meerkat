package logical

type Node interface {
	GetParent() Node
	SetParent(n Node)
	AddChild(n Node)
	GetChildren() []Node
	String() string
}
