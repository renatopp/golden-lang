package syntax

import (
	"fmt"

	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/lang"
)

func Parse(tokens []*lang.Token) (*core.AstNode, error) {
	scanner := lang.NewTokenScanner(tokens)
	parser := &Parser{
		Parser:      lang.NewParser(scanner),
		ValueSolver: lang.NewPrattSolver[*core.AstNode](),
		TypeSolver:  lang.NewPrattSolver[*core.AstNode](),
	}

	parser.ValueSolver.SetPrecedenceFn(parser.valuePrecedence)
	parser.registerValueExpressions()
	parser.TypeSolver.SetPrecedenceFn(parser.typePrecedence)
	parser.registerTypeExpressions()

	module := parser.Parse()
	if parser.Scanner.HasErrors() || parser.HasErrors() {
		return nil, lang.NewErrorList(append(parser.Errors(), parser.Scanner.Errors()...))
	}

	return module, nil
}

type Parser struct {
	*lang.Parser
	ValueSolver *lang.PrattSolver[*core.AstNode]
	TypeSolver  *lang.PrattSolver[*core.AstNode]
}

func (p *Parser) Parse() *core.AstNode {
	defer func() {
		r := recover()
		if r == nil {
			return
		} else if err, ok := r.(lang.Error); ok {
			p.RegisterError(err)
		} else {
			p.RegisterError(lang.NewError(lang.Loc{}, "unknown error", fmt.Sprintf("%v", r)))
		}
	}()

	return p.parseModule()
}

func (p *Parser) ExpectTokens(kind ...string) {
	if !p.IsNextTokens(kind...) {
		p.Error(p.PeekToken().Loc, "unexpected token", "expected %s, got %s", kind, p.PeekToken().Kind)
	}
}

func (p *Parser) ExpectLiterals(lit ...string) {
	if !p.IsNextLiterals(lit...) {
		p.Error(p.PeekToken().Loc, "unexpected literal", "expected %s, got %s", lit, p.PeekToken().Literal)
	}
}

func (p *Parser) ExpectLiteralsOf(kind string, lit ...string) {
	p.ExpectTokens(kind)
	p.ExpectLiterals(lit...)
}

func (p *Parser) Error(loc lang.Loc, kind, msg string, args ...any) {
	panic(lang.NewError(loc, kind, fmt.Sprintf(msg, args...)))
}

func (p *Parser) SkipNewlines() {
	p.Skip(core.TNewline)
}

func (p *Parser) SkipSeparator(kind ...string) {
	p.SkipNewlines()
	p.SkipN(1, kind...)
	p.SkipNewlines()
}

func (p *Parser) parseModule() *core.AstNode {
	imports := []*ast.ModuleImport{}
	types := []*core.AstNode{}
	functions := []*core.AstNode{}
	variables := []*core.AstNode{}

	first := p.PeekToken()
	p.Skip(core.TNewline)
	for {
		switch {
		case p.IsNextLiteralsOf(core.TKeyword, core.KImport):
			imports = append(imports, p.parseImport())
		case p.IsNextLiteralsOf(core.TKeyword, core.KData):
			types = append(types, p.parseTypeExpression())
		case p.IsNextLiteralsOf(core.TKeyword, core.KFn):
			functions = append(functions, p.parseValueExpression())
		case p.IsNextLiteralsOf(core.TKeyword, core.KLet):
			variables = append(variables, p.parseValueExpression())
		case p.IsNextTokens(core.TComment):
			p.EatToken()
		case p.IsNextTokens(core.TEof):
			// EOF
		default:
			p.Error(
				p.PeekToken().Loc,
				"unexpected token",
				"expected import, let, data or fn, got %s:%s",
				p.PeekToken().Kind,
				p.PeekToken().Literal,
			)
		}

		p.Skip(core.TNewline)
		if p.IsNextTokens(core.TEof) {
			break
		}
	}

	return core.NewNode(first, &ast.Module{
		Imports:   imports,
		Types:     types,
		Functions: functions,
		Variables: variables,
	})
}

func (p *Parser) parseImport() *ast.ModuleImport {
	p.ExpectLiteralsOf(core.TKeyword, core.KImport)
	p.EatToken()
	path := p.EatToken().Literal
	alias := ""
	if p.IsNextLiteralsOf(core.TKeyword, core.KAs) {
		p.EatToken()
		p.ExpectTokens(core.TVarIdent)
		alias = p.EatToken().Literal
	}
	return &ast.ModuleImport{Path: path, Alias: alias}
}
