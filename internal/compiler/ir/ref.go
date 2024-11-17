package ir

import (
	"github.com/renatopp/golden/internal/core"
)

type Ref struct {
	Package    *core.Package
	Module     *core.Module
	Identifier string
	SsaCount   int
}

func R(pkg *core.Package, mod *core.Module, ident string, ssaCount int) *Ref {
	return &Ref{
		Package:    pkg,
		Module:     mod,
		Identifier: ident,
		SsaCount:   ssaCount,
	}
}

// func (r *Ref) PreName() string {
// 	// TODO: improve naming
// 	// return fmt.Sprintf("%s_%s", r.Module.Name, r.Identifier)
// 	return r.Identifier
// }

// func (r *Ref) Name() string {
// 	// TODO: improve naming
// 	// return fmt.Sprintf("%s_%s", r.Module.Name, r.Identifier)
// 	return fmt.Sprintf("%s%d", r.Identifier, r.SsaCount)
// }
