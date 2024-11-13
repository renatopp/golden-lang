package build

import (
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/fs"
)

type Options struct {
	InputFilePath  string
	OutputFilePath string
	NumWorkers     int
	Debug          bool
}

func Build(opts Options) error {
	return errors.WithRecovery(func() { build(opts) })
}

func build(opts Options) {
	inputPath := opts.InputFilePath
	ext := fs.GetFileExtension(inputPath)
	if ext == "" {
		inputPath += ".gold"
	}

	if err := fs.CheckFileExists(inputPath); err != nil {
		errors.Throw(errors.InvalidFileError, "file %s not found", err)
	}

	if !fs.IsFileExtension(inputPath, ".gold", false) {
		errors.Throw(errors.InvalidFileError, "invalid file extension, expected '.gold', but received '%s'", fs.GetFileExtension(inputPath))
	}

	if err := fs.CheckFilePermissions(inputPath); err != nil {
		errors.Throw(errors.InvalidFileError, "file '%s 'does not have the read permission", inputPath)
	}

	path, err := fs.GetAbsolutePath(inputPath)
	if err != nil {
		errors.Throw(errors.InvalidFileError, "could not get the absolute path")
	}

	if name := fs.ModulePath_To_ModuleName(path); !fs.IsModuleNameValid(name) {
		errors.Throw(errors.InvalidFileError, "invalid module name, expected a valid module name, but received '%s'", name)
	}

	ctx := NewContext()
	ctx.Options = opts
	ctx.EntryModulePath = path
	ctx.ScheduleDiscoverPackage(path)

	steps := createSteps(ctx)
	for i := 0; i < opts.NumWorkers; i++ {
		go startWorker(steps, ctx)
	}

	<-ctx.Done
}
