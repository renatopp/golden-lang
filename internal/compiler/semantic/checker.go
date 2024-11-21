package semantic

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/env"
	"github.com/renatopp/golden/internal/compiler/types"
	"github.com/renatopp/golden/internal/helpers/ds"
)

var _ ast.Visitor = &TypeChecker{}

type TypeChecker struct {
	scopeStack *ds.Stack[env.Scope]
}

func NewTypeChecker() *TypeChecker {
	return &TypeChecker{
		scopeStack: ds.NewStack[env.Scope](),
	}
}

func (c *TypeChecker) pushScope(scope *env.Scope) {
	c.scopeStack.Push(scope)
}

func (c *TypeChecker) popScope() *env.Scope {
	return c.scopeStack.Pop()
}

func (c *TypeChecker) scope() *env.Scope {
	return c.scopeStack.Top()
}

func (c *TypeChecker) PreResolve(node *ast.Module) {
	c.pushScope(node.Type().(*types.Module).Scope)
	defer c.popScope()

	for _, fn := range node.Functions {
		c.scope().Values.Set(fn.Name.Unwrap().Literal, env.B(nil))
	}

	for _, v := range node.Variables {
		c.scope().Values.Set(v.Name.Literal, env.B(nil))
	}
}

func (c *TypeChecker) Resolve(node *ast.Module) {}

func (c *TypeChecker) VisitModule(node *ast.Module)               {}
func (c *TypeChecker) VisitImport(node *ast.Import)               {}
func (c *TypeChecker) VisitInt(node *ast.Int)                     {}
func (c *TypeChecker) VisitFloat(node *ast.Float)                 {}
func (c *TypeChecker) VisitString(node *ast.String)               {}
func (c *TypeChecker) VisitBool(node *ast.Bool)                   {}
func (c *TypeChecker) VisitVarIdent(node *ast.VarIdent)           {}
func (c *TypeChecker) VisitVarDecl(node *ast.VarDecl)             {}
func (c *TypeChecker) VisitBlock(node *ast.Block)                 {}
func (c *TypeChecker) VisitUnaryOp(node *ast.UnaryOp)             {}
func (c *TypeChecker) VisitBinaryOp(node *ast.BinaryOp)           {}
func (c *TypeChecker) VisitAccess(node *ast.Access)               {}
func (c *TypeChecker) VisitTypeIdent(node *ast.TypeIdent)         {}
func (c *TypeChecker) VisitFuncType(node *ast.FuncType)           {}
func (c *TypeChecker) VisitFuncTypeParam(node *ast.FuncTypeParam) {}
func (c *TypeChecker) VisitFuncDecl(node *ast.FuncDecl)           {}
func (c *TypeChecker) VisitFuncDeclParam(node *ast.FuncDeclParam) {}
func (c *TypeChecker) VisitAppl(node *ast.Appl)                   {}
func (c *TypeChecker) VisitApplArg(node *ast.ApplArg)             {}
