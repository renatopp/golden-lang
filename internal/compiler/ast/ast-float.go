package ast

import "github.com/renatopp/golden/lang"

type Float struct {
	*BaseNode
	Literal float64
}

func NewFloat(token *lang.Token, literal float64) *Float {
	return &Float{
		BaseNode: NewBaseNode(
			ValueExpressionKind,
			token,
		),
		Literal: literal,
	}
}

func (n *Float) Accept(v Visitor) { v.VisitFloat(n) }
