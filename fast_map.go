package ungo

type FastMap[K comparable, V any] struct {
	indexes []K
	values  []V
}

func (fm *FastMap[K, V]) Set(key K, value V) {
	for i := 0; i < len(fm.indexes); i++ {
		if fm.indexes[i] == key {
			fm.values[i] = value
			return
		}
	}
	fm.indexes = append(fm.indexes, key)
	fm.values = append(fm.values, value)
}

func (fm *FastMap[K, V]) Get(key K) (V, bool) {
	idx := fm.indexes
	for i := 0; i < len(idx); i++ {
		if idx[i] == key {
			return fm.values[i], true
		}
	}
	var zero V
	return zero, false
}

func (fm *FastMap[K, V]) Delete(key K) {
	for i := 0; i < len(fm.indexes); i++ {
		if fm.indexes[i] == key {
			lastIdx := len(fm.indexes) - 1

			copy(fm.indexes[i:], fm.indexes[i+1:])
			copy(fm.values[i:], fm.values[i+1:])

			var zeroK K
			var zeroV V
			fm.indexes[lastIdx] = zeroK
			fm.values[lastIdx] = zeroV

			fm.indexes = fm.indexes[:lastIdx]
			fm.values = fm.values[:lastIdx]
			return
		}
	}
}

func (fm *FastMap[K, V]) Size() int {
	return len(fm.indexes)
}

func (fm *FastMap[K, V]) ForEach(f func(key K, value V)) {
	for i := 0; i < len(fm.indexes); i++ {
		f(fm.indexes[i], fm.values[i])
	}
}

func (fm *FastMap[K, V]) Keys() []K {
	return fm.indexes
}

func (fm *FastMap[K, V]) Values() []V {
	return fm.values
}

func (fm *FastMap[K, V]) Clear() {
	var zeroK K
	var zeroV V
	for i := range fm.indexes {
		fm.indexes[i] = zeroK
		fm.values[i] = zeroV
	}
	fm.indexes = fm.indexes[:0]
	fm.values = fm.values[:0]
}

func (fm *FastMap[K, V]) Contains(key K) bool {
	idx := fm.indexes
	for i := 0; i < len(idx); i++ {
		if idx[i] == key {
			return true
		}
	}
	return false
}
