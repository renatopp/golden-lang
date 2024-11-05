package internal

import (
	"fmt"

	"github.com/renatopp/golden/lang"
)

func Parse(tokens []*lang.Token) (*Node, error) {
	scanner := lang.NewTokenScanner(tokens)
	parser := &parser{
		Parser: lang.NewParser(scanner),
		Pratt:  lang.NewPrattSolver[*Node](),
	}

	parser.Pratt.SetPrecedenceFn(parser.precedence)

	node := parser.Parse()
	if parser.Scanner.HasErrors() || parser.HasErrors() {
		return nil, lang.NewErrorList(append(parser.Errors(), parser.Scanner.Errors()...))
	}

	return node, nil
}

type parser struct {
	*lang.Parser
	Pratt *lang.PrattSolver[*Node]
}

func (p *parser) Parse() (out *Node) {
	defer func() {
		if r := recover(); r != nil {
			if r == nil {
				return
			} else if err, ok := r.(lang.Error); ok {
				p.RegisterError(err)
			} else {
				p.RegisterError(lang.NewError(lang.Loc{}, "unknown error", fmt.Sprintf("%v", r)))
			}
		}
	}()

	return p.parseModule()
}

// Overriding methods from lang.Parser

func (p *parser) ExpectToken(kinds ...string) {
	if !p.Parser.ExpectToken(kinds...) {
		panic(nil)
	}
}
func (p *parser) ExpectLiteral(literals ...string) {
	if !p.Parser.ExpectLiteral(literals...) {
		panic(nil)
	}
}
func (p *parser) Expect(kind string, literals ...string) {
	if !p.Parser.Expect(kind, literals...) {
		panic(nil)
	}
}
func (p *parser) ExpectSkipToken1(kinds ...string) {
	if !p.Parser.ExpectSkipToken1(kinds...) {
		panic(nil)
	}
}
func (p *parser) ExpectSkipTokenAll(kinds ...string) {
	if !p.Parser.ExpectSkipTokenAll(kinds...) {
		panic(nil)
	}
}
func (p *parser) ExpectSkipLiteral1(literals ...string) {
	if !p.Parser.ExpectSkipLiteral1(literals...) {
		panic(nil)
	}
}
func (p *parser) ExpectSkipLiteralAll(literals ...string) {
	if !p.Parser.ExpectSkipLiteralAll(literals...) {
		panic(nil)
	}
}
func (p *parser) ExpectSkip1(kind string, literals ...string) {
	if !p.Parser.ExpectSkip1(kind, literals...) {
		panic(nil)
	}
}
func (p *parser) ExpectSkipAll(kind string, literals ...string) {
	if !p.Parser.ExpectSkipAll(kind, literals...) {
		panic(nil)
	}
}

// Custom methods

func (p *parser) precedence(t *lang.Token) int {
	return 0
}

func (p *parser) parseModule() *Node {
	first := p.PeekToken()
	imports := []*Node{}
	types := []*Node{}
	functions := []*Node{}
	variables := []*Node{}

	for {
		p.Skip(TNewline)
		stmt := p.parseStatement()
		if stmt == nil {
			break
		}
		switch stmt.Data.(type) {
		// case *AstImport { imports = append(imports, stmt) }
		case *AstFunctionDecl:
			functions = append(functions, stmt)
		}
	}

	return NewNode(first, &AstModule{
		Imports:   imports,
		Types:     types,
		Functions: functions,
		Variables: variables,
	})
}

func (p *parser) parseStatement() *Node {
	switch {
	case p.IsNext(TKeyword, KFn):
		return p.parseFunctionDecl()
	default:
		return nil
	}
}

func (p *parser) parseFunctionDecl() *Node {
	p.Expect(TKeyword, KFn)
	fn := p.EatToken()

	p.Expect(TVarIdent)
	name := p.EatToken()

	p.Expect(TLparen)
	p.EatToken()

	p.Expect(TRparen)
	p.EatToken()

	p.Expect(TLbrace)
	p.EatToken()

	p.Expect(TRbrace)
	p.EatToken()

	return NewNode(fn, &AstFunctionDecl{
		Name: name.Literal,
		Body: NewNode(fn, nil),
	})
}
