package ast

import "github.com/renatopp/golden/lang"

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
