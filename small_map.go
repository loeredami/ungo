package ungo

import (
	"unsafe"
)

type entry[K comparable, V any] struct {
	key   K
	value V
	used  bool
}

type SmallMap[K comparable, V any] struct {
	data []entry[K, V]
	size int
	mask uintptr
}

func NewSmallMap[K comparable, V any](capacity int) *SmallMap[K, V] {
	pow2 := 8
	for pow2 < capacity {
		pow2 <<= 1
	}
	return &SmallMap[K, V]{
		data: make([]entry[K, V], pow2),
		mask: uintptr(pow2 - 1),
	}
}

func (fm *SmallMap[K, V]) hash(key K) uintptr {
	u := *(*uintptr)(unsafe.Pointer(&key))
	return u * 0xdeece66d
}

func (fm *SmallMap[K, V]) Set(key K, value V) {
	if fm.size*2 >= len(fm.data) {
		fm.grow()
	}

	idx := fm.hash(key) & fm.mask
	for {
		if !fm.data[idx].used {
			fm.data[idx] = entry[K, V]{key: key, value: value, used: true}
			fm.size++
			return
		}
		if fm.data[idx].key == key {
			fm.data[idx].value = value
			return
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *SmallMap[K, V]) Get(key K) (V, bool) {
	if len(fm.data) == 0 {
		var zero V
		return zero, false
	}

	idx := fm.hash(key) & fm.mask
	for {
		e := &fm.data[idx]
		if !e.used {
			var zero V
			return zero, false
		}
		if e.key == key {
			return e.value, true
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *SmallMap[K, V]) Delete(key K) {
	idx := fm.hash(key) & fm.mask
	for {
		if !fm.data[idx].used {
			return
		}
		if fm.data[idx].key == key {
			fm.data[idx].used = false
			fm.size--
			fm.rehashChain(idx)
			return
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *SmallMap[K, V]) rehashChain(i uintptr) {
	j := i
	for {
		j = (j + 1) & fm.mask
		if !fm.data[j].used {
			return
		}
		r := fm.hash(fm.data[j].key) & fm.mask
		if (j > i && (r <= i || r > j)) || (j < i && (r <= i && r > j)) {
			fm.data[i] = fm.data[j]
			i = j
			fm.data[i].used = false
		}
	}
}

func (fm *SmallMap[K, V]) grow() {
	oldData := fm.data
	newCap := len(oldData) * 2
	fm.data = make([]entry[K, V], newCap)
	fm.mask = uintptr(newCap - 1)
	fm.size = 0

	for i := range oldData {
		if oldData[i].used {
			fm.Set(oldData[i].key, oldData[i].value)
		}
	}
}

func (fm *SmallMap[K, V]) Size() int {
	return fm.size
}

func (fm *SmallMap[K, V]) ForEach(f func(key K, value V)) {
	for i := range fm.data {
		if fm.data[i].used {
			f(fm.data[i].key, fm.data[i].value)
		}
	}
}

func (fm *SmallMap[K, V]) Keys() []K {
	res := make([]K, 0, fm.size)
	for i := range fm.data {
		if fm.data[i].used {
			res = append(res, fm.data[i].key)
		}
	}
	return res
}

func (fm *SmallMap[K, V]) Values() []V {
	res := make([]V, 0, fm.size)
	for i := range fm.data {
		if fm.data[i].used {
			res = append(res, fm.data[i].value)
		}
	}
	return res
}

func (fm *SmallMap[K, V]) Clear() {
	for i := range fm.data {
		fm.data[i] = entry[K, V]{}
	}
	fm.size = 0
}

func (fm *SmallMap[K, V]) Contains(key K) bool {
	idx := fm.hash(key) & fm.mask
	for {
		if !fm.data[idx].used {
			return false
		}
		if fm.data[idx].key == key {
			return true
		}
		idx = (idx + 1) & fm.mask
	}
}
