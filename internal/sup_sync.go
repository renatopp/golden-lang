package internal

import "sync"

type SyncMap[K comparable, V any] struct {
	m   map[K]V
	mtx sync.RWMutex
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{m: make(map[K]V)}
}

func (m *SyncMap[K, V]) Get(key K) (value V, ok bool) {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	value, ok = m.m[key]
	return
}

func (m *SyncMap[K, V]) GetOr(key K, def V) V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	value, ok := m.Get(key)
	if !ok {
		return def
	}
	return value
}

func (m *SyncMap[K, V]) Set(key K, value V) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.m[key] = value
}

func (m *SyncMap[K, V]) Delete(key K) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	delete(m.m, key)
}

func (m *SyncMap[K, V]) Has(key K) bool {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	_, ok := m.m[key]
	return ok
}

func (m *SyncMap[K, V]) Len() int {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	return len(m.m)
}

func (m *SyncMap[K, V]) Values() []V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	values := make([]V, 0, len(m.m))
	for _, value := range m.m {
		values = append(values, value)
	}
	return values
}

func (m *SyncMap[K, V]) Keys() []K {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	keys := make([]K, 0, len(m.m))
	for key := range m.m {
		keys = append(keys, key)
	}
	return keys
}

func (m *SyncMap[K, V]) Items() map[K]V {
	m.mtx.RLock()
	defer m.mtx.RUnlock()
	items := make(map[K]V, len(m.m))
	for key, value := range m.m {
		items[key] = value
	}
	return items
}

func (m *SyncMap[K, V]) Clear() {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.m = make(map[K]V)
}

type SyncList[T comparable] struct {
	l   []T
	mtx sync.RWMutex
}

func NewSyncList[T comparable]() *SyncList[T] {
	return &SyncList[T]{l: make([]T, 0)}
}

func (l *SyncList[T]) Append(value T) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	l.l = append(l.l, value)
}

func (l *SyncList[T]) Get(index int) (value T, ok bool) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	if index < 0 || index >= len(l.l) {
		return
	}
	return l.l[index], true
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

func (l *SyncList[T]) Remove(index int) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if index < 0 || index >= len(l.l) {
		return
	}
	l.l = append(l.l[:index], l.l[index+1:]...)
}

func (l *SyncList[T]) RemoveValue(value T) {
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

func (l *SyncList[T]) Contains(value T) bool {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	for _, v := range l.l {
		if v == value {
			return true
		}
	}
	return false
}
