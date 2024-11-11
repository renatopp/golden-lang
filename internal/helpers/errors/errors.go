package errors

import (
	"runtime/debug"

	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/lang"
)

const (
	InternalError core.ErrorKind = iota
)

func WithRecovery(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(*core.Error); ok {
				err = e
			} else {
				err = core.NewError(InternalError, "%v", r)
				debug.PrintStack()
			}
		}
	}()
	f()
	return
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
