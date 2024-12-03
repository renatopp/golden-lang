package env

import "github.com/renatopp/golden/internal/compiler/ast"

var _bindingId uint64 = 0

type ValueBinding struct {
	Id             uint64
	DefinitionNode ast.Node
	Assignments    []ast.Node
	References     []ast.Node
	LastNode       ast.Node
	Type           ast.Type
}

func NewValueBinding(n ast.Node, t ast.Type) *ValueBinding {
	_bindingId++
	return &ValueBinding{
		Id:             _bindingId,
		DefinitionNode: n,
		Assignments:    []ast.Node{n},
		References:     []ast.Node{},
		LastNode:       n,
		Type:           t,
	}
}

func (b *ValueBinding) Assign(n ast.Node) {
	b.Assignments = append(b.Assignments, n)
	b.LastNode = n
}

func (b *ValueBinding) Reference(n ast.Node) {
	b.References = append(b.References, n)
}

func (b *ValueBinding) IsSolved() bool {
	return b.Type != nil
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

func (b *TypeBinding) IsSolved() bool {
	return b.Type != nil
}

var TB = NewTypeBinding
