package semantic

import "github.com/renatopp/golden/internal/compiler/ast"

var _ ast.Visitor = &Checker{}

type Checker struct{}

func (c *Checker) VisitModule(node *ast.Module)               {}
func (c *Checker) VisitImport(node *ast.Import)               {}
func (c *Checker) VisitInt(node *ast.Int)                     {}
func (c *Checker) VisitFloat(node *ast.Float)                 {}
func (c *Checker) VisitString(node *ast.String)               {}
func (c *Checker) VisitBool(node *ast.Bool)                   {}
func (c *Checker) VisitVarIdent(node *ast.VarIdent)           {}
func (c *Checker) VisitVarDecl(node *ast.VarDecl)             {}
func (c *Checker) VisitBlock(node *ast.Block)                 {}
func (c *Checker) VisitUnaryOp(node *ast.UnaryOp)             {}
func (c *Checker) VisitBinaryOp(node *ast.BinaryOp)           {}
func (c *Checker) VisitAccess(node *ast.Access)               {}
func (c *Checker) VisitTypeIdent(node *ast.TypeIdent)         {}
func (c *Checker) VisitFuncType(node *ast.FuncType)           {}
func (c *Checker) VisitFuncTypeParam(node *ast.FuncTypeParam) {}
func (c *Checker) VisitFuncDecl(node *ast.FuncDecl)           {}
func (c *Checker) VisitFuncDeclParam(node *ast.FuncDeclParam) {}
func (c *Checker) VisitAppl(node *ast.Appl)                   {}
func (c *Checker) VisitApplArg(node *ast.ApplArg)             {}
