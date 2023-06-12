package structure

import "sync"

type BiMap[K comparable, V comparable] struct {
	mutex   sync.RWMutex
	Forward map[K]V
	Reverse map[V]K
}

func NewBiMap[K comparable, V comparable]() *BiMap[K, V] {
	return &BiMap[K, V]{
		mutex:   sync.RWMutex{},
		Forward: make(map[K]V),
		Reverse: make(map[V]K),
	}
}

func (m *BiMap[K, V]) Put(key K, value V) {
	m.mutex.Lock()
	m.Forward[key] = value
	m.Reverse[value] = key
	m.mutex.Unlock()
}

func (m *BiMap[K, V]) Get(key any) any {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if value, ok := m.Forward[key]; ok {
		return value
	}
	if value, ok := m.Reverse[key]; ok {
		return value
	}
	return nil
}

func (m *BiMap[K, V]) Delete(key K, value V) {
	m.mutex.Lock()
	delete(m.Forward, key)
	delete(m.Reverse, value)
	m.mutex.Unlock()
}
