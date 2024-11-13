package build

import (
	"fmt"

	"github.com/renatopp/golden/internal/compiler/semantic"
	"github.com/renatopp/golden/internal/compiler/semantic/types"
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/debug"
	"github.com/renatopp/golden/internal/helpers/errors"
)

type StepResolveBindings struct {
	ctx *Context
}

func NewStepResolveBindings(ctx *Context) *StepResolveBindings {
	return &StepResolveBindings{ctx: ctx}
}

func (s *StepResolveBindings) Process(packages []*core.Package) {
	for _, pkg := range packages {
		modules := pkg.Modules.Values()
		s.createScope(modules)
		s.attachPackageScopes(modules)
		s.preResolveTypes(modules)
		s.preResolveFunctions(modules)
		s.preResolveVariables(modules)
		s.resolve(modules)
	}

	if s.ctx.Options.Debug {
		module, _ := s.ctx.Modules.Get(s.ctx.EntryModulePath)
		s.debugPrintAst(module.Node)
	}
}

func (s *StepResolveBindings) createScope(modules []*core.Module) {
	for _, module := range modules {
		module.Scope = s.ctx.GlobalScope.New()
		module.Resolver = semantic.NewResolver(module)
		module.Node.WithType(types.NewModule(module.Name, module))
	}
}

func (s *StepResolveBindings) attachPackageScopes(modules []*core.Module) {
	for _, module := range modules {
		for _, other := range modules {
			if module == other {
				continue
			}
			module.Scope.Values.Set(other.Name, core.BindValue(other.Node))
		}
	}
}

func (s *StepResolveBindings) preResolveTypes(modules []*core.Module) {
	for _, module := range modules {
		node := module.Node.Data().(*ast.Module)
		for _, tp := range node.Types {
			if err := module.Resolver.PreResolve(tp); err != nil {
				errors.Rethrow(err)
			}
		}
	}
}

func (s *StepResolveBindings) preResolveFunctions(modules []*core.Module) {
	for _, module := range modules {
		node := module.Node.Data().(*ast.Module)
		for _, tp := range node.Functions {
			if err := module.Resolver.PreResolve(tp); err != nil {
				errors.Rethrow(err)
			}
		}
	}
}

func (s *StepResolveBindings) preResolveVariables(modules []*core.Module) {
	for _, module := range modules {
		node := module.Node.Data().(*ast.Module)
		for _, tp := range node.Variables {
			if err := module.Resolver.PreResolve(tp); err != nil {
				errors.Rethrow(err)
			}
		}
	}
}

func (s *StepResolveBindings) resolve(modules []*core.Module) {
	for _, module := range modules {
		if err := module.Resolver.Resolve(module.Node); err != nil {
			errors.Rethrow(err)
		}
	}
}

func (s *StepResolveBindings) debugPrintAst(root *core.AstNode) {
	fmt.Printf("[ANNOTATED AST]\n")
	debug.PrintAst(root)
}
