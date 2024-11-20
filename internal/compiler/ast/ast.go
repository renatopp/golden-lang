package ast

import (
	"github.com/renatopp/golden/internal/helpers/safe"
	"github.com/renatopp/golden/lang"
)

//
//
//

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

//
//
//

type Import struct {
	*BaseNode
	Path  *String
	Alias safe.Optional[*VarIdent] // Optional
}

func NewImport(token *lang.Token, path *String, alias safe.Optional[*VarIdent]) *Import {
	return &Import{
		BaseNode: NewBaseNode(
			token,
		),
		Path:  path,
		Alias: alias,
	}
}

func (n *Import) Accept(v Visitor) { v.VisitImport(n) }

//
//
//

type String struct {
	*BaseNode
	Literal string
}

func NewString(token *lang.Token, literal string) *String {
	return &String{
		BaseNode: NewBaseNode(
			token,
		),
		Literal: literal,
	}
}

func (n *String) Accept(v Visitor) { v.VisitString(n) }

//
//
//

type VarIdent struct {
	*BaseNode
	Literal string
}

func NewVarIdent(token *lang.Token, literal string) *VarIdent {
	return &VarIdent{
		BaseNode: NewBaseNode(
			token,
		),
		Literal: literal,
	}
}

func (n *VarIdent) Accept(v Visitor) { v.VisitVarIdent(n) }
