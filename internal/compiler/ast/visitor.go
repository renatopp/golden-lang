package ast

import (
	"github.com/renatopp/golden/internal/helpers/iter"
	"github.com/renatopp/golden/internal/helpers/safe"
)

type Visitor interface {
	VisitModule(Module) Node
	VisitConst(Const) Node
	VisitInt(Int) Node
	VisitFloat(Float) Node
	VisitString(String) Node
	VisitBool(Bool) Node
	VisitVarIdent(VarIdent) Node
	VisitTypeIdent(TypeIdent) Node
	VisitBinOp(BinOp) Node
	VisitUnaryOp(UnaryOp) Node
	VisitBlock(Block) Node
}

type ReplacerVisitor struct {
}

func (v *ReplacerVisitor) VisitModule(node Module) Node {
	node.Consts = iter.Map(node.Consts, func(e Const) Const { return e.Visit(v).(Const) })
	return node
}
func (v *ReplacerVisitor) VisitConst(node Const) Node {
	node.Name = node.Name.Visit(v).(VarIdent)
	if node.TypeExpr.Has() {
		node.TypeExpr = safe.Some(node.TypeExpr.Unwrap().Visit(v))
	}
	node.ValueExpr = node.ValueExpr.Visit(v)
	return node
}
func (v *ReplacerVisitor) VisitInt(node Int) Node {
	return node
}
func (v *ReplacerVisitor) VisitFloat(node Float) Node {
	return node
}
func (v *ReplacerVisitor) VisitString(node String) Node {
	return node
}
func (v *ReplacerVisitor) VisitBool(node Bool) Node {
	return node
}
func (v *ReplacerVisitor) VisitVarIdent(node VarIdent) Node {
	return node
}
func (v *ReplacerVisitor) VisitTypeIdent(node TypeIdent) Node {
	return node
}
func (v *ReplacerVisitor) VisitBinOp(node BinOp) Node {
	node.Left = node.Left.Visit(v)
	node.Right = node.Right.Visit(v)
	return node
}
func (v *ReplacerVisitor) VisitUnaryOp(node UnaryOp) Node {
	node.Right = node.Right.Visit(v)
	return node
}
func (v *ReplacerVisitor) VisitBlock(node Block) Node {
	node.Expressions = iter.Map(node.Expressions, func(e Node) Node { return e.Visit(v) })
	return node
}
