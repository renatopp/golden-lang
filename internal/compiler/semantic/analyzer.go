// type checking, inference, binding checks, cyclic reference check, argument validations, control flow checks, visibility checks, mutability, etc.
package semantic

import (
	"fmt"
	"slices"
	"strings"

	"github.com/renatopp/golden/internal/compiler/semantic/types"
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/lang"
)

var Int, Float, String, Bool, Void core.TypeData

func init() {
	Int = types.NewPrimitive("Int", func() (core.AstData, error) { return &ast.Int{Value: 0}, nil })
	Float = types.NewPrimitive("Float", func() (core.AstData, error) { return &ast.Float{Value: 0}, nil })
	String = types.NewPrimitive("String", func() (core.AstData, error) { return &ast.String{Value: ""}, nil })
	Bool = types.NewPrimitive("Bool", func() (core.AstData, error) { return &ast.Bool{Value: false}, nil })
	Void = types.NewVoid()
}

type Analyzer struct {
	*lang.ErrorData
	module          *core.Module
	scope           *core.Scope
	moduleScope     *core.Scope
	scopeStack      []*core.Scope
	resolutionStack []*core.AstNode
}

func NewAnalyzer(module *core.Module) *Analyzer {
	return &Analyzer{
		ErrorData:       lang.NewErrorData(),
		module:          module,
		scope:           module.Scope,
		moduleScope:     module.Scope,
		scopeStack:      []*core.Scope{module.Scope},
		resolutionStack: []*core.AstNode{},
	}
}

// Pre-analyze the module, adding types and function signatures to the module
// scope, so they can be used later in the analysis
func (a *Analyzer) PreAnalyzeTypes() error {
	a.WithRecovery(a.preAnalyzeFunctions)

	if a.HasErrors() {
		return lang.NewErrorList(a.Errors())
	}
	return nil
}

func (a *Analyzer) PreAnalyzeFunctions() error {
	a.WithRecovery(a.preAnalyzeFunctions)

	if a.HasErrors() {
		return lang.NewErrorList(a.Errors())
	}
	return nil
}

func (a *Analyzer) PreAnalyzeVariables() error {
	a.WithRecovery(a.preAnalyzeVariables)

	if a.HasErrors() {
		return lang.NewErrorList(a.Errors())
	}
	return nil
}

func (a *Analyzer) Analyze() error {
	a.WithRecovery(a.analyze)

	if a.HasErrors() {
		return lang.NewErrorList(a.Errors())
	}
	return nil
}

func (a *Analyzer) Error(loc lang.Loc, kind, msg string, args ...any) {
	panic(lang.NewError(loc, kind, fmt.Sprintf(msg, args...)))
}

func (a *Analyzer) GetValueOrError(name string) *core.AstNode {
	node := a.scope.GetValue(name)
	if node == nil {
		a.Error(lang.Loc{}, "undefined", "undefined identifier %s", name)
	}
	return node
}

func (a *Analyzer) GetTypeOrError(name string) core.TypeData {
	tp := a.scope.GetType(name)
	if tp == nil {
		a.Error(lang.Loc{}, "undefined", "undefined type identifier %s", name)
	}
	return tp
}

func (a *Analyzer) ExpectMatchingTypes(nodes ...*core.AstNode) {
	if len(nodes) == 0 {
		return
	}
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			loc := lang.Loc{}
			if nodes[j].Token() == nil {
				loc = nodes[i].Token().Loc
			}

			if nodes[i].Type() == nil && nodes[j].Type() == nil {
				continue
			}

			if nodes[i].Type() == nil {
				a.Error(loc, "type", "expected type %s, got %s", nodes[j].Type().Signature(), "Void")
			}

			if nodes[j].Type() == nil {
				a.Error(loc, "type", "expected type %s, got %s", "Void", nodes[i].Type().Signature())
			}

			if !nodes[i].Type().Accepts(nodes[j].Type()) {
				a.Error(loc, "type", "expected type %s, got %s", nodes[i].Type().Signature(), nodes[j].Type().Signature())
			}
		}
	}
}

func (a *Analyzer) ExpectTypeToBeAnyOf(base *core.AstNode, nodes ...core.TypeData) {
	for _, node := range nodes {
		if node.Accepts(base.Type()) {
			return
		}
	}

	names := []string{}
	for _, node := range nodes {
		names = append(names, node.Signature())
	}
	if base.Type() == nil {
		a.Error(base.Token().Loc, "type", "expected types %s, got Void", strings.Join(names, ", "))
	}
	a.Error(base.Token().Loc, "type", "expected types %s, got %s", strings.Join(names, ", "), base.Type().Signature())
}

func (a *Analyzer) pushScope(scope *core.Scope) {
	a.scopeStack = append(a.scopeStack, a.scope)
	a.scope = scope
}

func (a *Analyzer) popScope() *core.Scope {
	if len(a.scopeStack) == 1 {
		panic("no scope to pop")
	}
	a.scope = a.scopeStack[len(a.scopeStack)-1]
	a.scopeStack = a.scopeStack[:len(a.scopeStack)-1]
	return a.scope
}

func (a *Analyzer) preAnalyzeTypes() {
	//for _, node := range a.module.Ast.Types {
	//a.preResolve(node)
	//}
}

func (a *Analyzer) preAnalyzeFunctions() {
	for _, node := range (a.module.Node.Data().(*ast.Module)).Functions {
		ast := node.Data().(*ast.FunctionDecl)

		returnType := Void
		if ast.ReturnType != nil {
			returnType = a.resolveType(ast.ReturnType).Type()
		}

		params := []core.TypeData{}
		for _, param := range ast.Params {
			param.Type = a.resolveType(param.Type)
			params = append(params, param.Type.Type())
		}

		tp := types.NewFunction(params, returnType)
		node.WithType(tp)
		a.scope.SetValue(ast.Name, node.WithType(tp))
	}
}

func (a *Analyzer) preAnalyzeVariables() {
	for _, node := range a.module.Node.Data().(*ast.Module).Variables {
		ast := node.Data().(*ast.VariableDecl)

		a.scope.SetValue(ast.Name, ast.Value)
	}
}

func (a *Analyzer) analyze() {
	for _, node := range a.module.Node.Data().(*ast.Module).Functions {
		a.resolveValue(node)
	}

	for _, node := range a.module.Node.Data().(*ast.Module).Variables {
		a.resolveValue(node)
	}
}

// ---------------------------------------------------------------------

func (a *Analyzer) resolveType(node *core.AstNode) *core.AstNode {
	switch ast := node.Data().(type) {
	case *ast.TypeIdent:
		a.resolveTypeIdentAsType(node, ast)

	case *ast.FunctionType:
		a.resolveFunctionType(node, ast)

	default:
		a.Error(node.Token().Loc, "unknown", "unknown node %v", node)
	}

	return node
}

func (a *Analyzer) resolveTypeIdentAsType(node *core.AstNode, ast *ast.TypeIdent) {
	node.WithType(a.GetTypeOrError(ast.Name))
}

func (a *Analyzer) resolveFunctionType(node *core.AstNode, ast *ast.FunctionType) {
	params := []core.TypeData{}
	for _, param := range ast.Params {
		params = append(params, a.resolveType(param.Type).Type())
	}

	ret := Void
	if ast.ReturnType != nil {
		ret = a.resolveType(ast.ReturnType).Type()
	}

	node.WithType(&types.Function{Parameters: params, Return: ret})
}

// Anonymous types
func (a *Analyzer) resolveTypeDecl(node *core.AstNode, ast *ast.DataDecl) {
	// for _, c := range ast.Constructors {
	// 	for _, f := range c.Fields {
	// 		a.resolveType(f.Type)
	// 	}
	// }

	// NewDataType("Anonymous")
}

// ---------------------------------------------------------------------

func (a *Analyzer) ResolveValue(node *core.AstNode) *core.AstNode { return a.resolveValue(node) } // TODO: remove

func (a *Analyzer) resolveValue(node *core.AstNode) *core.AstNode {
	if slices.Contains(a.resolutionStack, node) {
		a.Error(node.Token().Loc, "circular", "circular reference detected")
	}
	a.resolutionStack = append(a.resolutionStack, node)
	defer func() {
		a.resolutionStack = a.resolutionStack[:len(a.resolutionStack)-1]
	}()

	switch ast := node.Data().(type) {
	case *ast.Block:
		a.resolveBlock(node, ast)

	case *ast.Bool:
		a.resolveBool(node, ast)

	case *ast.Int:
		a.resolveInt(node, ast)

	case *ast.Float:
		a.resolveFloat(node, ast)

	case *ast.String:
		a.resolveString(node, ast)

	case *ast.UnaryOp:
		a.resolveUnaryOp(node, ast)

	case *ast.BinaryOp:
		a.resolveBinaryOp(node, ast)

	case *ast.TypeIdent:
		a.resolveTypeIdentAsValue(node, ast)

	case *ast.VarIdent:
		a.resolveVarIdent(node, ast)

	case *ast.VariableDecl:
		a.resolveVariableDecl(node, ast)

	case *ast.FunctionDecl:
		a.resolveFunctionDecl(node, ast)

	case *ast.Apply:
		a.resolveApply(node, ast)

	case *ast.Access:
		a.resolveAccessValue(node, ast)

	default:
		a.Error(node.Token().Loc, "unknown", "unknown node %v", node)
	}

	return node
}

func (a *Analyzer) resolveBlock(node *core.AstNode, ast *ast.Block) {
	node.WithType(Void)
	for _, expr := range ast.Expressions {
		a.resolveValue(expr)
		node.WithType(expr.Type())
	}
}

func (a *Analyzer) resolveBool(node *core.AstNode, ast *ast.Bool) {
	node.WithType(Bool)
}

func (a *Analyzer) resolveInt(node *core.AstNode, ast *ast.Int) {
	node.WithType(Int)
}

func (a *Analyzer) resolveFloat(node *core.AstNode, ast *ast.Float) {
	node.WithType(Float)
}

func (a *Analyzer) resolveString(node *core.AstNode, ast *ast.String) {
	node.WithType(String)
}

func (a *Analyzer) resolveUnaryOp(node *core.AstNode, ast *ast.UnaryOp) {
	a.resolveValue(ast.Right)

	switch ast.Operator {
	case "-", "+":
		a.ExpectTypeToBeAnyOf(ast.Right, Int, Float)
		node.WithType(ast.Right.Type())
	case "!":
		a.ExpectTypeToBeAnyOf(ast.Right, Bool)
		node.WithType(Bool)
	default:
		a.Error(node.Token().Loc, "unknown", "unknown unary operator %s", ast.Operator)
	}
}

func (a *Analyzer) resolveBinaryOp(node *core.AstNode, ast *ast.BinaryOp) {
	a.resolveValue(ast.Left)
	a.resolveValue(ast.Right)

	switch ast.Operator {
	case "+":
		a.ExpectTypeToBeAnyOf(ast.Right, Int, Float, String)
		a.ExpectMatchingTypes(ast.Left, ast.Right)
		node.WithType(ast.Left.Type())

	case "-", "*", "/", "%":
		a.ExpectTypeToBeAnyOf(ast.Right, Int, Float)
		a.ExpectMatchingTypes(ast.Left, ast.Right)
		node.WithType(ast.Left.Type())

	case "==", "!=":
		node.WithType(Bool)

	case "<", "<=", ">", ">=":
		a.ExpectTypeToBeAnyOf(ast.Right, Int, Float)
		a.ExpectMatchingTypes(ast.Left, ast.Right)
		node.WithType(Bool)

	case "<=>":
		a.ExpectTypeToBeAnyOf(ast.Right, Int, Float)
		a.ExpectMatchingTypes(ast.Left, ast.Right)
		node.WithType(Int)

	case "and", "or", "xor":
		a.ExpectMatchingTypes(ast.Left, ast.Right)
		a.ExpectTypeToBeAnyOf(ast.Left, Bool)
		node.WithType(Bool)

	default:
		a.Error(node.Token().Loc, "unknown", "unknown binary operator %s", ast.Operator)
	}
}

func (a *Analyzer) resolveTypeIdentAsValue(node *core.AstNode, ast *ast.TypeIdent) {
	node.WithType(a.GetValueOrError(ast.Name).Type())
}

func (a *Analyzer) resolveVarIdent(node *core.AstNode, ast *ast.VarIdent) {
	val := a.GetValueOrError(ast.Name)

	if val.Type() == nil {
		a.resolveValue(val)
	}

	node.WithType(val.Type())
}

func (a *Analyzer) resolveVariableDecl(node *core.AstNode, ast *ast.VariableDecl) {
	if ast.Type != nil {
		a.resolveType(ast.Type)
	}

	if ast.Value == nil {
		def, err := ast.Type.Type().Default()
		if err != nil {
			a.Error(node.Token().Loc, "type", err.Error())
		}
		ast.Value = node.Copy().WithData(def)
	}

	a.resolveValue(ast.Value)

	if ast.Type != nil {
		a.ExpectMatchingTypes(ast.Type, ast.Value)
	}
	node.WithType(ast.Value.Type())
	a.scope.SetValue(ast.Name, ast.Value)
}

func (a *Analyzer) resolveFunctionDecl(node *core.AstNode, ast *ast.FunctionDecl) {
	var tp *types.Function
	// If not already pre-analyzed
	if node.Type() == nil {
		returnType := Void
		if ast.ReturnType != nil {
			returnType = a.resolveType(ast.ReturnType).Type()
		}

		params := []core.TypeData{}
		for _, param := range ast.Params {
			param.Type = a.resolveType(param.Type)
			params = append(params, param.Type.Type())
		}

		tp = types.NewFunction(params, returnType)
	} else {
		tp = node.Type().(*types.Function)
	}

	a.pushScope(a.scope.New())
	for _, param := range ast.Params {
		a.scope.SetValue(param.Name, param.Type)
	}

	a.resolveValue(ast.Body)
	a.ExpectTypeToBeAnyOf(ast.Body, tp.Return)
	a.popScope()

	node.WithType(tp)
	a.scope.SetValue(ast.Name, node)
}

func (a *Analyzer) resolveApply(node *core.AstNode, ast *ast.Apply) {
	if ast.Target == nil {
		a.resolveAnonymousApply(node, ast)
	} else {
		a.resolveTargetApply(node, ast)
	}
}

func (a *Analyzer) resolveAnonymousApply(node *core.AstNode, ast *ast.Apply) {
}

func (a *Analyzer) resolveTargetApply(node *core.AstNode, ast *ast.Apply) {
	a.resolveValue(ast.Target)

	tp, ok := ast.Target.Type().(core.ApplicableTypeData)
	if !ok {
		a.Error(ast.Target.Token().Loc, "type", "expected type to be applicable, got %s", ast.Target.Type().Signature())
	}

	args := []core.TypeData{}
	for _, arg := range ast.Args {
		a.resolveValue(arg.Value)
		args = append(args, arg.Value.Type())
	}

	ret, err := tp.Apply(args)
	if err != nil {
		a.Error(node.Token().Loc, "type", err.Error())
	}

	node.WithType(ret)
}

func (a *Analyzer) resolveAccessValue(node *core.AstNode, ast *ast.Access) {
	a.resolveValue(ast.Target)

	tp, ok := ast.Target.Type().(core.AccessibleTypeData)
	if !ok {
		a.Error(ast.Target.Token().Loc, "type", "expected type to be accessible, got %s", ast.Target.Type().Signature())
	}

	val, err := tp.AccessValue(ast.Accessor)
	if err != nil {
		a.Error(node.Token().Loc, "type", err.Error())
	}

	if val.Type() == nil {
		a.resolveValue(val)
	}

	node.WithType(val.Type())
}
