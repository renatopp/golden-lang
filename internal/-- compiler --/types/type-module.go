package types

import (
	"fmt"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/env"
)

type Module struct {
	*BaseType
	Path  string
	Scope *env.Scope
}

func NewModule(node *ast.Module, path string, scope *env.Scope) *Module {
	return &Module{
		BaseType: NewBaseType(node),
		Path:     path,
		Scope:    scope,
	}
}

func (m *Module) Signature() string {
	return fmt.Sprintf("Module('%s')", m.Path)
}

func (m *Module) Compatible(other ast.Type) bool {
	return other != nil && m.Id() == other.Id()
}

func (m *Module) Default() (ast.Node, error) {
	return nil, fmt.Errorf("module type does not have a default value")
}
