package build

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/renatopp/golden/internal/compiler/semantic"
	"github.com/renatopp/golden/internal/compiler/semantic/types"
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/logger"
)

type StepResolveBindings struct {
	ctx *Context
}

func NewStepResolveBindings(ctx *Context) *StepResolveBindings {
	return &StepResolveBindings{ctx: ctx}
}

func (s *StepResolveBindings) Process(packages []*core.Package) {
	for _, pkg := range packages {
		logger.Debug("[worker:analyze] analyzing package: %s", pkg.Path)
		mods := pkg.Modules.Values()

		// Create module scopes
		for _, module := range mods {
			module.Scope = s.ctx.GlobalScope.New()
			module.Resolver = semantic.NewResolver(module)
			module.Node.WithType(types.NewModule(module.Name, module))
		}

		// Attach scopes to each module
		for _, module := range mods {
			for _, other := range mods {
				if module == other {
					continue
				}
				module.Scope.Values.Set(other.Name, core.BindValue(other.Node))
			}
		}

		// Pre-resolve types
		for _, module := range mods {
			for _, tp := range module.Node.Data().(*ast.Module).Types {
				if err := module.Resolver.PreResolve(tp); err != nil {
					panic(err)
				}
			}
		}

		// Pre-resolve functions
		for _, module := range mods {
			for _, tp := range module.Node.Data().(*ast.Module).Functions {
				if err := module.Resolver.PreResolve(tp); err != nil {
					panic(err)
				}
			}
		}

		// Pre-resolve variables
		for _, module := range mods {
			for _, tp := range module.Node.Data().(*ast.Module).Variables {
				if err := module.Resolver.PreResolve(tp); err != nil {
					panic(err)
				}
			}
		}

		// for _, module := range mods {
		// 	println("# SCOPE OF", module.Path)
		// 	println(module.Scope.String())
		// 	println()
		// }

		// Resolve everything
		for _, module := range mods {
			if err := module.Resolver.Resolve(module.Node); err != nil {
				panic(err)
			}
		}
	}

	if s.ctx.Options.Debug {
		module, _ := s.ctx.Modules.Get(s.ctx.EntryModulePath)

		fmt.Printf("[%s:ANNOTATED AST]\n", module.Path)
		module.Node.Traverse(func(node *core.AstNode, level int) {
			ident := strings.Repeat("  ", level)
			line := "    " + ident + node.Signature()
			comment := " -- " + ident + node.Tag()

			size := utf8.RuneCountInString(line)
			if size < 50 {
				println(line, strings.Repeat(" ", 50-utf8.RuneCountInString(line)), comment)
			} else {
				println(line, comment)
			}
		})
		println()
	}

}
