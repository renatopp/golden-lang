package syntax

import (
	"strconv"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/token"
	"github.com/renatopp/golden/internal/helpers/errors"
)

type Parser struct {
	*BaseParser
}

func NewParser(tokens []token.Token) *Parser {
	p := &Parser{
		BaseParser: NewBaseParser(tokens),
	}

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

func (p *Parser) Parse() (res ast.Module, err error) {
	err = errors.WithRecovery(func() {
		res = p.parseModule()
	})
	return res, err
}

func (p *Parser) parseModule() ast.Module {
	consts := []ast.Const{}
	for {
		p.Skip(token.TNewline, token.TComment)
		if p.Peek().Is(token.TEof) {
			break
		}
		consts = append(consts, p.parseConst())
	}
	return ast.Module{
		Consts: consts,
	}
}

// const <var-ident> <type-expr>? = <value-expr>
func (p *Parser) parseConst() ast.Const {
	// const
	p.Expect(token.TConst)
	tok := p.Eat()

	// var-ident
	name := p.parseVarIdent()

	// type-expr

	// =
	p.Expect(token.TAssign)
	p.Eat()

	// value-expr
	value := p.parseValueExpression(0)
	return ast.Const{
		Token:     tok,
		Name:      name,
		ValueExpr: value,
	}
}

func (p *Parser) parseValueExpression(prec int) ast.Node {
	return p.ValueSolver.SolveExpression(prec)
}

func (p *Parser) parseVarIdent() ast.VarIdent {
	p.Expect(token.TVarIdent)
	tok := p.Eat()
	return ast.VarIdent{
		Token: tok,
		Value: tok.Value,
	}
}

func (p *Parser) parseInt() ast.Node {
	p.Expect(token.TInt)
	tok := p.Eat()
	val, err := strconv.ParseInt(tok.Value, 10, 64)
	if err != nil {
		errors.ThrowAtToken(tok, errors.ParserError, "invalid integer '%s'", tok.Value)
	}
	return ast.Int{
		Token: tok,
		Value: val,
	}
}

func (p *Parser) parseHex() ast.Node {
	p.Expect(token.THex)
	tok := p.Eat()
	val, err := strconv.ParseInt(tok.Value, 16, 64)
	if err != nil {
		errors.ThrowAtToken(tok, errors.ParserError, "invalid hex '%s'", tok.Value)
	}
	return ast.Int{
		Token: tok,
		Value: val,
	}
}

func (p *Parser) parseOctal() ast.Node {
	p.Expect(token.TOctal)
	tok := p.Eat()
	val, err := strconv.ParseInt(tok.Value, 8, 64)
	if err != nil {
		errors.ThrowAtToken(tok, errors.ParserError, "invalid octal '%s'", tok.Value)
	}
	return ast.Int{
		Token: tok,
		Value: val,
	}
}

func (p *Parser) parseBinary() ast.Node {
	p.Expect(token.TBinary)
	tok := p.Eat()
	val, err := strconv.ParseInt(tok.Value, 2, 64)
	if err != nil {
		errors.ThrowAtToken(tok, errors.ParserError, "invalid bin '%s'", tok.Value)
	}
	return ast.Int{
		Token: tok,
		Value: val,
	}
}

func (p *Parser) parseFloat() ast.Node {
	p.Expect(token.TFloat)
	tok := p.Eat()
	val, err := strconv.ParseFloat(tok.Value, 64)
	if err != nil {
		errors.ThrowAtToken(tok, errors.ParserError, "invalid float '%s'", tok.Value)
	}
	return ast.Float{
		Token: tok,
		Value: val,
	}
}

func (p *Parser) parseString() ast.Node {
	p.Expect(token.TString)
	tok := p.Eat()
	return ast.String{
		Token: tok,
		Value: tok.Value,
	}
}

func (p *Parser) parseBool() ast.Node {
	p.Expect(token.TTrue, token.TFalse)
	tok := p.Eat()
	val, err := strconv.ParseBool(tok.Value)
	if err != nil {
		errors.ThrowAtToken(tok, errors.ParserError, "invalid bool '%s'", tok.Value)
	}
	return ast.Bool{
		Token: tok,
		Value: val,
	}
}

func (p *Parser) parseUnaryOp() ast.Node {
	tok := p.Eat()
	right := p.parseValueExpression(0)
	return ast.UnaryOp{
		Token: tok,
		Op:    tok.Value,
		Right: right,
	}
}

func (p *Parser) parseBinOp(left ast.Node) ast.Node {
	tok := p.Eat()
	right := p.parseValueExpression(p.ValuePrecedence(tok))
	return ast.BinOp{
		Token: tok,
		Op:    tok.Value,
		Left:  left,
		Right: right,
	}
}

func (p *Parser) parseBlock() ast.Node {
	p.Expect(token.TLeftBrace)

	tok := p.Eat()
	nodes := []ast.Node{}
	p.SkipNewlines()
	for !p.IsNext(token.TRightBrace) {
		node := p.parseValueExpression(0)
		if node == nil {
			errors.ThrowAtToken(p.Peek(), errors.ParserError, "expected value expression, but none was found")
		}
		nodes = append(nodes, node)
		p.SkipSeparator(token.TSemicolon)
	}

	p.Expect(token.TRightBrace)
	p.Eat()
	return ast.Block{
		Token:       tok,
		Expressions: nodes,
	}
}
