package internal

import (
	"sync"
	"sync/atomic"
)

// BuildPipeline is the structure that holds the global state of the build pipeline
type BuildPipeline struct {
	GlobalScope        *Scope
	EntryModulePath    string
	Packages           *SyncMap[string, *Package]
	Modules            *SyncMap[string, *Module]
	PendingModuleCount atomic.Int64

	mtx        sync.Mutex
	toDiscover chan string
	toPrepare  chan string
	done       chan any
}

func NewBuildPipeline() *BuildPipeline {
	scope := NewScope()
	scope.SetValue("$_scope", NewEmptyNode().WithType(String).WithData(&AstString{Value: "global"}))

	return &BuildPipeline{
		GlobalScope:        scope,
		Packages:           NewSyncMap[string, *Package](),
		Modules:            NewSyncMap[string, *Module](),
		PendingModuleCount: atomic.Int64{},

		mtx:        sync.Mutex{},
		toDiscover: make(chan string),
		toPrepare:  make(chan string),
		done:       make(chan any),
	}
}

func (p *BuildPipeline) Discover(modulePath string) {
	p.PendingModuleCount.Add(1)
	p.toDiscover <- modulePath
}

func (p *BuildPipeline) Prepare(modulePath string) {
	p.PendingModuleCount.Add(1)
	p.toPrepare <- modulePath
}

func (p *BuildPipeline) AckModule() {
	p.PendingModuleCount.Add(-1)
}

func (p *BuildPipeline) PreRegisterModule(modulePath string) bool {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if p.Modules.Has(modulePath) {
		return false
	}
	p.Modules.Set(modulePath, nil)
	return true
}

func (p *BuildPipeline) RegisterModule(module *Module) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.Modules.Set(module.Path, module)
}

func (p *BuildPipeline) CreateOrGetPackage(packageName, packagePath string) *Package {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	if pkg, ok := p.Packages.Get(packagePath); ok {
		return pkg
	}

	pkg := NewPackage()
	pkg.Name = packageName
	pkg.Path = packagePath
	p.Packages.Set(packagePath, pkg)
	return pkg
}
