package ast

import (
	"github.com/renatopp/golden/internal/helpers/iter"
	"github.com/renatopp/golden/internal/helpers/safe"
)

type Visitor interface {
	VisitModule(*Module) Node
	VisitVarDecl(*VarDecl) Node
	VisitInt(*Int) Node
	VisitFloat(*Float) Node
	VisitString(*String) Node
	VisitBool(*Bool) Node
	VisitVarIdent(*VarIdent) Node
	VisitTypeIdent(*TypeIdent) Node
	VisitBinOp(*BinOp) Node
	VisitUnaryOp(*UnaryOp) Node
	VisitBlock(*Block) Node

	VisitFnDecl(*FnDecl) Node
	VisitFnDeclParam(*FnDeclParam) Node
	VisitTypeFn(*TypeFn) Node
	VisitApplication(*Application) Node
	VisitReturn(*Return) Node
}

// Use it to replace nodes in the AST.
type Visiter struct {
	self Visitor
}

func NewVisiter(self Visitor) *Visiter {
	return &Visiter{self: self}
}

func (v *Visiter) VisitModule(node *Module) Node {
	node.Exprs = iter.Map(node.Exprs, func(e Node) Node { return e.Visit(v.self) })
	return node
}
func (v *Visiter) VisitVarDecl(node *VarDecl) Node {
	node.Name = node.Name.Visit(v.self).(*VarIdent)
	node.TypeExpr = safe.Map(node.TypeExpr, func(n Node) Node { return n.Visit(v.self) })
	node.ValueExpr = node.ValueExpr.Visit(v.self)
	return node
}
func (v *Visiter) VisitInt(node *Int) Node             { return node }
func (v *Visiter) VisitFloat(node *Float) Node         { return node }
func (v *Visiter) VisitString(node *String) Node       { return node }
func (v *Visiter) VisitBool(node *Bool) Node           { return node }
func (v *Visiter) VisitVarIdent(node *VarIdent) Node   { return node }
func (v *Visiter) VisitTypeIdent(node *TypeIdent) Node { return node }
func (v *Visiter) VisitBinOp(node *BinOp) Node {
	node.LeftExpr = node.LeftExpr.Visit(v.self)
	node.RightExpr = node.RightExpr.Visit(v.self)
	return node
}
func (v *Visiter) VisitUnaryOp(node *UnaryOp) Node {
	node.RightExpr = node.RightExpr.Visit(v.self)
	return node
}
func (v *Visiter) VisitBlock(node *Block) Node {
	node.Exprs = iter.Map(node.Exprs, func(e Node) Node { return e.Visit(v.self) })
	return node
}

func (v *Visiter) VisitFnDecl(node *FnDecl) Node {
	node.Name = safe.Map(node.Name, func(n *VarIdent) *VarIdent { return n.Visit(v.self).(*VarIdent) })
	node.Params = iter.Map(node.Params, func(n *FnDeclParam) *FnDeclParam { return n.Visit(v.self).(*FnDeclParam) })
	node.TypeExpr = node.TypeExpr.Visit(v.self)
	node.ValueExpr = node.ValueExpr.Visit(v.self).(*Block)
	return node
}
func (v *Visiter) VisitFnDeclParam(node *FnDeclParam) Node {
	node.Name = node.Name.Visit(v.self).(*VarIdent)
	node.TypeExpr = node.TypeExpr.Visit(v.self)
	return node
}
func (v *Visiter) VisitTypeFn(node *TypeFn) Node {
	node.Parameters = iter.Map(node.Parameters, func(n Node) Node { return n.Visit(v.self) })
	node.ReturnExpr = node.ReturnExpr.Visit(v.self)
	return node
}
func (v *Visiter) VisitApplication(node *Application) Node {
	node.Target = node.Target.Visit(v.self)
	node.Args = iter.Map(node.Args, func(n Node) Node { return n.Visit(v.self) })
	return node
}
func (v *Visiter) VisitReturn(node *Return) Node {
	node.ValueExpr = node.ValueExpr.Visit(v.self)
	return node
}
