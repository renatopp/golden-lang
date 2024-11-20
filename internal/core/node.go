package core

import "github.com/renatopp/golden/lang"

type Node interface {
	Id() uint64
	Token() *lang.Token
}

var _genericNodeId uint64

// Generic Node as basis for all nodes
type GenericNode struct {
	id    uint64
	token *lang.Token
}

func NewGenericNode(token *lang.Token) *GenericNode {
	_genericNodeId++
	return &GenericNode{id: _genericNodeId, token: token}
}

func (n *GenericNode) Id() uint64         { return n.id }
func (n *GenericNode) Token() *lang.Token { return n.token }
