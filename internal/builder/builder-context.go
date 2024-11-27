package builder

import (
	"github.com/renatopp/golden/internal/compiler/env"
	"github.com/renatopp/golden/internal/helpers/ds"
)

type BuildContext struct {
	Options         *BuildOptions
	PackageRegistry *ds.SyncMap[string, *Package]
	ModuleRegistry  *ds.SyncMap[string, *Module]
	EntryPackage    *Package
	EntryModule     *Module
	DependencyOrder []*Package
	GlobalScope     *env.Scope
}
