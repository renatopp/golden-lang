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

// Initialization checker
// Type checker
// Type inference
// Scope checker
// Const checker
// Const folding
// Const propagation
type Checker struct {
	scopeStack          *ds.Stack[*env.Scope]
	initializationStack *ds.Stack[ast.Node]
	root                ast.Module
}

func NewChecker(root ast.Module) *Checker {
	return &Checker{
		scopeStack:          ds.NewStack[*env.Scope](),
		initializationStack: ds.NewStack[ast.Node](),
		root:                root,
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

func (c *Checker) declare(name string, node ast.Node, tp ast.Type) {
	scope := c.scope().Values
	bind := scope.GetLocal(name, nil)
	if bind != nil && bind.IsSolved() {
		errors.ThrowAtNode(node, errors.NameAlreadyDefined, "name '%s' already defined", name)
	}

	if bind != nil {
		bind.Type = tp
	} else {
		scope.Set(name, env.VB(node, tp))
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

func (c *Checker) expectCompatibleNodeTypes(a, b ast.Node) {
	aWrappedType := a.GetType()
	bWrappedType := b.GetType()

	if !aWrappedType.Has() {
		errors.ThrowAtNode(a, errors.InternalError, "expression has 'unknown' type")
	}

	if !bWrappedType.Has() {
		errors.ThrowAtNode(b, errors.InternalError, "expression has 'unknown' type")
	}

	aType := aWrappedType.Unwrap()
	bType := bWrappedType.Unwrap()

	if !aType.IsCompatible(bType) {
		errors.ThrowAtNode(a, errors.TypeError, "expected type '%s', but got '%s'", bType.GetSignature(), aType.GetSignature())
	}
}

// Interface

func (c *Checker) Check() (res ast.Module, err error) {
	err = errors.WithRecovery(func() {
		c.preCheck()
		res = c.VisitModule(c.root).(ast.Module)
	})
	return res, err
}

func (c *Checker) preCheck() {
	tp := c.root.GetType().Unwrap().(*types.Module)
	c.pushScope(tp.Scope)
	defer c.popScope()

	for _, e := range c.root.Exprs {
		switch n := e.(type) {
		case ast.Const:
			v := n.Name.Value
			c.scope().Values.Set(v, env.VB(n, nil))
		}
	}

}

func (c *Checker) VisitModule(node ast.Module) ast.Node {
	node.Exprs = iter.Map(node.Exprs, func(e ast.Node) ast.Node { return e.Visit(c) })
	return node
}

func (c *Checker) VisitConst(node ast.Const) ast.Node {
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
	node.BaseNode = ast.SetType(node.BaseNode, tp)
	c.declare(node.Name.Value, node, tp)
	return node
}

func (c *Checker) VisitInt(node ast.Int) ast.Node {
	node.BaseNode = ast.SetType(node.BaseNode, types.Int)
	return node
}

func (c *Checker) VisitFloat(node ast.Float) ast.Node {
	node.BaseNode = ast.SetType(node.BaseNode, types.Float)
	return node
}

func (c *Checker) VisitString(node ast.String) ast.Node {
	node.BaseNode = ast.SetType(node.BaseNode, types.String)
	return node
}

func (c *Checker) VisitBool(node ast.Bool) ast.Node {
	node.BaseNode = ast.SetType(node.BaseNode, types.Bool)
	return node
}

func (c *Checker) VisitVarIdent(node ast.VarIdent) ast.Node {
	name := node.Value
	bind := c.scope().Values.Get(name, nil)
	if bind == nil {
		errors.ThrowAtNode(node, errors.NameNotFound, "variable '%s' not defined", name)
	}
	// if !bind.IsSolved() {
	// TODO: how to mutate the node outside the node structure?

	return node
}

func (c *Checker) VisitTypeIdent(node ast.TypeIdent) ast.Node {
	return node
}

func (c *Checker) VisitBinOp(node ast.BinOp) ast.Node {
	return node
}

func (c *Checker) VisitUnaryOp(node ast.UnaryOp) ast.Node {
	return node
}

func (c *Checker) VisitBlock(node ast.Block) ast.Node {
	return node
}
