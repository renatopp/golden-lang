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

func (l *loader) discover(modulePath string) {
	l.pending.Add(1)
	go l.loadModule(modulePath)
}

func (l *loader) loadModule(modulePath string) {
	defer l.pending.Done()

	// Create the file
	file := NewFile(
		fs.ModulePath2ModuleName(modulePath),
		modulePath,
		fs.ModulePath2ModuleFileName(modulePath),
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
	lexer := syntax.NewLexer(modulePath, bytes)
	tokens, err := lexer.Lex()
	if err != nil {
		l.errors.Add(err)
		return
	}
	l.ctx.Options.OnTokensReady.Emit(file, tokens)

	// Convert tokens to AST
	// parser := syntax.NewParser(tokens, modulePath)
	// root, err := syntax.Parse(tokens, modulePath)
	// if err != nil {
	// 	l.errors.Add(err)
	// 	return
	// }
	// module.Root = root
	// l.ctx.Options.OnAstReady.Emit(module, root)

	// // Add the module to the package
	// pkg.Modules.Add(module)
	// l.ctx.ModuleRegistry.Set(modulePath, module)

	// // Discover imports
	// for _, a := range root.Imports {
	// 	path := fs.ImportName2ModulePath(a.Path.Literal)
	// 	alias := fs.ModulePath2ModuleName(path)
	// 	if a.Alias.Has() {
	// 		alias = a.Alias.Unwrap().Literal
	// 	}

	// 	module.Imports = append(module.Imports, &ModuleImport{
	// 		Path:  path,
	// 		Alias: alias,
	// 	})

	// 	packagePath := fs.ModulePath2PackagePath(path)
	// 	module.Package.Imports.AddUnique(packagePath)

	// 	if fs.CheckFolderExists(packagePath) != nil {
	// 		l.errors.Add(errors.NewError(errors.InvalidFolderError, "could not find package '%s'", packagePath).WithNode(a.Path))
	// 	}

	// 	if fs.CheckFileExists(path) != nil {
	// 		l.errors.Add(errors.NewError(errors.InvalidFileError, "could not find module '%s'", path).WithNode(a.Path))
	// 	}

	// 	l.discover(path)
	// }
}

// func (l *loader) discover(modulePath string) {
// 	if l.errors.Len() > 0 {
// 		return
// 	}

// 	packagePath := fs.ModulePath2PackagePath(modulePath)
// 	if ok := l.ctx.PackageRegistry.SetFirst(packagePath, nil); !ok {
// 		return
// 	}

// 	l.pending.Add(1)
// 	go l.loadPackage(packagePath)
// }

// func (l *loader) loadPackage(packagePath string) {
// 	defer l.pending.Done()

// 	pkg := NewPackage(
// 		fs.PackagePath2PackageName(packagePath),
// 		packagePath,
// 	)
// 	l.ctx.PackageRegistry.Set(packagePath, pkg)
// 	files := fs.DiscoverModules(packagePath)
// 	for _, modulePath := range files {
// 		l.pending.Add(1)
// 		go l.loadModule(pkg, modulePath)
// 	}
// }

// func (l *loader) loadModule(pkg *Package, modulePath string) {
// 	defer l.pending.Done()

// 	// Create the module
// 	module := NewFile(
// 		fs.ModulePath2ModuleName(modulePath),
// 		modulePath,
// 		fs.ModulePath2ModuleFileName(modulePath),
// 		pkg,
// 	)

// 	// Read the bytes
// 	bytes, err := os.ReadFile(modulePath)
// 	if err != nil {
// 		l.errors.Add(
// 			errors.NewError(errors.InvalidFileError, "could not read module '%s', reason: %v", modulePath, err),
// 		)
// 		return
// 	}

// 	// Convert bytes to tokens
// 	tokens, err := syntax.Lex(bytes, modulePath)
// 	if err != nil {
// 		l.errors.Add(err)
// 		return
// 	}
// 	l.ctx.Options.OnTokensReady.Emit(module, tokens)

// 	// Convert tokens to AST
// 	root, err := syntax.Parse(tokens, modulePath)
// 	if err != nil {
// 		l.errors.Add(err)
// 		return
// 	}
// 	module.Root = root
// 	l.ctx.Options.OnAstReady.Emit(module, root)

// 	// Add the module to the package
// 	pkg.Modules.Add(module)
// 	l.ctx.ModuleRegistry.Set(modulePath, module)

// 	// Discover imports
// 	for _, a := range root.Imports {
// 		path := fs.ImportName2ModulePath(a.Path.Literal)
// 		alias := fs.ModulePath2ModuleName(path)
// 		if a.Alias.Has() {
// 			alias = a.Alias.Unwrap().Literal
// 		}

// 		module.Imports = append(module.Imports, &ModuleImport{
// 			Path:  path,
// 			Alias: alias,
// 		})

// 		packagePath := fs.ModulePath2PackagePath(path)
// 		module.Package.Imports.AddUnique(packagePath)

// 		if fs.CheckFolderExists(packagePath) != nil {
// 			l.errors.Add(errors.NewError(errors.InvalidFolderError, "could not find package '%s'", packagePath).WithNode(a.Path))
// 		}

// 		if fs.CheckFileExists(path) != nil {
// 			l.errors.Add(errors.NewError(errors.InvalidFileError, "could not find module '%s'", path).WithNode(a.Path))
// 		}

// 		l.discover(path)
// 	}
// }
