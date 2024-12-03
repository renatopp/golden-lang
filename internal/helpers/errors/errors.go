package errors

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/renatopp/golden/internal/compiler/ast"
	"github.com/renatopp/golden/internal/compiler/token"
	"github.com/renatopp/golden/internal/helpers/safe"
)

type ErrorCode uint64

const (
	InternalError ErrorCode = iota
	NotImplemented
	InvalidFileError
	InvalidFolderError
	CircularReferenceError
	ParserError
	TypeError
	NameNotFound
	NameAlreadyDefined
	InvalidEntryFile
	TemporaryImplementationError
)

var codeToName = map[ErrorCode]string{
	InternalError:                "internal error",
	NotImplemented:               "not implemented",
	InvalidFileError:             "invalid file error",
	InvalidFolderError:           "invalid folder error",
	CircularReferenceError:       "circular reference error",
	ParserError:                  "parser error",
	TypeError:                    "type error",
	NameNotFound:                 "name not found",
	NameAlreadyDefined:           "name already defined",
	InvalidEntryFile:             "invalid entry file",
	TemporaryImplementationError: "temporary implementation error",
}

//
//
//

// GoldenError is a custom error type that contains information about the error
type GoldenError struct {
	Loc   safe.Optional[*token.Span]
	Token safe.Optional[*token.Token]
	Node  safe.Optional[ast.Node]
	Code  ErrorCode
	Msg   string
	Stack string
}

func NewError(code ErrorCode, msg string, args ...any) GoldenError {
	return GoldenError{
		Code: code,
		Msg:  fmt.Sprintf(msg, args...),
	}
}

func (e GoldenError) Error() string { return e.Msg }

func (e GoldenError) WithLoc(loc *token.Span) GoldenError {
	e.Loc = safe.Some(loc)
	return e
}

func (e GoldenError) WithToken(token *token.Token) GoldenError {
	e.Loc = safe.Some(token.Loc)
	e.Token = safe.Some(token)
	return e
}

func (e GoldenError) WithNode(node ast.Node) GoldenError {
	e.Token = safe.Some(node.GetToken())
	e.Loc = safe.Some(e.Token.Unwrap().Loc)
	e.Node = safe.Some(node)
	return e
}

func (e GoldenError) WithCode(code ErrorCode) GoldenError {
	e.Code = code
	return e
}

func (e GoldenError) WithMessage(msg string, args ...any) GoldenError {
	e.Msg = fmt.Sprintf(msg, args...)
	return e
}

func (e GoldenError) WithStack(stack string) GoldenError {
	e.Stack = stack
	return e
}

//
//
//

func ToGoldenError(e any) GoldenError {
	if e, ok := e.(GoldenError); ok {
		return e
	}
	if e, ok := e.(*GoldenError); ok && e != nil {
		return *e
	}

	return NewError(InternalError, "%v", e).WithStack(string(debug.Stack()))
}

func WithRecovery(f func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ToGoldenError(r)
			// err = gerr.WithStack(string(debug.Stack()))
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

func ThrowAtLocation(loc *token.Span, code ErrorCode, msg string, args ...any) {
	panic(NewError(code, msg, args...).WithLoc(loc))
}

func ThrowAtToken(token *token.Token, code ErrorCode, msg string, args ...any) {
	panic(NewError(code, msg, args...).WithToken(token))
}

func ThrowAtNode(node ast.Node, code ErrorCode, msg string, args ...any) {
	panic(NewError(code, msg, args...).WithNode(node))
}

func Throw(code ErrorCode, msg string, args ...any) {
	panic(NewError(code, msg, args...))
}

//
//
//

func PrettyPrint(e error) {
	switch e := e.(type) {
	case GoldenError:
		prettyGoldenError(e)
	case *GoldenError:
		prettyGoldenError(*e)
	default:
		prettySimpleError(e)
	}
}

func prettySimpleError(e error) {
	fmt.Printf("Error: %s", e.Error())
}

func prettyGoldenError(e GoldenError) {
	if !e.Loc.Has() {
		prettySimpleError(e)
		return
	}

	loc := e.Loc.Unwrap()
	filePath := loc.Filename
	source, err := os.ReadFile(filePath)
	if err != nil {
		prettySimpleError(e)
		return
	}

	lines := strings.Split(string(source), "\n")
	fromLine := loc.FromLine
	fromColumn := loc.FromColumn
	// toLine := loc.ToLine
	toColumn := loc.ToColumn
	columnSpan := max(1, toColumn-fromColumn)

	targetLine := lines[fromLine-1]

	code := codeToName[e.Code]

	fmt.Printf("[%v] %s at line:%d, column:%d\n", code, filePath, fromLine, fromColumn)
	fmt.Printf("\n")
	fmt.Printf("    %s\n", targetLine)
	fmt.Printf("    %s\n", (strings.Repeat(" ", fromColumn-1) + strings.Repeat("^", columnSpan)))
	fmt.Printf("\n")
	fmt.Printf("Error: %s", e.Msg)

	if e.Stack != "" {
		fmt.Printf("\n%s\n", e.Stack)
	}
}
