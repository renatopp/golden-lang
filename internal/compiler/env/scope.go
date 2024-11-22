package env

import "github.com/renatopp/golden/internal/compiler/ast"

type Binding struct {
	Node ast.Node // The node that the binding refers to
	Type ast.Type // The type of the binding
}

func NewBinding(t ast.Type, n ast.Node) *Binding {
	return &Binding{Node: n, Type: t}
}

func NewSimpleBinding(t ast.Type) *Binding {
	return &Binding{Type: t}
}

var B = NewSimpleBinding
var BN = NewBinding

//
//
//

type ScopeMap struct {
	Parent   *ScopeMap
	Bindings map[string]*Binding
}

func (s *ScopeMap) Get(key string) *Binding {
	if binding, ok := s.Bindings[key]; ok {
		return binding
	}
	if s.Parent != nil {
		return s.Parent.Get(key)
	}
	return nil
}

func (s *ScopeMap) GetLocal(key string) *Binding {
	return s.Bindings[key]
}

func (s *ScopeMap) Set(key string, binding *Binding) {
	s.Bindings[key] = binding
}

func (s *ScopeMap) Clear() {
	s.Bindings = map[string]*Binding{}
}

//
//
//

type Scope struct {
	Depth  int
	Parent *Scope
	Types  *ScopeMap
	Values *ScopeMap
}

func NewScope() *Scope {
	return &Scope{
		Depth:  0,
		Parent: nil,
		Types:  &ScopeMap{Bindings: map[string]*Binding{}},
		Values: &ScopeMap{Bindings: map[string]*Binding{}},
	}
}

func (s *Scope) New() *Scope {
	return &Scope{
		Depth:  s.Depth + 1,
		Parent: s,
		Types:  &ScopeMap{Parent: s.Types, Bindings: map[string]*Binding{}},
		Values: &ScopeMap{Parent: s.Values, Bindings: map[string]*Binding{}},
	}
}
