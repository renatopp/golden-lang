package ast

import (
	"github.com/renatopp/golden/internal/helpers/iter"
	"github.com/renatopp/golden/internal/helpers/safe"
)

type Visitor interface {
	VisitModule(*Module) Node
	VisitConst(*Const) Node
	VisitInt(*Int) Node
	VisitFloat(*Float) Node
	VisitString(*String) Node
	VisitBool(*Bool) Node
	VisitVarIdent(*VarIdent) Node
	VisitTypeIdent(*TypeIdent) Node
	VisitBinOp(*BinOp) Node
	VisitUnaryOp(*UnaryOp) Node
	VisitBlock(*Block) Node
}

// Use it to replace nodes in the AST.
type Visiter struct {
}

func (v *Visiter) VisitModule(node *Module) Node {
	node.Exprs = iter.Map(node.Exprs, func(e Node) Node { return e.Visit(v) })
	return node
}
func (v *Visiter) VisitConst(node *Const) Node {
	node.Name = node.Name.Visit(v).(*VarIdent)
	node.TypeExpr = safe.Map(node.TypeExpr, func(n Node) Node { return n.Visit(v) })
	node.ValueExpr = node.ValueExpr.Visit(v)
	return node
}
func (v *Visiter) VisitInt(node *Int) Node             { return node }
func (v *Visiter) VisitFloat(node *Float) Node         { return node }
func (v *Visiter) VisitString(node *String) Node       { return node }
func (v *Visiter) VisitBool(node *Bool) Node           { return node }
func (v *Visiter) VisitVarIdent(node *VarIdent) Node   { return node }
func (v *Visiter) VisitTypeIdent(node *TypeIdent) Node { return node }
func (v *Visiter) VisitBinOp(node *BinOp) Node {
	node.LeftExpr = node.LeftExpr.Visit(v)
	node.RightExpr = node.RightExpr.Visit(v)
	return node
}
func (v *Visiter) VisitUnaryOp(node *UnaryOp) Node {
	node.RightExpr = node.RightExpr.Visit(v)
	return node
}
func (v *Visiter) VisitBlock(node *Block) Node {
	node.Exprs = iter.Map(node.Exprs, func(e Node) Node { return e.Visit(v) })
	return node
}
