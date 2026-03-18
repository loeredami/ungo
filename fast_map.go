package ungo

// FastMap is a map that uses a slice of keys and a slice of values for fast lookups.
type FastMap[K comparable, V any] struct {
	indexes []K
	values  []V
}

func (fm *FastMap[K, V]) Set(key K, value V) {
	for i, k := range fm.indexes {
		if k == key {
			fm.values[i] = value
			return
		}
	}
	fm.indexes = append(fm.indexes, key)
	fm.values = append(fm.values, value)
}

func (fm *FastMap[K, V]) Get(key K) (V, bool) {
	for i, k := range fm.indexes {
		if k == key {
			return fm.values[i], true
		}
	}
	var zero V
	return zero, false
}

func (fm *FastMap[K, V]) Delete(key K) {
	for i, k := range fm.indexes {
		if k == key {
			fm.indexes = append(fm.indexes[:i], fm.indexes[i+1:]...)
			fm.values = append(fm.values[:i], fm.values[i+1:]...)
			return
		}
	}
}

func (fm *FastMap[K, V]) Size() int {
	return len(fm.indexes)
}

func (fm *FastMap[K, V]) ForEach(f func(key K, value V)) {
	for i, k := range fm.indexes {
		f(k, fm.values[i])
	}
}

func (fm *FastMap[K, V]) Keys() []K {
	return fm.indexes
}

func (fm *FastMap[K, V]) Values() []V {
	return fm.values
}

func (fm *FastMap[K, V]) Clear() {
	fm.indexes = make([]K, 0)
	fm.values = make([]V, 0)
}

func (fm *FastMap[K, V]) Contains(key K) bool {
	for _, k := range fm.indexes {
		if k == key {
			return true
		}
	}
	return false
}
