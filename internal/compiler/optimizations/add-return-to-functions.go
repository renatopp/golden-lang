package optimizations

import "github.com/renatopp/golden/internal/compiler/ast"

type AddReturnToFunctions struct {
	*ast.Visiter
}

func NewAddReturnToFunctions() *AddReturnToFunctions {
	opt := &AddReturnToFunctions{}
	opt.Visiter = ast.NewVisiter(opt)
	return opt
}

func (a *AddReturnToFunctions) VisitFnDecl(node *ast.FnDecl) ast.Node {
	defer a.Visiter.VisitFnDecl(node)

	block := node.ValueExpr
	if len(block.Exprs) == 0 {
		return node
	}

	lastIdx := len(block.Exprs) - 1
	lastExpr := block.Exprs[lastIdx]
	if _, ok := lastExpr.(*ast.Return); ok {
		return node
	}

	block.Exprs[lastIdx] = &ast.Return{
		ValueExpr: lastExpr,
	}
	block.Exprs[lastIdx].SetType(lastExpr.GetType().Unwrap())

	return node
}
