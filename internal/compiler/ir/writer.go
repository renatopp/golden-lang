package ir

import (
	"strings"

	"github.com/renatopp/golden/internal/compiler/ir/comp"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/ds"
	"github.com/renatopp/golden/internal/helpers/errors"
)

// Should be used per package
type GirWriter struct {
	Package     *core.Package
	Types       []*core.IrComp
	Functions   []*core.IrComp
	ScopeStack  *ds.Stack[GirScope]
	ModuleStack *ds.Stack[core.Module]
}

var _ core.IrWriter = &GirWriter{}

func NewGirWriter(pkg *core.Package) *GirWriter {
	w := &GirWriter{
		Package:     pkg,
		Types:       []*core.IrComp{},
		Functions:   []*core.IrComp{},
		ScopeStack:  ds.NewStack[GirScope](),
		ModuleStack: ds.NewStack[core.Module](),
	}

	w.ScopeStack.Push(NewGirScope())

	return w
}

func (w *GirWriter) module() *core.Module {
	top := w.ModuleStack.Top()
	if top == nil {
		errors.Throw(errors.InternalError, "module stack is empty")
	}
	return top
}

func (w *GirWriter) scope() *GirScope {
	top := w.ScopeStack.Top()
	if top == nil {
		errors.Throw(errors.InternalError, "scope stack is empty")
	}
	return top
}

func (w *GirWriter) EnterModule(module *core.Module) {
	w.ModuleStack.Push(module)
	ident := strings.Repeat("  ", w.ModuleStack.Len())
	println(ident, module.Path)
}

func (w *GirWriter) ExitModule() {
	if w.ModuleStack.Pop() == nil {
		errors.Throw(errors.InternalError, "trying to pop a module from and empty stack")
	}
}

func (w *GirWriter) Declare(identifier string, comp core.IrComp, node *core.AstNode) core.IrComp {
	// Convert name into Ref(package, module, identifier)
	// Register in scope
	// Convert name to SSA

	// ref := R(w.Package, w.module(), identifier)
	// name := ref.Name()
	// scope := w.scope()
	// ssa := scope.Inc(name)
	// w.scope().Set()

	return nil
}

func (w *GirWriter) NewInt(value int64, node *core.AstNode) core.IrComp {
	return &comp.Int{
		Base:  *comp.NewBase(node),
		Value: value,
	}
}

func (w *GirWriter) NewFloat(value float64, node *core.AstNode) core.IrComp {
	return &comp.Float{
		Base:  *comp.NewBase(node),
		Value: value,
	}
}

func (w *GirWriter) NewBool(value bool, node *core.AstNode) core.IrComp {
	return &comp.Bool{
		Base:  *comp.NewBase(node),
		Value: value,
	}
}

func (w *GirWriter) NewString(value string, node *core.AstNode) core.IrComp {
	return &comp.String{
		Base:  *comp.NewBase(node),
		Value: value,
	}
}

// func (w *GirWriter) BeginFunction(name string, node *core.AstNode) *GirFunctionWriter { return nil }
// func (w *GirWriter) EndFunction()                                                     {}

// func (w *GirWriter) OpenBlock(node *core.AstNode) {}
// func (w *GirWriter) CloseBlock()                  {}

// func (w *GirWriter) NewVarIdent(name string, node *core.AstNode)                           {}
// func (w *GirWriter) NewTypeIdent(name string, node *core.AstNode)                          {}
// func (w *GirWriter) NewInteger(value int, node *core.AstNode) {}
// func (w *GirWriter) NewDeclare(name string, value *core.IrValue, node *core.AstNode)            {}
// func (w *GirWriter) NewBinOp(op string, left *core.IrValue, right *core.IrValue, node *core.AstNode) {}
// func (w *GirWriter) NewUnOp(op string, value *core.IrValue, node *core.AstNode)                 {}
// func (w *GirWriter) NewCall(name string, args []*core.IrValue, node *core.AstNode)              {} // map??

// type GirFunctionWriter struct{}

// func (w *GirFunctionWriter) WithParam(name string, type_ *Ref, node *core.AstNode) {}
// func (w *GirFunctionWriter) WithReturn(type_ *Ref)                                 {}
// func (w *GirFunctionWriter) OpenBlock(node *core.AstNode)                          {}
// func (w *GirFunctionWriter) CloseBlock()                                           {}

// type Ref struct {
// 	Package    *core.Package
// 	Module     *core.Module
// 	Node       *core.AstNode
// 	Identifier string
// }
