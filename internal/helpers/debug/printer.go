package debug

import (
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/helpers/iter"
	"github.com/renatopp/golden/internal/helpers/str"
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

func (p *AstPrinter) inc()           { p.depth++ }
func (p *AstPrinter) dec()           { p.depth-- }
func (p *AstPrinter) indent() string { return strings.Repeat("  ", p.depth-1) }
func (p *AstPrinter) print(n ast.Node, s string, args ...any) {
	first := fmt.Sprintf(p.indent()+s, args...)
	fmt.Print(first)
	n.GetType().IfElse(func(tp ast.Type) {
		second := fmt.Sprintf(" → %s", tp.GetSignature())
		println(str.Repeat(" ", 50-len(first)), second)
	}, func() {
		println()
	})
}

func (p *AstPrinter) VisitModule(node *ast.Module) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[module]")
	iter.Each(node.Exprs, func(e ast.Node) { e.Visit(p) })
	return node
}

func (p *AstPrinter) VisitVarDecl(node *ast.VarDecl) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[let]")
	node.Name.Visit(p)
	node.TypeExpr.If(func(n ast.Node) { n.Visit(p) })
	node.ValueExpr.Visit(p)
	return node
}

func (p *AstPrinter) VisitInt(node *ast.Int) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[int:%d]", node.Value)
	return node
}

func (p *AstPrinter) VisitFloat(node *ast.Float) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[float:%f]", node.Value)
	return node
}

func (p *AstPrinter) VisitString(node *ast.String) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[string:'%s']", Escape(node.Value))
	return node
}

func (p *AstPrinter) VisitBool(node *ast.Bool) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[bool:%t]", node.Value)
	return node
}

func (p *AstPrinter) VisitVarIdent(node *ast.VarIdent) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[var-ident:%s]", node.Value)
	return node
}

func (p *AstPrinter) VisitTypeIdent(node *ast.TypeIdent) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[type-ident:%s]", node.Value)
	return node
}

func (p *AstPrinter) VisitBinOp(node *ast.BinOp) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[bin-op:%s]", node.Op)
	node.LeftExpr.Visit(p)
	node.RightExpr.Visit(p)
	return node
}

func (p *AstPrinter) VisitUnaryOp(node *ast.UnaryOp) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[unary-op:%s]", node.Op)
	node.RightExpr.Visit(p)
	return node
}

func (p *AstPrinter) VisitBlock(node *ast.Block) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[block]")
	iter.Each(node.Exprs, func(e ast.Node) { e.Visit(p) })
	return node
}

func (p *AstPrinter) VisitFnDecl(node *ast.FnDecl) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[fn-decl]")
	node.Name.If(func(n *ast.VarIdent) { n.Visit(p) })
	iter.Each(node.Params, func(n *ast.FnDeclParam) { n.Visit(p) })
	node.TypeExpr.Visit(p)
	node.ValueExpr.Visit(p)
	return node
}

func (p *AstPrinter) VisitFnDeclParam(node *ast.FnDeclParam) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[fn-decl-param]")
	node.Name.Visit(p)
	node.TypeExpr.Visit(p)
	return node
}

func (p *AstPrinter) VisitTypeFn(node *ast.TypeFn) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[type-fn]")
	iter.Each(node.Parameters, func(n ast.Node) { n.Visit(p) })
	node.ReturnExpr.Visit(p)
	return node
}

func (p *AstPrinter) VisitApplication(node *ast.Application) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[application]")
	node.Target.Visit(p)
	iter.Each(node.Args, func(n ast.Node) { n.Visit(p) })
	return node
}

func (p *AstPrinter) VisitReturn(node *ast.Return) ast.Node {
	p.inc()
	defer p.dec()
	p.print(node, "[return]")
	node.ValueExpr.If(func(n ast.Node) { n.Visit(p) })
	return node
}
