package syntax

import (
	"strconv"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/token"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/safe"
)

type Parser struct {
	*BaseParser
}

func NewParser(tokens []*token.Token) *Parser {
	p := &Parser{
		BaseParser: NewBaseParser(tokens),
	}

	p.ValueSolver.RegisterPrefixFn(token.TVarIdent, p.parseVarIdent)
	p.ValueSolver.RegisterPrefixFn(token.TInt, p.parseInt)
	p.ValueSolver.RegisterPrefixFn(token.THex, p.parseHex)
	p.ValueSolver.RegisterPrefixFn(token.TOctal, p.parseOctal)
	p.ValueSolver.RegisterPrefixFn(token.TBinary, p.parseBinary)
	p.ValueSolver.RegisterPrefixFn(token.TFloat, p.parseFloat)
	p.ValueSolver.RegisterPrefixFn(token.TString, p.parseString)
	p.ValueSolver.RegisterPrefixFn(token.TTrue, p.parseBool)
	p.ValueSolver.RegisterPrefixFn(token.TFalse, p.parseBool)
	p.ValueSolver.RegisterPrefixFn(token.TPlus, p.parseUnaryOp)
	p.ValueSolver.RegisterPrefixFn(token.TMinus, p.parseUnaryOp)
	p.ValueSolver.RegisterPrefixFn(token.TBang, p.parseUnaryOp)
	p.ValueSolver.RegisterPrefixFn(token.TLeftBrace, p.parseBlock)

	p.ValueSolver.RegisterInfixFn(token.TPlus, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TMinus, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TStar, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TSlash, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TPercent, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TSpaceShip, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TEqual, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TNotEqual, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TGreater, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TGreaterEqual, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TLess, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TLessEqual, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TAnd, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TOr, p.parseBinOp)
	p.ValueSolver.RegisterInfixFn(token.TXor, p.parseBinOp)

	return p
}

func (p *Parser) Parse() (res *ast.Module, err error) {
	err = errors.WithRecovery(func() {
		res = p.parseModule()
	})
	return res, err
}

func (p *Parser) parseValueExpression(prec int) safe.Optional[ast.Node] {
	return p.ValueSolver.SolveExpression(prec)
}

// file
func (p *Parser) parseModule() *ast.Module {
	exprs := []ast.Node{}
	first := p.Peek()
	for {
		p.Skip(token.TNewline, token.TComment)
		if p.Peek().Is(token.TEof) {
			break
		}
		exprs = append(exprs, p.parseConst())
	}
	return ast.NewModule(first, exprs)
}

// const <var-ident> <type-expr>? = <value-expr>
func (p *Parser) parseConst() *ast.Const {
	tok := p.ExpectAndEat(token.TConst)       // const
	name := p.parseVarIdent().(*ast.VarIdent) // var-ident
	// TODO: add type expression parsing // type-expr
	assign := p.ExpectAndEat(token.TAssign) // =
	val := p.parseValueExpression(0)        // value-expr
	if !val.Has() {
		errors.ThrowAtToken(assign, errors.ParserError, "expected value expression after assignment, but none was found")
	}
	return ast.NewConst(tok, name, safe.None[ast.Node](), val.Unwrap())
}

// foo, bar, _bar, _1, a_1, ...
func (p *Parser) parseVarIdent() ast.Node {
	tok := p.ExpectAndEat(token.TVarIdent)
	return ast.NewVarIdent(tok, tok.Literal)
}

// 1, 1_000, ...
func (p *Parser) parseInt() ast.Node {
	tok := p.ExpectAndEat(token.TInt)
	val, err := strconv.ParseInt(tok.Literal, 10, 64)
	if err != nil {
		errors.ThrowAtToken(tok, errors.ParserError, "invalid integer literal '%s'", tok.Literal)
	}
	return ast.NewInt(tok, val)
}

// 0x1, 0x1_000, ...
func (p *Parser) parseHex() ast.Node {
	tok := p.ExpectAndEat(token.THex)
	val, err := strconv.ParseInt(tok.Literal, 16, 64)
	if err != nil {
		errors.ThrowAtToken(tok, errors.ParserError, "invalid hexadecimal literal '%s'", tok.Literal)
	}
	return ast.NewInt(tok, val)
}

// 0o1, 0o1_000, ...
func (p *Parser) parseOctal() ast.Node {
	tok := p.ExpectAndEat(token.TOctal)
	val, err := strconv.ParseInt(tok.Literal, 8, 64)
	if err != nil {
		errors.ThrowAtToken(tok, errors.ParserError, "invalid octal literal '%s'", tok.Literal)
	}
	return ast.NewInt(tok, val)
}

// 0b1, 0b1_000, ...
func (p *Parser) parseBinary() ast.Node {
	tok := p.ExpectAndEat(token.TBinary)
	val, err := strconv.ParseInt(tok.Literal, 2, 64)
	if err != nil {
		errors.ThrowAtToken(tok, errors.ParserError, "invalid binary literal '%s'", tok.Literal)
	}
	return ast.NewInt(tok, val)
}

// 1.0, 1.0_000, 1e10, ...
func (p *Parser) parseFloat() ast.Node {
	tok := p.ExpectAndEat(token.TFloat)
	val, err := strconv.ParseFloat(tok.Literal, 64)
	if err != nil {
		errors.ThrowAtToken(tok, errors.ParserError, "invalid float '%s'", tok.Literal)
	}
	return ast.NewFloat(tok, val)
}

// 'asdfasdf'
func (p *Parser) parseString() ast.Node {
	tok := p.ExpectAndEat(token.TString)
	return ast.NewString(tok, tok.Literal)
}

// true, false
func (p *Parser) parseBool() ast.Node {
	tok := p.ExpectAndEat(token.TTrue, token.TFalse)
	return ast.NewBool(tok, tok.Is(token.TTrue))
}

// <op><value-expr>
func (p *Parser) parseUnaryOp() ast.Node {
	tok := p.Eat()
	right := p.parseValueExpression(0)
	if !right.Has() {
		errors.ThrowAtToken(tok, errors.ParserError, "expected value expression after unary operator, but none was found")
	}
	return ast.NewUnaryOp(tok, tok.Literal, right.Unwrap())
}

// <value-expr><op><value-expr>
func (p *Parser) parseBinOp(left ast.Node) ast.Node {
	tok := p.Eat()
	right := p.parseValueExpression(p.ValuePrecedence(tok))
	if !right.Has() {
		errors.ThrowAtToken(tok, errors.ParserError, "expected value expression after binary operator, but none was found")
	}
	return ast.NewBinOp(tok, tok.Literal, left, right.Unwrap())
}

// { ... }
func (p *Parser) parseBlock() ast.Node {
	tok := p.ExpectAndEat(token.TLeftBrace)
	exprs := []ast.Node{}
	p.SkipNewlines()
	for !p.IsNext(token.TRightBrace) {
		node := p.parseValueExpression(0)
		if !node.Has() {
			errors.ThrowAtToken(p.Peek(), errors.ParserError, "expected value expression inside the block, but none was found")
		}
		exprs = append(exprs, node.Unwrap())
		p.SkipSeparator(token.TSemicolon)
	}
	p.ExpectAndEat(token.TRightBrace)
	return ast.NewBlock(tok, exprs)
}
