package ast

import "github.com/renatopp/golden/lang"

type Int struct {
	*BaseNode
	Literal int64
}

func NewInt(token *lang.Token, literal int64) *Int {
	return &Int{
		BaseNode: NewBaseNode(
			token,
		),
		Literal: literal,
	}
}

func (n *Int) Accept(v Visitor) { v.VisitInt(n) }
