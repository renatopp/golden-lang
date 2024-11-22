package ast

import "github.com/renatopp/golden/lang"

type Block struct {
	*BaseNode
	Expressions []Node
}

func NewBlock(token *lang.Token, expressions []Node) *Block {
	return &Block{
		BaseNode: NewBaseNode(
			ValueExpressionKind,
			token,
		),
		Expressions: expressions,
	}
}

func (n *Block) Accept(v Visitor) { v.VisitBlock(n) }
