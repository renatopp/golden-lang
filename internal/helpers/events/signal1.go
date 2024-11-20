package events

import "reflect"

type Signal1[T any] struct {
	subscribers     []func(T)
	onceSubscribers []func(T)
}

func NewSignal1[T any]() *Signal1[T] {
	return &Signal1[T]{
		subscribers:     []func(T){},
		onceSubscribers: []func(T){},
	}
}

func (s *Signal1[T]) Subscribe(fn func(T)) {
	s.subscribers = append(s.subscribers, fn)
}

func (s *Signal1[T]) SubscribeOnce(fn func(T)) {
	s.onceSubscribers = append(s.onceSubscribers, fn)
}

func (s *Signal1[T]) Unsubscribe(fn func(T)) {
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

func (s *Signal1[T]) Emit(a T) {
	for _, fn := range s.subscribers {
		fn(a)
	}

	for _, fn := range s.onceSubscribers {
		fn(a)
	}
	s.onceSubscribers = []func(T){}
}

func (s *Signal1[T]) Clear() {
	s.subscribers = []func(T){}
	s.onceSubscribers = []func(T){}
}
