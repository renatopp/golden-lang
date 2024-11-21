package syntax

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/tokens"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/safe"
)

func (p *parser) parseModule() *ast.Module {
	imports := []*ast.Import{}
	functions := []*ast.FuncDecl{}
	variables := []*ast.VarDecl{}

	first := p.PeekToken()
	p.Skip(tokens.TNewline)
	for {
		switch {
		case p.IsNextTokens(tokens.TImport):
			imports = append(imports, p.parserImport())

		// case p.IsNextTokens(tokens.TData):
		// types = append(types, p.parseTypeExpression())

		case p.IsNextTokens(tokens.TFn):
			functions = append(functions, p.parseValueExpression().(*ast.FuncDecl))

		case p.IsNextTokens(tokens.TLet):
			variables = append(variables, p.parseValueExpression().(*ast.VarDecl))

		case p.IsNextTokens(tokens.TComment):
			p.EatToken()

		case p.IsNextTokens(tokens.TEof):
			// pass
		default:
			tok := p.PeekToken()
			errors.ThrowAtToken(tok, errors.ParserError, "unexpected token '%s' at module level", tok.Literal)
		}

		p.Skip(tokens.TNewline)
		if p.IsNextTokens(tokens.TEof) {
			break
		}
	}

	return ast.NewModule(first, p.ModulePath, imports, functions, variables)
}

func (p *parser) parserImport() *ast.Import {
	p.ExpectTokens(tokens.TImport)
	first := p.EatToken()

	p.ExpectTokens(tokens.TString)
	tok := p.EatToken()
	path := ast.NewString(tok, tok.Literal)

	alias := safe.None[*ast.VarIdent]()
	if p.IsNextTokens(tokens.TAs) {
		p.EatToken()
		p.ExpectTokens(tokens.TVarIdent)
		ident := p.EatToken()
		alias = safe.Some(ast.NewVarIdent(ident, ident.Literal))
	}

	return ast.NewImport(first, path, alias)
}
