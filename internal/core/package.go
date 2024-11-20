package core

import "github.com/renatopp/golden/internal/helpers/ds"

// Represents a folder containing .gold files.
type Package struct {
	Name    string                // Name of the package, ex: `@/foo/bar`
	Path    string                // Absolute path of the package in the file system, ex: `/d/project/foo/bar`
	Modules *ds.SyncList[*Module] // Modules in the package
}
