package builder

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/helpers/safe"
)

// Represents a folder containing .gold files.
// type Package struct {
// 	Name    string               // Name of the package, ex: `@/foo/bar`
// 	Path    string               // Absolute path of the package in the file system, ex: `/d/project/foo/bar`
// 	Modules *ds.SyncList[*File]  // Modules in the package
// 	Imports *ds.SyncList[string] // Absolute paths of other packages that this package imports
// }

// func NewPackage(name, path string) *Package {
// 	return &Package{
// 		Name:    name,
// 		Path:    path,
// 		Modules: ds.NewSyncList[*File](),
// 		Imports: ds.NewSyncList[string](),
// 	}
// }

// Represents a .gold file.
type File struct {
	Name     string                    // Name of the module, ex: `hello`
	Path     string                    // Absolute path of the module in the file system, ex: `/d/project/foo/bar/hello.gold`
	FileName string                    // Name of the file, ex: `hello.gold`
	Root     safe.Optional[ast.Module] // Root node of the module, type is `ast.Module`
	// Package  *Package        // Package that contains the module
	// Imports []*ModuleImport // Modules that this module imports
}

func NewFile(name, path, fileName string) *File {
	return &File{
		Name:     name,
		Path:     path,
		FileName: fileName,
		Root:     safe.None[ast.Module](),
		// Package:  pkg,
		// Imports: make([]*ModuleImport, 0),
	}
}

// func (m *File) Scope() *env.Scope {
// 	return m.Root.Type().(*types.Module).Scope
// }

// Represents the import from one module to another.
type ModuleImport struct {
	Path  string // Absolute path of the module in the file system, eg: /d/project/foo/bar/hello.gold``
	Alias string // Alias of the module, eg: `hello`
}
