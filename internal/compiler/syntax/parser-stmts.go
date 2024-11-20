package syntax

import (
	"os"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/tokens"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/safe"
)

func (p *parser) parseModule() *ast.Module {
	imports := []*ast.Import{}

	first := p.PeekToken()
	p.Skip(tokens.TNewline)
	for {
		switch {
		case p.IsNextTokens(tokens.TImport):
			stmt := p.parserImport()
			imports = append(imports, stmt)

		case p.IsNextTokens(tokens.TData):
			println("TODO: data")
			os.Exit(0)
		case p.IsNextTokens(tokens.TFn):
			println("TODO: fn")
			os.Exit(0)
		case p.IsNextTokens(tokens.TLet):
			println("TODO: let")
			os.Exit(0)
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

	return ast.NewModule(first, p.ModulePath, imports)
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
