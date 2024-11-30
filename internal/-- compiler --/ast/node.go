package ast

import "github.com/renatopp/golden/lang"

type ExpressionKind int

const (
	UnknownExpressionKind ExpressionKind = iota
	TypeExpressionKind
	ValueExpressionKind
)

type Type interface {
	Id() uint64
	Definition() Node
	Signature() string
	Compatible(Type) bool
	Default() (Node, error)
}

type Node interface {
	Id() uint64
	Token() *lang.Token
	Accept(Visitor)
	Type() Type
	SetConstant()
	IsConstant() bool
	SetType(Type)
	ExpressionKind() ExpressionKind
	SetExpressionKind(ExpressionKind)
}

//
// Base Node
//

var _baseNodeId uint64

type BaseNode struct {
	id             uint64
	token          *lang.Token
	type_          Type
	constant       bool
	expressionKind ExpressionKind
}

func NewBaseNode(expressionKind ExpressionKind, token *lang.Token) *BaseNode {
	_baseNodeId++
	return &BaseNode{
		id:             _baseNodeId,
		token:          token,
		expressionKind: expressionKind,
	}
}

func (n *BaseNode) Id() uint64                            { return n.id }
func (n *BaseNode) Token() *lang.Token                    { return n.token }
func (n *BaseNode) Type() Type                            { return n.type_ }
func (n *BaseNode) SetConstant()                          { n.constant = true }
func (n *BaseNode) IsConstant() bool                      { return n.constant }
func (n *BaseNode) SetType(t Type)                        { n.type_ = t }
func (n *BaseNode) ExpressionKind() ExpressionKind        { return n.expressionKind }
func (n *BaseNode) SetExpressionKind(kind ExpressionKind) { n.expressionKind = kind }
