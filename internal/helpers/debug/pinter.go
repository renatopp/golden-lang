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
	fmt.Printf(p.indent()+s+"\n", args...)
}

func (p *AstPrinter) VisitModule(node *ast.Module) {
	p.inc()
	defer p.dec()

	p.print("- [module %s]", node.ModulePath)
	for _, stmt := range node.Imports {
		stmt.Accept(p)
	}
}

func (p *AstPrinter) VisitImport(node *ast.Import) {
	p.inc()
	defer p.dec()

	p.print("- [import]")
	node.Path.Accept(p)
	node.Alias.If(func(alias *ast.VarIdent) {
		alias.Accept(p)
	})
}

func (p *AstPrinter) VisitString(node *ast.String) {
	p.inc()
	defer p.dec()

	p.print("- [string %s]", node.Literal)
}

func (p *AstPrinter) VisitVarIdent(node *ast.VarIdent) {
	p.inc()
	defer p.dec()

	p.print("- [var-ident %s]", node.Literal)
}
