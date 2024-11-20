package syntax

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/lang"
)

func Parse(tokens []*lang.Token, module *core.Module) (*ast.Module, error) {
	scanner := lang.NewTokenScanner(tokens)
	parser := &parser{
		Module:      module,
		Parser:      lang.NewParser(scanner),
		ValueSolver: lang.NewPrattSolver[core.Node](),
		TypeSolver:  lang.NewPrattSolver[core.Node](),
	}

	parser.ValueSolver.SetPrecedenceFn(parser.valuePrecedence)
	parser.TypeSolver.SetPrecedenceFn(parser.typePrecedence)
	parser.registerValueExpressions()
	parser.registerTypeExpressions()

	return parser.Parse()
}

type parser struct {
	*lang.Parser
	Module      *core.Module
	ValueSolver *lang.PrattSolver[core.Node]
	TypeSolver  *lang.PrattSolver[core.Node]
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
	p.Skip(core.TNewline)
}

func (p *parser) SkipSeparator(kind ...string) {
	p.SkipNewlines()
	p.SkipN(1, kind...)
	p.SkipNewlines()
}
