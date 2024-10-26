package lang

import (
	"slices"
)

// Represents a character in the source code, with its position.
type Char struct {
	Rune rune
	Size int // Size in bytes
	Span Span
}

// Checks if the character rune is one of the given runes.
func (p Char) Is(runes ...rune) bool {
	return slices.Contains(runes, p.Rune)
}

// Creates an error at the char position.
func (p Char) AsError(kind, msg string) Error {
	return NewError(
		Loc{Start: p.Span, End: p.Span},
		kind,
		msg,
	)
}
