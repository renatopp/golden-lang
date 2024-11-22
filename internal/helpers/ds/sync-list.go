package ds

import "sync"

type SyncList[T comparable] struct {
	list []T
	mtx  sync.RWMutex
}

func NewSyncList[T comparable]() *SyncList[T] {
	return &SyncList[T]{list: make([]T, 0)}
}

func (l *SyncList[T]) Get(index int) (value T, ok bool) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	if index < 0 || index >= len(l.list) {
		return
	}
	return l.list[index], true
}

func (l *SyncList[T]) Add(value T) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.list = append(l.list, value)
}

func (l *SyncList[T]) AddUnique(value T) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	for _, v := range l.list {
		if v == value {
			return
		}
	}
	l.list = append(l.list, value)
}

func (l *SyncList[T]) Delete(index int) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if index < 0 || index >= len(l.list) {
		return
	}
	l.list = append(l.list[:index], l.list[index+1:]...)
}

func (l *SyncList[T]) DeleteValue(value T) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	for i, v := range l.list {
		if v == value {
			l.list = append(l.list[:i], l.list[i+1:]...)
			return
		}
	}
}

func (l *SyncList[T]) IndexOf(value T) int {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	for i, v := range l.list {
		if v == value {
			return i
		}
	}
	return -1
}

func (l *SyncList[T]) Has(value T) bool {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	for _, v := range l.list {
		if v == value {
			return true
		}
	}
	return false
}

func (l *SyncList[T]) Len() int {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	return len(l.list)
}

func (l *SyncList[T]) Values() []T {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	values := make([]T, len(l.list))
	copy(values, l.list)
	return values
}

func (l *SyncList[T]) Clear() {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.list = make([]T, 0)
}
