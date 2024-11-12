package builder

import (
	"sync"
	"sync/atomic"

	"github.com/renatopp/golden/internal/compiler/semantic"
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/syncds"
)

// Pipeline is the structure that holds the global state of the build pipeline
type Pipeline struct {
	GlobalScope        *core.Scope
	EntryModulePath    string
	Packages           *syncds.SyncMap[string, *core.Package]
	Modules            *syncds.SyncMap[string, *core.Module]
	PendingModuleCount atomic.Int64

	mtx        sync.Mutex
	toDiscover chan string
	toPrepare  chan string
	done       chan any
}

func NewPipeline() *Pipeline {
	scope := core.NewScope()
	scope.SetValue("$_scope", core.NewEmptyNode().WithType(semantic.String).WithData(&ast.String{Value: "global"}))
	scope.SetType("Int", semantic.Int)
	scope.SetType("Float", semantic.Float)
	scope.SetType("String", semantic.String)
	scope.SetType("Bool", semantic.Bool)
	println("NewPipeline", semantic.Int)

	return &Pipeline{
		GlobalScope:        scope,
		Packages:           syncds.NewSyncMap[string, *core.Package](),
		Modules:            syncds.NewSyncMap[string, *core.Module](),
		PendingModuleCount: atomic.Int64{},

		mtx:        sync.Mutex{},
		toDiscover: make(chan string),
		toPrepare:  make(chan string),
		done:       make(chan any),
	}
}

func (p *Pipeline) Discover(modulePath string) {
	p.PendingModuleCount.Add(1)
	p.toDiscover <- modulePath
}

func (p *Pipeline) Prepare(modulePath string) {
	p.PendingModuleCount.Add(1)
	p.toPrepare <- modulePath
}

func (p *Pipeline) AckModule() {
	p.PendingModuleCount.Add(-1)
}

func (p *Pipeline) PreRegisterModule(modulePath string) bool {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if p.Modules.Has(modulePath) {
		return false
	}
	p.Modules.Set(modulePath, nil)
	return true
}

func (p *Pipeline) RegisterModule(module *core.Module) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.Modules.Set(module.Path, module)
}

func (p *Pipeline) CreateOrGetPackage(packageName, packagePath string) *core.Package {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if pkg, ok := p.Packages.Get(packagePath); ok {
		return pkg
	}

	pkg := core.NewPackage()
	pkg.Name = packageName
	pkg.Path = packagePath
	p.Packages.Set(packagePath, pkg)
	return pkg
}
