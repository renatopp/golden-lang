package errors

import (
	"runtime/debug"

	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/lang"
)

const (
	InternalError          core.ErrorKind = iota // For unexpected errors, which should never happen
	InvalidFileError                             // For invalid file paths
	UndefinedVariableError                       // For referencing variables identifiers that are not defined
	UndefinedTypeError                           // For referencing types identifiers that are not defined
	CircularReferenceError                       // For circular references in initialization
	ExpressionError                              // For wrong expression results
	TypeError                                    // For type mismatch
)

func toGoldenError(e any) *core.Error {
	if e, ok := e.(*core.Error); ok {
		return e
	}
	return core.NewError(InternalError, "%v", e)
}

func WithRecovery(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = toGoldenError(r).WithStack(string(debug.Stack()))
		}
	}()
	f()
	return
}

func RethrowWith(e error, kind core.ErrorKind, msg string, args ...any) {
	panic(toGoldenError(e).WithKind(kind).WithMessage(msg, args...))
}

func Rethrow(e error) {
	panic(e)
}

func ThrowAtLocation(loc lang.Loc, kind core.ErrorKind, msg string, args ...any) {
	panic(core.NewError(kind, msg, args...).WithLoc(&loc))
}

func ThrowAtToken(token lang.Token, kind core.ErrorKind, msg string, args ...any) {
	panic(core.NewError(kind, msg, args...).WithToken(&token))
}

func ThrowAtNode(node *core.AstNode, kind core.ErrorKind, msg string, args ...any) {
	panic(core.NewError(kind, msg, args...).WithNode(node))
}

func Throw(kind core.ErrorKind, msg string, args ...any) {
	panic(core.NewError(kind, msg, args...))
}
