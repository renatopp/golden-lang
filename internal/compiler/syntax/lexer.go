package syntax

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/renatopp/golden/internal/core"
	"github.com/renatopp/golden/internal/helpers/errors"
	"github.com/renatopp/golden/lang"
	"github.com/renatopp/golden/lang/eaters"
	"github.com/renatopp/golden/lang/runes"
)

var typeRegex = regexp.MustCompile(`^_*[A-Z][a-zA-Z0-9_]*$`)

func Lex(input []byte, modulePath string) ([]*lang.Token, error) {
	scanner := lang.NewByteScanner(input)
	tokens := []*lang.Token{}
	for {
		t := nextToken(scanner)
		tokens = append(tokens, t.WithFile(modulePath))
		if t.Kind == core.TEof {
			break
		}
	}

	if scanner.HasErrors() {
		err := lang.NewErrorList(scanner.Errors())
		return tokens, errors.ToGoldenError(err)
	}

	return tokens, nil
}

func nextToken(scanner *lang.ByteScanner) *lang.Token {
	for {
		c0 := scanner.PeekCharAt(0)
		c1 := scanner.PeekCharAt(1)
		c2 := scanner.PeekCharAt(2)

		s1 := string(c0.Rune)
		s2 := s1 + string(c1.Rune)
		s3 := s2 + string(c2.Rune)

		if scanner.TotalErrors() >= 10 {
			return lang.NewToken(core.TEof, "").WithChars(c0, c1)
		}

		switch {
		// EOF
		case c0.Is(0):
			return lang.NewToken(core.TEof, "").WithChars(c0, c1)

		// Whitespaces
		case c0.Is(' ', '\t', '\r'):
			eaters.EatSpaces(scanner)
			continue

			// Newlines
		case c0.Is('\n'):
			return eaters.EatNewlines(scanner).WithType(core.TNewline)

			// Comments
		case s2 == "--":
			t := eaters.EatUntilEndOfLine(scanner).WithType(core.TComment)
			eaters.EatSpaces(scanner) // \r
			scanner.EatChar()         // \n
			return t

		// // Literals, Identifiers and Keywords
		case runes.IsAlpha(c0.Rune) || c0.Is('_'):
			t := eaters.EatIdentifier(scanner)
			tok := core.LiteralToToken(t.Literal)
			if tok != core.TUnknown {
				return t.WithType(tok)
			}
			if typeRegex.MatchString(t.Literal) {
				return t.WithType(core.TTypeIdent)
			}
			return t.WithType(core.TVarIdent)

		// Numbers
		case runes.IsNumeric(c0.Rune):
			switch {
			case c1.Is('x', 'X'):
				return eaters.EatHexadecimal(scanner).WithType(core.THex)
			case c1.Is('o', 'O'):
				return eaters.EatOctal(scanner).WithType(core.TOctal)
			case c1.Is('b', 'B'):
				return eaters.EatBinary(scanner).WithType(core.TBinary)
			default:
				num := eaters.EatNumber(scanner)
				if strings.Contains(num.Literal, ".") || strings.Contains(num.Literal, "e") {
					return num.WithType(core.TFloat)
				}
				return num.WithType(core.TInteger)
			}

		// Strings
		case c0.Is('\''):
			return eaters.EatString(scanner).WithType(core.TString)

		case c0.Is('`'):
			return eaters.EatRawString(scanner).WithType(core.TString)

		default:
			if tok := core.LiteralToToken(s3); tok != core.TUnknown {
				chars := scanner.EatChars(3)
				return lang.NewToken(tok, s3).WithChars(chars[0], chars[2])
			}
			if tok := core.LiteralToToken(s2); tok != core.TUnknown {
				chars := scanner.EatChars(2)
				return lang.NewToken(tok, s2).WithChars(chars[0], chars[1])
			}
			if tok := core.LiteralToToken(s1); tok != core.TUnknown {
				chars := scanner.EatChars(1)
				return lang.NewToken(tok, s1).WithChars(chars[0], chars[0])
			}

			scanner.EatChar()
			scanner.RegisterError(c0.AsError(eaters.ErrSyntax, fmt.Sprintf("unexpected character '%v'", s1)))
			return lang.NewToken(core.TInvalid, s1).WithChars(c0, c1)
		}
	}
}
