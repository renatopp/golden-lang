// Parser for value expressions.
package syntax

import (
	"github.com/renatopp/golden/internal/compiler/tokens"
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
	// p.ValueSolver.RegisterPrefixFn(tokens.TKeyword, p.parseValueKeyword)
	// p.ValueSolver.RegisterPrefixFn(tokens.TLbrace, p.parseBlock)
	// p.ValueSolver.RegisterPrefixFn(tokens.TPlus, p.parseUnaryOperator)
	// p.ValueSolver.RegisterPrefixFn(tokens.TMinus, p.parseUnaryOperator)
	// p.ValueSolver.RegisterPrefixFn(tokens.TBang, p.parseUnaryOperator)
	// p.ValueSolver.RegisterPrefixFn(tokens.TInteger, p.parseInteger)
	// p.ValueSolver.RegisterPrefixFn(tokens.THex, p.parseInteger)
	// p.ValueSolver.RegisterPrefixFn(tokens.TOctal, p.parseInteger)
	// p.ValueSolver.RegisterPrefixFn(tokens.TBinary, p.parseInteger)
	// p.ValueSolver.RegisterPrefixFn(tokens.TFloat, p.parseFloat)
	// p.ValueSolver.RegisterPrefixFn(tokens.TBool, p.parseBool)
	// p.ValueSolver.RegisterPrefixFn(tokens.TString, p.parseString)
	// p.ValueSolver.RegisterPrefixFn(tokens.TVarIdent, p.parseVarIdent)
	// p.ValueSolver.RegisterPrefixFn(tokens.TTypeIdent, p.parseTypeExpressionAsValue)
	// p.ValueSolver.RegisterPrefixFn(tokens.TLparen, p.parseAnonymousDataApply)

	// p.ValueSolver.RegisterInfixFn(tokens.TPlus, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TMinus, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TStar, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TSlash, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TPercent, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TSpaceship, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TEqual, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TNequal, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TAnd, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TOr, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TXor, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TLt, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TLte, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TGt, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TGte, p.parseBinaryOperator)
	// p.ValueSolver.RegisterInfixFn(tokens.TDot, p.parseAccess)
	// p.ValueSolver.RegisterInfixFn(tokens.TLparen, p.parseApply)
}

// // Nullable
// func (p *Parser) parseValueExpression(precedence ...int) *tokens.AstNode {
// 	pr := 0
// 	if len(precedence) > 0 {
// 		pr = precedence[0]
// 	}
// 	return p.ValueSolver.SolveExpression(p.Scanner, pr)
// }

// func (p *Parser) parseValueKeyword() *tokens.AstNode {
// 	switch {
// 	case p.IsNextLiteralsOf(tokens.TKeyword, tokens.KFn):
// 		return p.parseFunctionDecl()

// 	case p.IsNextLiteralsOf(tokens.TKeyword, tokens.KLet):
// 		return p.parseVariableDecl()
// 	}

// 	errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected value expression keyword, got '%s' instead", p.PeekToken().Literal)
// 	return nil
// }

// func (p *Parser) parseVariableDecl() *tokens.AstNode {
// 	p.ExpectLiteralsOf(tokens.TKeyword, tokens.KLet)
// 	let := p.EatToken()

// 	p.ExpectTokens(tokens.TVarIdent)
// 	ident := p.EatToken()

// 	tp := p.parseTypeExpression()

// 	var value *tokens.AstNode
// 	if p.IsNextTokens(tokens.TAssign) {
// 		p.EatToken()
// 		value = p.parseValueExpression()
// 		if value == nil {
// 			errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected value expression, but none was found")
// 		}
// 	}

// 	return tokens.NewNode(let, &ast.VariableDecl{
// 		Name:  ident.Literal,
// 		Type:  tp,
// 		Value: value,
// 	})
// }

// func (p *Parser) parseFunctionDecl() *tokens.AstNode {
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

// func (p *Parser) parseFunctionParams() []*ast.FunctionDeclParam {
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

// // Parse an integer literal with support for different bases.
// func (p *Parser) parseInteger() *tokens.AstNode {
// 	p.ExpectTokens(tokens.TInteger, tokens.THex, tokens.TOctal, tokens.TBinary)

// 	token := p.EatToken()
// 	base := 10
// 	switch token.Kind {
// 	case tokens.THex:
// 		base = 16
// 	case tokens.TOctal:
// 		base = 8
// 	case tokens.TBinary:
// 		base = 2
// 	}

// 	value, err := strconv.ParseInt(token.Literal, base, 64)
// 	if err != nil {
// 		errors.ThrowAtToken(token, errors.ParserError, "invalid integer literal '%s'", token.Literal)
// 	}

// 	return tokens.NewNode(token, &ast.Int{Value: value})
// }

// // Parse a float literal.
// func (p *Parser) parseFloat() *tokens.AstNode {
// 	p.ExpectTokens(tokens.TFloat)
// 	token := p.EatToken()
// 	value, err := strconv.ParseFloat(token.Literal, 64)
// 	if err != nil {
// 		errors.ThrowAtToken(token, errors.ParserError, "invalid float literal '%s'", token.Literal)
// 	}
// 	return tokens.NewNode(token, &ast.Float{Value: value})
// }

// // Parse a boolean literal.
// func (p *Parser) parseBool() *tokens.AstNode {
// 	p.ExpectTokens(tokens.TBool)
// 	p.ExpectLiterals("true", "false")
// 	token := p.EatToken()
// 	value := token.Literal == "true"
// 	return tokens.NewNode(token, &ast.Bool{Value: value})
// }

// // Parse a string literal.
// func (p *Parser) parseString() *tokens.AstNode {
// 	p.ExpectTokens(tokens.TString)
// 	token := p.EatToken()
// 	value := strings.ReplaceAll(token.Literal, "\r", "")
// 	return tokens.NewNode(token, &ast.String{Value: value})
// }

// // Parse a variable identifier.
// func (p *Parser) parseVarIdent() *tokens.AstNode {
// 	p.ExpectTokens(tokens.TVarIdent)
// 	token := p.EatToken()
// 	return tokens.NewNode(token, &ast.VarIdent{Name: token.Literal})
// }

// // Parse a block expression. Example: `{ ... }`
// func (p *Parser) parseBlock() *tokens.AstNode {
// 	p.ExpectTokens(tokens.TLbrace)
// 	lbrace := p.EatToken()

// 	nodes := []*tokens.AstNode{}
// 	p.SkipNewlines()
// 	for !p.IsNextTokens(tokens.TRbrace) {
// 		node := p.parseValueExpression()
// 		if node == nil {
// 			errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected value expression, but none was found")
// 		}
// 		nodes = append(nodes, node)

// 		p.SkipSeparator(tokens.TSemicolon)
// 	}

// 	p.ExpectTokens(tokens.TRbrace)
// 	p.EatToken()

// 	return tokens.NewNode(lbrace, &ast.Block{Expressions: nodes})
// }

// // Parse a unary operator. Example: `+x`
// func (p *Parser) parseUnaryOperator() *tokens.AstNode {
// 	op := p.EatToken()
// 	right := p.parseValueExpression(p.valuePrecedence(op))
// 	return tokens.NewNode(op, &ast.UnaryOp{
// 		Operator: op.Literal,
// 		Right:    right,
// 	})
// }

// func (p *Parser) parseTypeExpressionAsValue() *tokens.AstNode {
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

// // Parse a binary operator expression. Example: `x + y`
// func (p *Parser) parseBinaryOperator(left *tokens.AstNode) *tokens.AstNode {
// 	op := p.EatToken()
// 	right := p.parseValueExpression(p.valuePrecedence(op))

// 	if right == nil {
// 		errors.ThrowAtToken(op, errors.ParserError, "expecting value expression after operator, but none was found")
// 	}

// 	return tokens.NewNode(op, &ast.BinaryOp{
// 		Operator: op.Literal,
// 		Left:     left,
// 		Right:    right,
// 	})
// }

// // Parse an assignment expression. Example: `x = y`
// func (p *Parser) parseAccess(left *tokens.AstNode) *tokens.AstNode {
// 	op := p.EatToken()
// 	p.ExpectTokens(tokens.TVarIdent, tokens.TTypeIdent, tokens.TInteger)
// 	accessor := p.EatToken()
// 	return tokens.NewNode(op, &ast.Access{
// 		Target:   left,
// 		Accessor: accessor.Literal,
// 	})
// }

// // Parse an application. Example: `f(x, y)` or `F(x=2, y=3)`
// func (p *Parser) parseApply(left *tokens.AstNode) *tokens.AstNode {
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
// func (p *Parser) parseAnonymousDataApply() *tokens.AstNode {
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
// func (p *Parser) parseApplyArguments() (shape string, args []*ast.ApplyArgument) {
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
