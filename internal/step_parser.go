package internal

import (
	"fmt"

	"github.com/renatopp/golden/lang"
)

func Parse(tokens []*lang.Token) (*Node, error) {
	scanner := lang.NewTokenScanner(tokens)
	parser := &parser{
		Parser:      lang.NewParser(scanner),
		ValueSolver: lang.NewPrattSolver[*Node](),
		TypeSolver:  lang.NewPrattSolver[*Node](),
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

type parser struct {
	*lang.Parser
	ValueSolver *lang.PrattSolver[*Node]
	TypeSolver  *lang.PrattSolver[*Node]
}

func (p *parser) Parse() *Node {
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

func (p *parser) ExpectTokens(kind ...string) {
	if !p.IsNextTokens(kind...) {
		p.Error(p.PeekToken().Loc, "unexpected token", "expected %s, got %s", kind, p.PeekToken().Kind)
	}
}

func (p *parser) ExpectLiterals(lit ...string) {
	if !p.IsNextLiterals(lit...) {
		p.Error(p.PeekToken().Loc, "unexpected literal", "expected %s, got %s", lit, p.PeekToken().Literal)
	}
}

func (p *parser) ExpectLiteralsOf(kind string, lit ...string) {
	p.ExpectTokens(kind)
	p.ExpectLiterals(lit...)
}

func (p *parser) Error(loc lang.Loc, kind, msg string, args ...any) {
	panic(lang.NewError(loc, kind, fmt.Sprintf(msg, args...)))
}

func (p *parser) SkipNewlines() {
	p.Skip(TNewline)
}

func (p *parser) SkipSeparator(kind ...string) {
	p.SkipNewlines()
	p.SkipN(1, kind...)
	p.SkipNewlines()
}

func (p *parser) parseModule() *Node {
	imports := []*AstModuleImport{}
	types := []*Node{}
	functions := []*Node{}
	variables := []*Node{}

	first := p.PeekToken()
	p.Skip(TNewline)
	for {
		switch {
		case p.IsNextLiteralsOf(TKeyword, KImport):
			imports = append(imports, p.parseImport())
		case p.IsNextLiteralsOf(TKeyword, KData):
			types = append(types, p.parseTypeExpression())
		case p.IsNextLiteralsOf(TKeyword, KFn):
			functions = append(functions, p.parseValueExpression())
		case p.IsNextLiteralsOf(TKeyword, KLet):
			variables = append(variables, p.parseValueExpression())
		case p.IsNextTokens(TEof):
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

		p.Skip(TNewline)
		if p.IsNextTokens(TEof) {
			break
		}
	}

	return NewNode(first, &AstModule{
		Imports:   imports,
		Types:     types,
		Functions: functions,
		Variables: variables,
	})
}

func (p *parser) parseImport() *AstModuleImport {
	p.ExpectLiteralsOf(TKeyword, KImport)
	p.EatToken()
	path := p.EatToken().Literal
	alias := ""
	if p.IsNextLiteralsOf(TKeyword, KAs) {
		p.EatToken()
		p.ExpectTokens(TVarIdent)
		alias = p.EatToken().Literal
	}
	return &AstModuleImport{Path: path, Alias: alias}
}
