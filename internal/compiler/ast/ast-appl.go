package ast

import "github.com/renatopp/golden/lang"

type Appl struct {
	*BaseNode
	Target    Node
	Arguments []*ApplArg
}

func NewAppl(token *lang.Token, target Node, args []*ApplArg) *Appl {
	return &Appl{
		BaseNode: NewBaseNode(
			token,
		),
		Target:    target,
		Arguments: args,
	}
}

func (n *Appl) Accept(v Visitor) { v.VisitAppl(n) }

type ApplArg struct {
	*BaseNode
	Index     int
	ValueExpr Node
}

func NewApplArg(token *lang.Token, index int, val Node) *ApplArg {
	return &ApplArg{
		BaseNode: NewBaseNode(
			token,
		),
		Index:     index,
		ValueExpr: val,
	}
}

func (n *ApplArg) Accept(v Visitor) { v.VisitApplArg(n) }
