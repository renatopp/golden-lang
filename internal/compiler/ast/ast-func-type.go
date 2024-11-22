package ast

import (
	"github.com/renatopp/golden/internal/helpers/safe"
	"github.com/renatopp/golden/lang"
)

type FuncType struct {
	*BaseNode
	Params []*FuncTypeParam
	Return safe.Optional[Node]
}

func NewFuncType(token *lang.Token, params []*FuncTypeParam, ret safe.Optional[Node]) *FuncType {
	return &FuncType{
		BaseNode: NewBaseNode(
			ValueExpressionKind,
			token,
		),
		Params: params,
		Return: ret,
	}
}

func (n *FuncType) Accept(v Visitor) { v.VisitFuncType(n) }

type FuncTypeParam struct {
	*BaseNode
	Index    int
	TypeExpr Node
}

func NewFuncTypeParam(token *lang.Token, index int, tpExpr Node) *FuncTypeParam {
	return &FuncTypeParam{
		BaseNode: NewBaseNode(
			ValueExpressionKind,
			token,
		),
		Index:    index,
		TypeExpr: tpExpr,
	}
}

func (n *FuncTypeParam) Accept(v Visitor) { v.VisitFuncTypeParam(n) }
