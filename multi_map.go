package ungo

type MultiMap[K comparable, V any] struct {
	items map[K][]V
}

func NewMultiMap[K comparable, V any]() *MultiMap[K, V] {
	return &MultiMap[K, V]{
		items: make(map[K][]V),
	}
}

func (m *MultiMap[K, V]) Add(key K, value V) {
	m.items[key] = append(m.items[key], value)
}

func (m *MultiMap[K, V]) Get(key K) []V {
	return m.items[key]
}

func (m *MultiMap[K, V]) Delete(key K) {
	delete(m.items, key)
}

func (m *MultiMap[K, V]) Clear() {
	m.items = make(map[K][]V)
}

func (m *MultiMap[K, V]) Keys() []K {
	keys := make([]K, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	return keys
}
func (m *MultiMap[K, V]) Values() [][]V {
	values := make([][]V, 0, len(m.items))
	for _, v := range m.items {
		values = append(values, v)
	}
	return values
}

func (m *MultiMap[K, V]) Len() int {
	return len(m.items)
}

func (m *MultiMap[K, V]) Contains(key K) bool {
	_, ok := m.items[key]
	return ok
}
