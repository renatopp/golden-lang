package syntax

import (
	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/token"
)

type peekFn func() token.Token
type precedenceFn func(token.Token) int
type prefixFn func() ast.Node
type infixFn func(ast.Node) ast.Node
type postfixFn func(ast.Node) ast.Node

type PrattSolver struct {
	peek       peekFn
	precedence precedenceFn
	prefixFns  map[token.TokenKind]prefixFn
	infixFns   map[token.TokenKind]infixFn
	postfixFns map[token.TokenKind]postfixFn
}

func NewPrattSolver(peek peekFn, prec precedenceFn) *PrattSolver {
	return &PrattSolver{
		peek:       peek,
		precedence: prec,
		prefixFns:  map[token.TokenKind]prefixFn{},
		infixFns:   map[token.TokenKind]infixFn{},
		postfixFns: map[token.TokenKind]postfixFn{},
	}
}

func (p *PrattSolver) RegisterPrefixFn(kind token.TokenKind, fn prefixFn) *PrattSolver {
	p.prefixFns[kind] = fn
	return p
}

func (p *PrattSolver) RegisterInfixFn(kind token.TokenKind, fn infixFn) *PrattSolver {
	p.infixFns[kind] = fn
	return p
}

func (p *PrattSolver) RegisterPostfixFn(kind token.TokenKind, fn postfixFn) *PrattSolver {
	p.postfixFns[kind] = fn
	return p
}

func (p *PrattSolver) SetPrecedenceFn(fn precedenceFn) *PrattSolver {
	p.precedence = fn
	return p
}

func (p *PrattSolver) SolveExpression(precedence int) ast.Node {
	prefix := p.prefixFns[p.peek().Kind]
	if prefix == nil {
		return nil
	}
	left := prefix()
	if left == nil {
		return nil
	}

	cur := p.peek()
	for {
		starting := cur

		if precedence < p.precedence(cur) {
			infix := p.infixFns[cur.Kind]
			if infix != nil {
				left = infix(left)
				cur = p.peek()
			}
		}

		for {
			postfix := p.postfixFns[cur.Kind]
			if postfix == nil {
				break
			}
			newLeft := postfix(left)
			if newLeft == nil {
				break
			}
			left = newLeft
			cur = p.peek()
		}

		// Didn't find any infix or postfix function
		if starting == cur {
			break
		}
	}

	return left
}
