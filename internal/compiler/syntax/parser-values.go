// Parser for value expressions.
package syntax

import (
	"strconv"
	"strings"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/tokens"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/safe"
	"github.com/renatopp/golden/lang"
)

func (p *parser) valuePrecedence(t *lang.Token) int {
	switch {
	case t.IsKind(tokens.TAssign):
		return 10
	case t.IsKind(tokens.TPipe):
		return 20
	case t.IsKind(tokens.TOr):
		return 40
	case t.IsKind(tokens.TXor):
		return 45
	case t.IsKind(tokens.TAnd):
		return 50
	case t.IsKind(tokens.TEqual, tokens.TNequal):
		return 70
	case t.IsKind(tokens.TLt, tokens.TGt, tokens.TLte, tokens.TGte):
		return 80
	case t.IsKind(tokens.TPlus, tokens.TMinus):
		return 90
	case t.IsKind(tokens.TStar, tokens.TSlash):
		return 100
	case t.IsKind(tokens.TSpaceship):
		return 110
	case t.IsKind(tokens.TPercent):
		return 120
	case t.IsKind(tokens.TLparen):
		return 130
	case t.IsKind(tokens.TDot):
		return 140
	}
	return 0
}

func (p *parser) registerValueExpressions() {
	p.ValueSolver.RegisterPrefixFn(tokens.TLet, p.parseVarDecl)
	p.ValueSolver.RegisterPrefixFn(tokens.TInteger, p.parseInt)
	p.ValueSolver.RegisterPrefixFn(tokens.THex, p.parseInt)
	p.ValueSolver.RegisterPrefixFn(tokens.TOctal, p.parseInt)
	p.ValueSolver.RegisterPrefixFn(tokens.TBinary, p.parseInt)
	p.ValueSolver.RegisterPrefixFn(tokens.TFloat, p.parseFloat)
	p.ValueSolver.RegisterPrefixFn(tokens.TBool, p.parseBool)
	p.ValueSolver.RegisterPrefixFn(tokens.TString, p.parseString)
	p.ValueSolver.RegisterPrefixFn(tokens.TVarIdent, p.parseVarIdent)
	p.ValueSolver.RegisterPrefixFn(tokens.TLbrace, p.parseBlock)
	p.ValueSolver.RegisterPrefixFn(tokens.TPlus, p.parseUnaryOp)
	p.ValueSolver.RegisterPrefixFn(tokens.TMinus, p.parseUnaryOp)
	p.ValueSolver.RegisterPrefixFn(tokens.TBang, p.parseUnaryOp)
	// p.ValueSolver.RegisterPrefixFn(tokens.TTypeIdent, p.parseTypeExpressionAsValue)
	// p.ValueSolver.RegisterPrefixFn(tokens.TLparen, p.parseAnonymousDataApply)

	p.ValueSolver.RegisterInfixFn(tokens.TPlus, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TMinus, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TStar, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TSlash, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TPercent, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TSpaceship, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TEqual, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TNequal, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TAnd, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TOr, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TXor, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TLt, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TLte, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TGt, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TGte, p.parseBinaryOp)
	p.ValueSolver.RegisterInfixFn(tokens.TDot, p.parseAccess)
	// p.ValueSolver.RegisterInfixFn(tokens.TLparen, p.parseApply)
}

// // Nullable
func (p *parser) parseValueExpression(precedence ...int) ast.Node {
	pr := 0
	if len(precedence) > 0 {
		pr = precedence[0]
	}
	return p.ValueSolver.SolveExpression(p.Scanner, pr)
}

func (p *parser) parseVarDecl() ast.Node {
	p.ExpectTokens(tokens.TLet)
	let := p.EatToken()

	p.ExpectTokens(tokens.TVarIdent)
	toc := p.EatToken()
	name := ast.NewVarIdent(toc, toc.Literal)

	tp := safe.None[ast.Node]()
	tpExpr := p.parseTypeExpression()
	if tpExpr != nil {
		tp = safe.Some(tpExpr)
	}

	val := safe.None[ast.Node]()
	if p.IsNextTokens(tokens.TAssign) {
		p.EatToken()
		valExpr := p.parseValueExpression()
		if valExpr == nil {
			errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected value expression, but none was found")
		}
		val = safe.Some(valExpr)
	}

	return ast.NewVarDecl(let, name, tp, val)
}

// Parse an integer literal with support for different bases.
func (p *parser) parseInt() ast.Node {
	p.ExpectTokens(tokens.TInteger, tokens.THex, tokens.TOctal, tokens.TBinary)

	token := p.EatToken()
	base := 10
	switch token.Kind {
	case tokens.THex:
		base = 16
	case tokens.TOctal:
		base = 8
	case tokens.TBinary:
		base = 2
	}

	value, err := strconv.ParseInt(token.Literal, base, 64)
	if err != nil {
		errors.ThrowAtToken(token, errors.ParserError, "invalid integer literal '%s'", token.Literal)
	}
	return ast.NewInt(token, value)
}

// Parse a float literal.
func (p *parser) parseFloat() ast.Node {
	p.ExpectTokens(tokens.TFloat)
	token := p.EatToken()
	value, err := strconv.ParseFloat(token.Literal, 64)
	if err != nil {
		errors.ThrowAtToken(token, errors.ParserError, "invalid float literal '%s'", token.Literal)
	}
	return ast.NewFloat(token, value)
}

// Parse a boolean literal.
func (p *parser) parseBool() ast.Node {
	p.ExpectTokens(tokens.TBool)
	token := p.EatToken()
	value := token.Literal == "true"
	return ast.NewBool(token, value)
}

// Parse a string literal.
func (p *parser) parseString() ast.Node {
	p.ExpectTokens(tokens.TString)
	token := p.EatToken()
	value := strings.ReplaceAll(token.Literal, "\r", "")
	return ast.NewString(token, value)
}

// Parse a variable identifier.
func (p *parser) parseVarIdent() ast.Node {
	p.ExpectTokens(tokens.TVarIdent)
	token := p.EatToken()
	return ast.NewVarIdent(token, token.Literal)
}

// Parse a block expression. Example: `{ ... }`
func (p *parser) parseBlock() ast.Node {
	p.ExpectTokens(tokens.TLbrace)
	lbrace := p.EatToken()

	nodes := []ast.Node{}
	p.SkipNewlines()
	for !p.IsNextTokens(tokens.TRbrace) {
		node := p.parseValueExpression()
		if node == nil {
			errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected value expression, but none was found")
		}
		nodes = append(nodes, node)
		p.SkipSeparator(tokens.TSemicolon)
	}

	p.ExpectTokens(tokens.TRbrace)
	p.EatToken()
	return ast.NewBlock(lbrace, nodes)
}

// Parse a unary operator. Example: `+x`
func (p *parser) parseUnaryOp() ast.Node {
	op := p.EatToken()
	right := p.parseValueExpression(p.valuePrecedence(op))
	return ast.NewUnaryOp(op, op.Literal, right)
}

// Parse a binary operator expression. Example: `x + y`
func (p *parser) parseBinaryOp(left ast.Node) ast.Node {
	op := p.EatToken()
	right := p.parseValueExpression(p.valuePrecedence(op))
	if right == nil {
		errors.ThrowAtToken(op, errors.ParserError, "expecting value expression after operator, but none was found")
	}
	return ast.NewBinaryOp(op, op.Literal, left, right)
}

// Parse an assignment expression. Example: `x = y`
func (p *parser) parseAccess(left ast.Node) ast.Node {
	op := p.EatToken()
	p.ExpectTokens(tokens.TVarIdent, tokens.TTypeIdent, tokens.TInteger)
	tok := p.EatToken()
	accessor := ast.NewVarIdent(tok, tok.Literal)
	return ast.NewAccess(op, left, accessor)
}

// func (p *parser) parseValueKeyword() *tokens.AstNode {
// 	switch {
// 	case p.IsNextLiteralsOf(tokens.TKeyword, tokens.KFn):
// 		return p.parseFunctionDecl()

// 	case p.IsNextLiteralsOf(tokens.TKeyword, tokens.KLet):
// 		return p.parseVariableDecl()
// 	}

// 	errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected value expression keyword, got '%s' instead", p.PeekToken().Literal)
// 	return nil
// }

// func (p *parser) parseFunctionDecl() *tokens.AstNode {
// 	p.ExpectLiteralsOf(tokens.TKeyword, tokens.KFn)
// 	fn := p.EatToken()

// 	name := ""
// 	if p.IsNextTokens(tokens.TVarIdent) {
// 		name = p.EatToken().Literal
// 	}

// 	p.ExpectTokens(tokens.TLparen)
// 	p.EatToken()
// 	params := p.parseFunctionParams()
// 	p.ExpectTokens(tokens.TRparen)
// 	p.EatToken()

// 	tp := p.parseTypeExpression()

// 	p.ExpectTokens(tokens.TLbrace)
// 	body := p.parseBlock()

// 	node := tokens.NewNode(fn, &ast.FunctionDecl{
// 		Name:       name,
// 		Params:     params,
// 		ReturnType: tp,
// 		Body:       body,
// 	})

// 	return node
// }

// func (p *parser) parseFunctionParams() []*ast.FunctionDeclParam {
// 	params := []*ast.FunctionDeclParam{}

// 	p.SkipNewlines()
// 	for {
// 		if !p.IsNextTokens(tokens.TVarIdent) {
// 			break
// 		}

// 		name := p.EatToken().Literal
// 		tp := p.parseTypeExpression()
// 		if tp == nil {
// 			errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected type expression, but none was found")
// 		}
// 		params = append(params, &ast.FunctionDeclParam{
// 			Name: name,
// 			Type: tp,
// 		})

// 		if !p.IsNextTokens(tokens.TComma, tokens.TNewline) {
// 			break
// 		}
// 		p.SkipSeparator(tokens.TComma)
// 	}

// 	p.SkipNewlines()
// 	return params
// }

// func (p *parser) parseTypeExpressionAsValue() *tokens.AstNode {
// 	tp := p.parseTypeExpression()
// 	if tp == nil {
// 		errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected type expression, but none was found")
// 	}
// 	return tokens.NewNode(tp.Token(), &ast.Apply{
// 		Shape:  "unit",
// 		Args:   []*ast.ApplyArgument{},
// 		Target: tp,
// 	})
// }

// // Parse an application. Example: `f(x, y)` or `F(x=2, y=3)`
// func (p *parser) parseApply(left *tokens.AstNode) *tokens.AstNode {
// 	p.ExpectTokens(tokens.TLparen)
// 	first := p.EatToken()
// 	shape, args := p.parseApplyArguments()
// 	p.ExpectTokens(tokens.TRparen)
// 	p.EatToken()
// 	return tokens.NewNode(first, &ast.Apply{
// 		Shape:  shape,
// 		Target: left,
// 		Args:   args,
// 	})
// }

// // Parse an anonymous data application. Example: `(2, 4)`
// func (p *parser) parseAnonymousDataApply() *tokens.AstNode {
// 	p.ExpectTokens(tokens.TLparen)
// 	lparen := p.EatToken()
// 	shape, args := p.parseApplyArguments()
// 	p.ExpectTokens(tokens.TRparen)
// 	p.EatToken()

// 	return tokens.NewNode(lparen, &ast.Apply{
// 		Shape: shape,
// 		Args:  args,
// 	})
// }

// // Parse a sequence of type application arguments. Example: `1, 2` or `x=2, s=2`
// func (p *parser) parseApplyArguments() (shape string, args []*ast.ApplyArgument) {
// 	args = []*ast.ApplyArgument{}
// 	p.SkipNewlines()

// 	n0 := p.PeekToken()
// 	n1 := p.PeekTokenAt(1)
// 	switch {
// 	case n0.IsKind(tokens.TRparen):
// 		shape = "unit"

// 	case n0.IsKind(tokens.TVarIdent) && n1.IsKind(tokens.TAssign):
// 		shape = "record"
// 		for {
// 			p.ExpectTokens(tokens.TVarIdent)
// 			name := p.EatToken()
// 			p.ExpectTokens(tokens.TAssign)
// 			p.EatToken()
// 			expr := p.parseValueExpression()
// 			if expr == nil {
// 				errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected value expression, but none was found")
// 			}
// 			args = append(args, &ast.ApplyArgument{
// 				Token: name,
// 				Name:  name.Literal,
// 				Value: expr,
// 			})

// 			if !p.IsNextTokens(tokens.TComma, tokens.TNewline) {
// 				break
// 			}
// 			p.SkipSeparator(tokens.TComma)
// 		}

// 	default:
// 		shape = "tuple"
// 		for {
// 			first := p.PeekToken()
// 			expr := p.parseValueExpression()
// 			if expr == nil {
// 				break
// 			}
// 			args = append(args, &ast.ApplyArgument{
// 				Token: first,
// 				Value: expr,
// 			})

// 			if !p.IsNextTokens(tokens.TComma, tokens.TNewline) {
// 				break
// 			}
// 			p.SkipSeparator(tokens.TComma)
// 		}
// 	}

// 	p.SkipNewlines()
// 	return shape, args
// }
