// Parser for type expressions.
package internal

import (
	"strconv"

	"github.com/renatopp/golden/lang"
)

func (p *parser) registerTypeExpressions() {
	p.TypeSolver.RegisterPrefixFn(TKeyword, p.parseTypeKeyword)
	p.TypeSolver.RegisterPrefixFn(TTypeIdent, p.parseTypeIdent)
	p.TypeSolver.RegisterPrefixFn(TLparen, p.parseAnonymousDataDecl)
}

func (p *parser) typePrecedence(t *lang.Token) int {
	return 0
}

// Nullable
func (p *parser) parseTypeExpression(precedence ...int) *Node {
	pr := 0
	if len(precedence) > 0 {
		pr = precedence[0]
	}
	return p.TypeSolver.SolveExpression(p.Scanner, pr)
}

func (p *parser) parseTypeKeyword() *Node {
	switch {
	case p.IsNextLiteralsOf(TKeyword, KData):
		return p.parseDataDecl()

	case p.IsNextLiteralsOf(TKeyword, KFN):
		return p.parseFunctionType()
	}

	p.Error(p.PeekToken().Loc, "unexpected token", "expected type expression keyword, got %s", p.PeekToken().Kind)
	return nil
}

// Parse data declaration. Example: `data ...`
func (p *parser) parseDataDecl() *Node {
	first := p.EatToken()

	p.ExpectTokens(TTypeIdent)
	ident := p.EatToken()

	constructors := []*DataConstructor{}
	switch {
	case p.IsNextTokens(TLparen):
		// parse declarations like `data Name(...)`
		p.EatToken()
		shape, fields := p.parseDataConstructorFields()
		p.ExpectTokens(TRparen)
		p.EatToken()
		constructors = append(constructors, &DataConstructor{
			Token:  ident,
			Name:   ident.Literal,
			Shape:  shape,
			Fields: fields,
		})

	case p.IsNextTokens(TAssign):
		// parse declarations like `data Name = ...`
		p.EatToken()
		constructors = p.parseDataConstructors()

	default:
		// parse declarations like `data Name`
		constructors = append(constructors, &DataConstructor{
			Token:  ident,
			Name:   ident.Literal,
			Shape:  "unit",
			Fields: []*DataConstructorField{},
		})
	}

	return NewNode(first, &AstDataDecl{
		Name:         ident.Literal,
		Constructors: constructors,
	})
}

// Parse data constructor. Example: `A | B(Int) | Other`
func (p *parser) parseDataConstructors() []*DataConstructor {
	constructors := []*DataConstructor{}

	if p.IsNextTokens(TNewline) {
		p.SkipSeparator(TPipe)
	}

	for {
		p.ExpectTokens(TTypeIdent)
		ident := p.EatToken()

		shape := "unit"
		fields := []*DataConstructorField{}
		if p.IsNextTokens(TLparen) {
			p.EatToken()
			shape, fields = p.parseDataConstructorFields()
			p.ExpectTokens(TRparen)
			p.EatToken()
		}

		constructors = append(constructors, &DataConstructor{
			Token:  ident,
			Name:   ident.Literal,
			Shape:  shape,
			Fields: fields,
		})

		p.ExpectTokens(TPipe, TNewline, TEof)
		p.SkipNewlines()
		if !p.IsNextTokens(TPipe) {
			break
		}
		p.SkipSeparator(TPipe)
	}

	return constructors
}

// Parse fields from data constructor, without parenthesis. Example: `Int, String` or `x Int, y String`
func (p *parser) parseDataConstructorFields() (shape string, fields []*DataConstructorField) {
	shape = ""
	fields = []*DataConstructorField{}

	p.SkipNewlines()
	switch {
	case p.IsNextTokens(TRparen):
		shape = "unit"

	case p.IsNextTokens(TVarIdent):
		shape = "record"

		for {
			ident := p.EatToken()
			tp := p.parseTypeExpression()
			if tp == nil {
				p.Error(p.PeekToken().Loc, "unexpected token", "expected type expression, got %s", p.PeekToken().Kind)
			}

			fields = append(fields, &DataConstructorField{
				Token: ident,
				Name:  ident.Literal,
				Type:  tp,
			})

			p.ExpectTokens(TComma, TNewline, TRparen)
			p.SkipSeparator(TComma)
			if p.IsNextTokens(TRparen) {
				break
			}
		}

	default:
		shape = "tuple"

		i := 0
		for {
			tp := p.parseTypeExpression()
			if tp == nil {
				p.Error(p.PeekToken().Loc, "unexpected token", "expected type expression, got %s", p.PeekToken().Kind)
			}

			fields = append(fields, &DataConstructorField{
				Token: tp.Token,
				Name:  strconv.Itoa(i),
				Type:  tp,
			})
			i++

			p.ExpectTokens(TComma, TNewline, TRparen)
			p.SkipSeparator(TComma)
			if p.IsNextTokens(TRparen) {
				break
			}
		}

	}

	return shape, fields
}

// Parse function type. Example: `Fn (Int) String`
func (p *parser) parseFunctionType() *Node {
	p.ExpectLiteralsOf(TKeyword, KFN)
	fn := p.EatToken()

	p.ExpectTokens(TLparen)
	p.EatToken()

	parameters := []*FunctionTypeParam{}
	for {
		p.SkipNewlines()

		tp := p.parseTypeExpression()
		if tp == nil {
			break
		}

		parameters = append(parameters, &FunctionTypeParam{
			Type: tp,
		})

		p.ExpectTokens(TComma, TRparen, TNewline)
		p.SkipSeparator(TComma)
	}

	p.ExpectTokens(TRparen)
	p.EatToken()

	tp := p.parseTypeExpression()
	return NewNode(fn, &AstFunctionType{
		Params:     parameters,
		ReturnType: tp,
	})
}

// Parse data identifier.
func (p *parser) parseTypeIdent() *Node {
	p.ExpectTokens(TTypeIdent)
	ident := p.EatToken()

	return NewNode(ident, &AstTypeIdent{
		Name: ident.Literal,
	})
}

func (p *parser) parseAnonymousDataDecl() *Node {
	p.ExpectTokens(TLparen)
	p.EatToken()
	shape, fields := p.parseDataConstructorFields()
	p.ExpectTokens(TRparen)
	p.EatToken()

	return NewNode(nil, &AstDataDecl{
		Name: "",
		Constructors: []*DataConstructor{{
			Name:   "",
			Shape:  shape,
			Fields: fields,
		}},
	})
}
