package ast

import "github.com/renatopp/golden/lang"

type Node interface {
	Id() uint64
	Token() *lang.Token
	Accept(Visitor)
}

//
// Base Node
//

var _baseNodeId uint64

type BaseNode struct {
	id    uint64
	token *lang.Token
}

func NewBaseNode(token *lang.Token) *BaseNode {
	_baseNodeId++
	return &BaseNode{id: _baseNodeId, token: token}
}

func (n *BaseNode) Id() uint64         { return n.id }
func (n *BaseNode) Token() *lang.Token { return n.token }