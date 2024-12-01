package builder

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/env"
	"github.com/renatopp/golden/internal/compiler/types"
	"github.com/renatopp/golden/internal/helpers/safe"
)

// Represents a .gold file.
type File struct {
	Name     string                     // Name of the module, ex: `hello`
	Path     string                     // Absolute path of the module in the file system, ex: `/d/project/foo/bar/hello.gold`
	FileName string                     // Name of the file, ex: `hello.gold`
	Root     safe.Optional[*ast.Module] // Root node of the module, type is `ast.Module`
	// Imports []*ModuleImport // Modules that this module imports
}

func NewFile(name, path, fileName string) *File {
	return &File{
		Name:     name,
		Path:     path,
		FileName: fileName,
		Root:     safe.None[*ast.Module](),
		// Imports: make([]*ModuleImport, 0),
	}
}

func (m *File) Scope() *env.Scope {
	return m.Root.Unwrap().GetType().Unwrap().(*types.Module).Scope
}

// Represents the import from one module to another.
type ModuleImport struct {
	Path  string // Absolute path of the module in the file system, eg: /d/project/foo/bar/hello.gold``
	Alias string // Alias of the module, eg: `hello`
}
