package ast

import (
	"github.com/renatopp/golden/internal/helpers/safe"
	"github.com/renatopp/golden/lang"
)

type Import struct {
	*BaseNode
	Path  *String
	Alias safe.Optional[*VarIdent] // Optional
}

func NewImport(token *lang.Token, path *String, alias safe.Optional[*VarIdent]) *Import {
	return &Import{
		BaseNode: NewBaseNode(
			token,
		),
		Path:  path,
		Alias: alias,
	}
}

func (n *Import) Accept(v Visitor) { v.VisitImport(n) }
