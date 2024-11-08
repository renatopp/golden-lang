package internal

import (
	"os"
	"sync"
	"sync/atomic"

	"github.com/renatopp/golden/internal/fs"
	"github.com/renatopp/golden/internal/logger"
)

type BuildPipeline struct {
	Mtx                 sync.Mutex
	Modules             *SyncMap[string, *Module]
	PendingPrepareCount atomic.Int64
	ToDiscover          chan string
	ToPrepare           chan string
	Done                chan any
}

func NewBuildPipeline() *BuildPipeline {
	return &BuildPipeline{
		Mtx:                 sync.Mutex{},
		Modules:             NewSyncMap[string, *Module](),
		PendingPrepareCount: atomic.Int64{},
		ToDiscover:          make(chan string),
		ToPrepare:           make(chan string),
		Done:                make(chan any),
	}
}

type BuildWorker struct {
	id       int
	pipeline *BuildPipeline
}

func NewBuildWorker(id int, pipeline *BuildPipeline) *BuildWorker {
	return &BuildWorker{
		id:       id,
		pipeline: pipeline,
	}
}

func (w *BuildWorker) Start() {
	for {
		select {
		case path := <-w.pipeline.ToDiscover:
			w.discover(path)

		case path := <-w.pipeline.ToPrepare:
			w.prepare(path)
		}

		if w.pipeline.PendingPrepareCount.Load() > 0 {
			continue
		}

		w.analyze()
	}
}

// From a given module path, discover all the modules in the same package
func (w *BuildWorker) discover(modulePath string) {
	logger.Trace("[worker:discover] discovering package of: %s", modulePath)
	files := fs.DiscoverModules(modulePath)
	for _, file := range files {
		w.pipeline.Mtx.Lock()
		if w.pipeline.Modules.Has(file) {
			w.pipeline.Mtx.Unlock()
			continue
		}

		w.pipeline.PendingPrepareCount.Add(1)
		w.pipeline.Modules.Set(file, nil)
		w.pipeline.Mtx.Unlock()
		w.pipeline.ToPrepare <- file
	}
}

// For a given module path, prepare (lex, parse, pre-analyze) the file for the analysis
func (w *BuildWorker) prepare(modulePath string) {
	logger.Trace("[worker:prepare] preparing file: %s", modulePath)

	bytes, err := os.ReadFile(modulePath)
	if err != nil {
		panic(err)
	}

	logger.Trace("[worker:prepare] lexing: %s", modulePath)
	tokens, err := Lex(bytes)
	if err != nil {
		panic(err)
	}

	for _, token := range tokens {
		loc := token.Loc
		loc.Filename = modulePath
		token.Loc = loc
	}

	logger.Trace("[worker:prepare] parsing: %s", modulePath)
	module, err := Parse(tokens)
	if err != nil {
		panic(err)
	}

	module.Path = modulePath
	module.Name = fs.ModulePath_To_ModuleName(modulePath)
	module.FileName = fs.ModulePath_To_ModuleFileName(modulePath)
	module.PackageName = fs.ModulePath_To_PackageName(modulePath)
	module.PackagePath = fs.ModulePath_To_PackagePath(modulePath)

	w.pipeline.Modules.Set(modulePath, module)
	w.pipeline.PendingPrepareCount.Add(-1)
}

// Analyze all the modules in the pipeline
func (w *BuildWorker) analyze() {
	logger.Trace("[worker:analyze] analyzing modules")
	w.pipeline.Mtx.Lock()
	defer w.pipeline.Mtx.Unlock()

	w.pipeline.Done <- nil
}
