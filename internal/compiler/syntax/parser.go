package syntax

import (
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

	// p.TypeSolver.RegisterInfixFn(token.TAssign, nil)
	// p.ValueSolver.RegisterInfixFn(token.T)

	return p
}

func (p *Parser) Parse() (res ast.Node, err error) {
	err = errors.WithRecovery(func() {
		res = p.parse()
	})
	return res, err
}

func (p *Parser) parse() ast.Node {
	//
	return nil
}

// func (p *Parser) parseModule() ast.Module {}

// func (p *Parser) parseInt() ast.Int {}
// func (p *Parser) parseHex() ast.Int {}
// func (p *Parser) parseOctal() ast.Int {}
// func (p *Parser) parseBin() ast.Int {}
// func (p *Parser) parseFloat() ast.Float {}
// func (p *Parser) parseString() ast.String {}
// func (p *Parser) parseBool() ast.Bool {}
// func (p *Parser) parseLeftBrace() ast.Block {}
// func (p *Parser) parseUnaryOp() ast.UnaryOp {}

// func (p *Parser) parseBinOp(left ast.Node) ast.BinOp {}
