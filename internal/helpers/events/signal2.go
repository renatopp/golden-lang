package events

import "reflect"

type Signal2[T, R any] struct {
	subscribers     []func(T, R)
	onceSubscribers []func(T, R)
}

func NewSignal2[T, R any]() *Signal2[T, R] {
	return &Signal2[T, R]{
		subscribers:     []func(T, R){},
		onceSubscribers: []func(T, R){},
	}
}

func (s *Signal2[T, R]) Subscribe(fn func(T, R)) {
	s.subscribers = append(s.subscribers, fn)
}

func (s *Signal2[T, R]) SubscribeOnce(fn func(T, R)) {
	s.onceSubscribers = append(s.onceSubscribers, fn)
}

func (s *Signal2[T, R]) Unsubscribe(fn func(T, R)) {
	fnPointer := reflect.ValueOf(fn).Pointer()

	for i, f := range s.subscribers {
		fPointer := reflect.ValueOf(f).Pointer()
		if fPointer == fnPointer {
			s.subscribers = append(s.subscribers[:i], s.subscribers[i+1:]...)
			return
		}
	}

	for i, f := range s.onceSubscribers {
		fPointer := reflect.ValueOf(f).Pointer()
		if fPointer == fnPointer {
			s.onceSubscribers = append(s.onceSubscribers[:i], s.onceSubscribers[i+1:]...)
			return
		}
	}
}

func (s *Signal2[T, R]) Emit(a T, b R) {
	for _, fn := range s.subscribers {
		fn(a, b)
	}

	for _, fn := range s.onceSubscribers {
		fn(a, b)
	}
	s.onceSubscribers = []func(T, R){}
}

func (s *Signal2[T, R]) Clear() {
	s.subscribers = []func(T, R){}
	s.onceSubscribers = []func(T, R){}
}
