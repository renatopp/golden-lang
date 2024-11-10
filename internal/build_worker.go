package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/renatopp/golden/internal/fs"
	"github.com/renatopp/golden/internal/logger"
)

type BuildWorker struct {
	id       int
	pipeline *BuildPipeline
}

func NewBuildWorker(id int, pipeline *BuildPipeline) *BuildWorker {
	return &BuildWorker{
		id:       id,
		pipeline: pipeline,
	}
}

func (w *BuildWorker) Start() {
	for {
		select {
		case path := <-w.pipeline.toDiscover:
			w.discover(path)

		case path := <-w.pipeline.toPrepare:
			w.prepare(path)
		}

		if w.pipeline.PendingModuleCount.Load() > 0 {
			continue
		}

		w.analyze()
		w.codegen()
	}
}

// From a given module path, discover all the modules in the same package
func (w *BuildWorker) discover(modulePath string) {
	defer w.pipeline.AckModule()

	logger.Debug("[worker:discover] discovering package of: %s", modulePath)
	files := fs.DiscoverModules(modulePath)
	for _, file := range files {
		if !w.pipeline.PreRegisterModule(file) {
			continue
		}
		w.pipeline.Prepare(file)
	}
}

// For a given module path, prepare (lex, parse, pre-analyze) the file for the analysis
func (w *BuildWorker) prepare(modulePath string) {
	defer w.pipeline.AckModule()

	logger.Debug("[worker:prepare] preparing file: %s", modulePath)

	bytes, err := os.ReadFile(modulePath)
	if err != nil {
		panic(err)
	}

	logger.Trace("[worker:prepare] lexing: %s", modulePath)
	tokens, err := Lex(bytes)
	if err != nil {
		panic(err)
	}

	// Annotate the tokens with the file information
	for _, token := range tokens {
		loc := token.Loc
		loc.Filename = modulePath
		token.Loc = loc
	}

	logger.Trace("[worker:prepare] parsing: %s", modulePath)
	root, err := Parse(tokens)
	if err != nil {
		panic(err)
	}

	// Annotate the module with the package and file information
	module := NewModule()
	module.Node = root
	module.Ast = root.Data.(*AstModule)
	module.Path = modulePath
	module.Name = fs.ModulePath_To_ModuleName(modulePath)
	module.FileName = fs.ModulePath_To_ModuleFileName(modulePath)
	for _, imp := range module.Ast.Imports {
		module.Imports = append(module.Imports, &Import{
			Name:  imp.Path,
			Alias: imp.Alias,
		})
	}

	// Create the package reference
	packageName := fs.ModulePath_To_PackageName(modulePath)
	packagePath := fs.ModulePath_To_PackagePath(modulePath)
	pkg := w.pipeline.CreateOrGetPackage(packageName, packagePath)
	module.Package = pkg
	pkg.Modules.Append(module)
	w.pipeline.RegisterModule(module)

	// Schedule imports for discovery
	for _, imp := range module.Imports {
		modulePath := fs.ImportName_To_ModulePath(imp.Name)

		if err := fs.CheckFileExists(modulePath); err != nil {
			panic(fmt.Sprintf("file '%s' does not exist. Remember that module names must be lower snake case, including the extension.", modulePath))
		}

		moduleName := fs.ModulePath_To_ModuleName(imp.Name)
		if !fs.IsModuleNameValid(moduleName) {
			panic(fmt.Sprintf("invalid module name '%s'. Remember that module names must be lower snake case.", moduleName))
		}

		imp.Path = modulePath
		w.pipeline.Discover(modulePath)
		// TODO: check if import is a project package or a core package
		// if package starts with @, it is a project package
		// if package matches the core packages (from a map?), it is a core package
		// otherwise, search package in the GOLDENPATH
		// if package is not found, error
	}
}

// Analyze all the modules in the pipeline
func (w *BuildWorker) analyze() {
	logger.Debug("[worker:analyze] analyzing modules")

	logger.Trace("[worker:analyze] building dependency graph")
	entryModule, _ := w.pipeline.Modules.Get(w.pipeline.EntryModulePath)
	w.buildDependencyGraph(entryModule)

	logger.Trace("[worker:analyze] analyzing dependency cycles")
	orderedPackages := w.checkDependencyGraph(entryModule.Package)

	for _, pkg := range orderedPackages {
		logger.Debug("[worker:analyze] analyzing package: %s", pkg.Path)
		mods := pkg.Modules.Values()
		for _, module := range mods {
			module.Scope = w.pipeline.GlobalScope.New()
			module.Analyzer = NewAnalyzer(module)
			module.Node.WithType(NewModuleType(module.Name, module))
		}

		for _, module := range mods {
			for _, other := range mods {
				if module == other {
					continue
				}
				module.Scope.SetValue(other.Name, other.Node)
			}
		}

		for _, module := range mods {
			if err := module.Analyzer.PreAnalyzeTypes(); err != nil {
				panic(err)
			}
		}

		for _, module := range mods {
			if err := module.Analyzer.PreAnalyzeFunctions(); err != nil {
				panic(err)
			}
		}

		for _, module := range mods {
			if err := module.Analyzer.PreAnalyzeVariables(); err != nil {
				panic(err)
			}
		}

		//for _, module := range mods {
		//	println("# SCOPE OF", module.Path)
		//	println(module.Scope.String())
		//	println()
		//}
		for _, module := range mods {
			if err := module.Analyzer.Analyze(); err != nil {
				panic(err)
			}
		}

		//for _, module := range mods {
		//	println("# SCOPE OF", module.Path)
		//	println(module.Scope.String())
		//	println()
		//}
	}

	logger.Trace("[worker:analyze] checking main function")
	w.checkMainFunction()
	w.pipeline.done <- nil
}

func (w *BuildWorker) buildDependencyGraph(module *Module) {
	if module.DependsOn.Len() > 0 {
		return
	}

	for _, mod := range module.Package.Modules.Values() {
		if mod == module {
			continue
		}
		module.DependsOn.Set(mod.Path, mod)
	}

	for _, imp := range module.Imports {
		imp.Module, _ = w.pipeline.Modules.Get(imp.Path)
		imp.Package = imp.Module.Package

		if module == imp.Module {
			panic(fmt.Sprintf("module '%s' cannot import itself", module.Path))
		}

		module.DependsOn.Set(imp.Module.Path, imp.Module)
		if module.Package != imp.Package {
			module.Package.DependsOn.Set(imp.Package.Path, imp.Package)
		}

		w.buildDependencyGraph(imp.Module)
	}
}

func (w *BuildWorker) checkDependencyGraph(pkg *Package) []*Package {
	visited := map[string]bool{}
	stack := map[string]bool{}
	order := []*Package{}
	return w.checkDependencyGraphLoop(pkg, visited, stack, order)
}

func (w *BuildWorker) checkDependencyGraphLoop(pkg *Package, visited, stack map[string]bool, order []*Package) []*Package {
	visited[pkg.Path] = true
	stack[pkg.Path] = true
	for _, dep := range pkg.DependsOn.Values() {
		if !visited[dep.Path] {
			order = w.checkDependencyGraphLoop(dep, visited, stack, order)

		} else if stack[dep.Path] {
			panic(fmt.Sprintf("cyclic dependency detected: %s", dep.Path))
		}
	}
	stack[pkg.Path] = false
	return append(order, pkg)
}

func (w *BuildWorker) checkMainFunction() {
	main, _ := w.pipeline.Modules.Get(w.pipeline.EntryModulePath)

	mainFunc := main.Scope.GetValue("main")
	if mainFunc == nil {
		panic("function 'main' not found")
	}

	mainFuncType := mainFunc.Type.(*FunctionType)
	if mainFuncType.ret != Void {
		panic("function 'main' must not return any value")
	}
	if len(mainFuncType.args) > 0 {
		panic("function 'main' must not have any parameter")
	}
}

// Code Generation
func (w *BuildWorker) codegen() {
	pkgs := w.pipeline.Packages.Values()
	for _, pkg := range pkgs {
		code, err := CodeGen_C(pkg)
		if err != nil {
			panic(err)
		}

		name := strings.ReplaceAll(pkg.Name, "/", "_")
		name = strings.ReplaceAll(name, "@", "main")
		os.WriteFile(".out/"+name+".c", []byte(code), 0644)
	}
}
