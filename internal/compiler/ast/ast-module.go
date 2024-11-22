package ast

import "github.com/renatopp/golden/lang"

type Module struct {
	*BaseNode
	ModulePath string
	Imports    []*Import
	Functions  []*FuncDecl
	Variables  []*VarDecl
}

func NewModule(token *lang.Token, path string, imports []*Import, functions []*FuncDecl, variables []*VarDecl) *Module {
	return &Module{
		BaseNode: NewBaseNode(
			ValueExpressionKind,
			token,
		),
		ModulePath: path,
		Imports:    imports,
		Functions:  functions,
		Variables:  variables,
	}
}

func (n *Module) Accept(v Visitor) { v.VisitModule(n) }
