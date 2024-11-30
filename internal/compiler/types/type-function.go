package types

import (
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/compiler/ast"
)

var NoopFn = NewFunction(nil, []ast.Type{}, Void)

var _ ast.Type = &Function{}

type Function struct {
	*BaseType
	Params []ast.Type
	Return ast.Type
}

func NewFunction(def ast.Node, parameters []ast.Type, returnType ast.Type) *Function {
	return &Function{
		BaseType: NewBaseType(def),
		Params:   parameters,
		Return:   returnType,
	}
}

func (f *Function) GetSignature() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = p.GetSignature()
	}
	p := strings.Join(params, ", ")

	ret := ""
	if f.Return != nil {
		ret = " " + f.Return.GetSignature()
	}

	return fmt.Sprintf("Fn (%s)%s", p, ret)
}

func (f *Function) GetDefault() (ast.Node, error) {
	return nil, fmt.Errorf("functions cannot have default values")
}

func (f *Function) IsCompatible(t ast.Type) bool {
	fn, ok := t.(*Function)
	if !ok {
		return false
	}

	if len(f.Params) != len(fn.Params) {
		return false
	}

	for i, p := range f.Params {
		if !p.IsCompatible(fn.Params[i]) {
			return false
		}
	}

	if f.Return != nil && fn.Return != nil {
		return f.Return.IsCompatible(fn.Return)
	}

	return true
}
