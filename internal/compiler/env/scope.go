package env

import "github.com/renatopp/golden/internal/compiler/ast"

type Binding struct {
	Type ast.Type
}

func NewBinding(t ast.Type) *Binding {
	return &Binding{Type: t}
}

var B = NewBinding

//
//
//

type ScopeMap struct {
	parent   *ScopeMap
	bindings map[string]*Binding
}

func (s *ScopeMap) Get(key string) *Binding {
	if binding, ok := s.bindings[key]; ok {
		return binding
	}
	if s.parent != nil {
		return s.parent.Get(key)
	}
	return nil
}

func (s *ScopeMap) GetLocal(key string) *Binding {
	return s.bindings[key]
}

func (s *ScopeMap) Set(key string, binding *Binding) {
	s.bindings[key] = binding
}

func (s *ScopeMap) Clear() {
	s.bindings = map[string]*Binding{}
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
		Types:  &ScopeMap{bindings: map[string]*Binding{}},
		Values: &ScopeMap{bindings: map[string]*Binding{}},
	}
}

func (s *Scope) New() *Scope {
	return &Scope{
		Depth:  s.Depth + 1,
		Parent: s,
		Types:  &ScopeMap{bindings: map[string]*Binding{}},
		Values: &ScopeMap{bindings: map[string]*Binding{}},
	}
}
