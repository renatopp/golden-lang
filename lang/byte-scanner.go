package lang

import "unicode/utf8"

// ByteScanner is responsible for reading the source code input and returning the
// sequence of characters in that.
type ByteScanner struct {
	*ErrorData
	input  []byte
	cursor int
	line   int
	column int
	eof    *Char
	queue  []Char
}

// Creates a new scanner attached to the given input.
func NewByteScanner(input []byte) *ByteScanner {
	return &ByteScanner{
		ErrorData: NewErrorData(),
		input:     input,
		cursor:    0,
		line:      1,
		column:    1,
		eof:       nil,
		queue:     []Char{},
	}
}

// Return if the next token is the end of the file.
func (s *ByteScanner) IsEof() bool {
	return s.PeekChar().Is(0)
}

// Return the next character in the input and consume it.
func (s *ByteScanner) EatChar() Char {
	c := s.PeekChar()
	if len(s.queue) > 0 {
		s.queue = s.queue[1:]
	}
	return c
}

// Return the next N characters in the input and consume them.
func (s *ByteScanner) EatChars(n int) []Char {
	chars := make([]Char, 0)
	for i := 0; i < n; i++ {
		chars = append(chars, s.EatChar())
	}
	return chars
}

// Return the next character in the input without consuming it.
func (s *ByteScanner) PeekChar() Char {
	return s.PeekCharAt(0)
}

// Return the Nth character in the input without consuming it.
func (s *ByteScanner) PeekCharAt(offset int) Char {
	for len(s.queue) <= offset {
		char, err := s.next()
		if err != nil {
			s.RegisterError(NewError(
				Loc{Start: char.Span, End: char.Span},
				ErrIO,
				err.Error(),
			))
		}
		s.queue = append(s.queue, char)
	}
	return s.queue[offset]
}

func (s *ByteScanner) next() (Char, error) {
	if s.eof != nil {
		return *s.eof, nil
	}

	if s.cursor >= len(s.input) {
		s.eof = &Char{0, 0, Span{
			Offset: s.cursor,
			Line:   s.line,
			Column: s.column,
		}}
		return *s.eof, nil
	}

	r, size := utf8.DecodeRune(s.input[s.cursor:])
	if r == utf8.RuneError {
		s.cursor += size
		s.column++
		c := Char{0, 0, Span{
			Offset: s.cursor,
			Line:   s.line,
			Column: s.column,
		}}
		return c, c.AsError(ErrIO, errMsgInvalidChar)
	}

	c := Char{r, size, Span{
		Offset: s.cursor,
		Line:   s.line,
		Column: s.column,
	}}
	s.cursor += size
	s.column++
	if r == '\n' {
		s.line++
		s.column = 1
	}
	return c, nil
}
