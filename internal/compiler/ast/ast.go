package ast

import (
	"github.com/renatopp/golden/internal/compiler/token"
	"github.com/renatopp/golden/internal/helpers/safe"
)

type Node interface {
	GetNotes() Annotations
	GetToken() token.Token
	Visit(Visitor) Node
}

type Type interface{}

type Annotations struct {
	Type safe.Optional[Type]
}

func (a Annotations) WithType(t Type) Annotations { a.Type = safe.Some(t); return a }

//
//
//

type Module struct {
	Token  token.Token
	Notes  Annotations
	Consts []Const
}

func (n Module) GetNotes() Annotations { return n.Notes }
func (n Module) GetToken() token.Token { return n.Token }
func (n Module) Visit(v Visitor) Node  { return v.VisitModule(n) }

type Const struct {
	Token     token.Token
	Notes     Annotations
	Name      VarIdent
	TypeExpr  safe.Optional[Node]
	ValueExpr Node
}

func (n Const) GetNotes() Annotations { return n.Notes }
func (n Const) GetToken() token.Token { return n.Token }
func (n Const) Visit(v Visitor) Node  { return v.VisitConst(n) }

type Int struct {
	Token token.Token
	Notes Annotations
	Value int64
}

func (n Int) GetNotes() Annotations { return n.Notes }
func (n Int) GetToken() token.Token { return n.Token }
func (n Int) Visit(v Visitor) Node  { return v.VisitInt(n) }

type Float struct {
	Token token.Token
	Notes Annotations
	Value float64
}

func (n Float) GetNotes() Annotations { return n.Notes }
func (n Float) GetToken() token.Token { return n.Token }
func (n Float) Visit(v Visitor) Node  { return v.VisitFloat(n) }

type String struct {
	Token token.Token
	Notes Annotations
	Value string
}

func (n String) GetNotes() Annotations { return n.Notes }
func (n String) GetToken() token.Token { return n.Token }
func (n String) Visit(v Visitor) Node  { return v.VisitString(n) }

type Bool struct {
	Token token.Token
	Notes Annotations
	Value bool
}

func (n Bool) GetNotes() Annotations { return n.Notes }
func (n Bool) GetToken() token.Token { return n.Token }
func (n Bool) Visit(v Visitor) Node  { return v.VisitBool(n) }

type VarIdent struct {
	Token token.Token
	Notes Annotations
	Value string
}

func (n VarIdent) GetNotes() Annotations { return n.Notes }
func (n VarIdent) GetToken() token.Token { return n.Token }
func (n VarIdent) Visit(v Visitor) Node  { return v.VisitVarIdent(n) }

type TypeIdent struct {
	Token token.Token
	Notes Annotations
	Value string
}

func (n TypeIdent) GetNotes() Annotations { return n.Notes }
func (n TypeIdent) GetToken() token.Token { return n.Token }
func (n TypeIdent) Visit(v Visitor) Node  { return v.VisitTypeIdent(n) }

type BinOp struct {
	Token token.Token
	Notes Annotations
	Op    string
	Left  Node
	Right Node
}

func (n BinOp) GetNotes() Annotations { return n.Notes }
func (n BinOp) GetToken() token.Token { return n.Token }
func (n BinOp) Visit(v Visitor) Node  { return v.VisitBinOp(n) }

type UnaryOp struct {
	Token token.Token
	Notes Annotations
	Op    string
	Right Node
}

func (n UnaryOp) GetNotes() Annotations { return n.Notes }
func (n UnaryOp) GetToken() token.Token { return n.Token }
func (n UnaryOp) Visit(v Visitor) Node  { return v.VisitUnaryOp(n) }

type Block struct {
	Token       token.Token
	Notes       Annotations
	Expressions []Node
}

func (n Block) GetNotes() Annotations { return n.Notes }
func (n Block) GetToken() token.Token { return n.Token }
func (n Block) Visit(v Visitor) Node  { return v.VisitBlock(n) }
