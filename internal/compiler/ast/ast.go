package ast

import (
	"github.com/renatopp/golden/internal/compiler/token"
	"github.com/renatopp/golden/internal/helpers/safe"
)

type Node interface {
	IsEqual(n Node) bool
	GetId() uint64
	SetToken(tok *token.Token)
	GetToken() *token.Token
	SetType(Type)
	GetType() safe.Optional[Type]
	Visit(Visitor) Node
}

type Type interface {
	GetId() uint64
	GetDefinition() Node
	GetSignature() string
	GetDefault() (Node, error)
	IsCompatible(Type) bool
}

//
//
//

var _nodeId = uint64(0)

type BaseNode struct {
	Id    uint64
	Token *token.Token
	Type  safe.Optional[Type]
}

func NewBaseNode(tok *token.Token) BaseNode {
	_nodeId++
	return BaseNode{
		Id:    _nodeId,
		Token: tok,
		Type:  safe.None[Type](),
	}
}

func (n *BaseNode) IsEqual(other Node) bool      { return n.Id == other.GetId() }
func (n *BaseNode) GetId() uint64                { return n.Id }
func (n *BaseNode) SetToken(tok *token.Token)    { n.Token = tok }
func (n *BaseNode) GetToken() *token.Token       { return n.Token }
func (n *BaseNode) SetType(tp Type)              { n.Type = safe.Some(tp) }
func (n *BaseNode) GetType() safe.Optional[Type] { return n.Type }
func (n *BaseNode) Visit(v Visitor) Node {
	panic("base node does not have visitor")
}

//
//
//

type Module struct {
	BaseNode
	Exprs []Node
}

func NewModule(tok *token.Token, exprs []Node) *Module { return &Module{NewBaseNode(tok), exprs} }
func (n *Module) Visit(v Visitor) Node                 { return v.VisitModule(n) }

//
//
//

type Const struct {
	BaseNode
	Name      *VarIdent
	TypeExpr  safe.Optional[Node]
	ValueExpr Node
}

func NewConst(tok *token.Token, name *VarIdent, tpexpr safe.Optional[Node], valexpr Node) *Const {
	return &Const{NewBaseNode(tok), name, tpexpr, valexpr}
}
func (n *Const) Visit(v Visitor) Node { return v.VisitConst(n) }

//
//
//

type Int struct {
	BaseNode
	Value int64
}

func NewInt(tok *token.Token, val int64) *Int { return &Int{NewBaseNode(tok), val} }
func (n *Int) Visit(v Visitor) Node           { return v.VisitInt(n) }

//
//
//

type Float struct {
	BaseNode
	Value float64
}

func NewFloat(tok *token.Token, val float64) *Float { return &Float{NewBaseNode(tok), val} }
func (n *Float) Visit(v Visitor) Node               { return v.VisitFloat(n) }

//
//
//

type String struct {
	BaseNode
	Value string
}

func NewString(tok *token.Token, val string) *String { return &String{NewBaseNode(tok), val} }
func (n *String) Visit(v Visitor) Node               { return v.VisitString(n) }

//
//
//

type Bool struct {
	BaseNode
	Value bool
}

func NewBool(tok *token.Token, val bool) *Bool { return &Bool{NewBaseNode(tok), val} }
func (n *Bool) Visit(v Visitor) Node           { return v.VisitBool(n) }

//
//
//

type VarIdent struct {
	BaseNode
	Value string
}

func NewVarIdent(tok *token.Token, val string) *VarIdent { return &VarIdent{NewBaseNode(tok), val} }
func (n *VarIdent) Visit(v Visitor) Node                 { return v.VisitVarIdent(n) }

//
//
//

type TypeIdent struct {
	BaseNode
	Value string
}

func NewTypeIdent(tok *token.Token, val string) *TypeIdent { return &TypeIdent{NewBaseNode(tok), val} }
func (n *TypeIdent) Visit(v Visitor) Node                  { return v.VisitTypeIdent(n) }

//
//
//

type BinOp struct {
	BaseNode
	Op        string
	LeftExpr  Node
	RightExpr Node
}

func NewBinOp(tok *token.Token, op string, left, right Node) *BinOp {
	return &BinOp{NewBaseNode(tok), op, left, right}
}
func (n *BinOp) Visit(v Visitor) Node { return v.VisitBinOp(n) }

//
//
//

type UnaryOp struct {
	BaseNode
	Op        string
	RightExpr Node
}

func NewUnaryOp(tok *token.Token, op string, right Node) *UnaryOp {
	return &UnaryOp{NewBaseNode(tok), op, right}
}
func (n *UnaryOp) Visit(v Visitor) Node { return v.VisitUnaryOp(n) }

//
//
//

type Block struct {
	BaseNode
	Exprs []Node
}

func NewBlock(tok *token.Token, exprs []Node) *Block { return &Block{NewBaseNode(tok), exprs} }
func (n *Block) Visit(v Visitor) Node                { return v.VisitBlock(n) }
