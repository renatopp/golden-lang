package internal

import (
	"github.com/renatopp/golden/internal/fs"
	"github.com/renatopp/golden/internal/logger"
)

type BuildPipeline struct {
	ToDiscover chan string
	ToPrepare  chan string
	Done       chan any
}

func NewBuildPipeline() *BuildPipeline {
	return &BuildPipeline{
		ToDiscover: make(chan string),
		ToPrepare:  make(chan string),
		Done:       make(chan any),
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
			w.Prepare(path)
		}
	}
}

// From a given module path, discover all the modules in the same package
func (w *BuildWorker) discover(modulePath string) {
	logger.Trace("[worker:discover] discovering package of: %s", modulePath)
	files := fs.DiscoverModules(modulePath)
	for _, file := range files {
		w.pipeline.ToPrepare <- file
	}
}

// For a given module path, prepare (lex, parse, pre-analyze) the file for the analysis
func (w *BuildWorker) Prepare(modulePath string) {
	logger.Trace("[worker:prepare] preparing file: %s", modulePath)

}
