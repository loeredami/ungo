package ungo

type MultiMap[K comparable, V any] struct {
	items FastMap[K, []V]
}

func NewMultiMap[K comparable, V any]() *MultiMap[K, V] {
	return &MultiMap[K, V]{
		items: FastMap[K, []V]{},
	}
}

func (m *MultiMap[K, V]) Add(key K, value V) {
	v, ok := m.items.Get(key)
	if !ok {
		v = []V{}
	}
	m.items.Set(key, append(v, value))
}

func (m *MultiMap[K, V]) Get(key K) []V {
	k, ok := m.items.Get(key)
	if !ok {
		return []V{}
	}
	return k
}

func (m *MultiMap[K, V]) Delete(key K) {
	m.items.Delete(key)
}

func (m *MultiMap[K, V]) Clear() {
	m.items.Clear()
}

func (m *MultiMap[K, V]) Keys() []K {
	keys := make([]K, 0, m.items.Size())
	for _, k := range m.items.Keys() {
		keys = append(keys, k)
	}
	return keys
}
func (m *MultiMap[K, V]) Values() [][]V {
	values := make([][]V, 0, m.items.Size())
	for _, v := range m.items.Values() {
		values = append(values, v)
	}
	return values
}

func (m *MultiMap[K, V]) Len() int {
	return m.items.Size()
}

func (m *MultiMap[K, V]) Contains(key K) bool {
	return m.items.Contains(key)
}
