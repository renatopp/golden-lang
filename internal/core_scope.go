package internal

import "strings"

type Scope struct {
	Depth  int
	Parent *Scope
	Values map[string]*Node
	Types  map[string]RtType
}

func NewScope() *Scope {
	return &Scope{
		Depth:  0,
		Parent: nil,
		Values: map[string]*Node{},
		Types:  map[string]RtType{},
	}
}

func (s *Scope) New() *Scope {
	return &Scope{
		Depth:  s.Depth + 1,
		Parent: s,
		Values: map[string]*Node{},
		Types:  map[string]RtType{},
	}
}

func (s *Scope) SetValue(key string, value *Node) {
	s.Values[key] = value
}

func (s *Scope) SetType(key string, value RtType) {
	s.Types[key] = value
}

func (s *Scope) GetValueLocal(key string) *Node {
	return s.Values[key]
}

func (s *Scope) GetTypeLocal(key string) RtType {
	return s.Types[key]
}

func (s *Scope) GetValue(key string) *Node {
	if value, ok := s.Values[key]; ok {
		return value
	}
	if s.Parent != nil {
		return s.Parent.GetValue(key)
	}
	return nil
}

func (s *Scope) GetType(key string) RtType {
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
		r += ident + "T: " + k + " \u2192 " + v.Name() + "\n"
	}
	for k, v := range s.Values {
		r += ident + "V: " + k + " \u2192 " + oneline(v.String()) + "\n"
	}
	return r
}
