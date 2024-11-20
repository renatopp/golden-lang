package syntax

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/core"
)

func (p *parser) parseModule() *ast.Module {
	return &ast.Module{
		GenericNode: core.NewGenericNode(p.PeekToken()),
		Module:      p.Module,
	}
}
