package env

type scopeMap[T any] struct {
	Parent   *scopeMap[T]
	Bindings map[string]T
}

func (s *scopeMap[T]) Get(key string, or T) T {
	if binding, ok := s.Bindings[key]; ok {
		return binding
	}
	if s.Parent != nil {
		return s.Parent.Get(key, or)
	}
	return or
}

func (s *scopeMap[T]) GetLocal(key string, or T) T {
	if binding, ok := s.Bindings[key]; ok {
		return binding
	}
	return or
}

func (s *scopeMap[T]) Set(key string, binding T) {
	s.Bindings[key] = binding
}

func (s *scopeMap[T]) Clear() {
	s.Bindings = map[string]T{}
}

//
//
//

type Scope struct {
	Depth    int
	IsModule bool
	Parent   *Scope
	Types    *scopeMap[*TypeBinding]
	Values   *scopeMap[*ValueBinding]
}

func NewScope() *Scope {
	return &Scope{
		Depth:  0,
		Parent: nil,
		Types:  &scopeMap[*TypeBinding]{Bindings: map[string]*TypeBinding{}},
		Values: &scopeMap[*ValueBinding]{Bindings: map[string]*ValueBinding{}},
	}
}

func (s *Scope) New() *Scope {
	return &Scope{
		Depth:  s.Depth + 1,
		Parent: s,
		Types:  &scopeMap[*TypeBinding]{Parent: s.Types, Bindings: map[string]*TypeBinding{}},
		Values: &scopeMap[*ValueBinding]{Parent: s.Values, Bindings: map[string]*ValueBinding{}},
	}
}
