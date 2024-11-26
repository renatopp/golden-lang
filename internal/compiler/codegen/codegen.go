package codegen

import (
	"fmt"

	"github.com/renatopp/golden/internal/compiler/ast"
)

var _ ast.Visitor = &Codegen{}

type Codegen struct {
}

func NewCodegen() *Codegen {
	return &Codegen{}
}

func (c *Codegen) StartGeneration() {}
func (c *Codegen) EndGeneration()   {}

func (c *Codegen) StartPackage() {}
func (c *Codegen) EndPackage()   {}

func (c *Codegen) VisitModule(a *ast.Module) {
	fmt.Printf("// module %s\n", a.Path)
	for _, exp := range a.Imports {
		exp.Accept(c)
	}
	for _, exp := range a.Variables {
		exp.Accept(c)
	}
	for _, exp := range a.Functions {
		exp.Accept(c)
	}
}

func (c *Codegen) VisitImport(a *ast.Import) {
	fmt.Printf("import %s\n", a.Path.Literal)
}

func (c *Codegen) VisitInt(a *ast.Int) {
	fmt.Printf("%d", a.Literal)
}

func (c *Codegen) VisitFloat(a *ast.Float) {
	fmt.Printf("%f", a.Literal)
}

func (c *Codegen) VisitBool(a *ast.Bool) {
	fmt.Printf("%t", a.Literal)
}

func (c *Codegen) VisitString(a *ast.String) {
	fmt.Printf("%s", a.Literal)
}

func (c *Codegen) VisitVarIdent(a *ast.VarIdent) {
	fmt.Printf("%s", a.Literal)
}
func (c *Codegen) VisitVarDecl(a *ast.VarDecl) {
	fmt.Printf("var %s = ", a.Name.Literal)
	a.ValueExpr.Unwrap().Accept(c)
	fmt.Printf("\n")
}

func (c *Codegen) VisitBlock(a *ast.Block)                 {}
func (c *Codegen) VisitUnaryOp(a *ast.UnaryOp)             {}
func (c *Codegen) VisitBinaryOp(a *ast.BinaryOp)           {}
func (c *Codegen) VisitAccess(a *ast.Access)               {}
func (c *Codegen) VisitTypeIdent(a *ast.TypeIdent)         {}
func (c *Codegen) VisitFuncType(a *ast.FuncType)           {}
func (c *Codegen) VisitFuncTypeParam(a *ast.FuncTypeParam) {}
func (c *Codegen) VisitFuncDecl(a *ast.FuncDecl)           {}
func (c *Codegen) VisitFuncDeclParam(a *ast.FuncDeclParam) {}
func (c *Codegen) VisitAppl(a *ast.Appl)                   {}
func (c *Codegen) VisitApplArg(a *ast.ApplArg)             {}
