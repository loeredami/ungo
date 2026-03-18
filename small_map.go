package ungo

type entry[K comparable, V any] struct {
	key   K
	value V
}

type SmallMap[K comparable, V any] struct {
	data []entry[K, V]

	overflow map[K]V
}

func NewSmallMap[K comparable, V any](capacity int) *SmallMap[K, V] {
	return &SmallMap[K, V]{
		data: make([]entry[K, V], 0, capacity),
	}
}

func (fm *SmallMap[K, V]) Set(key K, value V) {
	if fm.overflow != nil {
		fm.overflow[key] = value
		return
	}

	d := fm.data
	for i := 0; i < len(d); i++ {
		if d[i].key == key {
			d[i].value = value
			return
		}
	}

	if len(d) >= 24 {
		fm.overflow = make(map[K]V, len(d)+1)
		for i := 0; i < len(d); i++ {
			fm.overflow[d[i].key] = d[i].value
		}
		fm.overflow[key] = value
		fm.data = nil
		return
	}

	fm.data = append(fm.data, entry[K, V]{key: key, value: value})
}

func (fm *SmallMap[K, V]) Get(key K) (V, bool) {
	if fm.overflow != nil {
		v, ok := fm.overflow[key]
		return v, ok
	}

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
	if fm.overflow != nil {
		delete(fm.overflow, key)
		return
	}

	d := fm.data
	for i := 0; i < len(d); i++ {
		if d[i].key == key {
			lastIdx := len(d) - 1
			d[i] = d[lastIdx]

			var zero entry[K, V]
			d[lastIdx] = zero

			fm.data = d[:lastIdx]
			return
		}
	}
}

func (fm *SmallMap[K, V]) Size() int {
	if fm.overflow != nil {
		return len(fm.overflow)
	}
	return len(fm.data)
}

func (fm *SmallMap[K, V]) ForEach(f func(key K, value V)) {
	if fm.overflow != nil {
		for k, v := range fm.overflow {
			f(k, v)
		}
		return
	}

	d := fm.data
	for i := 0; i < len(d); i++ {
		f(d[i].key, d[i].value)
	}
}

func (fm *SmallMap[K, V]) Keys() []K {
	if fm.overflow != nil {
		res := make([]K, 0, len(fm.overflow))
		for k := range fm.overflow {
			res = append(res, k)
		}
		return res
	}

	res := make([]K, len(fm.data))
	for i := 0; i < len(fm.data); i++ {
		res[i] = fm.data[i].key
	}
	return res
}

func (fm *SmallMap[K, V]) Values() []V {
	if fm.overflow != nil {
		res := make([]V, 0, len(fm.overflow))
		for _, v := range fm.overflow {
			res = append(res, v)
		}
		return res
	}

	res := make([]V, len(fm.data))
	for i := 0; i < len(fm.data); i++ {
		res[i] = fm.data[i].value
	}
	return res
}

func (fm *SmallMap[K, V]) Clear() {
	if fm.overflow != nil {
		for k := range fm.overflow {
			delete(fm.overflow, k)
		}
	}

	for i := range fm.data {
		fm.data[i] = entry[K, V]{}
	}
	fm.data = fm.data[:0]
}

func (fm *SmallMap[K, V]) Contains(key K) bool {
	if fm.overflow != nil {
		_, ok := fm.overflow[key]
		return ok
	}

	d := fm.data
	for i := 0; i < len(d); i++ {
		if d[i].key == key {
			return true
		}
	}
	return false
}
