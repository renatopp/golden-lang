package types

import (
	"fmt"

	"github.com/renatopp/golden/internal/core"
)

type Void struct{ baseType }

func NewVoid() *Void {
	return &Void{newBase()}
}

func (t *Void) Tag() string {
	return "Void"
}

func (t *Void) Signature() string {
	return "Void"
}

func (t *Void) Accepts(other core.TypeData) bool {
	if t == nil || other == nil {
		return false
	}
	return true
}

func (t *Void) Default() (core.AstData, error) {
	return nil, fmt.Errorf("Void does not have a default value")
}
