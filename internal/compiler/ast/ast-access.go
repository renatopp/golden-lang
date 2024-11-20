package ast

import "github.com/renatopp/golden/lang"

type Access struct {
	*BaseNode
	Target   Node
	Accessor *VarIdent
}

func NewAccess(token *lang.Token, target Node, accessor *VarIdent) *Access {
	return &Access{
		BaseNode: NewBaseNode(
			token,
		),
		Target:   target,
		Accessor: accessor,
	}
}

func (n *Access) Accept(v Visitor) { v.VisitAccess(n) }
