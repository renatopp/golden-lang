package syntax

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/tokens"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/lang"
)

func Parse(tokens []*lang.Token, modulePath string) (*ast.Module, error) {
	scanner := lang.NewTokenScanner(tokens)
	parser := &parser{
		ModulePath:  modulePath,
		Parser:      lang.NewParser(scanner),
		ValueSolver: lang.NewPrattSolver[ast.Node](),
		TypeSolver:  lang.NewPrattSolver[ast.Node](),
	}

	parser.ValueSolver.SetPrecedenceFn(parser.valuePrecedence)
	parser.TypeSolver.SetPrecedenceFn(parser.typePrecedence)
	parser.registerValueExpressions()
	parser.registerTypeExpressions()

	return parser.Parse()
}

type parser struct {
	*lang.Parser
	ModulePath  string
	ValueSolver *lang.PrattSolver[ast.Node]
	TypeSolver  *lang.PrattSolver[ast.Node]
}

func (p *parser) Parse() (a *ast.Module, err error) {
	err = errors.WithRecovery(func() {
		a = p.parseModule()
	})
	return a, err
}

func (p *parser) ExpectTokens(kinds ...string) {
	if !p.IsNextTokens(kinds...) {
		errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected token '%s', got '%s'", kinds, p.PeekToken().Kind)
	}
}

func (p *parser) SkipNewlines() {
	p.Skip(tokens.TNewline)
}

func (p *parser) SkipSeparator(kind ...string) {
	p.SkipNewlines()
	p.SkipN(1, kind...)
	p.SkipNewlines()
}
