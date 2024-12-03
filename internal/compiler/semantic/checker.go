package semantic

import (
	"fmt"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/env"
	"github.com/renatopp/golden/internal/compiler/types"
	"github.com/renatopp/golden/internal/helpers/ds"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/iter"
	"github.com/renatopp/golden/internal/helpers/safe"
	"github.com/renatopp/golden/internal/helpers/str"
)

var _ ast.Visitor = &Checker{}

// Checker is the main semantic analysis component. It performs the following procedures:
//
// [x] Circular Initialization checker
// [x] Type checker
// [x] Type inference
// [x] Scope checker
// [ ] Const folding
// [ ] Const propagation
type Checker struct {
	scopeStack          *ds.Stack[*env.Scope]
	initializationStack *ds.Stack[ast.Node]
}

func NewChecker() *Checker {
	return &Checker{
		scopeStack:          ds.NewStack[*env.Scope](),
		initializationStack: ds.NewStack[ast.Node](),
	}
}

// Scoping

func (c *Checker) pushScope(scope *env.Scope) {
	c.scopeStack.Push(scope)
}

func (c *Checker) popScope() *env.Scope {
	return c.scopeStack.Pop(nil)
}

func (c *Checker) scope() *env.Scope {
	scope := c.scopeStack.Top(nil)
	if scope == nil {
		errors.Throw(errors.InternalError, "no scope found")
	}
	return scope
}

func (c *Checker) declare(name ast.Node, node ast.Node, tp ast.Type) {
	scope := c.scope().Values
	lit := name.GetToken().Literal
	bind := scope.GetLocal(lit, nil)
	if bind != nil && bind.IsSolved() {
		errors.ThrowAtNode(name, errors.NameAlreadyDefined, "name '%s' already defined", lit)
	}

	if bind != nil {
		bind.Type = tp
	} else {
		scope.Set(lit, env.VB(node, tp))
	}
}

// Initialization Stack

func (c *Checker) pushInitialization(node ast.Node) {
	for _, e := range c.initializationStack.Iter() {
		if e.IsEqual(node) {
			// TODO: improve error message
			errors.ThrowAtNode(node, errors.CircularReferenceError, "circular initialization detected")
		}
	}
	c.initializationStack.Push(node)
}

func (c *Checker) popInitialization() ast.Node {
	return c.initializationStack.Pop(nil)
}

// Checks

func (c *Checker) expectNodeWithCompatibleType(node ast.Node, types ...ast.Type) {
	wrappedType := node.GetType()
	if !wrappedType.Has() {
		errors.ThrowAtNode(node, errors.InternalError, "expected type '%s', but got 'unknown'", types[0].GetSignature())
	}
	tp := wrappedType.Unwrap()

	for _, t := range types {
		if t.IsCompatible(tp) {
			return
		}
	}

	if len(types) == 1 {
		errors.ThrowAtNode(node, errors.TypeError, "expected type '%s', but got '%s'", types[0].GetSignature(), tp.GetSignature())
	}

	names := str.MapHumanList(types, func(t ast.Type) string {
		return fmt.Sprintf("'%s'", t.GetSignature())
	}, "or")
	errors.ThrowAtNode(node, errors.TypeError, "expected one of  %s, but got '%s'", names, tp.GetSignature())
}

func (c *Checker) expectCompatibleNodeTypes(receiver, giver ast.Node) {
	aWrappedType := receiver.GetType()
	bWrappedType := giver.GetType()

	if !aWrappedType.Has() {
		errors.ThrowAtNode(receiver, errors.InternalError, "expression has 'unknown' type")
	}

	if !bWrappedType.Has() {
		errors.ThrowAtNode(giver, errors.InternalError, "expression has 'unknown' type")
	}

	receiverType := aWrappedType.Unwrap()
	giverType := bWrappedType.Unwrap()

	if !receiverType.IsCompatible(giverType) {
		errors.ThrowAtNode(receiver, errors.TypeError, "expected type '%s', but got '%s'", receiverType.GetSignature(), giverType.GetSignature())
	}
}

// Interface

func (c *Checker) PreCheck(root *ast.Module) {
	tp := root.GetType().Unwrap().(*types.Module)
	c.pushScope(tp.Scope)
	defer c.popScope()

	for _, e := range root.Exprs {
		switch n := e.(type) {
		case *ast.VarDecl:
			v := n.Name.Value
			c.scope().Values.Set(v, env.VB(n, nil))
		}
	}
}

func (c *Checker) Check(root *ast.Module) (res *ast.Module, err error) {
	err = errors.WithRecovery(func() {
		res = c.VisitModule(root).(*ast.Module)
	})
	return res, err
}

func (c *Checker) VisitModule(node *ast.Module) ast.Node {
	c.pushScope(node.Type.Unwrap().(*types.Module).Scope)
	defer c.popScope()
	node.Exprs = iter.Map(node.Exprs, func(e ast.Node) ast.Node { return e.Visit(c) })
	return node
}

func (c *Checker) VisitVarDecl(node *ast.VarDecl) ast.Node {
	if node.Type.Has() {
		return node
	}
	c.pushInitialization(node)
	defer c.popInitialization()

	node.TypeExpr = safe.Map(node.TypeExpr, func(e ast.Node) ast.Node { return e.Visit(c) })
	node.ValueExpr = node.ValueExpr.Visit(c)
	if node.TypeExpr.Has() {
		c.expectCompatibleNodeTypes(node.TypeExpr.Unwrap(), node.ValueExpr)
	}

	tp := node.ValueExpr.GetType().Unwrap()
	node.SetType(tp)
	node.Name.SetType(tp)
	c.declare(node.Name, node, tp)
	return node
}

func (c *Checker) VisitInt(node *ast.Int) ast.Node {
	node.SetType(types.Int)
	return node
}

func (c *Checker) VisitFloat(node *ast.Float) ast.Node {
	node.SetType(types.Float)
	return node
}

func (c *Checker) VisitString(node *ast.String) ast.Node {
	node.SetType(types.String)
	return node
}

func (c *Checker) VisitBool(node *ast.Bool) ast.Node {
	node.SetType(types.Bool)
	return node
}

func (c *Checker) VisitVarIdent(node *ast.VarIdent) ast.Node {
	name := node.Value
	bind := c.scope().Values.Get(name, nil)
	if bind == nil {
		errors.ThrowAtNode(node, errors.NameNotFound, "variable '%s' not defined", name)
	}
	if !bind.IsSolved() {
		bind.LastNode.Visit(c)
		bind.Type = bind.LastNode.GetType().Unwrap()
	}
	node.SetType(bind.Type)

	return node
}

func (c *Checker) VisitBinOp(node *ast.BinOp) ast.Node {
	node.LeftExpr.Visit(c)
	node.RightExpr.Visit(c)

	switch node.Op {
	case "+":
		c.expectCompatibleNodeTypes(node.LeftExpr, node.RightExpr)
		node.SetType(node.LeftExpr.GetType().Unwrap())

	case "-", "*", "/":
		c.expectNodeWithCompatibleType(node.LeftExpr, types.Int, types.Float)
		c.expectNodeWithCompatibleType(node.RightExpr, types.Int, types.Float)
		c.expectCompatibleNodeTypes(node.LeftExpr, node.RightExpr)
		node.SetType(node.LeftExpr.GetType().Unwrap())

	case "==", "!=":
		node.SetType(types.Bool)

	case ">", "<", ">=", "<=":
		c.expectNodeWithCompatibleType(node.LeftExpr, types.Int, types.Float)
		c.expectNodeWithCompatibleType(node.RightExpr, types.Int, types.Float)
		c.expectCompatibleNodeTypes(node.LeftExpr, node.RightExpr)
		node.SetType(types.Bool)

	case "<=>":
		c.expectNodeWithCompatibleType(node.LeftExpr, types.Int, types.Float)
		c.expectNodeWithCompatibleType(node.RightExpr, types.Int, types.Float)
		c.expectCompatibleNodeTypes(node.LeftExpr, node.RightExpr)
		node.SetType(types.Int)

	case "and", "or", "xor":
		c.expectNodeWithCompatibleType(node.LeftExpr, types.Bool)
		c.expectNodeWithCompatibleType(node.RightExpr, types.Bool)
		node.SetType(types.Bool)

	default:
		errors.ThrowAtNode(node, errors.NotImplemented, "binary operator '%s' not implemented.", node.Op)
	}

	return node
}

func (c *Checker) VisitUnaryOp(node *ast.UnaryOp) ast.Node {
	node.RightExpr.Visit(c)

	switch node.Op {
	case "-", "+":
		c.expectNodeWithCompatibleType(node.RightExpr, types.Int, types.Float)
	case "!":
		c.expectNodeWithCompatibleType(node.RightExpr, types.Bool)
	default:
		errors.ThrowAtNode(node, errors.NotImplemented, "unary operator '%s' not implemented.", node.Op)
	}

	node.SetType(node.RightExpr.GetType().Unwrap())
	return node
}

func (c *Checker) VisitBlock(node *ast.Block) ast.Node {
	c.pushScope(c.scope().New())
	defer c.popScope()

	var tp ast.Type = types.Void
	for _, exp := range node.Exprs {
		exp.Visit(c)
		tp = exp.GetType().Unwrap()
	}
	node.SetType(tp)

	return node
}

func (c *Checker) VisitFnDecl(node *ast.FnDecl) ast.Node {
	if node.Type.Has() {
		return node
	}

	c.pushInitialization(node)
	defer c.popInitialization()

	node.TypeExpr = node.TypeExpr.Visit(c)
	node.Parameters = iter.Map(node.Parameters, func(p *ast.FnDeclParam) *ast.FnDeclParam { return p.Visit(c).(*ast.FnDeclParam) })

	c.pushScope(c.scope().New())
	tps := []ast.Type{}
	for _, p := range node.Parameters {
		tp := p.Type.Unwrap()
		tps = append(tps, tp)
		c.declare(p.Name, p, tp)
	}
	node.ValueExpr = node.ValueExpr.Visit(c)
	c.popScope()

	c.expectCompatibleNodeTypes(node.TypeExpr, node.ValueExpr)

	fnType := types.NewFunction(node, tps, node.TypeExpr.GetType().Unwrap())
	node.SetType(fnType)

	if node.Name.Has() {
		name := node.Name.Unwrap()
		name.SetType(fnType)
		c.declare(name, node, fnType)
	}

	return node
}

func (c *Checker) VisitFnDeclParam(node *ast.FnDeclParam) ast.Node {
	node.TypeExpr = node.TypeExpr.Visit(c)
	tp := node.TypeExpr.GetType().Unwrap()
	node.Name.SetType(tp)
	node.SetType(tp)
	return node
}

func (c *Checker) VisitTypeIdent(node *ast.TypeIdent) ast.Node {
	name := node.Value
	bind := c.scope().Types.Get(name, nil)
	if bind == nil {
		errors.ThrowAtNode(node, errors.NameNotFound, "type '%s' not defined", name)
	}
	if !bind.IsSolved() {
		bind.DefinitionNode.Visit(c)
		bind.Type = bind.DefinitionNode.GetType().Unwrap()
	}
	node.SetType(bind.Type)
	return node
}

func (c *Checker) VisitTypeFn(node *ast.TypeFn) ast.Node {
	node.Parameters = iter.Map(node.Parameters, func(p ast.Node) ast.Node { return p.Visit(c) })
	node.ReturnExpr = node.ReturnExpr.Visit(c)

	tps := []ast.Type{}
	for _, p := range node.Parameters {
		tps = append(tps, p.GetType().Unwrap())
	}

	node.SetType(types.NewFunction(node, tps, node.ReturnExpr.GetType().Unwrap()))
	return node
}
