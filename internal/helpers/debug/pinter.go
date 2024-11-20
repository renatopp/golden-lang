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
}
