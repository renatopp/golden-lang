package types

import (
	"fmt"

	"github.com/renatopp/golden/internal/core"
)

type Module struct {
	baseType
	Name   string
	Module *core.Module
}

func NewModule(name string, module *core.Module) *Module {
	return &Module{
		baseType: newBase(),
		Name:     name,
		Module:   module,
	}
}

func (t *Module) Tag() string {
	return t.Name
}

func (t *Module) Signature() string {
	return ""
}

func (t *Module) Accepts(other core.TypeData) bool {
	return false
}

func (t *Module) Default() (core.AstData, error) {
	return nil, fmt.Errorf("Module does not have a default value")
}

func (t *Module) AccessValue(name string) (*core.AstNode, error) {
	val := t.Module.Scope.GetValue(name)
	if val == nil {
		return nil, fmt.Errorf("value %s not found", name)
	}
	if val.Type() == nil {
		t.Module.Analyzer.ResolveValue(val)
	}
	return val, nil
}

func (t *Module) AccessType(name string) (core.TypeData, error) {
	val := t.Module.Scope.GetType(name)
	if val == nil {
		return nil, fmt.Errorf("type %s not found", name)
	}
	return val, nil
}
