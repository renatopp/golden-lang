package ast

import "github.com/renatopp/golden/lang"

type Module struct {
	*BaseNode
	ModulePath string
	Imports    []*Import
}

func NewModule(token *lang.Token, path string, imports []*Import) *Module {
	return &Module{
		BaseNode: NewBaseNode(
			token,
		),
		ModulePath: path,
		Imports:    imports,
	}
}

func (n *Module) Accept(v Visitor) { v.VisitModule(n) }
