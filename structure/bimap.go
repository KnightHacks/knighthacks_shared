package structure

type BiMap struct {
	Forward map[any]any
	Reverse map[any]any
}

func NewBiMap() *BiMap {
	return &BiMap{}
}

func (m *BiMap) Put(key any, value any) {
	m.Forward[key] = value
	m.Reverse[value] = key
}

func (m *BiMap) Get(key any) any {
	if value, ok := m.Forward[key]; ok {
		return value
	}
	return m.Reverse[key]
}

func (m *BiMap) Delete(key any, value any) {
	delete(m.Forward, key)
	delete(m.Forward, value)
}
