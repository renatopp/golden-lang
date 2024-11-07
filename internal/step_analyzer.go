package internal

import (
	"fmt"
	"runtime/debug"

	"github.com/renatopp/golden/lang"
)

// Type check the module
func Analyze(module *Module, scope *Scope) error {
	analyzer := &analyzer{
		ErrorData: lang.NewErrorData(),
		module:    module,
		scope:     scope,
	}
	analyzer.Analyze()
	if analyzer.HasErrors() {
		return lang.NewErrorList(analyzer.Errors())
	}
	return nil
}

type analyzer struct {
	*lang.ErrorData
	module *Module
	scope  *Scope
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

func (a *analyzer) resolveType(node *Node) *Node {
	switch ast := node.Data.(type) {
	case *AstTypeIdent:
		node.WithType(a.GetTypeOrError(ast.Name))

	case *AstFunctionType:
		params := []RtType{}
		for _, param := range ast.Params {
			params = append(params, a.resolveType(param.Type).Type)
		}

		ret := Void
		if ast.ReturnType != nil {
			ret = a.resolveType(ast.ReturnType).Type
		}

		node.WithType(&FunctionType{args: params, ret: ret})

	default:
		a.Error(node.Token.Loc, "unknown", "unknown node %s", node)
	}

	return node
}

func (a *analyzer) resolveValue(node *Node) *Node {
	switch ast := node.Data.(type) {
	case *AstBlock:
		node.WithType(Void)
		for _, expr := range ast.Expressions {
			a.resolveValue(expr)
			node.WithType(expr.Type)
		}

	case *AstBool:
		node.WithType(Bool)

	case *AstInt:
		node.WithType(Int)

	case *AstFloat:
		node.WithType(Float)

	case *AstString:
		node.WithType(String)

	case *AstUnaryOp:
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

	case *AstBinaryOp:
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

	case *AstTypeIdent:
		node.WithType(a.GetValueOrError(ast.Name).Type)

	case *AstVarIdent:
		node.WithType(a.GetValueOrError(ast.Name).Type)

	case *AstVariableDecl:
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

		scope := a.scope.New()
		for _, param := range ast.Params {
			scope.SetValue(param.Name, param.Type)
		}

		a.resolveValue(ast.Body)
		a.ExpectTypeToBeAnyOf(ast.Body, returnType)

		tp := &FunctionType{args: params, ret: returnType}

		node.WithType(tp)

	case *AstApply:
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

	default:
		a.Error(node.Token.Loc, "unknown", "unknown node %s", node)
	}

	return node
}
