package builder

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	// "github.com/renatopp/golden/internal/compiler/codegen"
	// "github.com/renatopp/golden/internal/compiler/env"
	// "github.com/renatopp/golden/internal/compiler/semantic"
	// "github.com/renatopp/golden/internal/compiler/types"

	"github.com/renatopp/golden/internal/compiler/env"
	"github.com/renatopp/golden/internal/helpers/ds"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/fs"
)

type BuildResult struct {
	Elapsed time.Duration
}

//
//
//

type Builder struct {
	ctx  *BuildContext
	opts *BuildOptions
}

func NewBuilder(opts *BuildOptions) *Builder {
	return &Builder{
		ctx:  nil,
		opts: opts,
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
	b.ctx = &BuildContext{
		Options:        b.opts,
		ModuleRegistry: ds.NewSyncMap[string, *File](),
		EntryModule:    nil,
	}

	fs.WorkingDir = b.ctx.Options.WorkingDir
	b.validateEntry()
	b.checkCacheFolders()
	b.loadModules()
	b.checkEntries()
	b.buildDependencyGraph()
	b.buildGlobalScope()
	b.semanticAnalysis()
	// b.checkMain()
	// b.generateCode()

	return res
}

func (b *Builder) validateEntry() {
	inputPath := b.ctx.Options.EntryFilePath

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

	b.ctx.Options.EntryFilePath = absPath
}

func (b *Builder) checkCacheFolders() {
	if err := fs.GuaranteeDirectoryExists(b.opts.GlobalCachePath); err != nil {
		errors.Throw(errors.InternalError, "could not create global cache path")
	}
	if err := fs.GuaranteeDirectoryExists(b.opts.GlobalTargetPath); err != nil {
		errors.Throw(errors.InternalError, "could not create global target path")
	}
	if err := fs.GuaranteeDirectoryExists(b.opts.LocalCachePath); err != nil {
		errors.Throw(errors.InternalError, "could not create local cache path")
	}
	if err := fs.GuaranteeDirectoryExists(b.opts.LocalTargetPath); err != nil {
		errors.Throw(errors.InternalError, "could not create local target path")
	}
	outputDir := filepath.Dir(b.opts.OutputFilePath)
	if err := fs.GuaranteeDirectoryExists(outputDir); err != nil {
		errors.Throw(errors.InternalError, "could not create output file path")
	}
}

func (b *Builder) loadModules() {
	l := &loader{
		ctx:     b.ctx,
		errors:  ds.NewSyncList[error](),
		pending: sync.WaitGroup{},
	}
	l.discover(b.ctx.Options.EntryFilePath)
	l.pending.Wait()

	if l.errors.Len() > 0 {
		e, _ := l.errors.Get(0)
		errors.Rethrow(e)
	}
}

func (b *Builder) checkEntries() {
	modulePath := b.ctx.Options.EntryFilePath
	b.ctx.EntryModule, _ = b.ctx.ModuleRegistry.Get(modulePath)
}

func (b *Builder) buildDependencyGraph() {
	registry := b.ctx.ModuleRegistry.Items()
	entry := b.ctx.EntryModule.Path

	visited := map[string]bool{}
	stack := map[string]bool{}
	order := []*File{}
	pkg := registry[entry]
	b.ctx.DependencyOrder = b.buildDependencyGraphLoop(registry, pkg, visited, stack, order)
	b.ctx.Options.OnDependencyGraphReady.Emit(b.ctx.DependencyOrder)
}

func (b *Builder) buildDependencyGraphLoop(registry map[string]*File, file *File, visited, stack map[string]bool, order []*File) []*File {
	visited[file.Path] = true
	stack[file.Path] = true
	for _, dep := range []string{} { // TODO: file.Imports.Values() {
		if !visited[dep] {
			p := registry[dep]
			order = b.buildDependencyGraphLoop(registry, p, visited, stack, order)

		} else if stack[dep] {
			names := []string{}
			for k := range stack {
				names = append(names, k)
			}
			deps := strings.Join(names, "\n- ")
			errors.Throw(errors.CircularReferenceError, "cyclic dependency detected importing packages: \n- %s", deps)
		}
	}
	stack[file.Path] = false
	return append(order, file)
}

func (b *Builder) buildGlobalScope() {
	b.ctx.GlobalScope = env.NewScope()
	// b.ctx.GlobalScope.Types.Set(types.Int.GetSignature(), env.B(types.Int))
	// b.ctx.GlobalScope.Types.Set(types.Float.GetSignature(), env.B(types.Float))
	// b.ctx.GlobalScope.Types.Set(types.Bool.GetSignature(), env.B(types.Bool))
	// b.ctx.GlobalScope.Types.Set(types.String.GetSignature(), env.B(types.String))
	// b.ctx.GlobalScope.Types.Set(types.Void.GetSignature(), env.B(types.Void))
}

func (b *Builder) semanticAnalysis() {
	// checker := semantic.NewTypeChecker()

	// for _, pkg := range b.ctx.DependencyOrder {
	// 	mods := pkg.Modules.Values()

	// 	// create type instances for all modules
	// 	for _, mod := range mods {
	// 		scope := b.ctx.GlobalScope.New()
	// 		scope.IsModule = true
	// 		mod.Root.SetType(types.NewModule(mod.Root, mod.Path, scope))
	// 	}

	// 	// attach type instances to the module scopes
	// 	for _, mod := range mods {
	// 		modType := mod.Root.Type().(*types.Module)

	// 		for _, other := range mods {
	// 			if mod == other {
	// 				continue
	// 			}

	// 			alias := fs.ModulePath2ModuleName(other.Path)
	// 			modType.Scope.Values.Set(alias, env.B(other.Root.Type()))
	// 		}
	// 	}

	// 	// pre-resolve all types, functions and module variables
	// 	for _, mod := range mods {
	// 		checker.PreResolve(mod.Root)
	// 	}

	// 	// resolve everything
	// 	for _, mod := range mods {
	// 		checker.Resolve(mod.Root)
	// 		b.ctx.Options.OnTypeCheckReady.Emit(mod, mod.Root, mod.Root.Type().(*types.Module).Scope)
	// 	}
	// }
}

// func (b *Builder) checkMain() {
// 	main := b.ctx.EntryModule.Scope().Values.Get("main")
// 	if main == nil {
// 		errors.Throw(errors.InvalidEntryFile, "entry module '%s' does not contain a 'main' function", b.ctx.EntryModule.Path)
// 	}

// 	if !types.NoopFn.Compatible(main.Type) {
// 		errors.Throw(errors.InvalidEntryFile, "entry module '%s' 'main' function has an invalid signature", b.ctx.EntryModule.Path)
// 	}
// }

// func (b *Builder) generateCode() {
// 	cg := codegen.NewCodegen(b.opts.LocalTargetPath)

// 	cg.StartGeneration()
// 	for _, pkg := range b.ctx.DependencyOrder {
// 		imports := pkg.Imports.Values()
// 		cg.StartPackage(pkg.Path, imports)
// 		mods := pkg.Modules.Values()
// 		for _, mod := range mods {
// 			mod.Root.Accept(cg)
// 		}
// 		cg.EndPackage()
// 	}
// 	cg.EndGeneration()
// }
