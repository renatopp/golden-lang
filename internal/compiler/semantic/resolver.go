package semantic

import (
	"github.com/renatopp/golden/internal/compiler/semantic/types"
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/lang"
)

type Resolver struct {
	*lang.ErrorData
	module              *core.Module
	ast                 *ast.Module
	scope               *core.Scope     // Current scope
	scopeStack          []*core.Scope   // Scope analysis
	initializationStack []*core.AstNode // Cyclic initialization detection
}

func NewResolver(module *core.Module) *Resolver {
	return &Resolver{
		module:              module,
		ast:                 module.Node.Data().(*ast.Module),
		scope:               module.Scope,
		scopeStack:          []*core.Scope{module.Scope},
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
	r.ir().EnterModule(r.module)
	defer r.ir().ExitModule()
	return errors.WithRecovery(func() {
		r.resolve(node)
	})
}

// INTERNAL HELPERS -----------------------------------------------------------

func (r *Resolver) ir() core.IrWriter {
	return r.module.Package.Ir
}

func (r *Resolver) pushScope(scope *core.Scope) {
	r.scopeStack = append(r.scopeStack, scope)
	r.scope = scope
}

func (r *Resolver) popScope() *core.Scope {
	if len(r.scopeStack) == 1 {
		// should never happen
		panic("cannot pop the root scope")
	}
	r.scope = r.scopeStack[len(r.scopeStack)-1]
	r.scopeStack = r.scopeStack[:len(r.scopeStack)-1]
	return r.scope
}

func (r *Resolver) getValueFromScope(name string, source *core.AstNode) *core.AstNode {
	binding := r.scope.Values.Get(name)
	if binding == nil {
		errors.ThrowAtNode(source, errors.UndefinedVariableError, "identifier '%s' not declared", name)
	}
	return binding.Node
}

func (r *Resolver) getTypeFromScope(name string, source *core.AstNode) core.TypeData {
	binding := r.scope.Types.Get(name)
	if binding == nil {
		errors.ThrowAtNode(source, errors.UndefinedTypeError, "type '%s' not declared", name)
	}
	return binding.Type
}

func (r *Resolver) getDefaultValue(source *core.AstNode, tp core.TypeData) *core.AstNode {
	def, err := tp.Default()
	if err != nil {
		errors.ThrowAtNode(source, errors.TypeError, err.Error())
	}
	return source.Copy().WithData(def)
}

func (r *Resolver) pushInitializationStack(node *core.AstNode) {
	for _, n := range r.initializationStack {
		if n.Id() == node.Id() {
			errors.ThrowAtNode(node, errors.CircularReferenceError, "circular initialization detected trying to resolve '%s'", node.Signature())
		}
	}
	r.initializationStack = append(r.initializationStack, node)
}

func (r *Resolver) popInitializationStack() *core.AstNode {
	node := r.initializationStack[len(r.initializationStack)-1]
	r.initializationStack = r.initializationStack[:len(r.initializationStack)-1]
	return node
}

func (r *Resolver) expectExpressionKind(node *core.AstNode, kind core.ExpressionKind) {
	if node.ExpressionKind() != kind {
		errors.ThrowAtNode(node, errors.ExpressionError, "expected %s expression but got %s expression instead", kind, node.ExpressionKind())
	}
}

func (r *Resolver) expectMatchingTypes(nodes ...*core.AstNode) {
	size := len(nodes)
	if size == 0 {
		return
	}

	// Check if all nodes have a matching type
	for i := 0; i < size; i++ {
		for j := i + 1; j < size; j++ {
			a, b := nodes[i], nodes[j]
			if a.Type() == nil && b.Type() == nil {
				continue
			}

			if a.Type() == nil {
				errors.ThrowAtNode(a, errors.TypeError, "expected type %s, got %s", b.Type().Signature(), "Void")
			}

			if b.Type() == nil {
				errors.ThrowAtNode(b, errors.TypeError, "expected type %s, got %s", a.Type().Signature(), "Void")
			}

			if !a.Type().Accepts(b.Type()) {
				errors.ThrowAtNode(b, errors.TypeError, "expected type %s, got %s", a.Type().Signature(), b.Type().Signature())
			}
		}
	}
}

func (r *Resolver) expectTypeToBeAnyOf(node *core.AstNode, types ...core.TypeData) core.TypeData {
	base := Void
	if node.Type() != nil {
		base = node.Type()
	}

	for _, tp := range types {
		// IMPORTANT: this checks which type is accepted by the base type
		if tp.Accepts(base) {
			return tp
		}
	}

	names := []string{}
	for _, tp := range types {
		names = append(names, tp.Signature())
	}
	if node.Type() == nil {
		errors.ThrowAtNode(node, errors.TypeError, "expected type Void to be any of %s", names)
	} else {
		errors.ThrowAtNode(node, errors.TypeError, "expected type '%s' to be any of %s", node.Type().Signature(), names)
	}
	return nil
}

// PRE RESOLVERS --------------------------------------------------------------

func (r *Resolver) preResolve(node *core.AstNode) *core.AstNode {
	if node == nil {
		return nil
	}
	node.WithModule(r.module)

	if node.Type() != nil {
		return node
	}

	switch ast := node.Data().(type) {
	case *ast.FunctionDecl:
		r.preResolveFunctionDecl(node, ast)

	case *ast.VariableDecl:
		r.preResolveVariableDecl(node, ast)

	}
	return node
}

func (r *Resolver) preResolveFunctionDecl(node *core.AstNode, ast *ast.FunctionDecl) {
	r.expectExpressionKind(node, core.ValueExpression)
	r.resolveFunctionSignature(node, ast)
	r.scope.Values.Set(ast.Name, core.BindValue(node))
}

func (r *Resolver) preResolveVariableDecl(node *core.AstNode, data *ast.VariableDecl) {
	r.expectExpressionKind(node, core.ValueExpression)

	switch sub := data.Value.Data().(type) {
	case *ast.FunctionDecl:
		sub.Name = data.Name
		r.preResolveFunctionDecl(data.Value, sub)
	}

	r.scope.Values.Set(data.Name, core.BindValue(data.Value))
}

// RESOLVERS ------------------------------------------------------------------
func (r *Resolver) resolve(node *core.AstNode) *core.AstNode {
	if node == nil {
		return nil
	}
	node.WithModule(r.module)

	defer func() {
		if node.Type() == nil {
			// errors.ThrowAtNode(node, errors.InternalError, "node type '%s' is nil after resolution", node.Signature())
		}
	}()

	switch node.ExpressionKind() {
	case core.TypeExpression:
		r.resolveTypeExpression(node)
	case core.ValueExpression:
		r.resolveValueExpression(node)
	default:
		errors.ThrowAtNode(node, errors.InternalError, "unknown expression kind %s", node.ExpressionKind())
	}
	return node
}

func (r *Resolver) resolveTypeExpression(node *core.AstNode) *core.AstNode {
	switch ast := node.Data().(type) {
	case *ast.TypeIdent:
		r.resolveTypeIdent(node, ast)
	case *ast.FunctionType:
		r.resolveFunctionType(node, ast)
	default:
		errors.ThrowAtNode(node, errors.InternalError, "unknown node %s", node.Signature())
	}

	return node
}

func (r *Resolver) resolveValueExpression(node *core.AstNode) *core.AstNode {
	r.pushInitializationStack(node)
	defer r.popInitializationStack()

	switch ast := node.Data().(type) {
	case *ast.Module:
		r.resolveModule(node, ast)

	case *ast.Block:
		r.resolveBlock(node, ast)

	case *ast.Bool:
		r.resolveBool(node, ast)

	case *ast.Int:
		r.resolveInt(node, ast)

	case *ast.Float:
		r.resolveFloat(node, ast)

	case *ast.String:
		r.resolveString(node, ast)

	case *ast.UnaryOp:
		r.resolveUnaryOp(node, ast)

	case *ast.BinaryOp:
		r.resolveBinaryOp(node, ast)

	case *ast.TypeIdent:
		r.resolveTypeIdent(node, ast)

	case *ast.VarIdent:
		r.resolveVarIdent(node, ast)

	case *ast.VariableDecl:
		r.resolveVariableDecl(node, ast)

	case *ast.FunctionDecl:
		r.resolveFunctionDecl(node, ast)

	case *ast.Apply:
		r.resolveApply(node, ast)

	case *ast.Access:
		r.resolveAccessValue(node, ast)

	default:
		errors.ThrowAtNode(node, errors.InternalError, "unknown node %s", node.Signature())
	}

	return node
}

func (r *Resolver) resolveTypeIdent(node *core.AstNode, data *ast.TypeIdent) {
	if node.ExpressionKind() == core.TypeExpression {
		node.WithType(r.getTypeFromScope(data.Name, node))
		return
	}

	node.WithType(r.getValueFromScope(data.Name, node).Type())
}

func (r *Resolver) resolveFunctionType(node *core.AstNode, data *ast.FunctionType) {
	r.expectExpressionKind(node, core.TypeExpression)

	params := []core.TypeData{}
	for _, p := range data.Params {
		params = append(params, r.resolve(p.Type).Type())
	}

	ret := Void
	if data.ReturnType != nil {
		ret = r.resolve(data.ReturnType).Type()
	}

	node.WithType(types.NewFunction(params, ret))
}

func (r *Resolver) resolveModule(node *core.AstNode, data *ast.Module) {
	r.expectExpressionKind(node, core.ValueExpression)
	for _, child := range data.Children() {
		r.resolve(child)
	}
}

func (r *Resolver) resolveBlock(node *core.AstNode, data *ast.Block) {
	r.expectExpressionKind(node, core.ValueExpression)
	node.WithType(Void)
	for _, child := range data.Children() {
		node.WithType(r.resolve(child).Type())
	}
}

func (r *Resolver) resolveInt(node *core.AstNode, a *ast.Int) {
	r.expectExpressionKind(node, core.ValueExpression)
	node.WithType(Int)
	r.ir().NewInt(a.Value, node)
}

func (r *Resolver) resolveFloat(node *core.AstNode, a *ast.Float) {
	r.expectExpressionKind(node, core.ValueExpression)
	node.WithType(Float)
	r.ir().NewFloat(a.Value, node)
}

func (r *Resolver) resolveBool(node *core.AstNode, a *ast.Bool) {
	r.expectExpressionKind(node, core.ValueExpression)
	node.WithType(Bool)
	r.ir().NewBool(a.Value, node)
}

func (r *Resolver) resolveString(node *core.AstNode, a *ast.String) {
	r.expectExpressionKind(node, core.ValueExpression)
	node.WithType(String)
	r.ir().NewString(a.Value, node)
}

func (r *Resolver) resolveUnaryOp(node *core.AstNode, data *ast.UnaryOp) {
	r.expectExpressionKind(node, core.ValueExpression)
	r.resolve(data.Right)

	switch data.Operator {
	case "-", "+":
		r.expectTypeToBeAnyOf(data.Right, Int, Float)
	case "!":
		r.expectTypeToBeAnyOf(data.Right, Bool)
	default:
		errors.ThrowAtNode(node, errors.InternalError, "unknown unary operator %s", data.Operator)
	}
}

func (r *Resolver) resolveBinaryOp(node *core.AstNode, data *ast.BinaryOp) {
	r.expectExpressionKind(node, core.ValueExpression)
	r.resolve(data.Left)
	r.resolve(data.Right)

	switch data.Operator {
	case "+":
		r.expectTypeToBeAnyOf(data.Left, Int, Float, String)
		r.expectMatchingTypes(data.Left, data.Right)
		node.WithType(data.Left.Type())

	case "-", "*", "/":
		r.expectTypeToBeAnyOf(data.Left, Int, Float)
		r.expectMatchingTypes(data.Left, data.Right)
		node.WithType(data.Left.Type())

	case "==", "!=":
		node.WithType(Bool)

	case "<", "<=", ">", ">=":
		r.expectTypeToBeAnyOf(data.Left, Int, Float)
		r.expectMatchingTypes(data.Left, data.Right)
		node.WithType(data.Left.Type())

	case "<=>":
		r.expectTypeToBeAnyOf(data.Right, Int, Float)
		r.expectMatchingTypes(data.Left, data.Right)
		node.WithType(Int)

	case "and", "or", "xor":
		r.expectTypeToBeAnyOf(data.Left, Bool)
		r.expectMatchingTypes(data.Left, data.Right)
		node.WithType(Bool)

	default:
		errors.ThrowAtNode(node, errors.InternalError, "unknown binary operator %s", data.Operator)
	}
}

func (r *Resolver) resolveVarIdent(node *core.AstNode, data *ast.VarIdent) {
	r.expectExpressionKind(node, core.ValueExpression)
	node.WithType(r.getValueFromScope(data.Name, node).Type())
}

func (r *Resolver) resolveVariableDecl(node *core.AstNode, data *ast.VariableDecl) {
	r.expectExpressionKind(node, core.ValueExpression)
	if data.Value == nil && data.Type == nil {
		errors.ThrowAtNode(node, errors.ExpressionError, "variable declaration must have a type or a value")
	}

	r.resolve(data.Type)
	if data.Value == nil {
		data.Value = r.getDefaultValue(node, data.Type.Type())
	}

	r.resolve(data.Value)
	if data.Type != nil {
		r.expectMatchingTypes(data.Type, data.Value)
	}

	node.WithType(data.Value.Type())
	r.scope.Values.Set(data.Name, core.BindValue(node))

	// r.ir().Declare(data.Name, node)
}

func (r *Resolver) resolveFunctionDecl(node *core.AstNode, data *ast.FunctionDecl) {
	r.expectExpressionKind(node, core.ValueExpression)

	// If not already pre-analyzed
	if node.Type() == nil {
		r.resolveFunctionSignature(node, data)
	}

	tp := node.Type().(*types.Function)
	tp.Scope = r.scope.New()
	r.pushScope(tp.Scope)
	for _, param := range data.Params {
		tp.Scope.Values.Set(param.Name, core.BindValue(param.Type))
	}
	r.resolve(data.Body)
	r.expectTypeToBeAnyOf(data.Body, tp.Return)
	r.popScope()

	node.WithType(tp)
	r.scope.Values.Set(data.Name, core.BindValue(node))
}

func (r *Resolver) resolveFunctionSignature(node *core.AstNode, data *ast.FunctionDecl) {
	ret := Void
	if data.ReturnType != nil {
		ret = r.resolve(data.ReturnType).Type()
	}

	params := []core.TypeData{}
	for _, param := range data.Params {
		param.Type = r.resolve(param.Type)
		params = append(params, param.Type.Type())
	}

	node.WithType(types.NewFunction(params, ret))
}

func (r *Resolver) resolveApply(node *core.AstNode, ast *ast.Apply) {
	r.expectExpressionKind(node, core.ValueExpression)
	if ast.Target == nil {
		r.resolveAnonymousApply(node, ast)
	} else {
		r.resolveTargetApply(node, ast)
	}
}

func (r *Resolver) resolveAnonymousApply(node *core.AstNode, data *ast.Apply) {
	r.expectExpressionKind(node, core.ValueExpression)

}

func (r *Resolver) resolveTargetApply(node *core.AstNode, data *ast.Apply) {
	r.expectExpressionKind(node, core.ValueExpression)
	r.resolve(data.Target)

	tp, ok := data.Target.Type().(core.ApplicableTypeData)
	if !ok {
		errors.ThrowAtNode(data.Target, errors.TypeError, "type '%s' is not applicable", data.Target.Type().Signature())
	}

	args := []core.TypeData{}
	for _, arg := range data.Args {
		val := arg.Value
		r.resolve(val)
		args = append(args, val.Type())
	}

	ret, err := tp.Apply(args)
	if err != nil {
		errors.ThrowAtNode(node, errors.TypeError, err.Error())
	}

	node.WithType(ret)
}

func (r *Resolver) resolveAccessValue(node *core.AstNode, data *ast.Access) {
	r.expectExpressionKind(node, core.ValueExpression)
	r.resolve(data.Target)

	tp, ok := data.Target.Type().(core.AccessibleTypeData)
	if !ok {
		errors.ThrowAtNode(data.Target, errors.TypeError, "type '%s' is not accessible", data.Target.Type().Signature())
	}

	val, err := tp.AccessValue(data.Accessor)
	if err != nil {
		errors.ThrowAtNode(node, errors.TypeError, "%s", err.Error())
	}

	if val.Type() == nil {
		r.resolve(val)
	}

	node.WithType(val.Type())
}
