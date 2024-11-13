package errors

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/lang"
)

const (
	InternalError            core.ErrorKind = iota // For unexpected errors, which should never happen
	InvalidFileError                               // For invalid file path
	UnexpectedCharacterError                       // For lexing  errors
	UndefinedVariableError                         // For referencing variables identifiers that are not defined
	UndefinedTypeError                             // For referencing types identifiers that are not defined
	CircularReferenceError                         // For circular references in initialization
	ExpressionError                                // For wrong expression results
	ParserError                                    // For type mismatch
	TypeError                                      // For type mismatch
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

func ThrowAtToken(token *lang.Token, kind core.ErrorKind, msg string, args ...any) {
	panic(core.NewError(kind, msg, args...).WithToken(token))
}

func ThrowAtNode(node *core.AstNode, kind core.ErrorKind, msg string, args ...any) {
	panic(core.NewError(kind, msg, args...).WithNode(node))
}

func Throw(kind core.ErrorKind, msg string, args ...any) {
	panic(core.NewError(kind, msg, args...))
}

func PrettyPrint(e error) {
	if e, ok := e.(*core.Error); ok {
		if e.Loc == nil {
			prettySimpleError(e)
			return
		}

		prettyGoldenError(e)
	} else {
		prettySimpleError(e)
	}
}

func prettySimpleError(e error) {
	println(e.Error())
}

func prettyGoldenError(e *core.Error) {
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

	fmt.Printf("[error %v] %s at line:%d, column:%d\n", e.Kind(), filePath, fromLine, fromColumn)
	fmt.Printf("\n")
	fmt.Printf("    %s\n", targetLine)
	fmt.Printf("    %s\n", (strings.Repeat(" ", fromColumn-1) + strings.Repeat("^", columnSpan)))
	fmt.Printf("\n")
	fmt.Printf("Error: %s", e.Message())
}
