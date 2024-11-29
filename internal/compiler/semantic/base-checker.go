package semantic

import (
	"fmt"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/env"
	"github.com/renatopp/golden/internal/helpers/ds"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/str"
)

type BaseChecker struct {
	scopeStack          *ds.Stack[*env.Scope]
	initializationStack *ds.Stack[ast.Node]
}

func NewBaseChecker() *BaseChecker {
	return &BaseChecker{
		scopeStack:          ds.NewStack[*env.Scope](),
		initializationStack: ds.NewStack[ast.Node](),
	}
}

// Scoping

func (c *BaseChecker) PushScope(scope *env.Scope) {
	c.scopeStack.Push(scope)
}

func (c *BaseChecker) PopScope() *env.Scope {
	return c.scopeStack.Pop(nil)
}

func (c *BaseChecker) Scope() *env.Scope {
	scope := c.scopeStack.Top(nil)
	if scope == nil {
		errors.Throw(errors.InternalError, "no scope found")
	}
	return scope
}

func (c *BaseChecker) DeclareValue(name string, node ast.Node) {
	// bind := c.Scope().Values.GetLocal(name, nil)
	// if bind != nil {
	// 	errors.ThrowAtNode(node, errors.NameAlreadyDefined, "name '%s' is already defined", name)
	// }

	// c.Scope().Values.Set(name, node)
}

// Initialization Stack

func (c *BaseChecker) PushInitialization(node ast.Node) {
	if c.initializationStack.Has(node) {
		// TODO: improve error message
		errors.ThrowAtNode(node, errors.CircularReferenceError, "circular initialization detected")
	}
	c.initializationStack.Push(node)
}

func (c *BaseChecker) PopInitialization() ast.Node {
	return c.initializationStack.Pop(nil)
}

// Checks

func (c *BaseChecker) ExpectNodeWithCompatibleType(node ast.Node, types ...ast.Type) {
	wrappedType := node.GetType()
	if !wrappedType.Has() {
		errors.ThrowAtNode(node, errors.InternalError, "expected type '%s', but got 'unknown'", types[0].Signature())
	}
	tp := wrappedType.Unwrap()

	for _, t := range types {
		if t.Compatible(tp) {
			return
		}
	}

	if len(types) == 1 {
		errors.ThrowAtNode(node, errors.TypeError, "expected type '%s', but got '%s'", types[0].Signature(), tp.Signature())
	}

	names := str.MapHumanList(types, func(t ast.Type) string {
		return fmt.Sprintf("'%s'", t.Signature())
	}, "or")
	errors.ThrowAtNode(node, errors.TypeError, "expected one of  %s, but got '%s'", names, tp.Signature())
}

func (c *BaseChecker) ExpectCompatibleNodeTypes(a, b ast.Node) {
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

	if !aType.Compatible(bType) {
		errors.ThrowAtNode(a, errors.TypeError, "expected type '%s', but got '%s'", bType.Signature(), aType.Signature())
	}
}
