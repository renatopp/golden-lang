package builder

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/env"
	"github.com/renatopp/golden/internal/compiler/optimizations"
	"github.com/renatopp/golden/internal/compiler/semantic"
	"github.com/renatopp/golden/internal/compiler/types"
	"github.com/renatopp/golden/internal/helpers/ds"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/fs"
	"github.com/renatopp/golden/internal/helpers/safe"
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
	if err == nil {
		return res, errors.WithRecovery(b.generateOutput)
	}
	return res, err
}

func (b *Builder) Run() (res *BuildResult, err error) {
	err = errors.WithRecovery(func() {
		start := time.Now()
		res = b.build()
		res.Elapsed = time.Since(start)
	})
	if err == nil {
		return res, errors.WithRecovery(b.run)
	}
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
	b.checkMain()
	b.applyOptimizations()
	b.generateCode()

	return res
}

func (b *Builder) run() {
	b.runCode()
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
	b.ctx.GlobalScope.Types.Set(types.Int.GetSignature(), env.TB(types.Int, nil))
	b.ctx.GlobalScope.Types.Set(types.Float.GetSignature(), env.TB(types.Float, nil))
	b.ctx.GlobalScope.Types.Set(types.Bool.GetSignature(), env.TB(types.Bool, nil))
	b.ctx.GlobalScope.Types.Set(types.String.GetSignature(), env.TB(types.String, nil))
	b.ctx.GlobalScope.Types.Set(types.Void.GetSignature(), env.TB(types.Void, nil))
}

func (b *Builder) semanticAnalysis() {
	checker := semantic.NewChecker()
	mods := b.ctx.DependencyOrder
	// create type instances for all modules
	for _, mod := range mods {
		root := mod.Root.Unwrap()
		scope := b.ctx.GlobalScope.New()
		scope.IsModule = true
		mod.Root.Unwrap().SetType(types.NewModule(root, mod.Path, scope))
	}

	// attach type instances to the module scopes
	for _, mod := range mods {
		root := mod.Root.Unwrap()
		modType := root.GetType().Unwrap().(*types.Module)

		for _, other := range mods {
			if mod == other {
				continue
			}

			alias := fs.ModulePath2ModuleName(other.Path)
			otherRoot := other.Root.Unwrap()
			modType.Scope.Values.Set(alias, env.VB(otherRoot, otherRoot.GetType().Unwrap()))
		}
	}

	// pre-resolve all types, functions and module variables
	for _, mod := range mods {
		root := mod.Root.Unwrap()
		checker.PreCheck(root)
	}

	// resolve everything
	for _, mod := range b.ctx.DependencyOrder {
		root := mod.Root.Unwrap()
		_, err := checker.Check(root)
		if err != nil {
			errors.Rethrow(err)
		}
		b.ctx.Options.OnTypeCheckReady.Emit(mod, root, root.GetType().Unwrap().(*types.Module).Scope)
	}
}

func (b *Builder) checkMain() {
	main := b.ctx.EntryModule.Scope().Values.Get("main", nil)
	if main == nil {
		errors.Throw(errors.InvalidEntryFile, "entry module '%s' does not contain a 'main' function", b.ctx.EntryModule.Path)
	}

	if !types.NoopFn.IsCompatible(main.Type) {
		errors.Throw(errors.InvalidEntryFile, "entry module '%s' 'main' function has an invalid signature", b.ctx.EntryModule.Path)
	}
}

func (b *Builder) applyOptimizations() {
	pipeline := optimizations.NewPipeline(
		optimizations.NewAddReturnToFunctions(),
	)
	for _, mod := range b.ctx.DependencyOrder {
		mod.Root = safe.Some(pipeline.Run(mod.Root.Unwrap()).(*ast.Module))
		b.ctx.Options.OnOptimizationReady.Emit(mod, mod.Root.Unwrap())
	}
}

func (b *Builder) generateCode() {
	backend := b.opts.OutputTarget
	backend.Initialize(b.opts.LocalTargetPath)

	backend.BeforeCodeGeneration()
	for _, mod := range b.ctx.DependencyOrder {
		backend.GenerateCode(mod.Path, mod.Root.Unwrap(), mod == b.ctx.EntryModule)
	}
	backend.AfterCodeGeneration()
	backend.Finalize()
}

func (b *Builder) generateOutput() {
	backend := b.opts.OutputTarget
	backend.Build(b.opts.OutputFilePath)
}

func (b *Builder) runCode() {
	backend := b.opts.OutputTarget
	backend.Run()
}
