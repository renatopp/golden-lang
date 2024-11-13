// Parser for type expressions.
package syntax

import (
	"github.com/renatopp/golden/internal/compiler/syntax/ast"
	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/lang"
)

func (p *Parser) registerTypeExpressions() {
	p.TypeSolver.RegisterPrefixFn(core.TKeyword, p.parseTypeKeyword)
	p.TypeSolver.RegisterPrefixFn(core.TTypeIdent, p.parseTypeIdent)
	// p.TypeSolver.RegisterPrefixFn(core.TLparen, p.parseAnonymousDataDecl)
}

func (p *Parser) typePrecedence(t *lang.Token) int {
	return 0
}

// Nullable
func (p *Parser) parseTypeExpression(precedence ...int) *core.AstNode {
	pr := 0
	if len(precedence) > 0 {
		pr = precedence[0]
	}
	return p.TypeSolver.SolveExpression(p.Scanner, pr)
}

func (p *Parser) parseTypeKeyword() *core.AstNode {
	switch {
	// case p.IsNextLiteralsOf(core.TKeyword, core.KData):
	// 	return p.parseDataDecl()

	case p.IsNextLiteralsOf(core.TKeyword, core.KFN):
		return p.parseFunctionType()
	}

	errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected type expression keyword, got '%s'", p.PeekToken().Literal)
	return nil
}

// Parse data declaration. Example: `data ...`
// func (p *Parser) parseDataDecl() *core.AstNode {
// 	first := p.EatToken()

// 	p.ExpectTokens(core.TTypeIdent)
// 	ident := p.EatToken()

// 	constructors := []*DataConstructor{}
// 	switch {
// 	case p.IsNextTokens(TLparen):
// 		// parse declarations like `data Name(...)`
// 		p.EatToken()
// 		shape, fields := p.parseDataConstructorFields()
// 		p.ExpectTokens(TRparen)
// 		p.EatToken()
// 		constructors = append(constructors, &DataConstructor{
// 			Token:  ident,
// 			Name:   ident.Literal,
// 			Shape:  shape,
// 			Fields: fields,
// 		})

// 	case p.IsNextTokens(TAssign):
// 		// parse declarations like `data Name = ...`
// 		p.EatToken()
// 		constructors = p.parseDataConstructors()

// 	default:
// 		// parse declarations like `data Name`
// 		constructors = append(constructors, &DataConstructor{
// 			Token:  ident,
// 			Name:   ident.Literal,
// 			Shape:  "unit",
// 			Fields: []*DataConstructorField{},
// 		})
// 	}

// 	return NewNode(first, &AstDataDecl{
// 		Name:         ident.Literal,
// 		Constructors: constructors,
// 	})
// }

// Parse data constructor. Example: `A | B(Int) | Other`
// func (p *Parser) parseDataConstructors() []*DataConstructor {
// 	constructors := []*DataConstructor{}

// 	if p.IsNextTokens(TNewline) {
// 		p.SkipSeparator(TPipe)
// 	}

// 	for {
// 		p.ExpectTokens(TTypeIdent)
// 		ident := p.EatToken()

// 		shape := "unit"
// 		fields := []*DataConstructorField{}
// 		if p.IsNextTokens(TLparen) {
// 			p.EatToken()
// 			shape, fields = p.parseDataConstructorFields()
// 			p.ExpectTokens(TRparen)
// 			p.EatToken()
// 		}

// 		constructors = append(constructors, &DataConstructor{
// 			Token:  ident,
// 			Name:   ident.Literal,
// 			Shape:  shape,
// 			Fields: fields,
// 		})

// 		p.ExpectTokens(TPipe, TNewline, TEof)
// 		p.SkipNewlines()
// 		if !p.IsNextTokens(TPipe) {
// 			break
// 		}
// 		p.SkipSeparator(TPipe)
// 	}

// 	return constructors
// }

// Parse fields from data constructor, without parenthesis. Example: `Int, String` or `x Int, y String`
// func (p *Parser) parseDataConstructorFields() (shape string, fields []*DataConstructorField) {
// 	shape = ""
// 	fields = []*DataConstructorField{}

// 	p.SkipNewlines()
// 	switch {
// 	case p.IsNextTokens(TRparen):
// 		shape = "unit"

// 	case p.IsNextTokens(TVarIdent):
// 		shape = "record"

// 		for {
// 			ident := p.EatToken()
// 			tp := p.parseTypeExpression()
// 			if tp == nil {
// 				p.Error(p.PeekToken().Loc, "unexpected token", "expected type expression, got %s", p.PeekToken().Kind)
// 			}

// 			fields = append(fields, &DataConstructorField{
// 				Token: ident,
// 				Name:  ident.Literal,
// 				Type:  tp,
// 			})

// 			p.ExpectTokens(TComma, TNewline, TRparen)
// 			p.SkipSeparator(TComma)
// 			if p.IsNextTokens(TRparen) {
// 				break
// 			}
// 		}

// 	default:
// 		shape = "tuple"

// 		i := 0
// 		for {
// 			tp := p.parseTypeExpression()
// 			if tp == nil {
// 				p.Error(p.PeekToken().Loc, "unexpected token", "expected type expression, got %s", p.PeekToken().Kind)
// 			}

// 			fields = append(fields, &DataConstructorField{
// 				Token: tp.Token,
// 				Name:  strconv.Itoa(i),
// 				Type:  tp,
// 			})
// 			i++

// 			p.ExpectTokens(TComma, TNewline, TRparen)
// 			p.SkipSeparator(TComma)
// 			if p.IsNextTokens(TRparen) {
// 				break
// 			}
// 		}

// 	}

// 	return shape, fields
// }

// Parse function type. Example: `Fn (Int) String`
func (p *Parser) parseFunctionType() *core.AstNode {
	p.ExpectLiteralsOf(core.TKeyword, core.KFN)
	fn := p.EatToken()

	p.ExpectTokens(core.TLparen)
	p.EatToken()

	parameters := []*ast.FunctionTypeParam{}
	for {
		p.SkipNewlines()

		tp := p.parseTypeExpression()
		if tp == nil {
			break
		}

		parameters = append(parameters, &ast.FunctionTypeParam{
			Type: tp,
		})

		p.ExpectTokens(core.TComma, core.TRparen, core.TNewline)
		p.SkipSeparator(core.TComma)
	}

	p.ExpectTokens(core.TRparen)
	p.EatToken()

	tp := p.parseTypeExpression()
	return core.NewNode(fn, &ast.FunctionType{
		Params:     parameters,
		ReturnType: tp,
	})
}

// Parse data identifier.
func (p *Parser) parseTypeIdent() *core.AstNode {
	p.ExpectTokens(core.TTypeIdent)
	ident := p.EatToken()

	return core.NewNode(ident, &ast.TypeIdent{
		Name: ident.Literal,
	})
}

// func (p *Parser) parseAnonymousDataDecl() *core.AstNode {
// 	p.ExpectTokens(TLparen)
// 	p.EatToken()
// 	shape, fields := p.parseDataConstructorFields()
// 	p.ExpectTokens(TRparen)
// 	p.EatToken()

// 	return NewNode(nil, &AstDataDecl{
// 		Name: "",
// 		Constructors: []*DataConstructor{{
// 			Name:   "",
// 			Shape:  shape,
// 			Fields: fields,
// 		}},
// 	})
// }
