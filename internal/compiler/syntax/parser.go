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
	p.ValueSolver.RegisterPrefixFn(token.TLeftParen, p.parseParen)
	p.ValueSolver.RegisterPrefixFn(token.TFn, p.parseFn)

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
	p.ValueSolver.RegisterInfixFn(token.TLeftParen, p.parseApplication)

	p.TypeSolver.RegisterPrefixFn(token.TTypeIdent, p.parseTypeIdentType)
	p.TypeSolver.RegisterPrefixFn(token.TFN, p.parseFnType)

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

func (p *Parser) parseTypeExpression(prec int) safe.Optional[ast.Node] {
	return p.TypeSolver.SolveExpression(prec)
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

		switch p.Peek().Kind {
		case token.TLet:
			exprs = append(exprs, p.parseLet())
		case token.TFn:
			exprs = append(exprs, p.parseFn())
		default:
			errors.ThrowAtToken(p.Peek(), errors.ParserError, "unexpected token '%s'", p.Peek().Literal)
		}

		p.SkipSeparator(token.TSemicolon)
	}
	return ast.NewModule(first, exprs)
}

// let <var-ident> <type-expr>? = <value-expr>
func (p *Parser) parseLet() *ast.VarDecl {
	tok := p.ExpectAndEat(token.TLet)         // let
	name := p.parseVarIdent().(*ast.VarIdent) // var-ident
	// TODO: add type expression parsing // type-expr
	p.ExpectAndEat(token.TAssign)    // =
	val := p.parseValueExpression(0) // value-expr
	if !val.Has() {
		p.ThrowExpectedValueExpression("after assignment")
	}
	return ast.NewVarDecl(tok, name, safe.None[ast.Node](), val.Unwrap())
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
		p.ThrowExpectedValueExpression("after unary operator '%s'", tok.Literal)
	}
	return ast.NewUnaryOp(tok, tok.Literal, right.Unwrap())
}

// (<value-expr>)
func (p *Parser) parseParen() ast.Node {
	p.ExpectAndEat(token.TLeftParen)
	p.SkipNewlines()
	node := p.parseValueExpression(0)
	if !node.Has() {
		p.ThrowExpectedValueExpression("inside the parentheses")
	}
	p.SkipNewlines()
	p.ExpectAndEat(token.TRightParen)
	return node.Unwrap()
}

// <value-expr><op><value-expr>
func (p *Parser) parseBinOp(left ast.Node) ast.Node {
	tok := p.Eat()
	right := p.parseValueExpression(p.ValuePrecedence(tok))
	if !right.Has() {
		p.ThrowExpectedValueExpression("after binary operator '%s'", tok.Literal)
	}
	return ast.NewBinOp(tok, tok.Literal, left, right.Unwrap())
}

// { ... }
func (p *Parser) parseBlock() ast.Node {
	tok := p.ExpectAndEat(token.TLeftBrace)
	exprs := []ast.Node{}
	p.SkipNewlines()
	for !p.IsNext(token.TRightBrace) {
		if p.IsNext(token.TReturn) {
			exprs = append(exprs, p.parseReturn())
		} else {
			node := p.parseValueExpression(0)
			if !node.Has() {
				p.ThrowExpectedValueExpression("inside the block")
			}
			exprs = append(exprs, node.Unwrap())
		}
		p.SkipSeparator(token.TSemicolon)
	}
	p.ExpectAndEat(token.TRightBrace)
	return ast.NewBlock(tok, exprs)
}

// fn <var-ident>?(<var-ident> <type-expr>, ...):<type-expr> = ...
func (p *Parser) parseFn() ast.Node {
	tok := p.ExpectAndEat(token.TFn)

	name := safe.None[*ast.VarIdent]()
	if p.IsNext(token.TVarIdent) {
		name = safe.Some(p.parseVarIdent().(*ast.VarIdent))
	}

	params := []*ast.FnDeclParam{}
	if p.IsNext(token.TLeftParen) {
		params = p.parseFnParams()
	}

	var returnExpr ast.Node = ast.NewTypeIdent(p.Peek(), "Void")
	if expr := p.parseTypeExpression(0); expr.Has() {
		returnExpr = expr.Unwrap()
	}

	p.SkipNewlines()
	p.Expect(token.TLeftBrace)
	val := p.parseBlock().(*ast.Block)
	return ast.NewFnDecl(tok, name, params, returnExpr, val)
}

// (<var-ident> <type-expr>, ...)
func (p *Parser) parseFnParams() []*ast.FnDeclParam {
	params := []*ast.FnDeclParam{}
	p.ExpectAndEat(token.TLeftParen)
	for {
		if p.IsNext(token.TRightParen) {
			break
		}
		p.Expect(token.TVarIdent)
		name := p.parseVarIdent().(*ast.VarIdent)
		typeExpr := p.parseTypeExpression(0)
		params = append(params, ast.NewFnDeclParam(name, typeExpr.Or(nil)))
		p.SkipSeparator(token.TComma)
	}
	last := p.ExpectAndEat(token.TRightParen)

	// treat backfilling of type expressions
	if len(params) > 0 {
		lastNode := params[len(params)-1]
		lastType := lastNode.TypeExpr
		if lastType == nil {
			errors.ThrowAtToken(last, errors.ParserError, "expected type expression after parameter name, but none was found")
			return nil
		}
		for i := len(params) - 1; i >= 0; i-- {
			p := params[i]
			if p.TypeExpr == nil {
				params[i].TypeExpr = lastType
			} else {
				lastType = p.TypeExpr
			}
		}
	}

	return params
}

func (p *Parser) parseReturn() ast.Node {
	tok := p.ExpectAndEat(token.TReturn)
	value := p.parseValueExpression(0)
	return ast.NewReturn(tok, value)
}

//
//
//

// Int, Float, ...
func (p *Parser) parseTypeIdentType() ast.Node {
	tok := p.ExpectAndEat(token.TTypeIdent)
	return ast.NewTypeIdent(tok, tok.Literal)
}

// Fn(<type-expr>, ...):<type-expr>
func (p *Parser) parseFnType() ast.Node {
	tok := p.ExpectAndEat(token.TFN)

	params := []ast.Node{}
	if p.IsNext(token.TLeftParen) {
		params = p.parseFnTypeParams()
	}

	var returnExpr ast.Node = ast.NewTypeIdent(p.Peek(), "Void")
	if expr := p.parseTypeExpression(0); expr.Has() {
		returnExpr = expr.Unwrap()
	}

	return ast.NewTypeFn(tok, params, returnExpr)
}

// (<type-expr>, ...)
func (p *Parser) parseFnTypeParams() []ast.Node {
	params := []ast.Node{}
	p.ExpectAndEat(token.TLeftParen)
	for {
		if p.IsNext(token.TRightParen) {
			break
		}
		typeExpr := p.parseTypeExpression(0)
		if !typeExpr.Has() {
			errors.ThrowAtToken(p.Peek(), errors.ParserError, "expected type expression, but none was found")
		}
		params = append(params, typeExpr.Unwrap())
		p.SkipSeparator(token.TComma)
	}
	p.ExpectAndEat(token.TRightParen)
	return params
}

// <target>(<value-expr>, ...)
func (p *Parser) parseApplication(left ast.Node) ast.Node {
	tok := p.ExpectAndEat(token.TLeftParen)
	args := []ast.Node{}
	for {
		if p.IsNext(token.TRightParen) {
			break
		}
		arg := p.parseValueExpression(0)
		if !arg.Has() {
			p.ThrowExpectedValueExpression("as argument")
		}
		args = append(args, arg.Unwrap())
		p.SkipSeparator(token.TComma)
	}
	p.ExpectAndEat(token.TRightParen)
	return ast.NewApplication(tok, left, args)
}
