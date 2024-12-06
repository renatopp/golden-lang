package golang

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/token"
	"github.com/renatopp/golden/internal/compiler/types"
	"github.com/renatopp/golden/internal/helpers/codegen"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/naming"
	"github.com/renatopp/golden/internal/helpers/tmpl"
)

//go:embed templates/module.go.tmpl
var raw_template_module string
var template_module, _ = template.New("module").Parse(raw_template_module)

var _ ast.Visitor = &Writer{}

type Writer struct {
	*ast.Visiter
	backend   *Golang
	stack     []string
	identer   *codegen.Identer
	funcLevel int
}

func NewWriter(backend *Golang) *Writer {
	w := &Writer{
		backend: backend,
		identer: codegen.NewIdenter(),
	}
	w.Visiter = ast.NewVisiter(w)
	return w
}

func (w *Writer) Push(s string) {
	w.stack = append(w.stack, s)
}

func (w *Writer) Pop() string {
	s := w.stack[len(w.stack)-1]
	w.stack = w.stack[:len(w.stack)-1]
	return s
}

func (w *Writer) Generate(packageName string, root *ast.Module) string {
	root.Visit(w)
	return tmpl.GenerateString(template_module, map[string]any{
		"PackageName": packageName,
		"Exprs":       w.Pop(),
	})
}

func (w *Writer) VisitModule(node *ast.Module) ast.Node {
	decls := []string{}
	for _, expr := range node.Exprs {
		expr.Visit(w)
		decls = append(decls, w.Pop())
	}

	w.Push(strings.Join(decls, "\n"))
	return node
}

func (w *Writer) VisitVarDecl(node *ast.VarDecl) ast.Node {
	node.Name.Visit(w)
	name := w.Pop()

	w.resolveType(node.Type.Unwrap())
	type_ := w.Pop()

	node.ValueExpr = node.ValueExpr.Visit(w)
	value := w.Pop()

	w.Push(fmt.Sprintf("var %s %s = %s", name, type_, value))
	return node
}

func (w *Writer) VisitVarIdent(node *ast.VarIdent) ast.Node {
	w.Push(w.name(node.Value))
	return node
}

func (w *Writer) VisitTypeIdent(node *ast.TypeIdent) ast.Node {
	errors.ThrowAtNode(node, errors.InternalError, "TypeIdent should not be visited, use resolveType instead")
	return node
}

func (w *Writer) VisitInt(node *ast.Int) ast.Node {
	w.Push(fmt.Sprintf("%d", node.Value))
	return node
}

func (w *Writer) VisitFloat(node *ast.Float) ast.Node {
	w.Push(fmt.Sprintf("%f", node.Value))
	return node
}

func (w *Writer) VisitString(node *ast.String) ast.Node {
	w.Push(fmt.Sprintf("%q", node.Value))
	return node
}

func (w *Writer) VisitBool(node *ast.Bool) ast.Node {
	w.Push(fmt.Sprintf("%t", node.Value))
	return node
}

func (w *Writer) VisitBinOp(node *ast.BinOp) ast.Node {
	node.LeftExpr.Visit(w)
	left := w.Pop()

	node.RightExpr.Visit(w)
	right := w.Pop()

	op := ""
	switch node.Op {
	case token.KindToLiteral(token.TAnd):
		op = "&&"
	case token.KindToLiteral(token.TOr):
		op = "||"
	case token.KindToLiteral(token.TEqual):
		op = "=="
	case token.KindToLiteral(token.TNotEqual):
		op = "!="
	case token.KindToLiteral(token.TLess):
		op = "<"
	case token.KindToLiteral(token.TLessEqual):
		op = "<="
	case token.KindToLiteral(token.TGreater):
		op = ">"
	case token.KindToLiteral(token.TGreaterEqual):
		op = ">="
	case token.KindToLiteral(token.TPlus):
		op = "+"
	case token.KindToLiteral(token.TMinus):
		op = "-"
	case token.KindToLiteral(token.TStar):
		op = "*"
	case token.KindToLiteral(token.TSlash):
		op = "/"
	case token.KindToLiteral(token.TPercent):
		op = "%"
	case token.KindToLiteral(token.TSpaceShip):
		errors.ThrowAtNode(node, errors.NotImplemented, "Spaceship operator not implemented yet")
		// term := fmt.Sprintf("((%s < %s) ? -1 : (%s > %s) ? 1 : 0)", left, right, left, right)
		// w.Push(term)
		return node
	}

	w.Push(fmt.Sprintf("(%s %s %s)", left, op, right))
	return node
}

func (w *Writer) VisitUnaryOp(node *ast.UnaryOp) ast.Node {
	node.RightExpr.Visit(w)
	right := w.Pop()

	w.Push(fmt.Sprintf("%s%s", node.Op, right))
	return node
}

func (w *Writer) VisitBlock(node *ast.Block) ast.Node {
	exprs := []string{}
	for _, expr := range node.Exprs {
		expr.Visit(w)
		exprs = append(exprs, w.Pop())
	}

	w.Push(strings.Join(exprs, "\n"))
	return node
}

func (w *Writer) VisitFnDecl(node *ast.FnDecl) ast.Node {
	name := ""
	if node.Name.Has() {
		name = w.name(node.Name.Unwrap().Value)
	}

	w.resolveType(node.TypeExpr.GetType().Unwrap())
	type_ := w.Pop()

	params := codegen.JoinList(", ", node.Params, func(p *ast.FnDeclParam) string {
		p.Visit(w)
		return w.Pop()
	})

	w.identer.Inc()
	w.funcLevel++
	node.ValueExpr.Visit(w)
	body := w.identer.Indent(w.Pop())
	w.funcLevel--
	w.identer.Dec()

	w.Push(fmt.Sprintf("func %s(%s) %s {\n%s\n}", name, params, type_, body))
	return node
}

func (w *Writer) VisitFnDeclParam(node *ast.FnDeclParam) ast.Node {
	node.Name.Visit(w)
	name := w.Pop()

	w.resolveType(node.TypeExpr.GetType().Unwrap())
	tp := w.Pop()

	w.Push(fmt.Sprintf("%s %s", name, tp))
	return node
}

func (w *Writer) VisitApplication(node *ast.Application) ast.Node {
	node.Target.Visit(w)
	target := w.Pop()

	args := codegen.JoinList(", ", node.Args, func(a ast.Node) string {
		a.Visit(w)
		return w.Pop()
	})

	w.Push(fmt.Sprintf("%s(%s)", target, args))
	return node
}

func (w *Writer) VisitReturn(node *ast.Return) ast.Node {
	node.ValueExpr.Visit(w)
	value := w.Pop()

	w.Push(fmt.Sprintf("return %s", value))
	return node
}

func (w *Writer) name(n string) string {
	if naming.IsPrivateName(n) {
		return strings.ToLower(n[:1]) + n[1:]
	} else {
		return strings.ToUpper(n[:1]) + n[1:]
	}
}

func (w *Writer) resolveType(tp ast.Type) {
	switch tp := tp.(type) {
	case *types.Primitive:
		switch tp.Name {
		case "Int":
			w.Push("int64")
		case "Float":
			w.Push("float64")
		case "String":
			w.Push("string")
		case "Bool":
			w.Push("bool")
		default:
			errors.Throw(errors.NotImplemented, "primitive type %v not implemented in go backend", tp)
		}

	case *types.Unit:
		w.Push("")

	case *types.Function:
		params := codegen.JoinList(", ", tp.Params, func(p ast.Type) string {
			w.resolveType(p)
			return w.Pop()
		})
		w.resolveType(tp.Return)
		returns := w.Pop()
		w.Push(fmt.Sprintf("func(%s) %s", params, returns))

	default:
		errors.Throw(errors.InternalError, "unknown type %s", tp.GetSignature())
	}
}
