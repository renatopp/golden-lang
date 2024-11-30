package types

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/token"
)

var (
	Zero        ast.Int
	One         ast.Int
	FZero       ast.Float
	FOne        ast.Float
	False       ast.Bool
	True        ast.Bool
	EmptyString ast.String

	Int    *Primitive
	Float  *Primitive
	Bool   *Primitive
	String *Primitive
)

func init() {
	Zero = ast.NewInt(token.Token{}, 0)
	One = ast.NewInt(token.Token{}, 1)
	FZero = ast.NewFloat(token.Token{}, 0)
	FOne = ast.NewFloat(token.Token{}, 1)
	False = ast.NewBool(token.Token{}, false)
	True = ast.NewBool(token.Token{}, true)
	EmptyString = ast.NewString(token.Token{}, "")

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

func (p *Primitive) GetSignature() string          { return p.Name }
func (p *Primitive) GetDefault() (ast.Node, error) { return p.DefaultFn() }
func (p *Primitive) IsCompatible(other ast.Type) bool {
	return other != nil && p.GetId() == other.GetId()
}
