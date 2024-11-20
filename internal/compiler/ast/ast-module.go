package ast

import "github.com/renatopp/golden/lang"

type Module struct {
	*BaseNode
	ModulePath string
	Imports    []*Import
	Variables  []*VarDecl
}

func NewModule(token *lang.Token, path string, imports []*Import, variables []*VarDecl) *Module {
	return &Module{
		BaseNode: NewBaseNode(
			token,
		),
		ModulePath: path,
		Imports:    imports,
		Variables:  variables,
	}
}

func (n *Module) Accept(v Visitor) { v.VisitModule(n) }
