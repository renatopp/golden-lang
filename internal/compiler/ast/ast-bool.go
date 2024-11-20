package ast

import "github.com/renatopp/golden/lang"

type Bool struct {
	*BaseNode
	Literal bool
}

func NewBool(token *lang.Token, literal bool) *Bool {
	return &Bool{
		BaseNode: NewBaseNode(
			token,
		),
		Literal: literal,
	}
}

func (n *Bool) Accept(v Visitor) { v.VisitBool(n) }
