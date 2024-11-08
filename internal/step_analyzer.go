package internal

import (
	"fmt"
	"runtime/debug"

	"github.com/renatopp/golden/lang"
)

// Type check the module
func Analyze(module *Module, scope *Scope) error {
	analyzer := &analyzer{
		ErrorData:   lang.NewErrorData(),
		module:      module,
		scope:       scope,
		moduleScope: scope,
		scopeStack:  []*Scope{scope},
	}
	analyzer.Analyze()
	if analyzer.HasErrors() {
		return lang.NewErrorList(analyzer.Errors())
	}
	return nil
}

// func PreAnalyze(module *Module, scope *Scope) error {
// 	analyzer := &analyzer{
// 		ErrorData:   lang.NewErrorData(),
// 		module:      module,
// 		scope:       scope,
// 		moduleScope: scope,
// 		scopeStack:  []*Scope{scope},
// 	}
// 	analyzer.Analyze()
// 	if analyzer.HasErrors() {
// 		return lang.NewErrorList(analyzer.Errors())
// 	}
// 	return nil
// }

type analyzer struct {
	*lang.ErrorData
	module      *Module
	scope       *Scope
	moduleScope *Scope
	scopeStack  []*Scope
}

func (a *analyzer) PreAnalyze() {

}

func (a *analyzer) Analyze() {
	defer func() {
		r := recover()
		if r == nil {
			return
		} else if err, ok := r.(lang.Error); ok {
			a.RegisterError(err)
		} else {
			a.RegisterError(lang.NewError(lang.Loc{}, "unknown error", fmt.Sprintf("%v", r)))
			debug.PrintStack()
		}
	}()

	a.resolveValue(a.module.Temp)
}

func (a *analyzer) Error(loc lang.Loc, kind, msg string, args ...any) {
	panic(lang.NewError(loc, kind, fmt.Sprintf(msg, args...)))
}

func (a *analyzer) GetValueOrError(name string) *Node {
	node := a.scope.GetValue(name)
	if node == nil {
		a.Error(lang.Loc{}, "undefined", "undefined identifier %s", name)
	}
	return node
}

func (a *analyzer) GetTypeOrError(name string) RtType {
	tp := a.scope.GetType(name)
	if tp == nil {
		a.Error(lang.Loc{}, "undefined", "undefined identifier %s", name)
	}
	return tp
}

func (a *analyzer) ExpectMatchingTypes(nodes ...*Node) {
	if len(nodes) == 0 {
		return
	}
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			loc := lang.Loc{}
			if nodes[j].Token == nil {
				loc = nodes[i].Token.Loc
			}

			if nodes[i].Type == nil && nodes[j].Type == nil {
				continue
			}

			if nodes[i].Type == nil {
				a.Error(loc, "type", "expected type %s, got %s", nodes[j].Type.Name(), "Void")
			}

			if nodes[j].Type == nil {
				a.Error(loc, "type", "expected type %s, got %s", "Void", nodes[i].Type.Name())
			}

			if !nodes[i].Type.Accepts(nodes[j].Type) {
				a.Error(loc, "type", "expected type %s, got %s", nodes[i].Type.Name(), nodes[j].Type.Name())
			}
		}
	}
}

func (a *analyzer) ExpectTypeToBeAnyOf(base *Node, nodes ...RtType) {
	for _, node := range nodes {
		if base.Type.Accepts(node) {
			return
		}
	}
	a.Error(base.Token.Loc, "type", "expected types %s, got %s", nodes, base.Type)
}

func (a *analyzer) pushScope(scope *Scope) {
	a.scopeStack = append(a.scopeStack, a.scope)
	a.scope = scope
}

func (a *analyzer) popScope() *Scope {
	if len(a.scopeStack) == 1 {
		panic("no scope to pop")
	}
	a.scope = a.scopeStack[len(a.scopeStack)-1]
	a.scopeStack = a.scopeStack[:len(a.scopeStack)-1]
	return a.scope
}

// ---------------------------------------------------------------------

func (a *analyzer) resolveType(node *Node) *Node {
	switch ast := node.Data.(type) {
	case *AstTypeIdent:
		a.resolveTypeIdentAsType(node, ast)

	case *AstFunctionType:
		a.resolveFunctionType(node, ast)

	default:
		a.Error(node.Token.Loc, "unknown", "unknown node %s", node)
	}

	return node
}

func (a *analyzer) resolveTypeIdentAsType(node *Node, ast *AstTypeIdent) {
	node.WithType(a.GetTypeOrError(ast.Name))
}

func (a *analyzer) resolveFunctionType(node *Node, ast *AstFunctionType) {
	params := []RtType{}
	for _, param := range ast.Params {
		params = append(params, a.resolveType(param.Type).Type)
	}

	ret := Void
	if ast.ReturnType != nil {
		ret = a.resolveType(ast.ReturnType).Type
	}

	node.WithType(&FunctionType{args: params, ret: ret})
}

// Anonymous types
func (a *analyzer) resolveTypeDecl(node *Node, ast *AstDataDecl) {
	// for _, c := range ast.Constructors {
	// 	for _, f := range c.Fields {
	// 		a.resolveType(f.Type)
	// 	}
	// }

	// NewDataType("Anonymous")
}

// ---------------------------------------------------------------------

func (a *analyzer) resolveValue(node *Node) *Node {
	switch ast := node.Data.(type) {
	case *AstBlock:
		a.resolveBlock(node, ast)

	case *AstBool:
		a.resolveBool(node, ast)

	case *AstInt:
		a.resolveInt(node, ast)

	case *AstFloat:
		a.resolveFloat(node, ast)

	case *AstString:
		a.resolveString(node, ast)

	case *AstUnaryOp:
		a.resolveUnaryOp(node, ast)

	case *AstBinaryOp:
		a.resolveBinaryOp(node, ast)

	case *AstTypeIdent:
		a.resolveTypeIdentAsValue(node, ast)

	case *AstVarIdent:
		a.resolveVarIdent(node, ast)

	case *AstVariableDecl:
		a.resolveVariableDecl(node, ast)

	case *AstFunctionDecl:
		a.resolveFunctionDecl(node, ast)

	case *AstApply:
		a.resolveApply(node, ast)

	default:
		a.Error(node.Token.Loc, "unknown", "unknown node %s", node)
	}

	return node
}

func (a *analyzer) resolveBlock(node *Node, ast *AstBlock) {
	node.WithType(Void)
	for _, expr := range ast.Expressions {
		a.resolveValue(expr)
		node.WithType(expr.Type)
	}
}

func (a *analyzer) resolveBool(node *Node, ast *AstBool) {
	node.WithType(Bool)
}

func (a *analyzer) resolveInt(node *Node, ast *AstInt) {
	node.WithType(Int)
}

func (a *analyzer) resolveFloat(node *Node, ast *AstFloat) {
	node.WithType(Float)
}

func (a *analyzer) resolveString(node *Node, ast *AstString) {
	node.WithType(String)
}

func (a *analyzer) resolveUnaryOp(node *Node, ast *AstUnaryOp) {
	a.resolveValue(ast.Right)

	switch ast.Operator {
	case "-", "+":
		a.ExpectTypeToBeAnyOf(ast.Right, Int, Float)
		node.WithType(ast.Right.Type)
	case "!":
		a.ExpectTypeToBeAnyOf(ast.Right, Bool)
		node.WithType(Bool)
	default:
		a.Error(node.Token.Loc, "unknown", "unknown unary operator %s", ast.Operator)
	}
}

func (a *analyzer) resolveBinaryOp(node *Node, ast *AstBinaryOp) {
	a.resolveValue(ast.Left)
	a.resolveValue(ast.Right)

	switch ast.Operator {
	case "+":
		a.ExpectTypeToBeAnyOf(ast.Right, Int, Float, String)
		a.ExpectMatchingTypes(ast.Left, ast.Right)
		node.WithType(ast.Left.Type)

	case "-", "*", "/", "%":
		a.ExpectTypeToBeAnyOf(ast.Right, Int, Float)
		a.ExpectMatchingTypes(ast.Left, ast.Right)
		node.WithType(ast.Left.Type)

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
		a.Error(node.Token.Loc, "unknown", "unknown binary operator %s", ast.Operator)
	}
}

func (a *analyzer) resolveTypeIdentAsValue(node *Node, ast *AstTypeIdent) {
	node.WithType(a.GetValueOrError(ast.Name).Type)
}

func (a *analyzer) resolveVarIdent(node *Node, ast *AstVarIdent) {
	node.WithType(a.GetValueOrError(ast.Name).Type)
}

func (a *analyzer) resolveVariableDecl(node *Node, ast *AstVariableDecl) {
	if ast.Type != nil {
		a.resolveType(ast.Type)
	}

	if ast.Value == nil {
		ast.Value = node.Copy().WithData(ast.Type.Type.Default())
	}

	a.resolveValue(ast.Value)

	if ast.Type != nil {
		a.ExpectMatchingTypes(ast.Type, ast.Value)
	}
	node.WithType(Void)
	a.scope.SetValue(ast.Name, ast.Value)
}

func (a *analyzer) resolveFunctionDecl(node *Node, ast *AstFunctionDecl) {
	returnType := Void
	if ast.ReturnType != nil {
		returnType = a.resolveType(ast.ReturnType).Type
	}

	params := []RtType{}
	for _, param := range ast.Params {
		param.Type = a.resolveType(param.Type)
		params = append(params, param.Type.Type)
	}

	a.pushScope(a.scope.New())
	for _, param := range ast.Params {
		a.scope.SetValue(param.Name, param.Type)
	}

	a.resolveValue(ast.Body)
	a.ExpectTypeToBeAnyOf(ast.Body, returnType)
	a.popScope()

	tp := &FunctionType{args: params, ret: returnType}

	node.WithType(tp)
}

func (a *analyzer) resolveApply(node *Node, ast *AstApply) {
	if ast.Target == nil {
		a.resolveAnonymousApply(node, ast)
	} else {
		a.resolveTargetApply(node, ast)
	}
}

func (a *analyzer) resolveAnonymousApply(node *Node, ast *AstApply) {

}

func (a *analyzer) resolveTargetApply(node *Node, ast *AstApply) {
	a.resolveValue(ast.Target)

	tp, ok := ast.Target.Type.(RtTypeApplicable)
	if !ok {
		a.Error(ast.Target.Token.Loc, "type", "expected type to be applicable, got %s", ast.Target.Type.Name())
	}

	args := []RtType{}
	for _, arg := range ast.Args {
		a.resolveValue(arg.Value)
		args = append(args, arg.Value.Type)
	}

	ret, err := tp.Apply(args)
	if err != nil {
		a.Error(node.Token.Loc, "type", err.Error())
	}

	node.WithType(ret)
}
