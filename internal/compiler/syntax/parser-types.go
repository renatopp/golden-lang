// Parser for type expressions.
package syntax

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/tokens"
	"github.com/renatopp/golden/internal/helpers/safe"
	"github.com/renatopp/golden/lang"
)

func (p *parser) typePrecedence(t *lang.Token) int {
	return 0
}

func (p *parser) registerTypeExpressions() {
	p.TypeSolver.RegisterPrefixFn(tokens.TTypeIdent, p.parseTypeIdent)
	p.TypeSolver.RegisterPrefixFn(tokens.TFN, p.parseFunctionType)
	// p.TypeSolver.RegisterPrefixFn(tokens.TLparen, p.parseAnonymousDataDecl)
}

// Nullable
func (p *parser) parseTypeExpression(precedence ...int) ast.Node {
	pr := 0
	if len(precedence) > 0 {
		pr = precedence[0]
	}
	return p.TypeSolver.SolveExpression(p.Scanner, pr)
}

// Parse type identifier.
func (p *parser) parseTypeIdent() ast.Node {
	p.ExpectTokens(tokens.TTypeIdent)
	ident := p.EatToken()

	node := ast.NewTypeIdent(ident, ident.Literal)
	node.SetExpressionKind(ast.TypeExpressionKind)
	return node
}

// Parse function type. Example: `Fn (Int) String`
func (p *parser) parseFunctionType() ast.Node {
	p.ExpectTokens(tokens.TFN)
	fn := p.EatToken()

	p.ExpectTokens(tokens.TLparen)
	p.EatToken()

	parameters := []*ast.FuncTypeParam{}
	for {
		p.SkipNewlines()

		tp := p.parseTypeExpression()
		if tp == nil {
			break
		}
		parameters = append(parameters, ast.NewFuncTypeParam(tp.Token(), len(parameters), tp))

		p.ExpectTokens(tokens.TComma, tokens.TRparen, tokens.TNewline)
		p.SkipSeparator(tokens.TComma)
	}

	p.ExpectTokens(tokens.TRparen)
	p.EatToken()

	ret := safe.None[ast.Node]()
	tok := p.parseTypeExpression()
	if tok != nil {
		ret = safe.Some(tok)
	}
	node := ast.NewFuncType(fn, parameters, ret)
	node.SetExpressionKind(ast.TypeExpressionKind)
	return node
}

// func (p *parser) parseTypeKeyword() ast.Node {
// 	switch {
// 	// case p.IsNextLiteralsOf(tokens.TKeyword, core.KData):
// 	// 	return p.parseDataDecl()

// 	case p.IsNextLiteralsOf(tokens.TKeyword, core.KFN):
// 		return p.parseFunctionType()
// 	}

// 	errors.ThrowAtToken(p.PeekToken(), errors.ParserError, "expected type expression keyword, got '%s'", p.PeekToken().Literal)
// 	return nil
// }

// // Parse data declaration. Example: `data ...`
// // func (p *parser) parseDataDecl() ast.Node {
// // 	first := p.EatToken()

// // 	p.ExpectTokens(tokens.TTypeIdent)
// // 	ident := p.EatToken()

// // 	constructors := []*DataConstructor{}
// // 	switch {
// // 	case p.IsNextTokens(TLparen):
// // 		// parse declarations like `data Name(...)`
// // 		p.EatToken()
// // 		shape, fields := p.parseDataConstructorFields()
// // 		p.ExpectTokens(TRparen)
// // 		p.EatToken()
// // 		constructors = append(constructors, &DataConstructor{
// // 			Token:  ident,
// // 			Name:   ident.Literal,
// // 			Shape:  shape,
// // 			Fields: fields,
// // 		})

// // 	case p.IsNextTokens(TAssign):
// // 		// parse declarations like `data Name = ...`
// // 		p.EatToken()
// // 		constructors = p.parseDataConstructors()

// // 	default:
// // 		// parse declarations like `data Name`
// // 		constructors = append(constructors, &DataConstructor{
// // 			Token:  ident,
// // 			Name:   ident.Literal,
// // 			Shape:  "unit",
// // 			Fields: []*DataConstructorField{},
// // 		})
// // 	}

// // 	return NewNode(first, &AstDataDecl{
// // 		Name:         ident.Literal,
// // 		Constructors: constructors,
// // 	})
// // }

// // Parse data constructor. Example: `A | B(Int) | Other`
// // func (p *parser) parseDataConstructors() []*DataConstructor {
// // 	constructors := []*DataConstructor{}

// // 	if p.IsNextTokens(TNewline) {
// // 		p.SkipSeparator(TPipe)
// // 	}

// // 	for {
// // 		p.ExpectTokens(TTypeIdent)
// // 		ident := p.EatToken()

// // 		shape := "unit"
// // 		fields := []*DataConstructorField{}
// // 		if p.IsNextTokens(TLparen) {
// // 			p.EatToken()
// // 			shape, fields = p.parseDataConstructorFields()
// // 			p.ExpectTokens(TRparen)
// // 			p.EatToken()
// // 		}

// // 		constructors = append(constructors, &DataConstructor{
// // 			Token:  ident,
// // 			Name:   ident.Literal,
// // 			Shape:  shape,
// // 			Fields: fields,
// // 		})

// // 		p.ExpectTokens(TPipe, TNewline, TEof)
// // 		p.SkipNewlines()
// // 		if !p.IsNextTokens(TPipe) {
// // 			break
// // 		}
// // 		p.SkipSeparator(TPipe)
// // 	}

// // 	return constructors
// // }

// // Parse fields from data constructor, without parenthesis. Example: `Int, String` or `x Int, y String`
// // func (p *parser) parseDataConstructorFields() (shape string, fields []*DataConstructorField) {
// // 	shape = ""
// // 	fields = []*DataConstructorField{}

// // 	p.SkipNewlines()
// // 	switch {
// // 	case p.IsNextTokens(TRparen):
// // 		shape = "unit"

// // 	case p.IsNextTokens(TVarIdent):
// // 		shape = "record"

// // 		for {
// // 			ident := p.EatToken()
// // 			tp := p.parseTypeExpression()
// // 			if tp == nil {
// // 				p.Error(p.PeekToken().Loc, "unexpected token", "expected type expression, got %s", p.PeekToken().Kind)
// // 			}

// // 			fields = append(fields, &DataConstructorField{
// // 				Token: ident,
// // 				Name:  ident.Literal,
// // 				Type:  tp,
// // 			})

// // 			p.ExpectTokens(TComma, TNewline, TRparen)
// // 			p.SkipSeparator(TComma)
// // 			if p.IsNextTokens(TRparen) {
// // 				break
// // 			}
// // 		}

// // 	default:
// // 		shape = "tuple"

// // 		i := 0
// // 		for {
// // 			tp := p.parseTypeExpression()
// // 			if tp == nil {
// // 				p.Error(p.PeekToken().Loc, "unexpected token", "expected type expression, got %s", p.PeekToken().Kind)
// // 			}

// // 			fields = append(fields, &DataConstructorField{
// // 				Token: tp.Token,
// // 				Name:  strconv.Itoa(i),
// // 				Type:  tp,
// // 			})
// // 			i++

// // 			p.ExpectTokens(TComma, TNewline, TRparen)
// // 			p.SkipSeparator(TComma)
// // 			if p.IsNextTokens(TRparen) {
// // 				break
// // 			}
// // 		}

// // 	}

// // 	return shape, fields
// // }

// // func (p *parser) parseAnonymousDataDecl() ast.Node {
// // 	p.ExpectTokens(TLparen)
// // 	p.EatToken()
// // 	shape, fields := p.parseDataConstructorFields()
// // 	p.ExpectTokens(TRparen)
// // 	p.EatToken()

// // 	return NewNode(nil, &AstDataDecl{
// // 		Name: "",
// // 		Constructors: []*DataConstructor{{
// // 			Name:   "",
// // 			Shape:  shape,
// // 			Fields: fields,
// // 		}},
// // 	})
// // }
