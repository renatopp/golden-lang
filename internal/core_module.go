package internal

type Import struct {
	Path  string
	Alias string
}

type Module struct {
	Name        string
	Path        string
	FileName    string
	PackageName string
	PackagePath string
	Scope       *Scope
	Imports     []*Import
	Types       []*Node
	Functions   []*Node
	Variables   []*Node
	Temp        *Node
}
