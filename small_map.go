package ungo

type SmallMap[K comparable, V any] struct {
	keys   []K
	values []V
}

func NewSmallMap[K comparable, V any](capacity int) *SmallMap[K, V] {
	return &SmallMap[K, V]{
		keys:   make([]K, 0, capacity),
		values: make([]V, 0, capacity),
	}
}

func (fm *SmallMap[K, V]) Set(key K, value V) {
	for i, k := range fm.keys {
		if k == key {
			fm.values[i] = value
			return
		}
	}
	fm.keys = append(fm.keys, key)
	fm.values = append(fm.values, value)
}

func (fm *SmallMap[K, V]) Get(key K) (V, bool) {
	for i, k := range fm.keys {
		if k == key {
			return fm.values[i], true
		}
	}
	var zero V
	return zero, false
}

func (fm *SmallMap[K, V]) Delete(key K) {
	for i, k := range fm.keys {
		if k == key {
			lastIdx := len(fm.keys) - 1

			// Swap with last to maintain O(1) delete
			fm.keys[i] = fm.keys[lastIdx]
			fm.values[i] = fm.values[lastIdx]

			// Clear last elements to avoid memory leaks
			var zeroK K
			var zeroV V
			fm.keys[lastIdx] = zeroK
			fm.values[lastIdx] = zeroV

			fm.keys = fm.keys[:lastIdx]
			fm.values = fm.values[:lastIdx]
			return
		}
	}
}

func (fm *SmallMap[K, V]) Size() int {
	return len(fm.keys)
}

func (fm *SmallMap[K, V]) ForEach(f func(key K, value V)) {
	for i, k := range fm.keys {
		f(k, fm.values[i])
	}
}

func (fm *SmallMap[K, V]) Keys() []K {
	res := make([]K, len(fm.keys))
	copy(res, fm.keys)
	return res
}

func (fm *SmallMap[K, V]) Values() []V {
	res := make([]V, len(fm.values))
	copy(res, fm.values)
	return res
}

func (fm *SmallMap[K, V]) Clear() {
	for i := range fm.keys {
		var zeroK K
		var zeroV V
		fm.keys[i] = zeroK
		fm.values[i] = zeroV
	}
	fm.keys = fm.keys[:0]
	fm.values = fm.values[:0]
}

func (fm *SmallMap[K, V]) Contains(key K) bool {
	for _, k := range fm.keys {
		if k == key {
			return true
		}
	}
	return false
}
