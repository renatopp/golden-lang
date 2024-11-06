package internal

type Import struct {
	Path  string
	Alias string
}

type Module struct {
	Imports   []*Import
	Types     []*Node
	Functions []*Node
	Variables []*Node
}
