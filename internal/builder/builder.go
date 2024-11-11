package builder

import (
	"fmt"
	"time"

	"github.com/renatopp/golden/internal/helpers/fs"
	"github.com/renatopp/golden/internal/helpers/logger"
)

type Builder struct {
	startTime time.Time
	pipeline  *Pipeline
}

func NewBuilder() *Builder {
	return &Builder{
		pipeline: NewPipeline(),
	}
}

type BuildOptions struct {
	InputFilePath  string
	OutputFilePath string
	NumWorkers     int
	Debug          bool
}

func (b *Builder) Build(opts BuildOptions) error {
	logger.Debug("[builder] starting building")
	logger.Trace("[builder] - input file path: %s", opts.InputFilePath)
	logger.Trace("[builder] - output file path: %s", opts.OutputFilePath)
	logger.Trace("[builder] - num workers: %d", opts.NumWorkers)
	b.startTime = time.Now()

	logger.Trace("[builder] checking if input file exists")
	if err := fs.CheckFileExists(opts.InputFilePath); err != nil {
		return fmt.Errorf("could not read the input file, reason: \n\n  %w", err)
	}

	logger.Trace("[builder] checking if input file has the correct extension")
	if !fs.IsFileExtension(opts.InputFilePath, ".gold", false) {
		return fmt.Errorf("invalid file extension, expected '.gold', but received '%s'", fs.GetFileExtension(opts.InputFilePath))
	}

	logger.Trace("[builder] checking if input file has the read permission")
	if err := fs.CheckFilePermissions(opts.InputFilePath); err != nil {
		return fmt.Errorf("input file does not have the read permission, reason: \n\n  %w", err)
	}

	path, err := fs.GetAbsolutePath(opts.InputFilePath)
	if err != nil {
		return fmt.Errorf("could not get the absolute path, reason: \n\n  %w", err)
	}

	logger.Trace("[builder] absolute path: %s", path)

	logger.Trace("[builder] checking if input file has valid name")
	if name := fs.ModulePath_To_ModuleName(path); !fs.IsModuleNameValid(name) {
		return fmt.Errorf("invalid module name, expected a valid module name, but received '%s'", name)
	}

	logger.Trace("[builder] starting %d workers", opts.NumWorkers)
	for i := 0; i < opts.NumWorkers; i++ {
		worker := NewBuildWorker(i, opts, b.pipeline)
		go worker.Start()
	}

	logger.Trace("[builder] scheduling discovery of file: %s", path)
	b.pipeline.EntryModulePath = path
	b.pipeline.Discover(path)

	<-b.pipeline.done

	logger.Trace("[builder] done!")
	logger.Debug("[builder] building finished in %s", time.Since(b.startTime))

	return nil
}
