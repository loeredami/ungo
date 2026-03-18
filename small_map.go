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
		// We allocate +1 for the sentinel record
		data: make([]entry[K, V], 0, capacity+1),
	}
}

func (fm *SmallMap[K, V]) Set(key K, value V) {
	d := fm.data
	for i := 0; i < len(d); i++ {
		if d[i].key == key {
			d[i].value = value
			return
		}
	}
	fm.data = append(fm.data, entry[K, V]{key: key, value: value})
}

func (fm *SmallMap[K, V]) Get(key K) (V, bool) {
	d := fm.data
	n := len(d)
	if n == 0 {
		var zero V
		return zero, false
	}

	lastOriginal := d[n-1]
	if lastOriginal.key == key {
		return lastOriginal.value, true
	}

	fm.data = fm.data[:n+1]
	fm.data[n].key = key

	i := 0
	for fm.data[i].key != key {
		i++
	}

	fm.data = fm.data[:n]

	if i < n {
		return d[i].value, true
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
	d := fm.data
	res := make([]K, len(d))
	for i := 0; i < len(d); i++ {
		res[i] = d[i].key
	}
	return res
}

func (fm *SmallMap[K, V]) Values() []V {
	d := fm.data
	res := make([]V, len(d))
	for i := 0; i < len(d); i++ {
		res[i] = d[i].value
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
	_, found := fm.Get(key)
	return found
}
