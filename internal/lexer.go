package internal

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/renatopp/golden/lang"
	"github.com/renatopp/golden/lang/eaters"
	"github.com/renatopp/golden/lang/runes"
)

var typeRegex = regexp.MustCompile(`^_*[A-Z][a-zA-Z0-9_]*$`)

func Lex(input []byte) ([]*lang.Token, error) {
	scanner := lang.NewByteScanner(input)
	lexer := &lexer{scanner: scanner}
	tokens := lexer.All()

	var err error = nil
	if scanner.HasErrors() {
		err = lang.NewErrorList(scanner.Errors())
	}
	return tokens, err
}

type lexer struct {
	scanner *lang.ByteScanner
}

func (l *lexer) All() []*lang.Token {
	tokens := []*lang.Token{}
	for {
		t := l.Next()
		tokens = append(tokens, t)
		if t.Kind == TEof {
			break
		}
	}
	return tokens
}

func (l *lexer) Next() *lang.Token {
	for {
		c0 := l.scanner.PeekCharAt(0)
		c1 := l.scanner.PeekCharAt(1)
		c2 := l.scanner.PeekCharAt(2)

		s1 := string(c0.Rune)
		s2 := s1 + string(c1.Rune)
		s3 := s2 + string(c2.Rune)

		if l.scanner.TotalErrors() >= 10 {
			return lang.NewToken(TEof, "").WithChars(c0, c1)
		}

		switch {
		// EOF
		case c0.Is(0):
			return lang.NewToken(TEof, "").WithChars(c0, c1)

		// Whitespaces
		case c0.Is(' ', '\t', '\r'):
			eaters.EatSpaces(l.scanner)
			continue

			// Newlines
		case c0.Is('\n'):
			return eaters.EatNewlines(l.scanner).WithType(TNewline)

			// Comments
		case s2 == "--":
			t := eaters.EatUntilEndOfLine(l.scanner).WithType(TComment)
			eaters.EatSpaces(l.scanner) // \r
			l.scanner.EatChar()         // \n
			return t

		// // Literals, Identifiers and Keywords
		case runes.IsAlpha(c0.Rune) || c0.Is('_'):
			t := eaters.EatIdentifier(l.scanner)

			switch {
			case t.IsLiteral("true", "false"):
				return t.WithType(TBool)
			case t.IsLiteral("and"):
				return t.WithType(TAnd)
			case t.IsLiteral("or"):
				return t.WithType(TOr)
			case t.IsLiteral("xor"):
				return t.WithType(TXor)
			case slices.Contains(Keywords, t.Literal):
				return t.WithType(TKeyword)
			default:
				if typeRegex.MatchString(t.Literal) {
					return t.WithType(TTypeIdent)
				}
				return t.WithType(TVarIdent)
			}

		// Numbers
		case runes.IsNumeric(c0.Rune):
			switch {
			case c1.Is('x', 'X'):
				return eaters.EatHexadecimal(l.scanner).WithType(THex)
			case c1.Is('o', 'O'):
				return eaters.EatOctal(l.scanner).WithType(TOctal)
			case c1.Is('b', 'B'):
				return eaters.EatBinary(l.scanner).WithType(TBinary)
			default:
				num := eaters.EatNumber(l.scanner)
				if strings.Contains(num.Literal, ".") || strings.Contains(num.Literal, "e") {
					return num.WithType(TFloat)
				}
				return num.WithType(TInteger)
			}

		// Strings
		case c0.Is('"'):
			return eaters.EatString(l.scanner).WithType(TString)

		case c0.Is('`'):
			return eaters.EatRawString(l.scanner).WithType(TString)

		default:
			if kind, ok := TripleCharTokens[s3]; ok {
				chars := l.scanner.EatChars(3)
				return lang.NewToken(kind, s3).WithChars(chars[0], chars[2])
			}

			if kind, ok := DoubleCharTokens[s2]; ok {
				chars := l.scanner.EatChars(2)
				return lang.NewToken(kind, s2).WithChars(chars[0], chars[1])
			}

			if kind, ok := SingleCharTokens[s1]; ok {
				c := l.scanner.EatChar()
				return lang.NewToken(kind, s1).WithChars(c0, c)
			}

			l.scanner.EatChar()
			l.scanner.RegisterError(c0.AsError(eaters.ErrSyntax, fmt.Sprintf("unexpected character '%v'", s1)))
			return lang.NewToken(TInvalid, s1).WithChars(c0, c1)
		}
	}
}
