package ir

import "github.com/renatopp/golden/internal/core"

type GirScope struct {
	parent      *GirScope
	depth       int
	nameCounter map[string]int
	values      map[string]core.IrComp
}

func NewGirScope() *GirScope {
	return &GirScope{
		parent:      nil,
		depth:       0,
		nameCounter: map[string]int{},
		values:      map[string]core.IrComp{},
	}
}

func (s *GirScope) New() *GirScope {
	return &GirScope{
		parent:      s,
		depth:       s.depth + 1,
		nameCounter: map[string]int{},
		values:      map[string]core.IrComp{},
	}
}

func (s *GirScope) Count(identifier string) int {
	if _, ok := s.nameCounter[identifier]; ok {
		return s.nameCounter[identifier]
	}
	return 0
}

func (s *GirScope) Incr(identifier string) int {
	if _, ok := s.nameCounter[identifier]; !ok {
		s.nameCounter[identifier] = 0
	}
	s.nameCounter[identifier]++
	return s.nameCounter[identifier]
}

func (s *GirScope) Set(uid string, value core.IrComp) {
	s.values[uid] = value
}

func (s *GirScope) Get(uid string) core.IrComp {
	if v, ok := s.values[uid]; ok {
		return v
	}
	if s.parent != nil {
		return s.parent.Get(uid)
	}
	return nil
}

func (s *GirScope) GetLocal(uid string) core.IrComp {
	if v, ok := s.values[uid]; ok {
		return v
	}
	return nil
}

func (s *GirScope) String() string {
	f := "scope"
	for k, v := range s.values {
		f += "\n\t" + k + " => " + v.Tag()
	}
	return f
}
