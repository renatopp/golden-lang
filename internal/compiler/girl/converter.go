package girl

import "github.com/renatopp/golden/internal/compiler/ast"

var _ ast.Visitor = &Converter{}

type Converter struct {
}

func NewConverter() *Converter {
	return &Converter{}
}

func (w *Converter) Process(modules []*ast.Module) {
	println("hello, world")
}

func (w *Converter) VisitModule(*ast.Module)               {}
func (w *Converter) VisitImport(*ast.Import)               {}
func (w *Converter) VisitInt(*ast.Int)                     {}
func (w *Converter) VisitFloat(*ast.Float)                 {}
func (w *Converter) VisitBool(*ast.Bool)                   {}
func (w *Converter) VisitString(*ast.String)               {}
func (w *Converter) VisitVarIdent(*ast.VarIdent)           {}
func (w *Converter) VisitVarDecl(*ast.VarDecl)             {}
func (w *Converter) VisitBlock(*ast.Block)                 {}
func (w *Converter) VisitUnaryOp(*ast.UnaryOp)             {}
func (w *Converter) VisitBinaryOp(*ast.BinaryOp)           {}
func (w *Converter) VisitAccess(*ast.Access)               {}
func (w *Converter) VisitTypeIdent(*ast.TypeIdent)         {}
func (w *Converter) VisitFuncType(*ast.FuncType)           {}
func (w *Converter) VisitFuncTypeParam(*ast.FuncTypeParam) {}
func (w *Converter) VisitFuncDecl(*ast.FuncDecl)           {}
func (w *Converter) VisitFuncDeclParam(*ast.FuncDeclParam) {}
func (w *Converter) VisitAppl(*ast.Appl)                   {}
func (w *Converter) VisitApplArg(*ast.ApplArg)             {}
