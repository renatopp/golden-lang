package build

import (
	"fmt"
	"os"
	"strings"

	"github.com/renatopp/golden/internal/compiler/syntax"
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/debug"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/fs"
	"github.com/renatopp/golden/internal/helpers/logger"
	"github.com/renatopp/golden/lang"
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

	// Read source code
	bytes, err := os.ReadFile(modulePath)
	if err != nil {
		errors.Rethrow(err)
	}

	// Extract tokens
	tokens, err := syntax.Lex(bytes)
	if err != nil {
		errors.Rethrow(err)
	}

	// Debug info!
	if s.ctx.Options.Debug && modulePath == s.ctx.EntryModulePath {
		s.debugPrintTokens(tokens)
	}

	// Annotate the tokens with the file information
	s.annotateTokens(tokens, modulePath)

	// Parse tokens into AST
	root, err := syntax.Parse(tokens)
	if err != nil {
		errors.Rethrow(err)
	}

	// Debug info!
	if s.ctx.Options.Debug && modulePath == s.ctx.EntryModulePath {
		s.debugPrintAst(root)
	}

	// Annotate the module with the package and file information
	module := s.createModule(root, modulePath)

	// Create the package reference
	s.createPackage(module)

	// Schedule imports for discovery
	s.scheduleNextImports(module)
}

func (s *StepPrepareAst) annotateTokens(tokens []*lang.Token, path string) {
	for _, token := range tokens {
		loc := token.Loc
		loc.Filename = path
		token.Loc = loc
	}
}

func (s *StepPrepareAst) createModule(root *core.AstNode, path string) *core.Module {
	module := core.NewModule()
	module.Node = root
	module.Path = path
	module.Name = fs.ModulePath_To_ModuleName(path)
	module.FileName = fs.ModulePath_To_ModuleFileName(path)

	node := root.Data().(*ast.Module)
	for _, imp := range node.Imports {
		a := imp.Data().(*ast.ModuleImport)
		module.Imports = append(module.Imports, &core.ModuleImport{
			Name:  a.Path,
			Alias: a.Alias,
			Node:  imp,
		})
	}
	return module
}

func (s *StepPrepareAst) createPackage(module *core.Module) *core.Package {
	packageName := fs.ModulePath_To_PackageName(module.Path)
	packagePath := fs.ModulePath_To_PackagePath(module.Path)
	module.Package = s.ctx.CreateOrGetPackage(packageName, packagePath)
	module.Package.Modules.Add(module)
	s.ctx.RegisterModule(module)
	return module.Package
}

func (s *StepPrepareAst) scheduleNextImports(module *core.Module) {
	for _, imp := range module.Imports {
		modulePath := fs.ImportName_To_ModulePath(imp.Name)
		packagePath := fs.ModulePath_To_PackagePath(modulePath)
		moduleName := fs.ModulePath_To_ModuleName(imp.Name)

		// check existence
		if err := fs.CheckFileExists(modulePath); err != nil {
			errors.ThrowAtNode(imp.Node, errors.InvalidFileError, "file '%s' does not exist", modulePath)
		}

		// check visibility
		if strings.HasPrefix(imp.Name, "_") && packagePath != module.Package.Path {
			errors.ThrowAtNode(imp.Node, errors.VisibilityError, "module '%s' is private and cannot be accessed from outside the package", imp.Name)
		}

		// check module name
		if !fs.IsModuleNameValid(moduleName) {
			errors.ThrowAtNode(imp.Node, errors.InvalidFileError, "invalid module name '%s'. Remember that module names must be lower snake case.", moduleName)
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

func (s *StepPrepareAst) debugPrintTokens(tokens []*lang.Token) {
	fmt.Printf("[TOKENS]\n")
	debug.PrintTokens(tokens)
}

func (s *StepPrepareAst) debugPrintAst(root *core.AstNode) {
	fmt.Printf("[AST]\n")
	debug.PrintAst(root)
}
