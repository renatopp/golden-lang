package types

import "github.com/renatopp/golden/internal/core"

type Primitive struct {
	baseType
	Name      string
	DefaultFn func() (core.AstData, error)
}

func NewPrimitive(name string, defaultFn func() (core.AstData, error)) *Primitive {
	return &Primitive{
		baseType:  newBase(),
		Name:      name,
		DefaultFn: defaultFn,
	}
}

func (t *Primitive) Tag() string {
	return t.Name
}

func (t *Primitive) Signature() string {
	return t.Name
}

func (t *Primitive) Accepts(other core.TypeData) bool {
	if t == nil || other == nil {
		return false
	}
	return t.Id() == other.Id()
}

func (t *Primitive) Default() (core.AstData, error) {
	return t.DefaultFn()
}
