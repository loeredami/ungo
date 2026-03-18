package ungo

type entry[K comparable, V any] struct {
	key   K
	value V
}

type SmallMap[K comparable, V any] struct {
	metadata []byte
	data     []entry[K, V]
	size     int
}

func NewSmallMap[K comparable, V any](capacity int) *SmallMap[K, V] {
	return &SmallMap[K, V]{
		metadata: make([]byte, capacity),
		data:     make([]entry[K, V], capacity),
	}
}

func (fm *SmallMap[K, V]) hash(key K) byte {
	h := uint8(0)
	return (h % 255) + 1
}

func (fm *SmallMap[K, V]) Set(key K, value V) {
	h := fm.hash(key)

	for i := 0; i < len(fm.metadata); i++ {
		if fm.metadata[i] == h {
			if fm.data[i].key == key {
				fm.data[i].value = value
				return
			}
		} else if fm.metadata[i] == 0 {
			fm.metadata[i] = h
			fm.data[i] = entry[K, V]{key: key, value: value}
			fm.size++
			return
		}
	}

	fm.grow()
	fm.Set(key, value)
}

func (fm *SmallMap[K, V]) Get(key K) (V, bool) {
	h := fm.hash(key)
	md := fm.metadata

	for i := 0; i < len(md); i++ {
		m := md[i]
		if m == 0 {
			break
		}
		if m == h && fm.data[i].key == key {
			return fm.data[i].value, true
		}
	}

	var zero V
	return zero, false
}

func (fm *SmallMap[K, V]) grow() {
	newCap := len(fm.data) * 2
	if newCap == 0 {
		newCap = 8
	}

	oldData := fm.data
	oldMD := fm.metadata

	fm.data = make([]entry[K, V], newCap)
	fm.metadata = make([]byte, newCap)
	fm.size = 0

	for i := 0; i < len(oldMD); i++ {
		if oldMD[i] != 0 {
			fm.Set(oldData[i].key, oldData[i].value)
		}
	}
}

func (fm *SmallMap[K, V]) Delete(key K) {
	h := fm.hash(key)
	for i := 0; i < len(fm.metadata); i++ {
		if fm.metadata[i] == h && fm.data[i].key == key {
			copy(fm.metadata[i:], fm.metadata[i+1:])
			copy(fm.data[i:], fm.data[i+1:])

			last := len(fm.metadata) - 1
			fm.metadata[last] = 0
			fm.data[last] = entry[K, V]{}
			fm.size--
			return
		}
	}
}

func (fm *SmallMap[K, V]) Size() int {
	return fm.size
}

func (fm *SmallMap[K, V]) ForEach(f func(key K, value V)) {
	for i, m := range fm.metadata {
		if m != 0 {
			f(fm.data[i].key, fm.data[i].value)
		}
	}
}

func (fm *SmallMap[K, V]) Keys() []K {
	res := make([]K, 0, fm.size)
	for i, m := range fm.metadata {
		if m != 0 {
			res = append(res, fm.data[i].key)
		}
	}
	return res
}

func (fm *SmallMap[K, V]) Values() []V {
	res := make([]V, 0, fm.size)
	for i, m := range fm.metadata {
		if m != 0 {
			res = append(res, fm.data[i].value)
		}
	}
	return res
}

func (fm *SmallMap[K, V]) Clear() {
	for i := range fm.metadata {
		fm.metadata[i] = 0
		fm.data[i] = entry[K, V]{}
	}
	fm.size = 0
}

func (fm *SmallMap[K, V]) Contains(key K) bool {
	_, found := fm.Get(key)
	return found
}
