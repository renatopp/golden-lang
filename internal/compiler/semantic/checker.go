package semantic

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/helpers/errors"
)

var _ ast.Visitor = &Checker{}

type Checker struct {
	root ast.Module
}

func NewChecker(root ast.Module) *Checker {
	return &Checker{
		root: root,
	}
}

func (c *Checker) Check() (res ast.Module, err error) {
	err = errors.WithRecovery(func() {
		res = c.VisitModule(c.root).(ast.Module)
	})
	return res, err
}

func (c *Checker) VisitModule(node ast.Module) ast.Node {
	// consts := []ast.Const{}
	// for _, c := range c.root.Consts {
	// 	consts = append(consts, c.Visit(c))
	// }
	// return ast.Module{
	// 	Token:  c.root.Token,
	// 	Consts: consts,
	// }
	return node
}

func (c *Checker) VisitConst(node ast.Const) ast.Node {
	return node
}

func (c *Checker) VisitInt(node ast.Int) ast.Node {
	return node
}

func (c *Checker) VisitFloat(node ast.Float) ast.Node {
	return node
}

func (c *Checker) VisitString(node ast.String) ast.Node {
	return node
}

func (c *Checker) VisitBool(node ast.Bool) ast.Node {
	return node
}

func (c *Checker) VisitVarIdent(node ast.VarIdent) ast.Node {
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
