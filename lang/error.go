package lang

import (
	"runtime/debug"
	"strconv"
	"strings"
)

var (
	ErrSyntax               = "syntax error"
	ErrIO                   = "IO error"
	errMsgInvalidChar       = "invalid UTF-8 encoding"
	errMsgUnexpectedToken   = "expect token(s): %s; got %s"
	errMsgUnexpectedLiteral = "expect literal(s): %s; got %s"
)

// An error that occurred during lexing or parsing, including the source code
// location.
type Error struct {
	Loc  Loc
	Kind string
	Msg  string
}

// Creates a new error.
func NewError(loc Loc, kind, msg string) Error {
	return Error{loc, kind, msg}
}

// Returns the error message.
func (e Error) Error() string {
	return e.Msg + " at " + strconv.Itoa(e.Loc.Start.Line) + ":" + strconv.Itoa(e.Loc.Start.Column)
}

// ---

type ErrorList struct {
	errors []error
}

func NewErrorList(errs []error) *ErrorList {
	return &ErrorList{errors: errs}
}

func (e ErrorList) Error() string {
	msgs := []string{}
	for _, err := range e.errors {
		msgs = append(msgs, "- "+err.Error())
	}
	return strings.Join(msgs, "\n")
}

//----

type ErrorData struct {
	errors []error
}

func NewErrorData() *ErrorData {
	return &ErrorData{errors: []error{}}
}

func (e *ErrorData) HasErrors() bool {
	return len(e.errors) > 0
}

func (e *ErrorData) TotalErrors() int {
	return len(e.errors)
}

func (e *ErrorData) Errors() []error {
	return e.errors
}

func (e *ErrorData) RegisterError(err error) {
	e.errors = append(e.errors, err)
}

func (e *ErrorData) WithRecovery(fn func()) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(Error); ok {
				e.RegisterError(err)
			} else {
				e.RegisterError(NewError(Loc{}, "unknown error", r.(string)))
				debug.PrintStack()
			}
		}
	}()
	fn()
}
