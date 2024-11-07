package internal

import "strings"

type Scope struct {
	Depth  int
	Parent *Scope
	Map    map[string]*Node
}

func NewScope() *Scope {
	return &Scope{
		Depth:  0,
		Parent: nil,
		Map:    map[string]*Node{},
	}
}

func (s *Scope) New() *Scope {
	return &Scope{
		Depth:  s.Depth + 1,
		Parent: s,
		Map:    map[string]*Node{},
	}
}

func (s *Scope) Set(key string, value *Node) {
	s.Map[key] = value
}

func (s *Scope) GetLocal(key string) *Node {
	return s.Map[key]
}

func (s *Scope) Get(key string) *Node {
	if value, ok := s.Map[key]; ok {
		return value
	}
	if s.Parent != nil {
		return s.Parent.Get(key)
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
	for k, v := range s.Map {
		r += ident + k + " \u2192 " + v.String() + "\n"
	}
	return r
}
