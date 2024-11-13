package build

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/renatopp/golden/internal/compiler/syntax"
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/fs"
	"github.com/renatopp/golden/internal/helpers/logger"
)

type StepPrepareAst struct {
	ctx *Context
}

func NewStepPrepareAst(ctx *Context) *StepPrepareAst {
	return &StepPrepareAst{ctx: ctx}
}

func (s *StepPrepareAst) Process(modulePath string) {
	logger.Trace("[worker:prepare] preparing: %s", modulePath)
	defer s.ctx.AckModule()

	bytes, err := os.ReadFile(modulePath)
	if err != nil {
		panic(err)
	}

	tokens, err := syntax.Lex(bytes)
	if err != nil {
		panic(err)
	}

	if s.ctx.Options.Debug && modulePath == s.ctx.EntryModulePath {
		fmt.Printf("[%s:TOKENS]\n", modulePath)
		for _, t := range tokens {
			fmt.Printf("    - %s: %q\n", t.Kind, t.Literal)
		}
		println()
	}

	// Annotate the tokens with the file information
	for _, token := range tokens {
		loc := token.Loc
		loc.Filename = modulePath
		token.Loc = loc
	}

	root, err := syntax.Parse(tokens)
	if err != nil {
		errors.Rethrow(err)
	}

	if s.ctx.Options.Debug && modulePath == s.ctx.EntryModulePath {
		fmt.Printf("[%s:AST]\n", modulePath)
		root.Traverse(func(node *core.AstNode, level int) {
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

	// Annotate the module with the package and file information
	module := core.NewModule()
	module.Node = root
	// module.Ast = root.Data.(*AstModule)
	module.Path = modulePath
	module.Name = fs.ModulePath_To_ModuleName(modulePath)
	module.FileName = fs.ModulePath_To_ModuleFileName(modulePath)
	for _, imp := range module.Node.Data().(*ast.Module).Imports {
		module.Imports = append(module.Imports, &core.ModuleImport{
			Name:  imp.Path,
			Alias: imp.Alias,
		})
	}

	// Create the package reference
	packageName := fs.ModulePath_To_PackageName(modulePath)
	packagePath := fs.ModulePath_To_PackagePath(modulePath)
	pkg := s.ctx.CreateOrGetPackage(packageName, packagePath)
	module.Package = pkg
	pkg.Modules.Add(module)
	s.ctx.RegisterModule(module)

	// Schedule imports for discovery
	for _, imp := range module.Imports {
		modulePath := fs.ImportName_To_ModulePath(imp.Name)

		// TODO: check if module is private

		if err := fs.CheckFileExists(modulePath); err != nil {
			panic(fmt.Sprintf("file '%s' does not exist. Remember that module names must be lower snake case, including the extension.", modulePath))
		}

		moduleName := fs.ModulePath_To_ModuleName(imp.Name)
		if !fs.IsModuleNameValid(moduleName) {
			panic(fmt.Sprintf("invalid module name '%s'. Remember that module names must be lower snake case.", moduleName))
		}

		imp.Path = modulePath

		logger.Info("[worker:prepare] scheduling import: %s", modulePath)
		s.ctx.ScheduleDiscoverPackage(modulePath)
		// TODO: check if import is a project package or a core package
		// if package starts with @, it is a project package
		// if package matches the core packages (from a map?), it is a core package
		// otherwise, search package in the GOLDENPATH
		// if package is not found, error
	}
}
