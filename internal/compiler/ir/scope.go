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
