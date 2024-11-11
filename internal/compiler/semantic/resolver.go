package semantic

import (
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/lang"
)

type Resolver struct {
	*lang.ErrorData
	module              *core.Module
	ast                 *ast.Module
	scopeStack          []*core.Scope   // Scope analysis
	initializationStack []*core.AstNode // Cyclic initialization detection
}

func NewResolver(module *core.Module) *Resolver {
	return &Resolver{
		module:              module,
		ast:                 module.Node.Data().(*ast.Module),
		scopeStack:          []*core.Scope{},
		initializationStack: []*core.AstNode{},
	}
}

// INTERFACE ------------------------------------------------------------------

func (r *Resolver) PreResolve(node *core.AstNode) error {
	return errors.WithRecovery(func() {
		r.preResolve(node)
	})
}

func (r *Resolver) Resolve(node *core.AstNode) error {
	return errors.WithRecovery(func() {
		r.resolve(node)
	})
}

func (r *Resolver) preResolve(node *core.AstNode) *core.AstNode {
	return node
}

func (r *Resolver) resolve(node *core.AstNode) *core.AstNode {
	return nil
}

// INTERNAL HELPERS -----------------------------------------------------------

// func (r *Resolver) pushScope(scope *core.Scope)

// func (r *Resolver) popScope() *core.Scope

// func (r *Resolver) getValueFromScope(name string) *core.AstNode

// func (r *Resolver) getTypeFromScope(name string) core.TypeData

// func (r *Resolver) pushInitializationStack(node *core.AstNode)

// func (r *Resolver) popInitializationStack() *core.AstNode

// func (r *Resolver) expectExpressionKind(node *core.AstNode, kind core.ExpressionKind)

// func (r *Resolver) expectMatchingTypes(nodes ...*core.AstNode)

// func (r *Resolver) expectTypeToBeAnyOf(node *core.AstNode, types ...core.AstData)

// RESOLVERS ------------------------------------------------------------------

// func (r *Resolver) resolveModule(node *core.AstNode, ast *ast.Module)

// func (r *Resolver) resolveBlock(node *core.AstNode, ast *ast.Block)

// func (r *Resolver) resolveBool(node *core.AstNode, ast *ast.Bool)

// func (r *Resolver) resolveInt(node *core.AstNode, ast *ast.Int)

// func (r *Resolver) resolveFloat(node *core.AstNode, ast *ast.Float)

// func (r *Resolver) resolveString(node *core.AstNode, ast *ast.String)

// func (r *Resolver) resolveUnaryOp(node *core.AstNode, ast *ast.UnaryOp)

// func (r *Resolver) resolveBinaryOp(node *core.AstNode, ast *ast.BinaryOp)

// func (r *Resolver) resolveTypeIdentAsValue(node *core.AstNode, ast *ast.TypeIdent)

// func (r *Resolver) resolveVarIdent(node *core.AstNode, ast *ast.VarIdent)

// func (r *Resolver) resolveVariableDecl(node *core.AstNode, ast *ast.VariableDecl)

// func (r *Resolver) resolveFunctionDecl(node *core.AstNode, ast *ast.FunctionDecl)

// func (r *Resolver) resolveApply(node *core.AstNode, ast *ast.Apply)

// func (r *Resolver) resolveAnonymousApply(node *core.AstNode, ast *ast.Apply)

// func (r *Resolver) resolveTargetApply(node *core.AstNode, ast *ast.Apply)

// func (r *Resolver) resolveAccessValue(node *core.AstNode, ast *ast.Access)
