package ds

import "iter"

type Stack[T any] struct {
	data []*T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{data: []*T{}}
}

func (s *Stack[T]) Push(value *T) {
	s.data = append(s.data, value)
}

func (s *Stack[T]) Pop() *T {
	if len(s.data) <= 0 {
		return nil
	}

	val := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return val
}

func (s *Stack[T]) Len() int {
	return len(s.data)
}

func (s *Stack[T]) Top() *T {
	if len(s.data) <= 0 {
		return nil
	}

	val := s.data[len(s.data)-1]
	return val
}

func (s *Stack[T]) Has(v *T) bool {
	for _, val := range s.data {
		if val == v {
			return true
		}
	}
	return false
}

func (s *Stack[T]) Iter() iter.Seq2[int, *T] {
	return func(yield func(int, *T) bool) {
		for i, v := range s.data {
			if !yield(i, v) {
				return
			}
		}
	}
}

func (s *Stack[T]) ReverseIter() iter.Seq2[int, *T] {
	return func(yield func(int, *T) bool) {
		size := len(s.data)
		for i := 0; i < size; i++ {
			if !yield(size-i, s.data[i]) {
				return
			}
		}
	}
}
