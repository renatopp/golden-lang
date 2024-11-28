package types

import (
	"fmt"

	"github.com/renatopp/golden/internal/compiler/ast"
)

var (
	Void *Unit
)

func init() {
	Void = NewUnit()
}

var _ ast.Type = &Unit{}

type Unit struct {
	*BaseType
}

func NewUnit() *Unit {
	return &Unit{
		BaseType: NewBaseType(nil),
	}
}

func (t *Unit) Signature() string              { return "Void" }
func (t *Unit) Compatible(other ast.Type) bool { return true }
func (t *Unit) Default() (ast.Node, error) {
	return nil, fmt.Errorf("cannot create a default value for void type")
}
