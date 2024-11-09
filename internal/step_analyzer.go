package internal

import (
	"fmt"

	"github.com/renatopp/golden/lang"
)

// Type check the module
func NewAnalyzer(module *Module) *Analyzer {
	return &Analyzer{
		ErrorData:   lang.NewErrorData(),
		module:      module,
		scope:       module.Scope,
		moduleScope: module.Scope,
		scopeStack:  []*Scope{module.Scope},
	}
}

type Analyzer struct {
	*lang.ErrorData
	module      *Module
	scope       *Scope
	moduleScope *Scope
	scopeStack  []*Scope
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

func (a *Analyzer) GetValueOrError(name string) *Node {
	node := a.scope.GetValue(name)
	if node == nil {
		a.Error(lang.Loc{}, "undefined", "undefined identifier %s", name)
	}
	return node
}

func (a *Analyzer) GetTypeOrError(name string) RtType {
	tp := a.scope.GetType(name)
	if tp == nil {
		a.Error(lang.Loc{}, "undefined", "undefined identifier %s", name)
	}
	return tp
}

func (a *Analyzer) ExpectMatchingTypes(nodes ...*Node) {
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

func (a *Analyzer) ExpectTypeToBeAnyOf(base *Node, nodes ...RtType) {
	for _, node := range nodes {
		if base.Type.Accepts(node) {
			return
		}
	}
	a.Error(base.Token.Loc, "type", "expected types %s, got %s", nodes, base.Type)
}

func (a *Analyzer) pushScope(scope *Scope) {
	a.scopeStack = append(a.scopeStack, a.scope)
	a.scope = scope
}

func (a *Analyzer) popScope() *Scope {
	if len(a.scopeStack) == 1 {
		panic("no scope to pop")
	}
	a.scope = a.scopeStack[len(a.scopeStack)-1]
	a.scopeStack = a.scopeStack[:len(a.scopeStack)-1]
	return a.scope
}

func (a *Analyzer) preAnalyzeFunctions() {
	for _, node := range a.module.Ast.Types {
		a.preResolve(node)
	}
	for _, node := range a.module.Ast.Functions {
		a.preResolve(node)
	}
}

func (a *Analyzer) analyze() {
	for _, node := range a.module.Ast.Functions {
		a.resolveValue(node)
	}

	for _, node := range a.module.Ast.Variables {
		a.resolveValue(node)
	}
}

// ---------------------------------------------------------------------

func (a *Analyzer) preResolve(node *Node) *Node {
	switch ast := node.Data.(type) {
	case *AstDataDecl:
		//

	case *AstFunctionDecl:
		returnType := Void
		if ast.ReturnType != nil {
			returnType = a.resolveType(ast.ReturnType).Type
		}

		params := []RtType{}
		for _, param := range ast.Params {
			param.Type = a.resolveType(param.Type)
			params = append(params, param.Type.Type)
		}

		tp := NewFunctionType(params, returnType)
		node.WithType(tp)
		a.scope.SetValue(ast.Name, node.WithType(tp))
	}

	return node
}

// ---------------------------------------------------------------------

func (a *Analyzer) resolveType(node *Node) *Node {
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

func (a *Analyzer) resolveTypeIdentAsType(node *Node, ast *AstTypeIdent) {
	node.WithType(a.GetTypeOrError(ast.Name))
}

func (a *Analyzer) resolveFunctionType(node *Node, ast *AstFunctionType) {
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
func (a *Analyzer) resolveTypeDecl(node *Node, ast *AstDataDecl) {
	// for _, c := range ast.Constructors {
	// 	for _, f := range c.Fields {
	// 		a.resolveType(f.Type)
	// 	}
	// }

	// NewDataType("Anonymous")
}

// ---------------------------------------------------------------------

func (a *Analyzer) resolveValue(node *Node) *Node {
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

	case *AstAccess:
		a.resolveAccessValue(node, ast)

	default:
		a.Error(node.Token.Loc, "unknown", "unknown node %s", node)
	}

	return node
}

func (a *Analyzer) resolveBlock(node *Node, ast *AstBlock) {
	node.WithType(Void)
	for _, expr := range ast.Expressions {
		a.resolveValue(expr)
		node.WithType(expr.Type)
	}
}

func (a *Analyzer) resolveBool(node *Node, ast *AstBool) {
	node.WithType(Bool)
}

func (a *Analyzer) resolveInt(node *Node, ast *AstInt) {
	node.WithType(Int)
}

func (a *Analyzer) resolveFloat(node *Node, ast *AstFloat) {
	node.WithType(Float)
}

func (a *Analyzer) resolveString(node *Node, ast *AstString) {
	node.WithType(String)
}

func (a *Analyzer) resolveUnaryOp(node *Node, ast *AstUnaryOp) {
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

func (a *Analyzer) resolveBinaryOp(node *Node, ast *AstBinaryOp) {
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

func (a *Analyzer) resolveTypeIdentAsValue(node *Node, ast *AstTypeIdent) {
	node.WithType(a.GetValueOrError(ast.Name).Type)
}

func (a *Analyzer) resolveVarIdent(node *Node, ast *AstVarIdent) {
	node.WithType(a.GetValueOrError(ast.Name).Type)
}

func (a *Analyzer) resolveVariableDecl(node *Node, ast *AstVariableDecl) {
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
	node.WithType(ast.Value.Type)
	a.scope.SetValue(ast.Name, ast.Value)
}

func (a *Analyzer) resolveFunctionDecl(node *Node, ast *AstFunctionDecl) {
	var tp *FunctionType
	// If not already pre-analyzed
	if node.Type != nil {
		returnType := Void
		if ast.ReturnType != nil {
			returnType = a.resolveType(ast.ReturnType).Type
		}

		params := []RtType{}
		for _, param := range ast.Params {
			param.Type = a.resolveType(param.Type)
			params = append(params, param.Type.Type)
		}

		tp = NewFunctionType(params, returnType)
	}

	a.pushScope(a.scope.New())
	for _, param := range ast.Params {
		a.scope.SetValue(param.Name, param.Type)
	}

	a.resolveValue(ast.Body)
	a.ExpectTypeToBeAnyOf(ast.Body, tp.ret)
	a.popScope()

	node.WithType(tp)
	a.scope.SetValue(ast.Name, node)
}

func (a *Analyzer) resolveApply(node *Node, ast *AstApply) {
	if ast.Target == nil {
		a.resolveAnonymousApply(node, ast)
	} else {
		a.resolveTargetApply(node, ast)
	}
}

func (a *Analyzer) resolveAnonymousApply(node *Node, ast *AstApply) {

}

func (a *Analyzer) resolveTargetApply(node *Node, ast *AstApply) {
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

func (a *Analyzer) resolveAccessValue(node *Node, ast *AstAccess) {
	a.resolveValue(ast.Target)

	tp, ok := ast.Target.Type.(RtTypeAccessible)
	if !ok {
		a.Error(ast.Target.Token.Loc, "type", "expected type to be accessible, got %s", ast.Target.Type.Name())
	}

	val, err := tp.AccessValue(ast.Accessor)
	if err != nil {
		a.Error(node.Token.Loc, "type", err.Error())
	}

	node.WithType(val.Type)
}
