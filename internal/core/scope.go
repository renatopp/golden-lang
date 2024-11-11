package core

import "strings"

type Scope struct {
	Depth  int
	Parent *Scope
	Values map[string]*AstNode
	Types  map[string]TypeData
}

func NewScope() *Scope {
	return &Scope{
		Depth:  0,
		Parent: nil,
		Values: map[string]*AstNode{},
		Types:  map[string]TypeData{},
	}
}

func (s *Scope) New() *Scope {
	return &Scope{
		Depth:  s.Depth + 1,
		Parent: s,
		Values: map[string]*AstNode{},
		Types:  map[string]TypeData{},
	}
}

func (s *Scope) SetValue(key string, value *AstNode) {
	s.Values[key] = value
}

func (s *Scope) SetType(key string, value TypeData) {
	s.Types[key] = value
}

func (s *Scope) GetValueLocal(key string) *AstNode {
	return s.Values[key]
}

func (s *Scope) GetTypeLocal(key string) TypeData {
	return s.Types[key]
}

func (s *Scope) GetValue(key string) *AstNode {
	if value, ok := s.Values[key]; ok {
		return value
	}
	if s.Parent != nil {
		return s.Parent.GetValue(key)
	}
	return nil
}

func (s *Scope) GetType(key string) TypeData {
	if tp, ok := s.Types[key]; ok {
		return tp
	}
	if s.Parent != nil {
		return s.Parent.GetType(key)
	}
	return nil
}

func (s *Scope) String() string {
	r := ""
	if s.Parent != nil {
		r += s.Parent.String()
	}

	ident := strings.Repeat("| ", s.Depth+1)
	r += strings.Repeat("\n| ", s.Depth) + "[scope]\n"
	for k, v := range s.Types {
		r += ident + "T: " + k + " = " + v.Signature() + "\n"
	}
	for k, v := range s.Values {
		r += ident + "V: " + k + " = " + v.Signature() + "\n"
	}
	return r
}
