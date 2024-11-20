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
}

type Module struct {
	Name     string   // Name of the module, ex: `hello`
	Path     string   // Absolute path of the module in the file system, ex: `/d/project/foo/bar/hello.gold`
	FileName string   // Name of the file, ex: `hello.gold`
	Package  *Package // Package that contains the module
	Root     ast.Node // Root node of the module, type is `*ast.Module`
}
