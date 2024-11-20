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

	pkg := &Package{
		Name:    fs.PackagePath2PackageName(packagePath),
		Path:    packagePath,
		Modules: ds.NewSyncList[*Module](),
	}
	l.ctx.PackageRegistry.Set(packagePath, pkg)
	files := fs.DiscoverModules(packagePath)
	for _, modulePath := range files {
		l.pending.Add(1)
		go l.loadModule(pkg, modulePath)
	}
}

func (l *loader) loadModule(pkg *Package, modulePath string) {
	defer l.pending.Done()

	module := &Module{
		Name:     fs.ModulePath2ModuleName(modulePath),
		Path:     modulePath,
		FileName: fs.ModulePath2ModuleFileName(modulePath),
		Package:  pkg,
		Root:     nil,
	}

	bytes, err := os.ReadFile(modulePath)
	if err != nil {
		l.errors.Add(
			errors.NewError(errors.InvalidFileError, "could not read module '%s', reason: %v", modulePath, err),
		)
		return
	}

	tokens, err := syntax.Lex(bytes, modulePath)
	if err != nil {
		l.errors.Add(err)
		return
	}

	if l.ctx.Options.OnTokensReady != nil {
		l.ctx.Options.OnTokensReady(module, tokens)
	}

	root, err := syntax.Parse(tokens, modulePath)
	if err != nil {
		l.errors.Add(err)
		return
	}
	module.Root = root

	if l.ctx.Options.OnAstReady != nil {
		l.ctx.Options.OnAstReady(module, root)
	}

	pkg.Modules.Add(module)
	l.ctx.ModuleRegistry.Set(modulePath, module)

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
