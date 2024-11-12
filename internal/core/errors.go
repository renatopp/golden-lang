package core

import (
	"fmt"

	"github.com/renatopp/golden/lang"
)

type ErrorKind uint64

type Error struct {
	loc   *lang.Loc
	token *lang.Token
	node  *AstNode
	kind  ErrorKind
	msg   string
	stack string
}

func NewError(kind ErrorKind, msg string, args ...any) *Error {
	return &Error{
		kind: kind,
		msg:  fmt.Sprintf(msg, args...),
	}
}

func NewEmptyError() *Error {
	return &Error{}
}

func (e *Error) Loc() *lang.Loc     { return e.loc }
func (e *Error) Token() *lang.Token { return e.token }
func (e *Error) Node() *AstNode     { return e.node }
func (e *Error) Kind() ErrorKind    { return e.kind }
func (e *Error) Message() string    { return e.msg }
func (e *Error) Stack() string      { return e.stack }
func (e *Error) Error() string      { return e.msg }

func (e *Error) WithLoc(loc *lang.Loc) *Error {
	e.loc = loc
	return e
}

func (e *Error) WithToken(token *lang.Token) *Error {
	e.token = token
	e.loc = &token.Loc
	return e
}

func (e *Error) WithNode(node *AstNode) *Error {
	e.node = node
	e.token = node.Token()
	e.loc = &e.token.Loc
	return e
}

func (e *Error) WithKind(kind ErrorKind) *Error {
	e.kind = kind
	return e
}

func (e *Error) WithMessage(msg string, args ...any) *Error {
	e.msg = fmt.Sprintf(msg, args...)
	return e
}

func (e *Error) WithStack(stack string) *Error {
	e.stack = stack
	return e
}
