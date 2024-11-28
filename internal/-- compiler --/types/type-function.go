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

func (f *Function) Signature() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = p.Signature()
	}
	p := strings.Join(params, ", ")

	ret := ""
	if f.Return != nil {
		ret = " " + f.Return.Signature()
	}

	return fmt.Sprintf("Fn (%s)%s", p, ret)
}

func (f *Function) Compatible(t ast.Type) bool {
	fn, ok := t.(*Function)
	if !ok {
		return false
	}

	if len(f.Params) != len(fn.Params) {
		return false
	}

	for i, p := range f.Params {
		if !p.Compatible(fn.Params[i]) {
			return false
		}
	}

	if f.Return != nil && fn.Return != nil {
		return f.Return.Compatible(fn.Return)
	}

	return true
}

func (f *Function) Default() (ast.Node, error) {
	return nil, fmt.Errorf("functions cannot have default values")
}
