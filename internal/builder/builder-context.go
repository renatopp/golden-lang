package builder

import (
	// "github.com/renatopp/golden/internal/compiler/env"
	"github.com/renatopp/golden/internal/helpers/ds"
)

type BuildContext struct {
	Options        *BuildOptions
	ModuleRegistry *ds.SyncMap[string, *File]
	EntryModule    *File
	// PackageRegistry *ds.SyncMap[string, *Package]
	// EntryPackage    *Package
	// DependencyOrder []*Package
	// GlobalScope     *env.Scope
}
