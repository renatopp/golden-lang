package syntax

type Scanner[T any] struct {
	source []T
	none   T   // Empty value if scanner is out of bounds
	cursor int // Cursor position
}

func NewScanner[T any](source []T, none T) *Scanner[T] {
	return &Scanner[T]{source: source, none: none, cursor: 0}
}

func (s *Scanner[T]) Reset() { s.cursor = 0 }
func (s *Scanner[T]) IsFinished() bool {
	return s.cursor >= len(s.source)
}
func (s *Scanner[T]) Eat() T {
	if s.IsFinished() {
		return s.none
	}
	value := s.source[s.cursor]
	s.cursor++
	return value
}
func (s *Scanner[T]) EatN(n int) []T {
	values := make([]T, n)
	for i := 0; i < n; i++ {
		values[i] = s.Eat()
	}
	return values
}
func (s *Scanner[T]) Peek() T { return s.PeekAt(0) }
func (s *Scanner[T]) PeekAt(offset int) T {
	if s.cursor+offset >= len(s.source) {
		return s.none
	}
	return s.source[s.cursor+offset]
}
