package debug

import (
	"fmt"
	"strings"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/helpers/iter"
	"github.com/renatopp/golden/internal/helpers/safe"
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

func (p *AstPrinter) inc()                        { p.depth++ }
func (p *AstPrinter) dec()                        { p.depth-- }
func (p *AstPrinter) indent() string              { return strings.Repeat("  ", p.depth-1) }
func (p *AstPrinter) print(s string, args ...any) { fmt.Printf(p.indent()+s, args...) }
func (p *AstPrinter) printType(tp safe.Optional[ast.Type]) {
	// tp.If(func(n ast.Node) {  })
	println()
}

func (p *AstPrinter) VisitModule(node ast.Module) ast.Node {
	p.inc()
	defer p.dec()
	p.print("[module]\n")
	iter.Each(node.Exprs, func(e ast.Node) { e.Visit(p) })
	return node
}

func (p *AstPrinter) VisitConst(node ast.Const) ast.Node {
	p.inc()
	defer p.dec()
	p.print("[const]\n")
	node.Name.Visit(p)
	node.TypeExpr.If(func(n ast.Node) { n.Visit(p) })
	node.ValueExpr.Visit(p)
	return node
}

func (p *AstPrinter) VisitInt(node ast.Int) ast.Node {
	p.inc()
	defer p.dec()
	p.print("[int:%d]\n", node.Value)
	return node
}

func (p *AstPrinter) VisitFloat(node ast.Float) ast.Node {
	p.inc()
	defer p.dec()
	p.print("[float:%f]\n", node.Value)
	return node
}

func (p *AstPrinter) VisitString(node ast.String) ast.Node {
	p.inc()
	defer p.dec()
	p.print("[string:'%s']\n", Escape(node.Value))
	return node
}

func (p *AstPrinter) VisitBool(node ast.Bool) ast.Node {
	p.inc()
	defer p.dec()
	p.print("[bool:%t]\n", node.Value)
	return node
}

func (p *AstPrinter) VisitVarIdent(node ast.VarIdent) ast.Node {
	p.inc()
	defer p.dec()
	p.print("[var-ident:%s]\n", node.Value)
	return node
}

func (p *AstPrinter) VisitTypeIdent(node ast.TypeIdent) ast.Node {
	p.inc()
	defer p.dec()
	p.print("[type-ident:%s]\n", node.Value)
	return node
}

func (p *AstPrinter) VisitBinOp(node ast.BinOp) ast.Node {
	p.inc()
	defer p.dec()
	p.print("[bin-op:%s]\n", node.Op)
	node.LeftExpr.Visit(p)
	node.RightExpr.Visit(p)
	return node
}

func (p *AstPrinter) VisitUnaryOp(node ast.UnaryOp) ast.Node {
	p.inc()
	defer p.dec()
	p.print("[unary-op:%s]\n", node.Op)
	node.RightExpr.Visit(p)
	return node
}

func (p *AstPrinter) VisitBlock(node ast.Block) ast.Node {
	p.inc()
	defer p.dec()
	p.print("[block]\n")
	iter.Each(node.Exprs, func(e ast.Node) { e.Visit(p) })
	return node
}
