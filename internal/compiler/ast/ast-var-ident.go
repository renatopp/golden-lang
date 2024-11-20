package ast

import "github.com/renatopp/golden/lang"

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
