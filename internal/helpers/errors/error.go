package errors

import (
	"fmt"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/lang"
)

type ErrorCode uint64

type Error struct {
	loc   *lang.Loc
	token *lang.Token
	node  ast.Node
	code  ErrorCode
	msg   string
	stack string
}

func NewError(code ErrorCode, msg string, args ...any) *Error {
	return &Error{
		code: code,
		msg:  fmt.Sprintf(msg, args...),
	}
}

func NewEmptyError() *Error {
	return &Error{}
}

func (e *Error) Loc() *lang.Loc     { return e.loc }
func (e *Error) Token() *lang.Token { return e.token }
func (e *Error) Node() ast.Node     { return e.node }
func (e *Error) Code() ErrorCode    { return e.code }
func (e *Error) Message() string    { return e.msg }
func (e *Error) Stack() string      { return e.stack }
func (e *Error) Error() string      { return e.msg }

func (e *Error) WithLoc(loc *lang.Loc) *Error {
	e.loc = loc
	return e
}

func (e *Error) WithToken(token *lang.Token) *Error {
	e.loc = token.Loc
	e.token = token
	return e
}

func (e *Error) WithNode(node ast.Node) *Error {
	e.token = node.Token()
	e.loc = e.token.Loc
	e.node = node
	return e
}

func (e *Error) WithCode(code ErrorCode) *Error {
	e.code = code
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
