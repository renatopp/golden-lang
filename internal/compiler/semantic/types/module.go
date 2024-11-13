package types

import (
	"fmt"
	"strings"

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
	return "Module(" + t.Name + ")"
}

func (t *Module) Signature() string {
	return "Module"
}

func (t *Module) Accepts(other core.TypeData) bool {
	return false
}

func (t *Module) Default() (core.AstData, error) {
	return nil, fmt.Errorf("Module does not have a default value")
}

func (t *Module) AccessValue(name string) (*core.AstNode, error) {
	if strings.HasPrefix(name, "_") {
		return nil, fmt.Errorf("value %s is private", name)
	}

	binding := t.Module.Scope.Values.Get(name)
	if binding == nil {
		return nil, fmt.Errorf("value %s not found", name)
	}
	val := binding.Node
	if val.Type() == nil {
		return val, t.Module.Resolver.Resolve(val)
	}
	return val, nil
}

func (t *Module) AccessType(name string) (core.TypeData, error) {
	if strings.HasPrefix(name, "_") {
		return nil, fmt.Errorf("type %s is private", name)
	}

	binding := t.Module.Scope.Types.Get(name)
	if binding == nil {
		return nil, fmt.Errorf("type %s not found", name)
	}
	return binding.Type, nil
}
