package syntax

import (
	"strings"

	"github.com/renatopp/golden/internal/compiler/token"
	"github.com/renatopp/golden/internal/helpers/errors"
)

type BaseParser struct {
	ValueSolver *PrattSolver
	TypeSolver  *PrattSolver
	Scanner     *Scanner[token.Token]
}

func NewBaseParser(tokens []token.Token) *BaseParser {
	p := &BaseParser{
		Scanner: NewScanner(tokens, token.Token{Kind: token.TEof}),
	}
	p.ValueSolver = NewPrattSolver(p.Peek, p.ValuePrecedence)
	p.TypeSolver = NewPrattSolver(p.Peek, p.TypePrecedence)
	return p
}

func (p *BaseParser) Eat() token.Token { return p.Scanner.Eat() }

func (p *BaseParser) EatN(n int) []token.Token { return p.Scanner.EatN(n) }

func (p *BaseParser) Peek() token.Token { return p.Scanner.Peek() }

func (p *BaseParser) PeekN(n int) token.Token { return p.Scanner.PeekAt(n) }

func (p *BaseParser) Skip(kinds ...token.TokenKind) []token.Token {
	res := []token.Token{}
	for p.IsNext(kinds...) {
		res = append(res, p.Eat())
	}
	return res
}

func (p *BaseParser) SkipN(n int, kinds ...token.TokenKind) []token.Token {
	res := []token.Token{}
	for i := 0; i < n && p.IsNext(kinds...); i++ {
		res = append(res, p.Eat())
	}
	return res
}

func (p *BaseParser) SkipNewlines() {
	p.Skip(token.TNewline)
}

func (p *BaseParser) SkipSeparator(kind ...token.TokenKind) {
	p.SkipNewlines()
	p.SkipN(1, kind...)
	p.SkipNewlines()
}

func (p *BaseParser) IsNext(kinds ...token.TokenKind) bool {
	next := p.Peek()
	for _, kind := range kinds {
		if next.Kind == kind {
			return true
		}
	}
	return false
}

func (p *BaseParser) IsNextLiteral(literals ...string) bool {
	next := p.Peek()
	for _, literal := range literals {
		if next.Literal == literal {
			return true
		}
	}
	return false
}

func (p *BaseParser) Expect(kinds ...token.TokenKind) {
	if !p.IsNext(kinds...) {
		names := []string{}
		for _, kind := range kinds {
			names = append(names, token.KindToLiteral(kind))
		}
		list := strings.Join(names, ", ")
		tok := p.Peek()
		errors.ThrowAtToken(tok, errors.ParserError, "expected token '%s', got '%s'", list, tok.Display())
	}
}

func (p *BaseParser) ExpectAndEat(kinds ...token.TokenKind) token.Token {
	p.Expect(kinds...)
	return p.Eat()
}

func (p *BaseParser) ValuePrecedence(t token.Token) int {
	switch {
	case t.Is(token.TAssign):
		return 10
	// case t.Is(token.TPipe):
	// 	return 20
	case t.Is(token.TOr):
		return 40
	case t.Is(token.TXor):
		return 45
	case t.Is(token.TAnd):
		return 50
	case t.Is(token.TEqual, token.TNotEqual):
		return 70
	case t.Is(token.TLess, token.TGreater, token.TLessEqual, token.TGreaterEqual):
		return 80
	case t.Is(token.TPlus, token.TMinus):
		return 90
	case t.Is(token.TStar, token.TSlash):
		return 100
	case t.Is(token.TSpaceShip):
		return 110
	case t.Is(token.TPercent):
		return 120
		// case t.Is(token.TLparen):
		// return 130
		// case t.Is(token.TDot):
		// return 140
	}
	return 0
}

func (p *BaseParser) TypePrecedence(t token.Token) int {
	return 0
}
