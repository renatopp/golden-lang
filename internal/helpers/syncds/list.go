package syncds

import "sync"

type SyncList[T comparable] struct {
	l   []T
	mtx sync.RWMutex
}

func NewSyncList[T comparable]() *SyncList[T] {
	return &SyncList[T]{l: make([]T, 0)}
}

func (l *SyncList[T]) Get(index int) (value T, ok bool) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	if index < 0 || index >= len(l.l) {
		return
	}
	return l.l[index], true
}

func (l *SyncList[T]) Add(value T) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.l = append(l.l, value)
}

func (l *SyncList[T]) Delete(index int) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if index < 0 || index >= len(l.l) {
		return
	}
	l.l = append(l.l[:index], l.l[index+1:]...)
}

func (l *SyncList[T]) DeleteValue(value T) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	for i, v := range l.l {
		if v == value {
			l.l = append(l.l[:i], l.l[i+1:]...)
			return
		}
	}
}

func (l *SyncList[T]) IndexOf(value T) int {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	for i, v := range l.l {
		if v == value {
			return i
		}
	}
	return -1
}

func (l *SyncList[T]) Has(value T) bool {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	for _, v := range l.l {
		if v == value {
			return true
		}
	}
	return false
}

func (l *SyncList[T]) Len() int {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	return len(l.l)
}

func (l *SyncList[T]) Values() []T {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	values := make([]T, len(l.l))
	copy(values, l.l)
	return values
}

func (l *SyncList[T]) Clear() {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.l = make([]T, 0)
}
