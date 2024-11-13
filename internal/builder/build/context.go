package build

import (
	"sync"
	"sync/atomic"

	"github.com/renatopp/golden/internal/compiler/semantic"
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/syncds"
)

type Context struct {
	Options         Options
	GlobalScope     *core.Scope
	EntryModulePath string
	Packages        *syncds.SyncMap[string, *core.Package]
	Modules         *syncds.SyncMap[string, *core.Module]

	toDiscoverPackage  chan string
	toPrepareAST       chan string
	pendingModuleCount atomic.Int64

	Done chan error

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
		toDiscoverPackage: make(chan string, 100),
		toPrepareAST:      make(chan string, 100),
		Done:              make(chan error),

		mtx: sync.Mutex{},
	}
}

func (c *Context) ScheduleDiscoverPackage(modulePath string) {
	c.pendingModuleCount.Add(1)
	c.toDiscoverPackage <- modulePath
}

func (c *Context) SchedulePrepareAST(modulePath string) {
	c.pendingModuleCount.Add(1)
	c.toPrepareAST <- modulePath
}

func (c *Context) AckModule() {
	c.pendingModuleCount.Add(-1)
}

func (c *Context) CanProceedToDependencyGraph() bool {
	return c.pendingModuleCount.Load() == 0
}

func (c *Context) PreRegisterModule(modulePath string) bool {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if c.Modules.Has(modulePath) {
		return false
	}
	c.Modules.Set(modulePath, nil)
	return true
}

func (c *Context) RegisterModule(module *core.Module) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.Modules.Set(module.Path, module)
}

func (c *Context) CreateOrGetPackage(packageName, packagePath string) *core.Package {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if pkg, ok := c.Packages.Get(packagePath); ok {
		return pkg
	}
	pkg := core.NewPackage()
	pkg.Name = packageName
	pkg.Path = packagePath
	c.Packages.Set(packagePath, pkg)
	return pkg
}
