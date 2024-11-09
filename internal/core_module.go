package internal

type Import struct {
	Name    string // eg: `@/foo/bar/hello``
	Alias   string
	Path    string // eg: `/d/project/foo/bar/hello.gold`
	Module  *Module
	Package *Package
}

type Module struct {
	Name      string
	Path      string
	FileName  string
	Package   *Package
	Scope     *Scope
	Imports   []*Import
	Ast       *AstModule
	Node      *Node
	DependsOn *SyncMap[string, *Module]
	Analyzer  *Analyzer
}

func NewModule() *Module {
	return &Module{
		Imports:   []*Import{},
		DependsOn: NewSyncMap[string, *Module](),
	}
}

type Package struct {
	Name      string
	Path      string
	Modules   *SyncList[*Module]
	DependsOn *SyncMap[string, *Package]
}

func NewPackage() *Package {
	return &Package{
		Modules:   NewSyncList[*Module](),
		DependsOn: NewSyncMap[string, *Package](),
	}
}
