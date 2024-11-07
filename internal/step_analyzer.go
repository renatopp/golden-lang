package internal

import (
	"fmt"

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
		}
	}()

	a.resolve(a.module.Temp)
}

func (a *analyzer) Error(loc lang.Loc, kind, msg string, args ...any) {
	panic(lang.NewError(loc, kind, fmt.Sprintf(msg, args...)))
}

func (a *analyzer) GetOrError(name string) *Node {
	node := a.scope.Get(name)
	if node == nil {
		a.Error(lang.Loc{}, "undefined", "undefined identifier %s", name)
	}
	return node
}

func (a *analyzer) GetTypeOrError(name string) RtType {
	node := a.GetOrError(name)
	if node.Type == nil {
		a.Error(lang.Loc{}, "undefined", "undefined type %s", name)
	}
	return node.Type
}

func (a *analyzer) ExpectMatchingTypes(nodes ...*Node) {
	if len(nodes) == 0 {
		return
	}
	for i := 0; i < len(nodes); i++ {
		for j := i + 1; j < len(nodes); j++ {
			// TODO: Implement Accepts method in RtType
			if nodes[i].Type.Name() != nodes[j].Type.Name() {
				// Should call type check function
				a.Error(nodes[j].Token.Loc, "type", "expected type %s, got %s", nodes[i].Type, nodes[j].Type)
			}
		}
	}
}

func (a *analyzer) ExpectTypeToBeAnyOf(base *Node, nodes ...*Node) {
	for _, node := range nodes {
		if node.Type == base.Type {
			return
		}
	}
	a.Error(base.Token.Loc, "type", "expected types %s, got %s", nodes, base.Type)
}

func (a *analyzer) resolve(node *Node) {
	switch ast := node.Data.(type) {
	case *AstBlock:
		node.WithType(Void.Type)
		for _, expr := range ast.Expressions {
			a.resolve(expr)
			node.WithType(expr.Type)
		}

	case *AstBool:
		node.WithType(Bool.Type)

	case *AstInt:
		node.WithType(Int.Type)

	case *AstFloat:
		node.WithType(Float.Type)

	case *AstString:
		node.WithType(String.Type)

	case *AstUnaryOp:
		a.resolve(ast.Right)

		switch ast.Operator {
		case "-", "+":
			a.ExpectTypeToBeAnyOf(ast.Right, Int, Float)
			node.WithType(ast.Right.Type)
		case "!":
			a.ExpectMatchingTypes(Bool, ast.Right)
			node.WithType(Bool.Type)
		default:
			a.Error(node.Token.Loc, "unknown", "unknown unary operator %s", ast.Operator)
		}

	case *AstBinaryOp:
		a.resolve(ast.Left)
		a.resolve(ast.Right)

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
			node.WithType(Bool.Type)

		case "<", "<=", ">", ">=":
			a.ExpectTypeToBeAnyOf(ast.Right, Int, Float)
			a.ExpectMatchingTypes(ast.Left, ast.Right)
			node.WithType(Bool.Type)

		case "<=>":
			a.ExpectTypeToBeAnyOf(ast.Right, Int, Float)
			a.ExpectMatchingTypes(ast.Left, ast.Right)
			node.WithType(Int.Type)

		case "and", "or", "xor":
			a.ExpectMatchingTypes(Bool, ast.Left, ast.Right)
			node.WithType(Bool.Type)

		default:
			a.Error(node.Token.Loc, "unknown", "unknown binary operator %s", ast.Operator)
		}

	case *AstVariableDecl:
		if ast.Value != nil {
			a.resolve(ast.Value)
		}

		if ast.Type != nil {
			a.resolve(ast.Type)
			a.ExpectMatchingTypes(ast.Type, ast.Value)
		}

		node.WithType(Void.Type)
		a.scope.Set(ast.Name, ast.Value)

	case *AstTypeIdent:
		node.WithType(a.GetOrError(ast.Name).Type)

	case *AstVarIdent:
		node.WithType(a.GetOrError(ast.Name).Type)

	default:
		a.Error(node.Token.Loc, "unknown", "unknown node %s", node)
	}

}
