package semantic

import (
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/env"
	"github.com/renatopp/golden/internal/compiler/types"
	"github.com/renatopp/golden/internal/helpers/ds"
	"github.com/renatopp/golden/internal/helpers/errors"
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

func (c *TypeChecker) expectTypeExpression(node ast.Node) {
	if node.ExpressionKind() != ast.TypeExpressionKind {
		errors.ThrowAtNode(node, errors.TypeError, "expected type expression")
	}
}

func (c *TypeChecker) expectValueExpression(node ast.Node) {
	if node.ExpressionKind() != ast.ValueExpressionKind {
		errors.ThrowAtNode(node, errors.TypeError, "expected value expression")
	}
}

func (c *TypeChecker) expectType(node ast.Node, types ...ast.Type) {
	for _, t := range types {
		if node.Type() == t {
			return
		}
	}

	if len(types) == 1 {
		errors.ThrowAtNode(node, errors.TypeError, "expected type '%s', but got '%s'", types[0].Signature(), node.Type().Signature())
	}

	tps := []string{}
	for _, t := range types {
		tps = append(tps, fmt.Sprintf("'%s'", t.Signature()))
	}
	names := strings.Join(tps[:len(tps)-1], ", ") + " or " + tps[len(tps)-1]
	errors.ThrowAtNode(node, errors.TypeError, "expected type %s, but got '%s'", names, node.Type().Signature())
}

func (c *TypeChecker) expectCompatible(a, b ast.Node) {
	if !a.Type().Compatible(b.Type()) {
		errors.ThrowAtNode(a, errors.TypeError, "mismatching types '%s' and '%s'", a.Type().Signature(), b.Type().Signature())
	}
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

func (c *TypeChecker) Resolve(node *ast.Module) {
	c.pushScope(node.Type().(*types.Module).Scope)
	defer c.popScope()

	for _, fn := range node.Functions {
		fn.Accept(c)
	}

	for _, v := range node.Variables {
		v.Accept(c)
	}
}

func (c *TypeChecker) VisitModule(node *ast.Module) {
	errors.Throw(errors.NotImplemented, "VisitModule not implemented.")
}

func (c *TypeChecker) VisitImport(node *ast.Import) {
	errors.Throw(errors.NotImplemented, "VisitImport not implemented.")
}

func (c *TypeChecker) VisitInt(node *ast.Int) {
	node.SetType(types.Int)
}

func (c *TypeChecker) VisitFloat(node *ast.Float) {
	node.SetType(types.Float)
}

func (c *TypeChecker) VisitString(node *ast.String) {
	node.SetType(types.String)
}

func (c *TypeChecker) VisitBool(node *ast.Bool) {
	node.SetType(types.Bool)
}

func (c *TypeChecker) VisitVarIdent(node *ast.VarIdent) {
	var name = node.Literal
	var binding = c.scope().Values.Get(name)
	if binding == nil {
		errors.ThrowAtNode(node, errors.NameNotFound, "variable '%s' not defined", name)
	}
	node.SetType(binding.Type)
}

func (c *TypeChecker) VisitVarDecl(node *ast.VarDecl) {
	var tp = node.TypeExpr.Or(nil)
	var val = node.ValueExpr.Or(nil)
	var err error

	node.TypeExpr.If(func(ast.Node) { tp.Accept(c) })

	// Get default value if type is defined and value is not
	if node.TypeExpr.Has() && !node.ValueExpr.Has() {
		val, err = tp.Type().Default()
		if err != nil {
			errors.ThrowAtNode(tp, errors.TypeError, "%s", err.Error())
		}
	}

	// Infer type from value
	val.Accept(c)
	if node.TypeExpr.Has() {
		if !val.Type().Compatible(tp.Type()) {
			errors.ThrowAtNode(node, errors.TypeError, "cannot assign type '%s' into a '%s' variable", val.Type().Signature(), tp.Type().Signature())
		}
	}
	node.SetType(types.Void)

	// Check for redeclaration
	name := node.Name.Literal
	if bind := c.scope().Values.GetLocal(name); bind != nil && bind.Type != nil {
		errors.ThrowAtNode(node, errors.NameAlreadyDefined, "variable '%s' already defined", name)
	}

	// Add to scope
	node.Name.SetType(val.Type())
	c.scope().Values.Set(name, env.BN(val.Type(), val))
}

func (c *TypeChecker) VisitBlock(node *ast.Block) {
	c.pushScope(c.scope().New())
	defer c.popScope()

	var tp ast.Type = types.Void
	for _, exp := range node.Expressions {
		exp.Accept(c)
		tp = exp.Type()
	}
	node.SetType(tp)
}

func (c *TypeChecker) VisitUnaryOp(node *ast.UnaryOp) {
	node.Right.Accept(c)

	switch node.Operator {
	case "-", "+":
		c.expectType(node.Right, types.Int, types.Float)
	case "!":
		c.expectType(node.Right, types.Bool)
	default:
		errors.ThrowAtNode(node, errors.NotImplemented, "unary operator '%s' not implemented.", node.Operator)
	}

	node.SetType(node.Right.Type())
}

func (c *TypeChecker) VisitBinaryOp(node *ast.BinaryOp) {
	node.Left.Accept(c)
	node.Right.Accept(c)

	switch node.Operator {
	case "+":
		c.expectType(node.Left, types.Int, types.Float, types.String)
		c.expectType(node.Right, types.Int, types.Float, types.String)
		c.expectCompatible(node.Left, node.Right)
		node.SetType(node.Left.Type())

	case "-", "*", "/":
		c.expectType(node.Left, types.Int, types.Float)
		c.expectType(node.Right, types.Int, types.Float)
		c.expectCompatible(node.Left, node.Right)
		node.SetType(node.Left.Type())

	case "==", "!=":
		node.SetType(types.Bool)

	case ">", "<", ">=", "<=":
		c.expectType(node.Left, types.Int, types.Float)
		c.expectType(node.Right, types.Int, types.Float)
		c.expectCompatible(node.Left, node.Right)
		node.SetType(types.Bool)

	case "<=>":
		c.expectType(node.Left, types.Int, types.Float)
		c.expectType(node.Right, types.Int, types.Float)
		c.expectCompatible(node.Left, node.Right)
		node.SetType(types.Int)

	case "and", "or", "xor":
		c.expectType(node.Left, types.Bool)
		c.expectType(node.Right, types.Bool)
		node.SetType(types.Bool)

	default:
		errors.ThrowAtNode(node, errors.NotImplemented, "binary operator '%s' not implemented.", node.Operator)
	}
}

func (c *TypeChecker) VisitAccess(node *ast.Access) {
	errors.Throw(errors.NotImplemented, "VisitAccess not implemented.")
}

func (c *TypeChecker) VisitTypeIdent(node *ast.TypeIdent) {
	if node.ExpressionKind() == ast.TypeExpressionKind {
		binding := c.scope().Types.Get(node.Literal)
		if binding == nil {
			errors.ThrowAtNode(node, errors.NameNotFound, "type '%s' not defined", node.Literal)
		}
		node.SetType(binding.Type)
		return
	}

	errors.Throw(errors.NotImplemented, "VisitTypeIdent as value not implemented.")
}

func (c *TypeChecker) VisitFuncType(node *ast.FuncType) {
	errors.Throw(errors.NotImplemented, "VisitFuncType not implemented.")
}

func (c *TypeChecker) VisitFuncTypeParam(node *ast.FuncTypeParam) {
	errors.Throw(errors.NotImplemented, "VisitFuncTypeParam not implemented.")
}

func (c *TypeChecker) VisitFuncDecl(node *ast.FuncDecl) {
	errors.Throw(errors.NotImplemented, "VisitFuncDecl not implemented.")
}

func (c *TypeChecker) VisitFuncDeclParam(node *ast.FuncDeclParam) {
	errors.Throw(errors.NotImplemented, "VisitFuncDeclParam not implemented.")
}

func (c *TypeChecker) VisitAppl(node *ast.Appl) {
	errors.Throw(errors.NotImplemented, "VisitAppl not implemented.")
}

func (c *TypeChecker) VisitApplArg(node *ast.ApplArg) {
	errors.Throw(errors.NotImplemented, "VisitApplArg not implemented.")
}
