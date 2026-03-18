package ungo

import (
	"fmt"
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
	pow2 := 16
	for pow2 < capacity {
		pow2 <<= 1
	}
	return &SmallMap[K, V]{
		data: make([]entry[K, V], pow2),
		mask: uintptr(pow2 - 1),
	}
}

func (fm *SmallMap[K, V]) hash(key K) uintptr {
	var h uintptr = 2166136261
	s := fmt.Sprintf("%v", key)
	for i := 0; i < len(s); i++ {
		h ^= uintptr(s[i])
		h *= 16777619
	}
	return h
}

func (fm *SmallMap[K, V]) Set(key K, value V) {
	if fm.size*10 >= len(fm.data)*7 {
		fm.grow()
	}

	h := fm.hash(key)
	mask := fm.mask
	idx := h & mask

	for i := uintptr(1); ; i++ {
		e := &fm.data[idx]
		if !e.used {
			*e = entry[K, V]{key: key, value: value, used: true}
			fm.size++
			return
		}
		if e.key == key {
			e.value = value
			return
		}
		idx = (idx + i) & mask
	}
}

func (fm *SmallMap[K, V]) Get(key K) (V, bool) {
	if fm.size == 0 {
		var zero V
		return zero, false
	}

	h := fm.hash(key)
	mask := fm.mask
	idx := h & mask

	for i := uintptr(1); ; i++ {
		e := &fm.data[idx]
		if !e.used {
			var zero V
			return zero, false
		}
		if e.key == key {
			return e.value, true
		}
		idx = (idx + i) & mask
		if i > mask {
			break
		}
	}
	var zero V
	return zero, false
}

func (fm *SmallMap[K, V]) Delete(key K) {
	h := fm.hash(key)
	mask := fm.mask
	idx := h & mask

	for i := uintptr(1); ; i++ {
		e := &fm.data[idx]
		if !e.used {
			return
		}
		if e.key == key {
			*e = entry[K, V]{}
			fm.size--
			return
		}
		idx = (idx + i) & mask
		if i > mask {
			break
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
	_, found := fm.Get(key)
	return found
}
