// Parser for value expressions.
package internal

import (
	"strconv"
	"strings"

	"github.com/renatopp/golden/lang"
)

func (p *parser) registerValueExpressions() {
	p.ValueSolver.RegisterPrefixFn(TKeyword, p.parseValueKeyword)
	p.ValueSolver.RegisterPrefixFn(TLbrace, p.parseBlock)
	p.ValueSolver.RegisterPrefixFn(TPlus, p.parseUnaryOperator)
	p.ValueSolver.RegisterPrefixFn(TMinus, p.parseUnaryOperator)
	p.ValueSolver.RegisterPrefixFn(TBang, p.parseUnaryOperator)
	p.ValueSolver.RegisterPrefixFn(TInteger, p.parseInteger)
	p.ValueSolver.RegisterPrefixFn(THex, p.parseInteger)
	p.ValueSolver.RegisterPrefixFn(TOctal, p.parseInteger)
	p.ValueSolver.RegisterPrefixFn(TBinary, p.parseInteger)
	p.ValueSolver.RegisterPrefixFn(TFloat, p.parseFloat)
	p.ValueSolver.RegisterPrefixFn(TBool, p.parseBool)
	p.ValueSolver.RegisterPrefixFn(TString, p.parseString)
	p.ValueSolver.RegisterPrefixFn(TVarIdent, p.parseVarIdent)
	p.ValueSolver.RegisterPrefixFn(TTypeIdent, p.parseTypeExpressionAsValue)
	p.ValueSolver.RegisterPrefixFn(TLparen, p.parseAnonymousDataApply)

	p.ValueSolver.RegisterInfixFn(TPlus, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TMinus, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TStar, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TSlash, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TPercent, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TSpaceship, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TEqual, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TNequal, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TAnd, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TOr, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TXor, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TLt, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TLte, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TGt, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TGte, p.parseBinaryOperator)
	p.ValueSolver.RegisterInfixFn(TDot, p.parseAccess)
	p.ValueSolver.RegisterInfixFn(TLparen, p.parseApply)

}

func (p *parser) valuePrecedence(t *lang.Token) int {
	switch {
	case t.IsKind(TAssign):
		return 10
	case t.IsKind(TPipe):
		return 20
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
	case t.IsKind(TDot):
		return 140
	}
	return 0
}

// Nullable
func (p *parser) parseValueExpression(precedence ...int) *Node {
	pr := 0
	if len(precedence) > 0 {
		pr = precedence[0]
	}
	return p.ValueSolver.SolveExpression(p.Scanner, pr)
}

func (p *parser) parseValueKeyword() *Node {
	switch {
	case p.IsNextLiteralsOf(TKeyword, KFn):
		return p.parseFunctionDecl()

	case p.IsNextLiteralsOf(TKeyword, KLet):
		return p.parseVariableDecl()
	}

	p.Error(p.PeekToken().Loc, "unexpected token", "expected value expression keyword, got %s", p.PeekToken().Kind)
	return nil
}

func (p *parser) parseVariableDecl() *Node {
	p.ExpectLiteralsOf(TKeyword, KLet)
	let := p.EatToken()

	p.ExpectTokens(TVarIdent)
	ident := p.EatToken()

	tp := p.parseTypeExpression()

	var value *Node
	if p.IsNextTokens(TAssign) {
		p.EatToken()
		value = p.parseValueExpression()
		if value == nil {
			p.Error(p.PeekToken().Loc, "unexpected token", "expected value expression")
		}
	}

	return NewNode(let, &AstVariableDecl{
		Name:  ident.Literal,
		Type:  tp,
		Value: value,
	})
}

func (p *parser) parseFunctionDecl() *Node {
	p.ExpectLiteralsOf(TKeyword, KFn)
	fn := p.EatToken()

	name := ""
	if p.IsNextTokens(TVarIdent) {
		name = p.EatToken().Literal
	}

	p.ExpectTokens(TLparen)
	p.EatToken()
	params := p.parseFunctionParams()
	p.ExpectTokens(TRparen)
	p.EatToken()

	tp := p.parseTypeExpression()

	p.ExpectTokens(TLbrace)
	body := p.parseBlock()

	node := NewNode(fn, &AstFunctionDecl{
		Name:       name,
		Params:     params,
		ReturnType: tp,
		Body:       body,
	})

	if name != "" {
		return NewNode(fn, &AstVariableDecl{
			Name:  name,
			Type:  nil,
			Value: node,
		})
	}

	return node
}

func (p *parser) parseFunctionParams() []*FunctionParam {
	params := []*FunctionParam{}

	p.SkipNewlines()
	for {
		if !p.IsNextTokens(TVarIdent) {
			break
		}

		name := p.EatToken().Literal
		tp := p.parseTypeExpression()
		if tp == nil {
			p.Error(p.PeekToken().Loc, "unexpected token", "expected type expression")
		}
		params = append(params, &FunctionParam{
			Name: name,
			Type: tp,
		})

		if !p.IsNextTokens(TComma, TNewline) {
			break
		}
		p.SkipSeparator(TComma)
	}

	p.SkipNewlines()
	return params
}

// Parse an integer literal with support for different bases.
func (p *parser) parseInteger() *Node {
	p.ExpectTokens(TInteger, THex, TOctal, TBinary)

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

// Parse a float literal.
func (p *parser) parseFloat() *Node {
	p.ExpectTokens(TFloat)
	token := p.EatToken()
	value, err := strconv.ParseFloat(token.Literal, 64)
	if err != nil {
		panic(lang.NewError(token.Loc, "invalid float", token.Literal))
	}
	return NewNode(token, &AstFloat{Value: value})
}

// Parse a boolean literal.
func (p *parser) parseBool() *Node {
	p.ExpectTokens(TBool)
	p.ExpectLiterals("true", "false")
	token := p.EatToken()
	value := token.Literal == "true"
	return NewNode(token, &AstBool{Value: value})
}

// Parse a string literal.
func (p *parser) parseString() *Node {
	p.ExpectTokens(TString)
	token := p.EatToken()
	value := strings.ReplaceAll(token.Literal, "\r", "")
	return NewNode(token, &AstString{Value: value})
}

// Parse a variable identifier.
func (p *parser) parseVarIdent() *Node {
	p.ExpectTokens(TVarIdent)
	token := p.EatToken()
	return NewNode(token, &AstVarIdent{Name: token.Literal})
}

// Parse a block expression. Example: `{ ... }`
func (p *parser) parseBlock() *Node {
	p.ExpectTokens(TLbrace)
	lbrace := p.EatToken()

	nodes := []*Node{}
	p.SkipNewlines()
	for !p.IsNextTokens(TRbrace) {
		node := p.parseValueExpression()
		if node == nil {
			p.Error(p.PeekToken().Loc, "unexpected token", "expected value expression")
		}
		nodes = append(nodes, node)

		p.SkipSeparator(TSemicolon)
	}

	p.ExpectTokens(TRbrace)
	p.EatToken()

	return NewNode(lbrace, &AstBlock{Expressions: nodes})
}

// Parse a unary operator. Example: `+x`
func (p *parser) parseUnaryOperator() *Node {
	op := p.EatToken()
	right := p.parseValueExpression(p.valuePrecedence(op))
	return NewNode(op, &AstUnaryOp{
		Operator: op.Literal,
		Right:    right,
	})
}

func (p *parser) parseTypeExpressionAsValue() *Node {
	tp := p.parseTypeExpression()
	if tp == nil {
		p.Error(p.PeekToken().Loc, "unexpected token", "expected type expression")
	}
	return NewNode(tp.Token, &AstApply{
		Shape:  "unit",
		Args:   []*ApplyArgument{},
		Target: tp,
	})
}

// Parse a binary operator expression. Example: `x + y`
func (p *parser) parseBinaryOperator(left *Node) *Node {
	op := p.EatToken()
	right := p.parseValueExpression(p.valuePrecedence(op))

	if right == nil {
		panic(lang.NewError(op.Loc, "expecting expression", ""))
	}

	return NewNode(op, &AstBinaryOp{
		Operator: op.Literal,
		Left:     left,
		Right:    right,
	})
}

// Parse an assignment expression. Example: `x = y`
func (p *parser) parseAccess(left *Node) *Node {
	op := p.EatToken()
	p.ExpectTokens(TVarIdent, TTypeIdent, TInteger)
	accessor := p.EatToken()
	return NewNode(op, &AstAccess{
		Target:   left,
		Accessor: accessor.Literal,
	})
}

// Parse an application. Example: `f(x, y)` or `F(x=2, y=3)`
func (p *parser) parseApply(left *Node) *Node {
	p.ExpectTokens(TLparen)
	first := p.EatToken()
	shape, args := p.parseApplyArguments()
	p.ExpectTokens(TRparen)
	p.EatToken()
	return NewNode(first, &AstApply{
		Shape:  shape,
		Target: left,
		Args:   args,
	})
}

// Parse an anonymous data application. Example: `(2, 4)`
func (p *parser) parseAnonymousDataApply() *Node {
	p.ExpectTokens(TLparen)
	lparen := p.EatToken()
	shape, args := p.parseApplyArguments()
	p.ExpectTokens(TRparen)
	p.EatToken()

	return NewNode(lparen, &AstApply{
		Shape: shape,
		Args:  args,
	})
}

// Parse a sequence of type application arguments. Example: `1, 2` or `x=2, s=2`
func (p *parser) parseApplyArguments() (shape string, args []*ApplyArgument) {
	args = []*ApplyArgument{}
	p.SkipNewlines()

	n0 := p.PeekToken()
	n1 := p.PeekTokenAt(1)
	switch {
	case n0.IsKind(TRparen):
		shape = "unit"

	case n0.IsKind(TVarIdent) && n1.IsKind(TAssign):
		shape = "record"
		for {
			p.ExpectTokens(TVarIdent)
			name := p.EatToken()
			p.ExpectTokens(TAssign)
			p.EatToken()
			expr := p.parseValueExpression()
			if expr == nil {
				p.Error(p.PeekToken().Loc, "unexpected token", "expected value expression")
			}
			args = append(args, &ApplyArgument{
				Token: name,
				Name:  name.Literal,
				Value: expr,
			})

			if !p.IsNextTokens(TComma, TNewline) {
				break
			}
			p.SkipSeparator(TComma)
		}

	default:
		shape = "tuple"
		for {
			first := p.PeekToken()
			expr := p.parseValueExpression()
			if expr == nil {
				break
			}
			args = append(args, &ApplyArgument{
				Token: first,
				Value: expr,
			})

			if !p.IsNextTokens(TComma, TNewline) {
				break
			}
			p.SkipSeparator(TComma)
		}
	}

	p.SkipNewlines()
	return shape, args
}
