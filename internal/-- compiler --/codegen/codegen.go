package codegen

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/codegen/core"
	"github.com/renatopp/golden/internal/compiler/codegen/golang"
)

var _ ast.Visitor = &Codegen{}

type Codegen struct {
	writer core.Writer
}

func NewCodegen(targetDirectory string) *Codegen {
	return &Codegen{
		writer: golang.NewGoWriter(targetDirectory),
	}
}

func (c *Codegen) StartGeneration() {
	c.writer.Start()
}

func (c *Codegen) EndGeneration() {
	c.writer.End()
}

func (c *Codegen) StartPackage(path string, imports []string) {
	c.writer.PushPackage(path, imports)
}

func (c *Codegen) EndPackage() {
	c.writer.Pop()
}

func (c *Codegen) VisitModule(a *ast.Module) {
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

func (c *Codegen) VisitImport(a *ast.Import) {}

func (c *Codegen) VisitInt(a *ast.Int) {}

func (c *Codegen) VisitFloat(a *ast.Float) {}

func (c *Codegen) VisitBool(a *ast.Bool) {}

func (c *Codegen) VisitString(a *ast.String) {}

func (c *Codegen) VisitVarIdent(a *ast.VarIdent) {}

func (c *Codegen) VisitVarDecl(a *ast.VarDecl) {}

func (c *Codegen) VisitBlock(a *ast.Block) {}

func (c *Codegen) VisitUnaryOp(a *ast.UnaryOp) {}

func (c *Codegen) VisitBinaryOp(a *ast.BinaryOp) {}

func (c *Codegen) VisitAccess(a *ast.Access) {}

func (c *Codegen) VisitTypeIdent(a *ast.TypeIdent) {}

func (c *Codegen) VisitFuncType(a *ast.FuncType) {}

func (c *Codegen) VisitFuncTypeParam(a *ast.FuncTypeParam) {}

func (c *Codegen) VisitFuncDecl(a *ast.FuncDecl) {}

func (c *Codegen) VisitFuncDeclParam(a *ast.FuncDeclParam) {}

func (c *Codegen) VisitAppl(a *ast.Appl) {}

func (c *Codegen) VisitApplArg(a *ast.ApplArg) {}
