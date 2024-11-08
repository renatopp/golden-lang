package internal

import (
	"fmt"
	"time"

	"github.com/renatopp/golden/internal/fs"
	"github.com/renatopp/golden/internal/logger"
)

type Builder struct {
	startTime time.Time
	pipeline  *BuildPipeline
}

func NewBuilder() *Builder {
	return &Builder{
		pipeline: NewBuildPipeline(),
	}
}

type BuildOptions struct {
	InputFilePath  string
	OutputFilePath string
	NumWorkers     int
}

func (b *Builder) Build(opts BuildOptions) error {
	logger.Debug("[builder] starting building")
	logger.Trace("[builder] - input file path: %s", opts.InputFilePath)
	logger.Trace("[builder] - output file path: %s", opts.OutputFilePath)
	logger.Trace("[builder] - num workers: %d", opts.NumWorkers)
	b.startTime = time.Now()

	logger.Trace("[builder] checking if input file exists")
	err := fs.CheckFileExists(opts.InputFilePath)
	if err != nil {
		return fmt.Errorf("could not read the input file, reason: \n\n  %w", err)
	}

	logger.Trace("[builder] checking if input file has the correct extension")
	if !fs.IsFileExtension(opts.InputFilePath, ".gold", false) {
		return fmt.Errorf("invalid file extension, expected '.gold', but received '%s'", fs.GetFileExtension(opts.InputFilePath))
	}

	logger.Trace("[builder] checking if input file has the read permission")
	err = fs.CheckFilePermissions(opts.InputFilePath)
	if err != nil {
		return fmt.Errorf("input file does not have the read permission, reason: \n\n  %w", err)
	}

	path, err := fs.GetAbsolutePath(opts.InputFilePath)
	if err != nil {
		return fmt.Errorf("could not get the absolute path, reason: \n\n  %w", err)
	}

	logger.Trace("[builder] absolute path: %s", path)

	logger.Trace("[builder] starting %d workers", opts.NumWorkers)
	for i := 0; i < opts.NumWorkers; i++ {
		worker := NewBuildWorker(i, b.pipeline)
		go worker.Start()
	}

	logger.Trace("[builder] scheduling discovery of file: %s", path)
	b.pipeline.ToDiscover <- path

	// <-b.pipeline.Done

	time.Sleep(1 * time.Second)
	logger.Trace("[builder] done!")
	logger.Debug("[builder] building finished in %s", time.Since(b.startTime))

	// Package discovery:
	// - Discover package
	// - Load all modules
	// - Schedule each module

	// Module processing:
	// - Lex
	// - Parse
	// - Pre Analyze
	// - Put in scope

	// Analysis:
	// - Analyze each module
	// - Start codegen

	// Codegen:
	// - Generate code

	return nil
}
