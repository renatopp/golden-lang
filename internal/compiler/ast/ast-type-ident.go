package ast

import "github.com/renatopp/golden/lang"

type TypeIdent struct {
	*BaseNode
	Literal string
}

func NewTypeIdent(token *lang.Token, literal string) *TypeIdent {
	return &TypeIdent{
		BaseNode: NewBaseNode(
			token,
		),
		Literal: literal,
	}
}

func (n *TypeIdent) Accept(v Visitor) { v.VisitTypeIdent(n) }
