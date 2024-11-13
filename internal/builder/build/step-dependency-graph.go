package build

import (
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/errors"
)

type StepDependencyGraph struct {
	ctx *Context
}

func NewStepDependencyGraph(ctx *Context) *StepDependencyGraph {
	return &StepDependencyGraph{ctx: ctx}
}

func (s *StepDependencyGraph) Process() []*core.Package {
	// Generate the dependency graph
	entryModule, _ := s.ctx.Modules.Get(s.ctx.EntryModulePath)
	s.buildDependencyGraph(entryModule)

	// Check validity of graph and returns the packages ordered by last dependency
	return s.checkDependencyGraph(entryModule.Package)
}

func (s *StepDependencyGraph) buildDependencyGraph(module *core.Module) {
	// Skip if already processed
	if module.DependsOn.Len() > 0 {
		return
	}

	// Setup implicit dependencies for modules in the same package
	for _, mod := range module.Package.Modules.Values() {
		if mod == module {
			continue
		}
		module.DependsOn.Set(mod.Path, mod)
	}

	// Setup explicit dependencies
	for _, imp := range module.Imports {
		imp.Module, _ = s.ctx.Modules.Get(imp.Path)
		imp.Package = imp.Module.Package

		// Check for circular dependencies
		if module == imp.Module {
			errors.ThrowAtNode(imp.Node, errors.CircularReferenceError, "module '%s' cannot import itself", module.Path)
		}

		// Add dependencies for both module and package
		module.DependsOn.Set(imp.Module.Path, imp.Module)
		if module.Package != imp.Package {
			module.Package.DependsOn.Set(imp.Package.Path, imp.Package)
		}

		// Recursively build dependency graph
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
			errors.Throw(errors.CircularReferenceError, "cyclic dependency detected: %s", dep.Path)
		}
	}
	stack[pkg.Path] = false
	return append(order, pkg)
}
