package core

import (
	"fmt"

	"github.com/renatopp/golden/internal/helpers/syncds"
)

// Represents a package (folder)
type Package struct {
	Name      string                            // eg: `@/foo/bar`
	Path      string                            // eg: `/d/project/foo/bar`
	Modules   *syncds.SyncList[*Module]         // the modules attached in this package
	DependsOn *syncds.SyncMap[string, *Package] // all packages that modules depends, including implicit ones
	Ir        IrWriter
}

func (p *Package) String() string {
	s := fmt.Sprintf("Package: %s\n", p.Path)
	for _, m := range p.Modules.Values() {
		s += fmt.Sprintf("  Module: %s\n", m.Path)
	}
	return s
}

func NewPackage() *Package {
	return &Package{
		Modules:   syncds.NewSyncList[*Module](),
		DependsOn: syncds.NewSyncMap[string, *Package](),
	}
}
