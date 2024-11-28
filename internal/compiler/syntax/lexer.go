package syntax

import (
	"strconv"
	"strings"

	"github.com/renatopp/golden/internal/compiler/token"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/internal/helpers/naming"
	"github.com/renatopp/golden/internal/helpers/runes"
)

type Lexer struct {
	filename   string
	line       int
	column     int
	fromLine   int
	fromColumn int
	scanner    *Scanner[rune]
}

func NewLexer(filename string, source []byte) *Lexer {
	return &Lexer{
		filename: filename,
		line:     1,
		column:   1,
		scanner:  NewScanner([]rune(string(source)), rune(0)),
	}
}

func (l *Lexer) Lex() (res []token.Token, err error) {
	err = errors.WithRecovery(func() {
		res = l.lex()
	})
	return res, err
}

func (l *Lexer) lex() []token.Token {
	tokens := []token.Token{}
	for {
		token, ok := l.next()
		if !ok {
			break
		}
		tokens = append(tokens, token)
	}
	return tokens
}

func (l *Lexer) next() (token.Token, bool) {
	for !l.scanner.IsFinished() {
		l.fromLine = l.line
		l.fromColumn = l.column

		c0 := l.scanner.PeekAt(0)
		c1 := l.scanner.PeekAt(1)
		c2 := l.scanner.PeekAt(2)

		s1 := string(c0)
		s2 := s1 + string(c1)
		s3 := s2 + string(c2)

		switch {
		// EOF
		case runes.IsEof(c0):
			return token.Token{}, false

		// Spaces
		case runes.IsSpace(c0):
			l.eatSpaces()
			continue

		// Newlines
		case runes.IsNewline(c0):
			return token.Token{
				Kind:  token.TNewline,
				Value: l.eatNewlines(),
				Loc:   l.span(),
			}, true

		// Comments
		case s2 == "--":
			return token.Token{
				Kind:  token.TComment,
				Value: l.eatComment(),
				Loc:   l.span(),
			}, true

		// Alpha literals, identifiers and keywords
		case runes.IsAlpha(c0) || runes.IsOneOf(c0, '_'):
			identifier := l.eatIdentifier()

			// Keywords or constants
			if kind := token.LiteralToKind(identifier); kind != token.TUnknown {
				return token.Token{
					Kind:  kind,
					Value: identifier,
					Loc:   l.span(),
				}, true
			}

			if naming.IsTypeName(identifier) {
				return token.Token{
					Kind:  token.TTypeIdent,
					Value: identifier,
					Loc:   l.span(),
				}, true
			}

			return token.Token{
				Kind:  token.TVarIdent,
				Value: identifier,
				Loc:   l.span(),
			}, true

		// Numeric literals
		case runes.IsDigit(c0) || c0 == '.' && runes.IsDigit(c1):
			switch {
			// Hex
			case runes.IsOneOf(c1, 'x', 'X'):
				return token.Token{
					Kind:  token.THex,
					Value: l.eatHexadecimal(),
					Loc:   l.span(),
				}, true

			// Octal
			case runes.IsOneOf(c1, 'o', 'O'):
				return token.Token{
					Kind:  token.TOctal,
					Value: l.eatOctal(),
					Loc:   l.span(),
				}, true

			// Binary
			case runes.IsOneOf(c1, 'b', 'B'):
				return token.Token{
					Kind:  token.TBinary,
					Value: l.eatBinary(),
					Loc:   l.span(),
				}, true

			// Float or integer
			default:
				num := l.eatNumber()
				if strings.Contains(num, ".") || strings.Contains(num, "e") {
					return token.Token{
						Kind:  token.TFloat,
						Value: num,
						Loc:   l.span(),
					}, true
				}

				return token.Token{
					Kind:  token.TInt,
					Value: num,
					Loc:   l.span(),
				}, true
			}

		// Strings
		case runes.IsOneOf(c0, '\''):
			return token.Token{
				Kind:  token.TString,
				Value: l.eatRawString(),
				Loc:   l.span(),
			}, true

		// Operators
		default:
			// 3-char operators
			if tok := token.LiteralToKind(s3); tok != token.TUnknown {
				l.eat()
				l.eat()
				l.eat()
				return token.Token{
					Kind:  tok,
					Value: s3,
					Loc:   l.span(),
				}, true
			}

			// 2-char operators
			if tok := token.LiteralToKind(s2); tok != token.TUnknown {
				l.eat()
				l.eat()
				return token.Token{
					Kind:  tok,
					Value: s2,
					Loc:   l.span(),
				}, true
			}

			// 1-char operators
			if tok := token.LiteralToKind(s1); tok != token.TUnknown {
				l.eat()
				return token.Token{
					Kind:  tok,
					Value: s1,
					Loc:   l.span(),
				}, true
			}

			// Unknown
			l.eat()
			errors.ThrowAtLocation(l.span(), errors.ParserError, "unexpected character: %s", s1)
		}
	}
	return token.Token{}, false
}

func (l *Lexer) span() token.Span {
	return token.Span{
		Filename:   l.filename,
		FromLine:   l.fromLine,
		FromColumn: l.fromColumn,
		ToLine:     l.line,
		ToColumn:   l.column,
	}
}

func (l *Lexer) eat() rune {
	c := l.scanner.Eat()
	if c == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
	return c
}

// Consumes all the characters that composes a space. Including space, tab and
// carriage return.
func (l *Lexer) eatSpaces() string {
	res := ""
	for {
		c := l.scanner.Peek()
		if !runes.IsSpace(c) {
			break
		}
		res += string(c)
		l.eat()
	}
	return res
}

// Consumes all the characters that composes a newline. Including only the
// newline character (without the carriage return).
func (l *Lexer) eatNewlines() string {
	res := ""
	for {
		c := l.scanner.Peek()
		if !runes.IsNewline(c) {
			break
		}
		res += string(c)
		l.eat()
	}
	return res
}

// Consumes all the characters until the end of the line or the end of file.
func (l *Lexer) eatComment() string {
	res := ""
	for {
		c := l.scanner.Peek()
		if runes.IsOneOf(c, '\r', '\n', 0) {
			break
		}

		res += string(c)
		l.eat()
	}
	return res
}

// Consumes all the characters that composes an common identifier. Including
// letters, digits and underscores.
func (l *Lexer) eatIdentifier() string {
	res := ""
	for {
		c := l.scanner.Peek()
		if !runes.IsAlphaNumeric(c) && c != '_' {
			break
		}

		res += string(c)
		l.eat()
	}
	return res
}

// Consumes all the characters that composes all common cases of numbers.
// Including:
//
// - Integers: `123_000`
// - Floats: `123.32`
// - Exponents: `123e32`
// - Floats with Exponents: `123.32e32`
// - Exponents with signal: `123e+32`
func (l *Lexer) eatNumber() string {
	res := ""
	dot := false
	exp := false
	for {
		c := l.scanner.Peek()
		switch {
		case c == '_':
			l.eat()
			continue

		case c == '.':
			if dot || exp {
				errors.ThrowAtLocation(l.span(), errors.ParserError, "unexpected dot")
				l.eat()
				continue
			}
			dot = true
			res += string(c)

		case runes.IsOneOf(c, 'e', 'E'):
			if exp {
				errors.ThrowAtLocation(l.span(), errors.ParserError, "unexpected e")
				l.eat()
				continue
			}
			exp = true
			res += string(c)

			next := l.scanner.PeekAt(1)
			if runes.IsOneOf(next, '+', '-') {
				l.eat()
				res += string(next)
			}

		case runes.IsDigit(c):
			res += string(c)

		default:
			return res
		}

		l.eat()
	}
}

// Consumes all the characters that composes a hexadecimal number. Considering
// `0x` or `0X` as optional prefix.
func (l *Lexer) eatHexadecimal() string {
	res := ""
	c0 := l.scanner.PeekAt(0)
	c1 := l.scanner.PeekAt(1)
	if runes.IsOneOf(c0, 'x', 'X') {
		l.eat()
	}
	if runes.IsOneOf(c1, 'x', 'X') {
		l.eat()
		l.eat()
	}
	for {
		c := l.scanner.Peek()
		if !runes.IsHexadecimal(c) {
			break
		}
		res += string(c)
		l.eat()
	}
	return res
}

// Consumes all the characters that composes an octal number. Considering `0`
// as optional prefix.
func (l *Lexer) eatOctal() string {
	res := ""
	c0 := l.scanner.PeekAt(0)
	c1 := l.scanner.PeekAt(1)
	if runes.IsOneOf(c0, 'o', 'O') {
		l.eat()
	}
	if runes.IsOneOf(c1, 'o', 'O') {
		l.eat()
		l.eat()
	}
	for {
		c := l.scanner.Peek()
		if !runes.IsOctal(c) {
			break
		}
		res += string(c)
		l.eat()
	}
	return res
}

// Consumes all the characters that composes a binary number. Considering `0b`
// or `0B` as optional prefix.
func (l *Lexer) eatBinary() string {
	res := ""
	c0 := l.scanner.PeekAt(0)
	c1 := l.scanner.PeekAt(1)
	if runes.IsOneOf(c0, 'b', 'B') {
		l.eat()
	}
	if runes.IsOneOf(c1, 'b', 'B') {
		l.eat()
		l.eat()
	}
	for {
		c := l.scanner.Peek()
		if !runes.IsBinary(c) {
			break
		}
		res += string(c)
		l.eat()
	}
	return res
}

// Consumes all the characters that composes a raw string. This function will
// consider the first character as the delimiter of the string, and will stop
// when it finds the same character again. If the string is not closed, it will
// register an error.
//
// This function will record the string as it was written, meaning that any
// escape character will be kept as it is, including new lines. If you need to
// ignore new lines, use `EatString` instead.
func (l *Lexer) eatRawString() string {
	res := ""
	escaping := false
	first := l.eat()
	for {
		c := l.scanner.Peek()
		println(string(c))

		if runes.IsOneOf(c, '\r') {
			l.eat()
			continue
		}

		if runes.IsEof(c) {
			errors.ThrowAtLocation(l.span(), errors.ParserError, "unexpected end of file")
			break
		}

		if !escaping && c == first {
			break
		}

		if !escaping && c == '\\' {
			escaping = true
			l.eat()
			continue
		}

		if escaping && c != first {
			escaping = false
			r, err := strconv.Unquote(`"\` + string(c) + `"`)
			if err != nil {
				errors.ThrowAtLocation(l.span(), errors.ParserError, "%v", err.Error())
			}
			c = []rune(r)[0]
		}

		res += string(c)
		l.eat()
	}
	l.eat()
	return res
}
