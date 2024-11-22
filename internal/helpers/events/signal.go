package events

import "reflect"

type Signal struct {
	subscribers     []func()
	onceSubscribers []func()
}

func NewSignal() *Signal {
	return &Signal{
		subscribers:     []func(){},
		onceSubscribers: []func(){},
	}
}

func (s *Signal) Subscribe(fn func()) {
	s.subscribers = append(s.subscribers, fn)
}

func (s *Signal) SubscribeOnce(fn func()) {
	s.onceSubscribers = append(s.onceSubscribers, fn)
}

func (s *Signal) Unsubscribe(fn func()) {
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

func (s *Signal) Emit() {
	for _, fn := range s.subscribers {
		fn()
	}

	for _, fn := range s.onceSubscribers {
		fn()
	}
	s.onceSubscribers = []func(){}
}

func (s *Signal) Clear() {
	s.subscribers = []func(){}
	s.onceSubscribers = []func(){}
}
