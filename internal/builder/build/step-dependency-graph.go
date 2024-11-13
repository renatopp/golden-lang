package build

import (
	"fmt"

	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/logger"
)

type StepDependencyGraph struct {
	ctx *Context
}

func NewStepDependencyGraph(ctx *Context) *StepDependencyGraph {
	return &StepDependencyGraph{ctx: ctx}
}

func (s *StepDependencyGraph) Process() []*core.Package {
	logger.Debug("step:dependency-graph] starting dependency graph")

	logger.Trace("[step:dependency-graph] building dependency graph")
	entryModule, _ := s.ctx.Modules.Get(s.ctx.EntryModulePath)
	s.buildDependencyGraph(entryModule)

	logger.Trace("[step:dependency-graph] checking dependency graph")
	orderedPackages := s.checkDependencyGraph(entryModule.Package)

	return orderedPackages
}

func (s *StepDependencyGraph) buildDependencyGraph(module *core.Module) {
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
		imp.Module, _ = s.ctx.Modules.Get(imp.Path)
		imp.Package = imp.Module.Package

		if module == imp.Module {
			panic(fmt.Sprintf("module '%s' cannot import itself", module.Path))
		}

		module.DependsOn.Set(imp.Module.Path, imp.Module)
		if module.Package != imp.Package {
			module.Package.DependsOn.Set(imp.Package.Path, imp.Package)
		}

		s.buildDependencyGraph(imp.Module)
	}
}

func (s *StepDependencyGraph) checkDependencyGraph(pkg *core.Package) []*core.Package {
	visited := map[string]bool{}
	stack := map[string]bool{}
	order := []*core.Package{}
	return s.checkDependencyGraphLoop(pkg, visited, stack, order)
}

func (s *StepDependencyGraph) checkDependencyGraphLoop(pkg *core.Package, visited, stack map[string]bool, order []*core.Package) []*core.Package {
	visited[pkg.Path] = true
	stack[pkg.Path] = true
	for _, dep := range pkg.DependsOn.Values() {
		if !visited[dep.Path] {
			order = s.checkDependencyGraphLoop(dep, visited, stack, order)

		} else if stack[dep.Path] {
			panic(fmt.Sprintf("cyclic dependency detected: %s", dep.Path))
		}
	}
	stack[pkg.Path] = false
	return append(order, pkg)
}
