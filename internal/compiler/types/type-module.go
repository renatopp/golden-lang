package types

import (
	"fmt"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/env"
)

var _ ast.Type = &Module{}

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

func (m *Module) GetSignature() string {
	return fmt.Sprintf("Module('%s')", m.Path)
}

func (m *Module) GetDefault() (ast.Node, error) {
	return nil, fmt.Errorf("module type does not have a default value")
}

func (m *Module) IsCompatible(other ast.Type) bool {
	return other != nil && m.GetId() == other.GetId()
}
