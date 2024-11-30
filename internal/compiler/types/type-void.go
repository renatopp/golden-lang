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

func (t *Unit) GetSignature() string { return "Void" }
func (t *Unit) GetDefault() (ast.Node, error) {
	return nil, fmt.Errorf("cannot create a default value for void type")
}
func (t *Unit) IsCompatible(other ast.Type) bool { return true }
