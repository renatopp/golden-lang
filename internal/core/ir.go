package core

import "fmt"

type IrComp interface {
	Node() *AstNode
	Tag() string
}

type IrWriter interface {
	EnterModule(*Module)
	ExitModule()

	// Declare(string, IrComp, *AstNode) IrComp
	Int(int64, *AstNode) *Ref
	Float(float64, *AstNode) *Ref
	Bool(bool, *AstNode) *Ref
	String(string, *AstNode) *Ref
}

type Ref struct {
	Package    *Package
	Module     *Module
	Identifier string
	Counter    int
}

func (r *Ref) FullName() string {
	// return r.Package.Name + "." + r.Module.Name + "." + r.Identifier
	return r.Identifier
}

func (r *Ref) IrName() string {
	return fmt.Sprintf("%s%d", r.Identifier, r.Counter)
}
