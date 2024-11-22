package debug

import (
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/compiler/ast"
)

var _ ast.Visitor = &AstPrinter{}

type AstPrinter struct {
	depth int
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{
		depth: 0,
	}
}

func (p *AstPrinter) inc() { p.depth++ }
func (p *AstPrinter) dec() { p.depth-- }
func (p *AstPrinter) indent() string {
	return strings.Repeat("  ", p.depth-1)
}

func (p *AstPrinter) print(s string, args ...any) {
	fmt.Printf(p.indent()+s, args...)
}

func (p *AstPrinter) printType(node ast.Node) {
	if tp := node.Type(); tp != nil {
		fmt.Printf(" → %s", tp.Signature())
	}
	println()
}

func (p *AstPrinter) VisitModule(node *ast.Module) {
	p.inc()
	defer p.dec()

	p.print("- [module %s]", node.ModulePath)
	p.printType(node)
	for _, stmt := range node.Imports {
		stmt.Accept(p)
	}
	for _, stmt := range node.Functions {
		stmt.Accept(p)
	}
	for _, stmt := range node.Variables {
		stmt.Accept(p)
	}
}

func (p *AstPrinter) VisitImport(node *ast.Import) {
	p.inc()
	defer p.dec()

	p.print("- [import]")
	p.printType(node)
	node.Path.Accept(p)
	node.Alias.If(func(alias *ast.VarIdent) {
		alias.Accept(p)
	})
}

func (p *AstPrinter) VisitInt(node *ast.Int) {
	p.inc()
	defer p.dec()

	p.print("- [int %d]", node.Literal)
	p.printType(node)
}

func (p *AstPrinter) VisitFloat(node *ast.Float) {
	p.inc()
	defer p.dec()

	p.print("- [float %f]", node.Literal)
	p.printType(node)
}

func (p *AstPrinter) VisitString(node *ast.String) {
	p.inc()
	defer p.dec()

	p.print("- [string %s]", node.Literal)
	p.printType(node)
}

func (p *AstPrinter) VisitBool(node *ast.Bool) {
	p.inc()
	defer p.dec()

	p.print("- [bool %t]", node.Literal)
	p.printType(node)
}

func (p *AstPrinter) VisitVarIdent(node *ast.VarIdent) {
	p.inc()
	defer p.dec()

	p.print("- [var-ident %s]", node.Literal)
	p.printType(node)
}

func (p *AstPrinter) VisitVarDecl(node *ast.VarDecl) {
	p.inc()
	defer p.dec()

	p.print("- [var-decl]")
	p.printType(node)
	node.Name.Accept(p)
	node.TypeExpr.If(func(expr ast.Node) { expr.Accept(p) })
	node.ValueExpr.If(func(expr ast.Node) { expr.Accept(p) })
}

func (p *AstPrinter) VisitBlock(node *ast.Block) {
	p.inc()
	defer p.dec()

	p.print("- [block]")
	p.printType(node)
	for _, stmt := range node.Expressions {
		stmt.Accept(p)
	}
}

func (p *AstPrinter) VisitUnaryOp(node *ast.UnaryOp) {
	p.inc()
	defer p.dec()

	p.print("- [unary-op %s]", node.Operator)
	p.printType(node)
	node.Right.Accept(p)
}

func (p *AstPrinter) VisitBinaryOp(node *ast.BinaryOp) {
	p.inc()
	defer p.dec()

	p.print("- [binary-op %s]", node.Operator)
	p.printType(node)
	node.Left.Accept(p)
	node.Right.Accept(p)
}

func (p *AstPrinter) VisitAccess(node *ast.Access) {
	p.inc()
	defer p.dec()

	p.print("- [access]")
	p.printType(node)
	node.Target.Accept(p)
	node.Accessor.Accept(p)
}

func (p *AstPrinter) VisitTypeIdent(node *ast.TypeIdent) {
	p.inc()
	defer p.dec()

	p.print("- [type-ident %s]", node.Literal)
	p.printType(node)
}

func (p *AstPrinter) VisitFuncType(node *ast.FuncType) {
	p.inc()
	defer p.dec()

	p.print("- [func-type]")
	p.printType(node)
	for _, param := range node.Params {
		param.Accept(p)
	}
	node.Return.If(func(expr ast.Node) { expr.Accept(p) })
}

func (p *AstPrinter) VisitFuncTypeParam(node *ast.FuncTypeParam) {
	p.inc()
	defer p.dec()

	p.print("- [func-type-param %d]", node.Index)
	p.printType(node)
	node.TypeExpr.Accept(p)
}

func (p *AstPrinter) VisitFuncDecl(node *ast.FuncDecl) {
	p.inc()
	defer p.dec()

	p.print("- [func-decl]")
	p.printType(node)
	node.Name.If(func(name *ast.VarIdent) { name.Accept(p) })
	for _, param := range node.Params {
		param.Accept(p)
	}
	node.Return.If(func(expr ast.Node) { expr.Accept(p) })
	node.Body.Accept(p)
}

func (p *AstPrinter) VisitFuncDeclParam(node *ast.FuncDeclParam) {
	p.inc()
	defer p.dec()

	p.print("- [func-decl-param %d %s]", node.Index, node.Name.Literal)
	p.printType(node)
	node.TypeExpr.Accept(p)
}

func (p *AstPrinter) VisitAppl(node *ast.Appl) {
	p.inc()
	defer p.dec()

	p.print("- [appl]")
	p.printType(node)
	node.Target.Accept(p)
	for _, arg := range node.Args {
		arg.Accept(p)
	}
}

func (p *AstPrinter) VisitApplArg(node *ast.ApplArg) {
	p.inc()
	defer p.dec()

	p.print("- [appl-arg %d]", node.Index)
	p.printType(node)
	node.ValueExpr.Accept(p)
}
