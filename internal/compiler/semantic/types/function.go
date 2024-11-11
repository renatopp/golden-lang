package types

import (
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/core"
)

type Function struct {
	baseType
	Parameters []core.TypeData
	Return     core.TypeData
}

func NewFunction(params []core.TypeData, returnType core.TypeData) *Function {
	return &Function{
		baseType:   newBase(),
		Parameters: params,
		Return:     returnType,
	}
}

func (t *Function) Tag() string {
	return "Fn"
}

func (t *Function) Signature() string {
	params := make([]string, len(t.Parameters))
	for i, p := range t.Parameters {
		params[i] = p.Signature()
	}
	return fmt.Sprintf("Fn(%s) %s", strings.Join(params, ", "), t.Return.Signature())
}

func (t *Function) Accepts(other core.TypeData) bool {
	if t == nil || other == nil {
		return false
	}

	fn, ok := other.(*Function)
	if !ok {
		return false
	}

	if len(t.Parameters) != len(fn.Parameters) {
		return false
	}

	for i, p := range t.Parameters {
		if !p.Accepts(fn.Parameters[i]) {
			return false
		}
	}

	return t.Return.Accepts(fn.Return)
}

func (t *Function) Default() (core.AstData, error) {
	return nil, fmt.Errorf("cannot create default value for function type")
}

func (t *Function) Apply(args []core.TypeData) (core.TypeData, error) {
	if len(t.Parameters) != len(args) {
		return nil, fmt.Errorf("expected %d arguments, got %d", len(t.Parameters), len(args))
	}

	for i, arg := range t.Parameters {
		if !arg.Accepts(args[i]) {
			return nil, fmt.Errorf("expected argument %d to be %s, got %s", i, arg.Signature(), args[i].Signature())
		}
	}

	return t.Return, nil

}
