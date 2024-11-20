package builder

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/helpers/ds"
)

// Represents a folder containing .gold files.
type Package struct {
	Name    string                // Name of the package, ex: `@/foo/bar`
	Path    string                // Absolute path of the package in the file system, ex: `/d/project/foo/bar`
	Modules *ds.SyncList[*Module] // Modules in the package
	Imports *ds.SyncList[string]  // Absolute paths of other packages that this package imports
}

type Module struct {
	Name     string          // Name of the module, ex: `hello`
	Path     string          // Absolute path of the module in the file system, ex: `/d/project/foo/bar/hello.gold`
	FileName string          // Name of the file, ex: `hello.gold`
	Package  *Package        // Package that contains the module
	Root     *ast.Module     // Root node of the module, type is `*ast.Module`
	Imports  []*ModuleImport // Modules that this module imports
}

type ModuleImport struct {
	Path  string // Absolute path of the module in the file system, eg: /d/project/foo/bar/hello.gold``
	Alias string // Alias of the module, eg: `hello`
}
