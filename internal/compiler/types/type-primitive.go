package types

import (
	"github.com/renatopp/golden/internal/compiler/ast"
)

var (
	Zero        *ast.Int
	One         *ast.Int
	FZero       *ast.Float
	FOne        *ast.Float
	False       *ast.Bool
	True        *ast.Bool
	EmptyString *ast.String

	Int    *Primitive
	Float  *Primitive
	Bool   *Primitive
	String *Primitive
)

func init() {
	Zero = ast.NewInt(nil, 0)
	One = ast.NewInt(nil, 1)
	FZero = ast.NewFloat(nil, 0)
	FOne = ast.NewFloat(nil, 1)
	False = ast.NewBool(nil, false)
	True = ast.NewBool(nil, true)
	EmptyString = ast.NewString(nil, "")

	Int = NewPrimitive("Int", func() (ast.Node, error) { return Zero, nil })
	Float = NewPrimitive("Float", func() (ast.Node, error) { return FZero, nil })
	Bool = NewPrimitive("Bool", func() (ast.Node, error) { return False, nil })
	String = NewPrimitive("String", func() (ast.Node, error) { return EmptyString, nil })
}

//
//
//

var _ ast.Type = &Primitive{}

type Primitive struct {
	*BaseType
	Name      string
	DefaultFn func() (ast.Node, error)
}

func NewPrimitive(name string, fn func() (ast.Node, error)) *Primitive {
	return &Primitive{
		BaseType:  NewBaseType(nil),
		Name:      name,
		DefaultFn: fn,
	}
}

func (p *Primitive) Signature() string {
	return p.Name
}

func (p *Primitive) Compatible(other ast.Type) bool {
	return other != nil && p.Id() == other.Id()
}

func (p *Primitive) Default() (ast.Node, error) {
	return p.DefaultFn()
}
