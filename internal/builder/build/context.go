package build

import (
	"sync"

	"github.com/renatopp/golden/internal/compiler/semantic"
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/syncds"
)

type Context struct {
	Options           Options
	GlobalScope       *core.Scope
	EntryModulePath   string
	Packages          *syncds.SyncMap[string, *core.Package]
	Modules           *syncds.SyncMap[string, *core.Module]
	ToDiscoverPackage chan string
	ToPrepareAST      chan string
	ToDependencyGraph chan string
	ToResolveBindings chan string
	ToFinish          chan string
	Done              chan any

	mtx sync.Mutex
}

func NewContext() *Context {
	scope := core.NewScope()
	scope.SetValue("$_scope", core.NewEmptyNode().WithType(semantic.String).WithData(&ast.String{Value: "global"}))
	scope.SetType("Int", semantic.Int)
	scope.SetType("Float", semantic.Float)
	scope.SetType("String", semantic.String)
	scope.SetType("Bool", semantic.Bool)

	return &Context{
		GlobalScope:       scope,
		Packages:          syncds.NewSyncMap[string, *core.Package](),
		Modules:           syncds.NewSyncMap[string, *core.Module](),
		ToDiscoverPackage: make(chan string, 100),
		ToPrepareAST:      make(chan string, 100),
		ToDependencyGraph: make(chan string, 100),
		ToResolveBindings: make(chan string, 100),
		ToFinish:          make(chan string, 100),
		Done:              make(chan any),

		mtx: sync.Mutex{},
	}
}

func (this *Context) PreRegisterModule(modulePath string) bool {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	if this.Modules.Has(modulePath) {
		return false
	}
	this.Modules.Set(modulePath, nil)
	return true
}

func (this *Context) RegisterModule(module *core.Module) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.Modules.Set(module.Path, module)
}

func (this *Context) CreateOrGetPackage(packageName, packagePath string) *core.Package {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	if pkg, ok := this.Packages.Get(packageName); ok {
		return pkg
	}
	pkg := core.NewPackage()
	pkg.Name = packageName
	pkg.Path = packagePath
	this.Packages.Set(packagePath, pkg)
	return pkg
}
