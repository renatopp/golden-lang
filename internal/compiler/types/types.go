package types

import "github.com/renatopp/golden/internal/compiler/ast"

var _typeId uint64

type BaseType struct {
	Id         uint64
	Definition ast.Node
}

func NewBaseType(def ast.Node) *BaseType {
	_typeId++
	return &BaseType{
		Id:         _typeId,
		Definition: def,
	}
}

func (t *BaseType) GetId() uint64           { return t.Id }
func (t *BaseType) GetDefinition() ast.Node { return t.Definition }
