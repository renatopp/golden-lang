package ast

import (
	"github.com/renatopp/golden/internal/helpers/safe"
	"github.com/renatopp/golden/lang"
)

type FuncDecl struct {
	*BaseNode
	Name   safe.Optional[*VarIdent]
	Params []*FuncDeclParam
	Return safe.Optional[Node]
	Body   *Block
}

func NewFuncDecl(
	token *lang.Token,
	name safe.Optional[*VarIdent],
	params []*FuncDeclParam,
	ret safe.Optional[Node],
	body *Block,
) *FuncDecl {
	return &FuncDecl{
		BaseNode: NewBaseNode(
			ValueExpressionKind,
			token,
		),
		Name:   name,
		Params: params,
		Return: ret,
		Body:   body,
	}
}

func (n *FuncDecl) Accept(v Visitor) { v.VisitFuncDecl(n) }

type FuncDeclParam struct {
	*BaseNode
	Index    int
	Name     *VarIdent
	TypeExpr Node
}

func NewFuncDeclParam(token *lang.Token, index int, name *VarIdent, typeExpr Node) *FuncDeclParam {
	return &FuncDeclParam{
		BaseNode: NewBaseNode(
			ValueExpressionKind,
			token,
		),
		Index:    index,
		Name:     name,
		TypeExpr: typeExpr,
	}
}

func (n *FuncDeclParam) Accept(v Visitor) { v.VisitFuncDeclParam(n) }
