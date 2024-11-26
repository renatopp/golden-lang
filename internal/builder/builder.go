package builder

import (
	"strings"
	"sync"
	"time"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/codegen"
	"github.com/renatopp/golden/internal/compiler/env"
	"github.com/renatopp/golden/internal/compiler/semantic"
	"github.com/renatopp/golden/internal/compiler/types"
	"github.com/renatopp/golden/internal/helpers/ds"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/events"
	"github.com/renatopp/golden/internal/helpers/fs"
	"github.com/renatopp/golden/lang"
)

//
//
//

type BuildContext struct {
	Options         *BuildOptions
	PackageRegistry *ds.SyncMap[string, *Package]
	ModuleRegistry  *ds.SyncMap[string, *Module]
	EntryPackage    *Package
	EntryModule     *Module
	DependencyOrder []*Package
	GlobalScope     *env.Scope
}

//
//
//

type BuildOptions struct {
	EntryFilePath          string // Absolute path of the entry file containing main function
	WorkingDir             string // Absolute path of the working directory
	OnTokensReady          *events.Signal2[*Module, []*lang.Token]
	OnAstReady             *events.Signal2[*Module, *ast.Module]
	OnDependencyGraphReady *events.Signal1[[]*Package]
	OnTypeCheckReady       *events.Signal3[*Module, *ast.Module, *env.Scope]
}

func NewBuildOptions(fileName string) *BuildOptions {
	return &BuildOptions{
		EntryFilePath:          fileName,
		WorkingDir:             fs.GetWorkingDir(),
		OnTokensReady:          events.NewSignal2[*Module, []*lang.Token](),
		OnAstReady:             events.NewSignal2[*Module, *ast.Module](),
		OnDependencyGraphReady: events.NewSignal1[[]*Package](),
		OnTypeCheckReady:       events.NewSignal3[*Module, *ast.Module, *env.Scope](),
	}
}

//
//
//

type BuildResult struct {
	Elapsed time.Duration
}

//
//
//

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
		PackageRegistry: ds.NewSyncMap[string, *Package](),
		ModuleRegistry:  ds.NewSyncMap[string, *Module](),
		EntryPackage:    nil,
		EntryModule:     nil,
	}

	fs.WorkingDir = ctx.Options.WorkingDir
	validateEntry(ctx)
	loadPackages(ctx)
	checkEntries(ctx)
	buildDependencyGraph(ctx)
	buildGlobalScope(ctx)
	semanticAnalysis(ctx)
	checkMain(ctx)
	generateCode(ctx)

	return res
}

func validateEntry(ctx *BuildContext) {
	inputPath := ctx.Options.EntryFilePath

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

	ctx.Options.EntryFilePath = absPath
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

func checkEntries(ctx *BuildContext) {
	modulePath := ctx.Options.EntryFilePath
	packagePath := fs.ModulePath2PackagePath(modulePath)
	ctx.EntryModule, _ = ctx.ModuleRegistry.Get(modulePath)
	ctx.EntryPackage, _ = ctx.PackageRegistry.Get(packagePath)
}

func buildDependencyGraph(ctx *BuildContext) {
	registry := ctx.PackageRegistry.Items()
	entry := ctx.EntryPackage.Path

	visited := map[string]bool{}
	stack := map[string]bool{}
	order := []*Package{}
	pkg := registry[entry]
	ctx.DependencyOrder = buildDependencyGraphLoop(registry, pkg, visited, stack, order)
	ctx.Options.OnDependencyGraphReady.Emit(ctx.DependencyOrder)
}

func buildDependencyGraphLoop(registry map[string]*Package, pkg *Package, visited, stack map[string]bool, order []*Package) []*Package {
	visited[pkg.Path] = true
	stack[pkg.Path] = true
	for _, dep := range pkg.Imports.Values() {
		if !visited[dep] {
			p := registry[dep]
			order = buildDependencyGraphLoop(registry, p, visited, stack, order)

		} else if stack[dep] {
			names := []string{}
			for k := range stack {
				names = append(names, k)
			}
			deps := strings.Join(names, "\n- ")
			errors.Throw(errors.CircularReferenceError, "cyclic dependency detected importing packages: \n- %s", deps)
		}
	}
	stack[pkg.Path] = false
	return append(order, pkg)
}

func buildGlobalScope(ctx *BuildContext) {
	ctx.GlobalScope = env.NewScope()
	ctx.GlobalScope.Types.Set(types.Int.Signature(), env.B(types.Int))
	ctx.GlobalScope.Types.Set(types.Float.Signature(), env.B(types.Float))
	ctx.GlobalScope.Types.Set(types.Bool.Signature(), env.B(types.Bool))
	ctx.GlobalScope.Types.Set(types.String.Signature(), env.B(types.String))
	ctx.GlobalScope.Types.Set(types.Void.Signature(), env.B(types.Void))
}

func semanticAnalysis(ctx *BuildContext) {
	checker := semantic.NewTypeChecker()

	for _, pkg := range ctx.DependencyOrder {
		mods := pkg.Modules.Values()

		// create type instances for all modules
		for _, mod := range mods {
			scope := ctx.GlobalScope.New()
			scope.IsModule = true
			mod.Root.SetType(types.NewModule(mod.Root, mod.Path, scope))
		}

		// attach type instances to the module scopes
		for _, mod := range mods {
			modType := mod.Root.Type().(*types.Module)

			for _, other := range mods {
				if mod == other {
					continue
				}

				alias := fs.ModulePath2ModuleName(other.Path)
				modType.Scope.Values.Set(alias, env.B(other.Root.Type()))
			}
		}

		// pre-resolve all types, functions and module variables
		for _, mod := range mods {
			checker.PreResolve(mod.Root)
		}

		// resolve everything
		for _, mod := range mods {
			checker.Resolve(mod.Root)
			ctx.Options.OnTypeCheckReady.Emit(mod, mod.Root, mod.Root.Type().(*types.Module).Scope)
		}
	}
}

func checkMain(ctx *BuildContext) {
	main := ctx.EntryModule.Scope().Values.Get("main")
	if main == nil {
		errors.Throw(errors.InvalidEntryFile, "entry module '%s' does not contain a 'main' function", ctx.EntryModule.Path)
	}

	if !types.NoopFn.Compatible(main.Type) {
		errors.Throw(errors.InvalidEntryFile, "entry module '%s' 'main' function has an invalid signature", ctx.EntryModule.Path)
	}
}

func generateCode(ctx *BuildContext) {
	cg := codegen.NewCodegen()

	cg.StartGeneration()
	for _, pkg := range ctx.DependencyOrder {
		cg.StartPackage()
		mods := pkg.Modules.Values()
		for _, mod := range mods {
			mod.Root.Accept(cg)
		}
		cg.EndPackage()
	}
	cg.EndGeneration()
}
