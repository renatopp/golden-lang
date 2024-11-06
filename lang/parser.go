package lang

import "slices"

type Parser struct {
	*ErrorData
	Scanner  *TokenScanner
	previous *Token
}

func NewParser(scanner *TokenScanner) *Parser {
	return &Parser{NewErrorData(), scanner, nil}
}

func (p *Parser) EatTokens(n int) []*Token {
	tks := p.Scanner.EatTokens(n)
	p.previous = tks[len(tks)-1]
	return tks
}

func (p *Parser) EatToken() *Token {
	tk := p.Scanner.EatToken()
	p.previous = tk
	return tk
}

func (p *Parser) PeekToken() *Token { return p.Scanner.PeekToken() }

func (p *Parser) PeekTokenAt(n int) *Token { return p.Scanner.PeekTokenAt(n) }

func (p *Parser) PreviousToken() *Token { return p.previous }

func (p *Parser) Skip(kinds ...string) []*Token {
	res := []*Token{}
	for p.IsNextTokens(kinds...) {
		res = append(res, p.EatToken())
	}
	return res
}

func (p *Parser) SkipN(n int, kinds ...string) []*Token {
	res := []*Token{}
	for i := 0; i < n && p.IsNextTokens(kinds...); i++ {
		res = append(res, p.EatToken())
	}
	return res
}

func (p *Parser) IsNextTokens(kinds ...string) bool {
	if len(kinds) == 0 {
		return true
	}
	return slices.Contains(kinds, p.PeekToken().Kind)
}

func (p *Parser) IsNextLiterals(literals ...string) bool {
	if len(literals) == 0 {
		return true
	}
	return slices.Contains(literals, p.PeekToken().Literal)
}

func (p *Parser) IsNextLiteralsOf(kind string, literals ...string) bool {
	return p.IsNextTokens(kind) && p.IsNextLiterals(literals...)
}
