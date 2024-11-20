package ast

import "github.com/renatopp/golden/lang"

type BinaryOp struct {
	*BaseNode
	Operator string
	Left     Node
	Right    Node
}

func NewBinaryOp(token *lang.Token, op string, left, right Node) *BinaryOp {
	return &BinaryOp{
		BaseNode: NewBaseNode(
			token,
		),
		Operator: op,
		Left:     left,
		Right:    right,
	}
}

func (n *BinaryOp) Accept(v Visitor) { v.VisitBinaryOp(n) }
