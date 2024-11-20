package syntax

import (
	"github.com/renatopp/golden/internal/compiler/ast"
)

func (p *parser) parseModule() *ast.Module {
	return ast.NewModule(p.PeekToken(), p.ModulePath)
}
