package env

import "github.com/renatopp/golden/internal/compiler/ast"

var _bindingId uint64 = 0

type ValueBinding struct {
	Id             uint64
	DefinitionNode ast.Node
	AssignedNodes  []ast.Node
	LastNode       ast.Node
	Type           ast.Type
}

func NewValueBinding(t ast.Type, n ast.Node) *ValueBinding {
	_bindingId++
	return &ValueBinding{
		Id:             _bindingId,
		DefinitionNode: n,
		Type:           t,
	}
}

func (b *ValueBinding) Assign(n ast.Node) {
	b.AssignedNodes = append(b.AssignedNodes, n)
	b.LastNode = n
}

var VB = NewValueBinding

type TypeBinding struct {
	Id             uint64
	DefinitionNode ast.Node
	Type           ast.Type
}

func NewTypeBinding(t ast.Type, n ast.Node) *TypeBinding {
	_bindingId++
	return &TypeBinding{
		Id:             _bindingId,
		DefinitionNode: n,
		Type:           t,
	}
}

var TB = NewTypeBinding
