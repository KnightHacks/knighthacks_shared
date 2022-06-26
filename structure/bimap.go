package structure

import "sync"

type BiMap struct {
	mutex   sync.RWMutex
	Forward map[any]any
	Reverse map[any]any
}

func NewBiMap() *BiMap {
	return &BiMap{}
}

func (m *BiMap) Put(key any, value any) {
	m.mutex.Lock()
	m.Forward[key] = value
	m.Reverse[value] = key
	m.mutex.Unlock()
}

func (m *BiMap) Get(key any) any {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if value, ok := m.Forward[key]; ok {
		return value
	}
	return m.Reverse[key]
}

func (m *BiMap) Delete(key any, value any) {
	m.mutex.Lock()
	delete(m.Forward, key)
	delete(m.Forward, value)
	m.mutex.Unlock()
}
