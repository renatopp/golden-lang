package errors

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/lang"
)

func ToGoldenError(e any) *Error {
	if e, ok := e.(*Error); ok {
		return e
	}
	return NewError(InternalError, "%v", e).WithStack(string(debug.Stack()))
}

func WithRecoveryCallback(f func(), e func(error)) {
	err := WithRecovery(f)
	e(err)
}

func WithRecovery(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ToGoldenError(r)
			// println(string(debug.Stack()))
		}
	}()
	f()
	return
}

func RethrowWith(e error, code ErrorCode, msg string, args ...any) {
	panic(ToGoldenError(e).WithCode(code).WithMessage(msg, args...))
}

func Rethrow(e error) {
	panic(e)
}

func ThrowAtLocation(loc lang.Loc, code ErrorCode, msg string, args ...any) {
	panic(NewError(code, msg, args...).WithLoc(&loc))
}

func ThrowAtToken(token *lang.Token, code ErrorCode, msg string, args ...any) {
	panic(NewError(code, msg, args...).WithToken(token))
}

func ThrowAtNode(node ast.Node, code ErrorCode, msg string, args ...any) {
	panic(NewError(code, msg, args...).WithNode(node))
}

func Throw(code ErrorCode, msg string, args ...any) {
	panic(NewError(code, msg, args...))
}

func PrettyPrint(e error) {
	if e, ok := e.(*Error); ok {
		if e.Loc() == nil {
			prettySimpleError(e)
			return
		}

		prettyGoldenError(e)
	} else {
		prettySimpleError(e)
	}
}

func prettySimpleError(e error) {
	fmt.Printf("Error: %s", e.Error())
}

func prettyGoldenError(e *Error) {
	filePath := e.Loc().Filename
	source, err := os.ReadFile(filePath)
	if err != nil {
		prettySimpleError(e)
		return
	}

	lines := strings.Split(string(source), "\n")
	fromLine := e.Loc().Start.Line
	fromColumn := e.Loc().Start.Column
	// toLine := e.Loc().End.Line
	toColumn := e.Loc().End.Column
	columnSpan := max(1, toColumn-fromColumn)

	targetLine := lines[fromLine-1]

	fmt.Printf("[error %v] %s at line:%d, column:%d\n", e.Code(), filePath, fromLine, fromColumn)
	fmt.Printf("\n")
	fmt.Printf("    %s\n", targetLine)
	fmt.Printf("    %s\n", (strings.Repeat(" ", fromColumn-1) + strings.Repeat("^", columnSpan)))
	fmt.Printf("\n")
	fmt.Printf("Error: %s", e.Message())

	if e.Stack() != "" {
		fmt.Printf("\n%s\n", e.Stack())
	}
}
