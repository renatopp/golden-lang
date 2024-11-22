package ast

import (
	"github.com/renatopp/golden/internal/helpers/safe"
	"github.com/renatopp/golden/lang"
)

type VarDecl struct {
	*BaseNode
	Name      *VarIdent
	TypeExpr  safe.Optional[Node]
	ValueExpr safe.Optional[Node]
}

func NewVarDecl(token *lang.Token, name *VarIdent, typeExpr, valueExpr safe.Optional[Node]) *VarDecl {
	return &VarDecl{
		BaseNode: NewBaseNode(
			ValueExpressionKind,
			token,
		),
		Name:      name,
		TypeExpr:  typeExpr,
		ValueExpr: valueExpr,
	}
}

func (n *VarDecl) Accept(v Visitor) { v.VisitVarDecl(n) }
