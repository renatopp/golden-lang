package semantic

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/env"
)

// Helper function to track if the checker is within a function and its scope.
type FunctionScope struct {
	Fn      *ast.FnDecl
	Scope   *env.Scope
	Returns []ast.Node
}

func NewFunctionScope(fn *ast.FnDecl, scope *env.Scope) *FunctionScope {
	return &FunctionScope{
		Fn:    fn,
		Scope: scope,
	}
}
