package events

import "reflect"

type Signal3[T, R, V any] struct {
	subscribers     []func(T, R, V)
	onceSubscribers []func(T, R, V)
}

func NewSignal3[T, R, V any]() *Signal3[T, R, V] {
	return &Signal3[T, R, V]{
		subscribers:     []func(T, R, V){},
		onceSubscribers: []func(T, R, V){},
	}
}

func (s *Signal3[T, R, V]) Subscribe(fn func(T, R, V)) {
	s.subscribers = append(s.subscribers, fn)
}

func (s *Signal3[T, R, V]) SubscribeOnce(fn func(T, R, V)) {
	s.onceSubscribers = append(s.onceSubscribers, fn)
}

func (s *Signal3[T, R, V]) Unsubscribe(fn func(T, R, V)) {
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

func (s *Signal3[T, R, V]) Emit(a T, b R, c V) {
	for _, fn := range s.subscribers {
		fn(a, b, c)
	}

	for _, fn := range s.onceSubscribers {
		fn(a, b, c)
	}
	s.onceSubscribers = []func(T, R, V){}
}

func (s *Signal3[T, R, V]) Clear() {
	s.subscribers = []func(T, R, V){}
	s.onceSubscribers = []func(T, R, V){}
}
