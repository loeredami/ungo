package ungo

type entry[K comparable, V any] struct {
	key   K
	value V
}

type SmallMap[K comparable, V any] struct {
	data []entry[K, V]
}

func NewSmallMap[K comparable, V any](capacity int) *SmallMap[K, V] {
	return &SmallMap[K, V]{
		data: make([]entry[K, V], 0, capacity),
	}
}

func (fm *SmallMap[K, V]) Set(key K, value V) {
	for i := 0; i < len(fm.data); i++ {
		if fm.data[i].key == key {
			fm.data[i].value = value
			return
		}
	}
	fm.data = append(fm.data, entry[K, V]{key: key, value: value})
}

func (fm *SmallMap[K, V]) Get(key K) (V, bool) {
	d := fm.data
	for i := 0; i < len(d); i++ {
		if d[i].key == key {
			return d[i].value, true
		}
	}
	var zero V
	return zero, false
}

func (fm *SmallMap[K, V]) Delete(key K) {
	d := fm.data
	for i := 0; i < len(d); i++ {
		if d[i].key == key {
			lastIdx := len(d) - 1
			d[i] = d[lastIdx]

			d[lastIdx] = entry[K, V]{}

			fm.data = d[:lastIdx]
			return
		}
	}
}

func (fm *SmallMap[K, V]) Size() int {
	return len(fm.data)
}

func (fm *SmallMap[K, V]) ForEach(f func(key K, value V)) {
	d := fm.data
	for i := 0; i < len(d); i++ {
		f(d[i].key, d[i].value)
	}
}

func (fm *SmallMap[K, V]) Keys() []K {
	res := make([]K, len(fm.data))
	for i := 0; i < len(fm.data); i++ {
		res[i] = fm.data[i].key
	}
	return res
}

func (fm *SmallMap[K, V]) Values() []V {
	res := make([]V, len(fm.data))
	for i := 0; i < len(fm.data); i++ {
		res[i] = fm.data[i].value
	}
	return res
}

func (fm *SmallMap[K, V]) Clear() {
	for i := range fm.data {
		fm.data[i] = entry[K, V]{}
	}
	fm.data = fm.data[:0]
}

func (fm *SmallMap[K, V]) Contains(key K) bool {
	d := fm.data
	for i := 0; i < len(d); i++ {
		if d[i].key == key {
			return true
		}
	}
	return false
}
