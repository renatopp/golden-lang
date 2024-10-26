package eaters

import (
	"strconv"

	"github.com/renatopp/golden/lang"
	"github.com/renatopp/golden/lang/runes"
)

var (
	ErrSyntax                 = "Syntax error"
	errMsgUnexpectedNewline   = "unexpected newline"
	errMsgUnexpectedEndOfFile = "unexpected end of file"
	errMsgUnexpectedDot       = "unexpected dot"
	errMsgUnexpectedE         = "unexpected exponent"
)

// Consumes all the characters that composes a string. This function will
// consider the first character as the delimiter of the string, and will stop
// when it finds the same character again. If the string is not closed, it will
// register an error. If the string contains a newline, it will register an
// error and ignore the newline. If you need to consume a string with a newline
// character, use `EatRawString` instead.
func EatString(s *lang.ByteScanner) *lang.Token {
	result := ""
	escaping := false
	first := s.EatChar()
	for {
		c := s.PeekChar()

		if c.Is('\n') {
			// Ignoring newlines
			s.RegisterError(c.AsError(ErrSyntax, errMsgUnexpectedNewline))
			s.EatChar()
			continue

		} else if s.IsEof() {
			// Stopping at EOF
			s.RegisterError(c.AsError(ErrSyntax, errMsgUnexpectedEndOfFile))
			break

		} else if !escaping && c.Is(first.Rune) {
			// End of string
			break

		} else if !escaping && c.Is('\\') {
			// Starting Escaping
			escaping = true
			s.EatChar()
			continue

		} else if escaping && !c.Is(first.Rune) {
			// Ending Escaping
			escaping = false
			r, err := strconv.Unquote(`"\` + string(c.Rune) + `"`)
			if err != nil {
				c := s.EatChar()
				s.RegisterError(c.AsError(ErrSyntax, err.Error()))
				continue
			}
			c.Rune = rune(r[0])
		}

		result += string(c.Rune)
		s.EatChar()
	}

	s.EatChar()
	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}

// Consumes all the characters that composes a raw string. This function will
// consider the first character as the delimiter of the string, and will stop
// when it finds the same character again. If the string is not closed, it will
// register an error.
//
// This function will record the string as it was written, meaning that any
// escape character will be kept as it is, including new lines. If you need to
// ignore new lines, use `EatString` instead.
func EatRawString(s *lang.ByteScanner) *lang.Token {
	result := ""
	escaping := false
	first := s.EatChar()
	for {
		c := s.PeekChar()

		if s.IsEof() {
			// Stopping at EOF
			s.RegisterError(c.AsError(ErrSyntax, errMsgUnexpectedEndOfFile))
			break

		} else if !escaping && c.Is(first.Rune) {
			// End of string
			break

		} else if !escaping && c.Is('\\') {
			// Starting Escaping
			escaping = true
			s.EatChar()
			continue

		} else if escaping && !c.Is(first.Rune) {
			// Ending Escaping
			escaping = false
			r, err := strconv.Unquote(`"\` + string(c.Rune) + `"`)
			if err != nil {
				c := s.EatChar()
				s.RegisterError(c.AsError(ErrSyntax, err.Error()))
				continue
			}
			c.Rune = rune(r[0])
		}

		result += string(c.Rune)
		s.EatChar()
	}

	s.EatChar()
	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}

// Consumes all the characters that composes all common cases of numbers.
// Including:
//
// - Integers: `123`
// - Floats: `123.32`
// - Exponents: `123e32`
// - Floats with Exponents: `123.32e32`
// - Exponents with signal: `123e+32`
func EatNumber(s *lang.ByteScanner) *lang.Token {
	result := ""
	dot := false
	exp := false

	first := s.PeekChar()
	for {
		c := s.PeekChar()

		switch {
		case c.Is('_'):
			s.EatChar()
			continue

		case c.Is('.'):
			if dot || exp {
				s.RegisterError(c.AsError(ErrSyntax, errMsgUnexpectedDot))
				s.EatChar()
				continue
			}

			dot = true
			result += string(c.Rune)

		case c.Is('e') || c.Is('E'):
			if exp {
				s.RegisterError(c.AsError(ErrSyntax, errMsgUnexpectedE))
				s.EatChar()
				continue
			}

			exp = true
			result += string(c.Rune)

			next := s.PeekCharAt(1)
			if next.Is('+') || next.Is('-') {
				s.EatChar()
				result += string(next.Rune)
			}

		case runes.IsDigit(c.Rune):
			result += string(c.Rune)

		default:
			return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
		}

		s.EatChar()
	}
}

// Consumes all the characters that composes an integer number.
func EatInteger(s *lang.ByteScanner) *lang.Token {
	result := ""

	first := s.PeekChar()
	for {
		c := s.PeekChar()
		if !runes.IsDigit(c.Rune) {
			break
		}

		result += string(c.Rune)
		s.EatChar()
	}

	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}

// Consumes all the characters that composes a hexadecimal number. Considering
// `0x` or `0X` as optional prefix.
func EatHexadecimal(s *lang.ByteScanner) *lang.Token {
	result := ""

	first := s.EatChar()
	next := s.PeekChar()
	if first.Is('0') && (next.Is('x') || next.Is('X')) {
		s.EatChar()
	} else {
		result += string(first.Rune)
	}

	for {
		c := s.PeekChar()
		if !runes.IsHexadecimal(c.Rune) {
			break
		}

		result += string(c.Rune)
		s.EatChar()
	}

	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}

// Consumes all the characters that composes an octal number. Considering `0`
// as optional prefix.
func EatOctal(s *lang.ByteScanner) *lang.Token {
	result := ""

	first := s.EatChar()
	next := s.PeekChar()
	if first.Is('0') && (next.Is('o') || next.Is('O')) {
		s.EatChar()
	} else if !first.Is('0') {
		result += string(first.Rune)
	}

	for {
		c := s.PeekChar()
		if !runes.IsOctal(c.Rune) {
			break
		}

		result += string(c.Rune)
		s.EatChar()
	}

	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}

// Consumes all the characters that composes a binary number. Considering `0b`
// or `0B` as optional prefix.
func EatBinary(s *lang.ByteScanner) *lang.Token {
	result := ""

	first := s.EatChar()
	next := s.PeekChar()
	if first.Is('0') && (next.Is('b') || next.Is('B')) {
		s.EatChar()
	} else {
		result += string(first.Rune)
	}

	for {
		c := s.PeekChar()
		if !runes.IsBinary(c.Rune) {
			break
		}

		result += string(c.Rune)
		s.EatChar()
	}

	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}

// Consumes all the characters that composes a whitespace. Including space,
// tab, newline and carriage return.
func EatWhitespaces(s *lang.ByteScanner) *lang.Token {
	result := ""

	first := s.PeekChar()
	for {
		c := s.PeekChar()
		if !runes.IsWhitespace(c.Rune) {
			break
		}

		result += string(c.Rune)
		s.EatChar()
	}

	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}

// Consumes all the characters that composes a space. Including space, tab and
// carriage return.
func EatSpaces(s *lang.ByteScanner) *lang.Token {
	result := ""

	first := s.PeekChar()
	for {
		c := s.PeekChar()
		if !runes.IsSpace(c.Rune) {
			break
		}

		result += string(c.Rune)
		s.EatChar()
	}

	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}

// Consumes all the characters that composes a newline. Including only the
// newline character (without the carriage return).
func EatNewlines(s *lang.ByteScanner) *lang.Token {
	result := ""

	first := s.PeekChar()
	for {
		c := s.PeekChar()
		if !c.Is('\n') {
			break
		}

		result += string(c.Rune)
		s.EatChar()
	}

	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}

// Consumes all the characters that composes an common identifier. Including
// letters, digits and underscores.
func EatIdentifier(s *lang.ByteScanner) *lang.Token {
	result := ""

	first := s.PeekChar()
	for {
		c := s.PeekChar()
		if !runes.IsAlphaNumeric(c.Rune) && !c.Is('_') {
			break
		}

		result += string(c.Rune)
		s.EatChar()
	}

	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}

// Consumes all the characters that composes an identifier. Including letters,
// digits and the characters passed as argument.
func EatIdentifierWith(s *lang.ByteScanner, allowedChars ...rune) *lang.Token {
	result := ""

	first := s.PeekChar()
	for {
		c := s.PeekChar()
		if !runes.IsAlphaNumeric(c.Rune) && !c.Is(allowedChars...) {
			break
		}

		result += string(c.Rune)
		s.EatChar()
	}

	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}

// Consumes all the characters that composes a word. Including letters and
// digits.
func EatWord(s *lang.ByteScanner) *lang.Token {
	result := ""

	first := s.PeekChar()
	for {
		c := s.PeekChar()
		if !runes.IsAlphaNumeric(c.Rune) {
			break
		}

		result += string(c.Rune)
		s.EatChar()
	}

	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}

// Consumes all the characters until the end of the line or the end of file.
func EatUntilEndOfLine(s *lang.ByteScanner) *lang.Token {
	result := ""

	first := s.PeekChar()
	for {
		c := s.PeekChar()
		if c.Is('\n', 0) {
			break
		}

		result += string(c.Rune)
		s.EatChar()
	}

	return lang.NewToken("unknown", result).WithChars(first, s.PeekChar())
}
