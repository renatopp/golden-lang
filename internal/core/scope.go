package core

import "strings"

// Represents a name bind, for both types and values
//
// If the expression kind is "type", then the node is the type declaration
// while the type is the type itself. If the expression kind is "value", then
// the node is the value declaration while the type is none.
type Binding struct {
	Kind ExpressionKind // type, value
	// Capability string // immutable, mutable
	Node *AstNode
	Type TypeData
}

func BindValue(node *AstNode) *Binding {
	return &Binding{
		Kind: ValueExpression,
		Node: node,
		Type: nil,
	}
}

func BindType(node *AstNode, t TypeData) *Binding {
	return &Binding{
		Kind: TypeExpression,
		Node: node,
		Type: t,
	}
}

// ----------------------------------------------------------------------------

type ScopeMap struct {
	parent   *ScopeMap
	bindings map[string]*Binding
}

func NewScopeMap(parent *ScopeMap) *ScopeMap {
	return &ScopeMap{
		parent:   parent,
		bindings: map[string]*Binding{},
	}
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

// ----------------------------------------------------------------------------

type Scope struct {
	Depth  int
	Parent *Scope
	Values *ScopeMap
	Types  *ScopeMap
}

func NewScope() *Scope {
	return &Scope{
		Depth:  0,
		Parent: nil,
		Values: NewScopeMap(nil),
		Types:  NewScopeMap(nil),
	}
}

func (s *Scope) New() *Scope {
	return &Scope{
		Depth:  s.Depth + 1,
		Parent: s,
		Values: NewScopeMap(s.Values),
		Types:  NewScopeMap(s.Types),
	}
}

func (s *Scope) String() string {
	r := ""
	if s.Parent != nil {
		r += s.Parent.String()
	}

	ident := strings.Repeat("| ", s.Depth+1)
	r += strings.Repeat("\n| ", s.Depth) + "[scope]\n"
	for k, v := range s.Types.bindings {
		r += ident + "T: " + k + " = " + v.Type.Signature() + "\n"
	}
	for k, v := range s.Values.bindings {
		r += ident + "V: " + k + " = " + v.Node.Signature() + "\n"
	}
	return r
}
