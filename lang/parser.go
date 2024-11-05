package lang

import (
	"fmt"
	"strings"
)

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
	for p.IsNextToken(kinds...) {
		res = append(res, p.EatToken())
	}
	return res
}

func (p *Parser) SkipN(n int, kinds ...string) []*Token {
	res := []*Token{}
	for i := 0; i < n && p.IsNextToken(kinds...); i++ {
		res = append(res, p.EatToken())
	}
	return res
}

func (p *Parser) ExpectToken(kinds ...string) bool {
	cur := p.PeekToken()
	for _, k := range kinds {
		if cur.Kind == k {
			return true
		}
	}

	expected := strings.Join(kinds, ", ")
	err := cur.AsError(ErrSyntax, fmt.Sprintf(errMsgUnexpectedToken, expected, cur.Kind))
	p.RegisterError(err)
	return false
}

func (p *Parser) ExpectLiteral(literals ...string) bool {
	if len(literals) == 0 {
		return true
	}

	cur := p.PeekToken()
	for _, lit := range literals {
		if cur.Literal == lit {
			return true
		}
	}

	expected := strings.Join(literals, ", ")
	err := cur.AsError(ErrSyntax, fmt.Sprintf(errMsgUnexpectedLiteral, expected, cur.Literal))
	p.RegisterError(err)
	return false
}

func (p *Parser) Expect(kind string, literals ...string) bool {
	return p.ExpectToken(kind) && p.ExpectLiteral(literals...)
}

func (p *Parser) ExpectSkipToken1(kinds ...string) bool {
	if p.ExpectToken(kinds...) {
		p.EatToken()
		return true
	}
	return false
}

func (p *Parser) ExpectSkipTokenAll(kinds ...string) bool {
	once := false
	for p.ExpectSkipToken1(kinds...) {
		once = true
	}
	return once
}

func (p *Parser) ExpectSkipLiteral1(literals ...string) bool {
	if p.ExpectLiteral(literals...) {
		p.EatToken()
		return true
	}
	return false
}

func (p *Parser) ExpectSkipLiteralAll(literals ...string) bool {
	once := false
	for p.ExpectSkipLiteral1(literals...) {
		once = true
	}
	return once
}

func (p *Parser) ExpectSkip1(kind string, literals ...string) bool {
	return p.Expect(kind, literals...) && p.ExpectSkipLiteral1(literals...)
}

func (p *Parser) ExpectSkipAll(kind string, literals ...string) bool {
	once := false
	for p.ExpectSkip1(kind, literals...) {
		once = true
	}
	return once
}

func (p *Parser) IsNextToken(kinds ...string) bool {
	cur := p.PeekToken()
	for _, k := range kinds {
		if cur.Kind == k {
			return true
		}
	}
	return false
}

func (p *Parser) IsNextLiteral(literals ...string) bool {
	cur := p.PeekToken()
	for _, lit := range literals {
		if cur.Literal == lit {
			return true
		}
	}
	return false
}

func (p *Parser) IsNext(kind string, literals ...string) bool {
	return p.IsNextToken(kind) && p.IsNextLiteral(literals...)
}
