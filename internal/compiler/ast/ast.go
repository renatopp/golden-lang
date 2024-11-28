package ast

import (
	"github.com/renatopp/golden/internal/compiler/token"
	"github.com/renatopp/golden/internal/helpers/safe"
)

type Node interface {
	GetToken() token.Token
	Visit(Visitor) Node
}

type Type interface{}

//
//
//

type Module struct {
	Token  token.Token
	Consts []Const
}

func (n Module) GetToken() token.Token { return n.Token }
func (n Module) Visit(v Visitor) Node  { return v.VisitModule(n) }

type Const struct {
	Token     token.Token
	Type      safe.Optional[Type]
	Name      VarIdent
	TypeExpr  safe.Optional[Node]
	ValueExpr Node
}

func (n Const) GetToken() token.Token { return n.Token }
func (n Const) Visit(v Visitor) Node  { return v.VisitConst(n) }

type Int struct {
	Token token.Token
	Type  safe.Optional[Type]
	Value int64
}

func (n Int) GetToken() token.Token { return n.Token }
func (n Int) Visit(v Visitor) Node  { return v.VisitInt(n) }

type Float struct {
	Token token.Token
	Type  safe.Optional[Type]
	Value float64
}

func (n Float) GetToken() token.Token { return n.Token }
func (n Float) Visit(v Visitor) Node  { return v.VisitFloat(n) }

type String struct {
	Token token.Token
	Type  safe.Optional[Type]
	Value string
}

func (n String) GetToken() token.Token { return n.Token }
func (n String) Visit(v Visitor) Node  { return v.VisitString(n) }

type Bool struct {
	Token token.Token
	Type  safe.Optional[Type]
	Value bool
}

func (n Bool) GetToken() token.Token { return n.Token }
func (n Bool) Visit(v Visitor) Node  { return v.VisitBool(n) }

type VarIdent struct {
	Token token.Token
	Type  safe.Optional[Type]
	Value string
}

func (n VarIdent) GetToken() token.Token { return n.Token }
func (n VarIdent) Visit(v Visitor) Node  { return v.VisitVarIdent(n) }

type TypeIdent struct {
	Token token.Token
	Type  safe.Optional[Type]
	Value string
}

func (n TypeIdent) GetToken() token.Token { return n.Token }
func (n TypeIdent) Visit(v Visitor) Node  { return v.VisitTypeIdent(n) }

type BinOp struct {
	Type  safe.Optional[Type]
	Op    token.Token
	Left  Node
	Right Node
}

func (n BinOp) GetToken() token.Token { return n.Op }
func (n BinOp) Visit(v Visitor) Node  { return v.VisitBinOp(n) }

type UnaryOp struct {
	Type  safe.Optional[Type]
	Op    token.Token
	Right Node
}

func (n UnaryOp) GetToken() token.Token { return n.Op }
func (n UnaryOp) Visit(v Visitor) Node  { return v.VisitUnaryOp(n) }

type Block struct {
	Token       token.Token
	Type        safe.Optional[Type]
	Expressions []Node
}

func (n Block) GetToken() token.Token { return n.Token }
func (n Block) Visit(v Visitor) Node  { return v.VisitBlock(n) }
