package ir

import (
	"fmt"

	"github.com/renatopp/golden/internal/core"
)

type Ref struct {
	Package    *core.Package
	Module     *core.Module
	Identifier string
}

func R(pkg *core.Package, mod *core.Module, ident string) *Ref {
	return &Ref{
		Package:    pkg,
		Module:     mod,
		Identifier: ident,
	}
}

func (r *Ref) Name() string {
	// TODO: improve naming
	return fmt.Sprintf("%s_%s", r.Module.Name, r.Identifier)
}
