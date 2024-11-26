package ast

import "github.com/renatopp/golden/lang"

type Module struct {
	*BaseNode
	Path      string
	Imports   []*Import
	Functions []*FuncDecl
	Variables []*VarDecl
}

func NewModule(token *lang.Token, path string, imports []*Import, functions []*FuncDecl, variables []*VarDecl) *Module {
	return &Module{
		BaseNode: NewBaseNode(
			ValueExpressionKind,
			token,
		),
		Path:      path,
		Imports:   imports,
		Functions: functions,
		Variables: variables,
	}
}

func (n *Module) Accept(v Visitor) { v.VisitModule(n) }
