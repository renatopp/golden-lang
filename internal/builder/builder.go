package builder

import (
	"time"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/ds"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/fs"
	"github.com/renatopp/golden/lang"
)

type BuildContext struct {
	Options         *BuildOptions
	PackageRegistry *ds.SyncMap[string, *core.Package]
	ModuleRegistry  *ds.SyncMap[string, *core.Module]
}

type BuildOptions struct {
	EntryFilePath string // Absolute path of the entry file containing main function
	OnTokensReady func(module *core.Module, tokens []*lang.Token)
	OnAstReady    func(module *core.Module, root *ast.Module)
}

type BuildResult struct {
	Elapsed time.Duration
}

type Builder struct {
	Options *BuildOptions
}

func NewBuilder(opts *BuildOptions) *Builder {
	return &Builder{
		Options: opts,
	}
}

func (b *Builder) Build() (res *BuildResult, err error) {
	err = errors.WithRecovery(func() {
		start := time.Now()
		res = b.build()
		res.Elapsed = time.Since(start)
	})

	return res, err
}

func (b *Builder) build() *BuildResult {
	res := &BuildResult{}
	ctx := &BuildContext{
		Options:         b.Options,
		PackageRegistry: ds.NewSyncMap[string, *core.Package](),
		ModuleRegistry:  ds.NewSyncMap[string, *core.Module](),
	}

	b.validateEntry()
	loadPackages(ctx)

	return res
}

func (b *Builder) validateEntry() {
	inputPath := b.Options.EntryFilePath

	extension := fs.GetFileExtension(inputPath)
	if extension == "" {
		inputPath += ".gold"
	}

	if err := fs.CheckFileExists(inputPath); err != nil {
		errors.Throw(errors.InvalidFileError, "input file '%s' not found", inputPath)
	}

	if !fs.IsFileExtension(inputPath, ".gold", false) {
		errors.Throw(errors.InvalidFileError, "input file '%s' must have a '.gold' extension", inputPath)
	}

	if err := fs.CheckFilePermissions(inputPath); err != nil {
		errors.Throw(errors.InvalidFileError, "input file '%s' does not have read permissions", inputPath)
	}

	absPath, err := fs.GetAbsolutePath(inputPath)
	if err != nil {
		errors.Throw(errors.InvalidFileError, "input file '%s' does not have a valid path", inputPath)
	}

	if name := fs.ModulePath2ModuleName(absPath); !fs.IsModuleNameValid(name) {
		errors.Throw(errors.InvalidFileError, "input file '%s' does not have a valid name", inputPath)
	}

	b.Options.EntryFilePath = absPath
}
