package girl

import "github.com/renatopp/golden/internal/helpers/ds"

// Reflects the golden scope to keep track of which golden names are mapped to
// which girl names.
type Scope struct {
	Parent *Scope
	// Names  map[string]*Ref // map golden names to girl names
}

type Value struct{}

type GirlWriter struct {
	scopeStack  *ds.Stack[Scope]
	ssaCounters map[string]int64
}

func NewGirlWriter() *GirlWriter {
	return &GirlWriter{
		scopeStack:  ds.NewStack[Scope](),
		ssaCounters: make(map[string]int64),
	}
}

// func DeclareVariable(package, module, name string, expr)

// func OpenScope --> push new scope to bind names
// func CloseScope  -->
