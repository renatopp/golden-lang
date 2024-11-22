package types

import "github.com/renatopp/golden/internal/compiler/ast"

var _baseTypeId uint64

type BaseType struct {
	id         uint64
	definition ast.Node
}

func NewBaseType(def ast.Node) *BaseType {
	_baseTypeId++
	return &BaseType{
		id:         _baseTypeId,
		definition: def,
	}
}

func (t *BaseType) Id() uint64           { return t.id }
func (t *BaseType) Definition() ast.Node { return t.definition }
