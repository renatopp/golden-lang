package ast

import "github.com/renatopp/golden/lang"

type UnaryOp struct {
	*BaseNode
	Operator string
	Right    Node
}

func NewUnaryOp(token *lang.Token, op string, right Node) *UnaryOp {
	return &UnaryOp{
		BaseNode: NewBaseNode(
			ValueExpressionKind,
			token,
		),
		Operator: op,
		Right:    right,
	}
}

func (n *UnaryOp) Accept(v Visitor) { v.VisitUnaryOp(n) }
