package internal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/renatopp/golden/lang"
)

func Parse(tokens []*lang.Token) (*Node, error) {
	scanner := lang.NewTokenScanner(tokens)
	parser := &parser{
		Parser: lang.NewParser(scanner),
		Pratt:  lang.NewPrattSolver[*Node](),
	}

	parser.Pratt.SetPrecedenceFn(parser.precedence)
	parser.Pratt.RegisterPrefixFn(TInteger, parser.parseInteger)
	parser.Pratt.RegisterPrefixFn(THex, parser.parseInteger)
	parser.Pratt.RegisterPrefixFn(TOctal, parser.parseInteger)
	parser.Pratt.RegisterPrefixFn(TBinary, parser.parseInteger)
	parser.Pratt.RegisterPrefixFn(TFloat, parser.parseFloat)
	parser.Pratt.RegisterPrefixFn(TBool, parser.parseBool)
	parser.Pratt.RegisterPrefixFn(TString, parser.parseString)
	parser.Pratt.RegisterPrefixFn(TLbrace, parser.parseBlock)
	parser.Pratt.RegisterPrefixFn(TPlus, parser.parseUnaryOperator)
	parser.Pratt.RegisterPrefixFn(TMinus, parser.parseUnaryOperator)
	parser.Pratt.RegisterPrefixFn(TBang, parser.parseUnaryOperator)

	parser.Pratt.RegisterInfixFn(TPlus, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TMinus, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TStar, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TSlash, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TPercent, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TSpaceship, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TEqual, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TNequal, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TAnd, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TOr, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TXor, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TLt, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TLte, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TGt, parser.parseBinaryOperator)
	parser.Pratt.RegisterInfixFn(TGte, parser.parseBinaryOperator)

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
	switch {
	case t.IsKind(TAssign):
		return 10
	case t.IsKind(TOr):
		return 40
	case t.IsKind(TXor):
		return 45
	case t.IsKind(TAnd):
		return 50
	case t.IsKind(TEqual, TNequal):
		return 70
	case t.IsKind(TLt, TGt, TLte, TGte):
		return 80
	case t.IsKind(TPlus, TMinus):
		return 90
	case t.IsKind(TStar, TSlash):
		return 100
	case t.IsKind(TSpaceship):
		return 110
	case t.IsKind(TPercent):
		return 120
	case t.IsKind(TLparen):
		return 130
	case t.IsKind(TLbrace):
		return 140
	}

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

		switch {
		case p.IsNext(TKeyword, KFn):
			fn := p.parseFunctionDecl()
			functions = append(functions, fn)
			continue
		}

		break
	}

	return NewNode(first, &AstModule{
		Imports:   imports,
		Types:     types,
		Functions: functions,
		Variables: variables,
	})
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

	body := p.parseBlock()

	return NewNode(fn, &AstFunctionDecl{
		Name: name.Literal,
		Body: body,
	})
}

func (p *parser) parseBlock() *Node {
	p.Expect(TLbrace)
	first := p.EatToken()
	p.Skip(TNewline, TSemicolon)

	expressions := []*Node{}
	for {
		expr := p.parseExpression()
		if expr == nil {
			break
		}

		expressions = append(expressions, expr)
		p.ExpectToken(TNewline, TSemicolon, TRbrace)
		p.Skip(TNewline, TSemicolon)
	}

	p.Expect(TRbrace)
	p.EatToken()

	return NewNode(first, &AstBlock{
		Expressions: expressions,
	})
}

// nullable
func (p *parser) parseExpression(precedence ...int) *Node {
	pr := 0
	if len(precedence) > 0 {
		pr = precedence[0]
	}
	return p.Pratt.SolveExpression(p.Scanner, pr)
}

func (p *parser) parseInteger() *Node {
	p.ExpectToken(TInteger, THex, TOctal, TBinary)

	token := p.EatToken()
	base := 10
	switch token.Kind {
	case THex:
		base = 16
	case TOctal:
		base = 8
	case TBinary:
		base = 2
	}

	value, err := strconv.ParseInt(token.Literal, base, 64)
	if err != nil {
		panic(lang.NewError(token.Loc, "invalid integer", token.Literal))
	}

	return NewNode(token, &AstInt{Value: value})
}

func (p *parser) parseFloat() *Node {
	p.ExpectToken(TFloat)
	token := p.EatToken()
	value, err := strconv.ParseFloat(token.Literal, 64)
	if err != nil {
		panic(lang.NewError(token.Loc, "invalid float", token.Literal))
	}
	return NewNode(token, &AstFloat{Value: value})
}

func (p *parser) parseBool() *Node {
	p.ExpectToken(TBool)
	p.ExpectLiteral("true", "false")
	token := p.EatToken()
	value := token.Literal == "true"
	return NewNode(token, &AstBool{Value: value})
}

func (p *parser) parseString() *Node {
	p.ExpectToken(TString)
	token := p.EatToken()
	value := strings.ReplaceAll(token.Literal, "\r", "")
	return NewNode(token, &AstString{Value: value})
}

func (p *parser) parseUnaryOperator() *Node {
	op := p.EatToken()
	right := p.parseExpression(p.precedence(op))
	return NewNode(op, &AstUnary{
		Op:    op.Literal,
		Right: right,
	})
}

func (p *parser) parseBinaryOperator(left *Node) *Node {
	op := p.EatToken()
	right := p.parseExpression(p.precedence(op))

	if right == nil {
		panic(lang.NewError(op.Loc, "expecting expression", ""))
	}

	return NewNode(op, &AstBinary{
		Op:    op.Literal,
		Left:  left,
		Right: right,
	})
}
