// Parser for value expressions.
package syntax

import (
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/lang"
)

func (p *parser) valuePrecedence(t *lang.Token) int {
	switch {
	case t.IsKind(core.TAssign):
		return 10
	case t.IsKind(core.TPipe):
		return 20
	case t.IsKind(core.TOr):
		return 40
	case t.IsKind(core.TXor):
		return 45
	case t.IsKind(core.TAnd):
		return 50
	case t.IsKind(core.TEqual, core.TNequal):
		return 70
	case t.IsKind(core.TLt, core.TGt, core.TLte, core.TGte):
		return 80
	case t.IsKind(core.TPlus, core.TMinus):
		return 90
	case t.IsKind(core.TStar, core.TSlash):
		return 100
	case t.IsKind(core.TSpaceship):
		return 110
	case t.IsKind(core.TPercent):
		return 120
	case t.IsKind(core.TLparen):
		return 130
	case t.IsKind(core.TDot):
		return 140
	}
	return 0
}

func (p *parser) registerValueExpressions() {
	// p.ValueSolver.RegisterPrefixFn(core.TKeyword, p.parseValueKeyword)
	// p.ValueSolver.RegisterPrefixFn(core.TLbrace, p.parseBlock)
	// p.ValueSolver.RegisterPrefixFn(core.TPlus, p.parseUnaryOperator)
	// p.ValueSolver.RegisterPrefixFn(core.TMinus, p.parseUnaryOperator)
	// p.ValueSolver.RegisterPrefixFn(core.TBang, p.parseUnaryOperator)
	// p.ValueSolver.RegisterPrefixFn(core.TInteger, p.parseInteger)
	// p.ValueSolver.RegisterPrefixFn(core.THex, p.parseInteger)
	// p.ValueSolver.RegisterPrefixFn(core.TOctal, p.parseInteger)
	// p.ValueSolver.RegisterPrefixFn(core.TBinary, p.parseInteger)
	// p.ValueSolver.RegisterPrefixFn(core.TFloat, p.parseFloat)
	// p.ValueSolver.RegisterPrefixFn(core.TBool, p.parseBool)
	// p.ValueSolver.RegisterPrefixFn(core.TString, p.parseString)
	// p.ValueSolver.RegisterPrefixFn(core.TVarIdent, p.parseVarIdent)
	// p.ValueSolver.RegisterPrefixFn(core.TTypeIdent, p.parseTypeExpressionAsValue)
	// p.ValueSolver.RegisterPrefixFn(core.TLparen, p.parseAnonymousDataApply)

	// p.ValueSolver.RegisterInfixFn(core.TPlus, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TMinus, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TStar, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TSlash, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TPercent, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TSpaceship, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TEqual, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TNequal, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TAnd, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TOr, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TXor, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TLt, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TLte, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TGt, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TGte, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(core.TDot, p.parseAccess)
	// p.ValueSolver.RegisterInfixFn(core.TLparen, p.parseApply)
}

// // Nullable
// func (p *Parser) parseValueExpression(precedence ...int) *core.AstNode {
// 	pr := 0
// 	if len(precedence) > 0 {
// 		pr = precedence[0]
// 	}
// 	return p.ValueSolver.SolveExpression(p.Scanner, pr)
// }

// func (p *Parser) parseValueKeyword() *core.AstNode {
// 	switch {
// 	case p.IsNextLiteralsOf(core.TKeyword, core.KFn):
// 		return p.parseFunctionDecl()

// 	case p.IsNextLiteralsOf(core.TKeyword, core.KLet):
// 		return p.parseVariableDecl()
// 	}

// 	errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected value expression keyword, got '%s' instead", p.PeekToken().Literal)
// 	return nil
// }

// func (p *Parser) parseVariableDecl() *core.AstNode {
// 	p.ExpectLiteralsOf(core.TKeyword, core.KLet)
// 	let := p.EatToken()

// 	p.ExpectTokens(core.TVarIdent)
// 	ident := p.EatToken()

// 	tp := p.parseTypeExpression()

// 	var value *core.AstNode
// 	if p.IsNextTokens(core.TAssign) {
// 		p.EatToken()
// 		value = p.parseValueExpression()
// 		if value == nil {
// 			errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected value expression, but none was found")
// 		}
// 	}

// 	return core.NewNode(let, &ast.VariableDecl{
// 		Name:  ident.Literal,
// 		Type:  tp,
// 		Value: value,
// 	})
// }

// func (p *Parser) parseFunctionDecl() *core.AstNode {
// 	p.ExpectLiteralsOf(core.TKeyword, core.KFn)
// 	fn := p.EatToken()

// 	name := ""
// 	if p.IsNextTokens(core.TVarIdent) {
// 		name = p.EatToken().Literal
// 	}

// 	p.ExpectTokens(core.TLparen)
// 	p.EatToken()
// 	params := p.parseFunctionParams()
// 	p.ExpectTokens(core.TRparen)
// 	p.EatToken()

// 	tp := p.parseTypeExpression()

// 	p.ExpectTokens(core.TLbrace)
// 	body := p.parseBlock()

// 	node := core.NewNode(fn, &ast.FunctionDecl{
// 		Name:       name,
// 		Params:     params,
// 		ReturnType: tp,
// 		Body:       body,
// 	})

// 	return node
// }

// func (p *Parser) parseFunctionParams() []*ast.FunctionDeclParam {
// 	params := []*ast.FunctionDeclParam{}

// 	p.SkipNewlines()
// 	for {
// 		if !p.IsNextTokens(core.TVarIdent) {
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

// 		if !p.IsNextTokens(core.TComma, core.TNewline) {
// 			break
// 		}
// 		p.SkipSeparator(core.TComma)
// 	}

// 	p.SkipNewlines()
// 	return params
// }

// // Parse an integer literal with support for different bases.
// func (p *Parser) parseInteger() *core.AstNode {
// 	p.ExpectTokens(core.TInteger, core.THex, core.TOctal, core.TBinary)

// 	token := p.EatToken()
// 	base := 10
// 	switch token.Kind {
// 	case core.THex:
// 		base = 16
// 	case core.TOctal:
// 		base = 8
// 	case core.TBinary:
// 		base = 2
// 	}

// 	value, err := strconv.ParseInt(token.Literal, base, 64)
// 	if err != nil {
// 		errors.ThrowAtToken(token, errors.ParserError, "invalid integer literal '%s'", token.Literal)
// 	}

// 	return core.NewNode(token, &ast.Int{Value: value})
// }

// // Parse a float literal.
// func (p *Parser) parseFloat() *core.AstNode {
// 	p.ExpectTokens(core.TFloat)
// 	token := p.EatToken()
// 	value, err := strconv.ParseFloat(token.Literal, 64)
// 	if err != nil {
// 		errors.ThrowAtToken(token, errors.ParserError, "invalid float literal '%s'", token.Literal)
// 	}
// 	return core.NewNode(token, &ast.Float{Value: value})
// }

// // Parse a boolean literal.
// func (p *Parser) parseBool() *core.AstNode {
// 	p.ExpectTokens(core.TBool)
// 	p.ExpectLiterals("true", "false")
// 	token := p.EatToken()
// 	value := token.Literal == "true"
// 	return core.NewNode(token, &ast.Bool{Value: value})
// }

// // Parse a string literal.
// func (p *Parser) parseString() *core.AstNode {
// 	p.ExpectTokens(core.TString)
// 	token := p.EatToken()
// 	value := strings.ReplaceAll(token.Literal, "\r", "")
// 	return core.NewNode(token, &ast.String{Value: value})
// }

// // Parse a variable identifier.
// func (p *Parser) parseVarIdent() *core.AstNode {
// 	p.ExpectTokens(core.TVarIdent)
// 	token := p.EatToken()
// 	return core.NewNode(token, &ast.VarIdent{Name: token.Literal})
// }

// // Parse a block expression. Example: `{ ... }`
// func (p *Parser) parseBlock() *core.AstNode {
// 	p.ExpectTokens(core.TLbrace)
// 	lbrace := p.EatToken()

// 	nodes := []*core.AstNode{}
// 	p.SkipNewlines()
// 	for !p.IsNextTokens(core.TRbrace) {
// 		node := p.parseValueExpression()
// 		if node == nil {
// 			errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected value expression, but none was found")
// 		}
// 		nodes = append(nodes, node)

// 		p.SkipSeparator(core.TSemicolon)
// 	}

// 	p.ExpectTokens(core.TRbrace)
// 	p.EatToken()

// 	return core.NewNode(lbrace, &ast.Block{Expressions: nodes})
// }

// // Parse a unary operator. Example: `+x`
// func (p *Parser) parseUnaryOperator() *core.AstNode {
// 	op := p.EatToken()
// 	right := p.parseValueExpression(p.valuePrecedence(op))
// 	return core.NewNode(op, &ast.UnaryOp{
// 		Operator: op.Literal,
// 		Right:    right,
// 	})
// }

// func (p *Parser) parseTypeExpressionAsValue() *core.AstNode {
// 	tp := p.parseTypeExpression()
// 	if tp == nil {
// 		errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected type expression, but none was found")
// 	}
// 	return core.NewNode(tp.Token(), &ast.Apply{
// 		Shape:  "unit",
// 		Args:   []*ast.ApplyArgument{},
// 		Target: tp,
// 	})
// }

// // Parse a binary operator expression. Example: `x + y`
// func (p *Parser) parseBinaryOperator(left *core.AstNode) *core.AstNode {
// 	op := p.EatToken()
// 	right := p.parseValueExpression(p.valuePrecedence(op))

// 	if right == nil {
// 		errors.ThrowAtToken(op, errors.ParserError, "expecting value expression after operator, but none was found")
// 	}

// 	return core.NewNode(op, &ast.BinaryOp{
// 		Operator: op.Literal,
// 		Left:     left,
// 		Right:    right,
// 	})
// }

// // Parse an assignment expression. Example: `x = y`
// func (p *Parser) parseAccess(left *core.AstNode) *core.AstNode {
// 	op := p.EatToken()
// 	p.ExpectTokens(core.TVarIdent, core.TTypeIdent, core.TInteger)
// 	accessor := p.EatToken()
// 	return core.NewNode(op, &ast.Access{
// 		Target:   left,
// 		Accessor: accessor.Literal,
// 	})
// }

// // Parse an application. Example: `f(x, y)` or `F(x=2, y=3)`
// func (p *Parser) parseApply(left *core.AstNode) *core.AstNode {
// 	p.ExpectTokens(core.TLparen)
// 	first := p.EatToken()
// 	shape, args := p.parseApplyArguments()
// 	p.ExpectTokens(core.TRparen)
// 	p.EatToken()
// 	return core.NewNode(first, &ast.Apply{
// 		Shape:  shape,
// 		Target: left,
// 		Args:   args,
// 	})
// }

// // Parse an anonymous data application. Example: `(2, 4)`
// func (p *Parser) parseAnonymousDataApply() *core.AstNode {
// 	p.ExpectTokens(core.TLparen)
// 	lparen := p.EatToken()
// 	shape, args := p.parseApplyArguments()
// 	p.ExpectTokens(core.TRparen)
// 	p.EatToken()

// 	return core.NewNode(lparen, &ast.Apply{
// 		Shape: shape,
// 		Args:  args,
// 	})
// }

// // Parse a sequence of type application arguments. Example: `1, 2` or `x=2, s=2`
// func (p *Parser) parseApplyArguments() (shape string, args []*ast.ApplyArgument) {
// 	args = []*ast.ApplyArgument{}
// 	p.SkipNewlines()

// 	n0 := p.PeekToken()
// 	n1 := p.PeekTokenAt(1)
// 	switch {
// 	case n0.IsKind(core.TRparen):
// 		shape = "unit"

// 	case n0.IsKind(core.TVarIdent) && n1.IsKind(core.TAssign):
// 		shape = "record"
// 		for {
// 			p.ExpectTokens(core.TVarIdent)
// 			name := p.EatToken()
// 			p.ExpectTokens(core.TAssign)
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

// 			if !p.IsNextTokens(core.TComma, core.TNewline) {
// 				break
// 			}
// 			p.SkipSeparator(core.TComma)
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

// 			if !p.IsNextTokens(core.TComma, core.TNewline) {
// 				break
// 			}
// 			p.SkipSeparator(core.TComma)
// 		}
// 	}

// 	p.SkipNewlines()
// 	return shape, args
// }
