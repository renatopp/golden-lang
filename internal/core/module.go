package core

import "github.com/renatopp/golden/internal/helpers/syncds"

type Resolver interface {
	PreResolve(*AstNode) error
	Resolve(*AstNode) error
}

// Represents a module (file)
type Module struct {
	Name      string                           // eg: `hello`
	Path      string                           // eg: `/d/project/foo/bar/hello.gold`
	FileName  string                           // eg: `hello.gold`
	Node      *AstNode                         // the ast node attached to this module
	Package   *Package                         // the package this module is attached
	Imports   []*ModuleImport                  // all import statements in this module
	DependsOn *syncds.SyncMap[string, *Module] // all modules imported in this one, including implicit modules

	Scope    *Scope
	Resolver Resolver
}

func NewModule() *Module {
	return &Module{
		Imports:   []*ModuleImport{},
		DependsOn: syncds.NewSyncMap[string, *Module](),
	}
}

// Represents the import of a module inside another module
type ModuleImport struct {
	Name    string // eg: `@/foo/bar/hello`
	Alias   string // eg: `x` in `import '...' as x`
	Path    string // eg: `/d/project/foo/bar/hello.gold`
	Module  *Module
	Package *Package
}

// Represents a package (folder)
type Package struct {
	Name      string                            // eg: `@/foo/bar`
	Path      string                            // eg: `/d/project/foo/bar`
	Modules   *syncds.SyncList[*Module]         // the modules attached in this package
	DependsOn *syncds.SyncMap[string, *Package] // all packages that modules depends, including implicit ones
}

func NewPackage() *Package {
	return &Package{
		Modules:   syncds.NewSyncList[*Module](),
		DependsOn: syncds.NewSyncMap[string, *Package](),
	}
}
