package lang

import (
	"fmt"
	"slices"
	"strings"
)

// Represents a token in the source code.
//
//	tokens.NewToken(EOF, "").WithChars(from, to)
type Token struct {
	Kind    string
	Literal string
	Loc     Loc
}

// Creates a new token.
func NewToken(t string, l string) *Token {
	return &Token{Kind: t, Literal: l}
}

// Creates a new error at the token location.
func (t *Token) AsError(kind, msg string) Error {
	return NewError(t.Loc, kind, msg)
}

// Creates a new token with the given char range. You can use this method as a
// chain.
func (t *Token) WithChars(from, to Char) *Token {
	t.Loc = Loc{Start: from.Span, End: to.Span}
	return t
}

// Creates a new token with the given type. You can use this method as a chain.
func (t *Token) WithType(tp string) *Token {
	t.Kind = tp
	return t
}

// Creates a new token with the given literal. You can use this method as a
// chain.
func (t *Token) WithLiteral(lit string) *Token {
	t.Literal = lit
	return t
}

// Returns true if the token is of the given type.
func (t *Token) IsType(tys ...string) bool {
	return slices.Contains(tys, t.Kind)
}

// Returns true if the token is the given literal.
func (t *Token) IsLiteral(lits ...string) bool {
	return slices.Contains(lits, t.Literal)
}

// Pretty string representation of the token.
func (t *Token) DebugString() string {
	literal := strings.ReplaceAll(t.Literal, "\n", "\\n")
	return fmt.Sprintf("[%d:%d] (%s) %s", t.Loc.Start.Line, t.Loc.Start.Column, t.Kind, literal)
}
