package semantic

import (
	"github.com/renatopp/golden/internal/compiler/semantic/types"
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
)

var Int, Float, String, Bool, Void core.TypeData

func init() {
	Int = types.NewPrimitive("Int", func() (core.AstData, error) { return &ast.Int{Value: 0}, nil })
	Float = types.NewPrimitive("Float", func() (core.AstData, error) { return &ast.Float{Value: 0}, nil })
	String = types.NewPrimitive("String", func() (core.AstData, error) { return &ast.String{Value: ""}, nil })
	Bool = types.NewPrimitive("Bool", func() (core.AstData, error) { return &ast.Bool{Value: false}, nil })
	Void = types.NewVoid()
}
