package builder

import (
	"os"
	"sync"

	"github.com/renatopp/golden/internal/compiler/syntax"
	"github.com/renatopp/golden/internal/helpers/ds"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/fs"
)

type loader struct {
	ctx     *BuildContext
	errors  *ds.SyncList[error]
	pending sync.WaitGroup
}

func loadPackages(ctx *BuildContext) {
	l := &loader{
		ctx:     ctx,
		errors:  ds.NewSyncList[error](),
		pending: sync.WaitGroup{},
	}
	l.discover(ctx.Options.EntryFilePath)
	l.pending.Wait()

	if l.errors.Len() > 0 {
		e, _ := l.errors.Get(0)
		errors.Rethrow(e)
	}
}

func (l *loader) discover(modulePath string) {
	if l.errors.Len() > 0 {
		return
	}

	packagePath := fs.ModulePath2PackagePath(modulePath)
	if ok := l.ctx.PackageRegistry.SetFirst(packagePath, nil); !ok {
		return
	}

	l.pending.Add(1)
	go l.loadPackage(packagePath)
}

func (l *loader) loadPackage(packagePath string) {
	defer l.pending.Done()

	pkg := NewPackage(
		fs.PackagePath2PackageName(packagePath),
		packagePath,
	)
	l.ctx.PackageRegistry.Set(packagePath, pkg)
	files := fs.DiscoverModules(packagePath)
	for _, modulePath := range files {
		l.pending.Add(1)
		go l.loadModule(pkg, modulePath)
	}
}

func (l *loader) loadModule(pkg *Package, modulePath string) {
	defer l.pending.Done()

	// Create the module
	module := NewModule(
		fs.ModulePath2ModuleName(modulePath),
		modulePath,
		fs.ModulePath2ModuleFileName(modulePath),
		pkg,
	)

	// Read the bytes
	bytes, err := os.ReadFile(modulePath)
	if err != nil {
		l.errors.Add(
			errors.NewError(errors.InvalidFileError, "could not read module '%s', reason: %v", modulePath, err),
		)
		return
	}

	// Convert bytes to tokens
	tokens, err := syntax.Lex(bytes, modulePath)
	if err != nil {
		l.errors.Add(err)
		return
	}

	if l.ctx.Options.OnTokensReady != nil {
		l.ctx.Options.OnTokensReady(module, tokens)
	}

	// Convert tokens to AST
	root, err := syntax.Parse(tokens, modulePath)
	if err != nil {
		l.errors.Add(err)
		return
	}
	module.Root = root

	if l.ctx.Options.OnAstReady != nil {
		l.ctx.Options.OnAstReady(module, root)
	}

	// Add the module to the package
	pkg.Modules.Add(module)
	l.ctx.ModuleRegistry.Set(modulePath, module)

	// Discover imports
	for _, a := range root.Imports {
		path := fs.ImportName2ModulePath(a.Path.Literal)
		alias := fs.ModulePath2ModuleName(path)
		if a.Alias.Has() {
			alias = a.Alias.Unwrap().Literal
		}

		module.Imports = append(module.Imports, &ModuleImport{
			Path:  path,
			Alias: alias,
		})

		packagePath := fs.ModulePath2PackagePath(path)
		module.Package.Imports.AddUnique(packagePath)

		l.discover(path)
	}
}
