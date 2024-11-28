package ast

import (
	"github.com/renatopp/golden/internal/compiler/token"
	"github.com/renatopp/golden/internal/helpers/safe"
)

type Node interface{}

type Type interface{}

type Module struct {
	Consts []Const
}

type Const struct {
	Token     token.Token
	Type      Type
	Name      VarIdent
	TypeExpr  safe.Optional[Node]
	ValueExpr Node
}

type Int struct {
	Token token.Token
	Type  Type
	Value int64
}

type Float struct {
	Token token.Token
	Type  Type
	Value float64
}

type String struct {
	Token token.Token
	Type  Type
	Value string
}

type Bool struct {
	Token token.Token
	Type  Type
	Value bool
}

type VarIdent struct {
	Token token.Token
	Type  Type
	Value string
}

type TypeIdent struct {
	Token token.Token
	Type  Type
	Value string
}

type BinOp struct {
	Type  Type
	Op    token.Token
	Left  Node
	Right Node
}

type UnaryOp struct {
	Type  Type
	Op    token.Token
	Right Node
}

type Block struct {
	Token       token.Token
	Type        Type
	Expressions []Node
}
