package ast

import "github.com/renatopp/golden/lang"

//
//
//

type Module struct {
	*BaseNode
	ModulePath string
}

func NewModule(token *lang.Token, path string) *Module {
	return &Module{
		BaseNode: NewBaseNode(
			token,
		),
		ModulePath: path,
	}
}

func (n *Module) Accept(v Visitor) {
	v.VisitModule(n)
}
