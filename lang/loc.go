package lang

// Loc represents a location in the source code.
type Loc struct {
	Filename string
	Start    Span
	End      Span
}

// Span represents a single point in the source code.
type Span struct {
	Offset int // Position in the source code byte array
	Line   int
	Column int
}
