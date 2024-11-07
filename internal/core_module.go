package internal

type Import struct {
	Path  string
	Alias string
}

type Module struct {
	Scope     *Scope
	Imports   []*Import
	Types     []*Node
	Functions []*Node
	Variables []*Node
	Temp      *Node
}
