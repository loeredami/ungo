package ungo

import (
	"unsafe"
)

type entry[K comparable, V any] struct {
	key   K
	value V
}

type SmallMap[K comparable, V any] struct {
	data  []entry[K, V]
	index []uint32
	mask  uintptr
}

func NewSmallMap[K comparable, V any](capacity int) *SmallMap[K, V] {
	pow2 := 8
	for pow2 < capacity*2 {
		pow2 <<= 1
	}
	return &SmallMap[K, V]{
		data:  make([]entry[K, V], 0, capacity),
		index: make([]uint32, pow2),
		mask:  uintptr(pow2 - 1),
	}
}

func (fm *SmallMap[K, V]) hash(key K) uintptr {
	u := *(*uintptr)(unsafe.Pointer(&key))
	return u * 0xdeece66d
}

func (fm *SmallMap[K, V]) Set(key K, value V) {
	h := fm.hash(key)
	idx := h & fm.mask

	for {
		dataIdxPlusOne := fm.index[idx]
		if dataIdxPlusOne == 0 {
			break
		}
		dataIdx := dataIdxPlusOne - 1
		if fm.data[dataIdx].key == key {
			fm.data[dataIdx].value = value
			return
		}
		idx = (idx + 1) & fm.mask
	}

	fm.data = append(fm.data, entry[K, V]{key: key, value: value})
	fm.index[idx] = uint32(len(fm.data))

	if len(fm.data)*10 > len(fm.index)*7 {
		fm.rehashIndex()
	}
}

func (fm *SmallMap[K, V]) Get(key K) (V, bool) {
	if len(fm.index) == 0 {
		var zero V
		return zero, false
	}

	idx := fm.hash(key) & fm.mask
	for {
		dataIdxPlusOne := fm.index[idx]
		if dataIdxPlusOne == 0 {
			var zero V
			return zero, false
		}
		dataIdx := dataIdxPlusOne - 1
		if fm.data[dataIdx].key == key {
			return fm.data[dataIdx].value, true
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *SmallMap[K, V]) rehashIndex() {
	newCap := len(fm.index) * 2
	fm.index = make([]uint32, newCap)
	fm.mask = uintptr(newCap - 1)

	for i := range fm.data {
		h := fm.hash(fm.data[i].key)
		idx := h & fm.mask
		for fm.index[idx] != 0 {
			idx = (idx + 1) & fm.mask
		}
		fm.index[idx] = uint32(i + 1)
	}
}

func (fm *SmallMap[K, V]) Delete(key K) {
	h := fm.hash(key)
	idx := h & fm.mask

	for {
		dataIdxPlusOne := fm.index[idx]
		if dataIdxPlusOne == 0 {
			return
		}
		dataIdx := dataIdxPlusOne - 1
		if fm.data[dataIdx].key == key {
			lastIdx := len(fm.data) - 1
			if int(dataIdx) != lastIdx {
				fm.data[dataIdx] = fm.data[lastIdx]
				fm.updateIndex(fm.data[dataIdx].key, uint32(dataIdx+1))
			}
			fm.data[lastIdx] = entry[K, V]{}
			fm.data = fm.data[:lastIdx]

			fm.rehashIndex()
			return
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *SmallMap[K, V]) updateIndex(key K, newIdxPlusOne uint32) {
	idx := fm.hash(key) & fm.mask
	for {
		if fm.data[fm.index[idx]-1].key == key {
			fm.index[idx] = newIdxPlusOne
			return
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *SmallMap[K, V]) Size() int {
	return len(fm.data)
}

func (fm *SmallMap[K, V]) ForEach(f func(key K, value V)) {
	for i := range fm.data {
		f(fm.data[i].key, fm.data[i].value)
	}
}

func (fm *SmallMap[K, V]) Keys() []K {
	res := make([]K, len(fm.data))
	for i := range fm.data {
		res[i] = fm.data[i].key
	}
	return res
}

func (fm *SmallMap[K, V]) Values() []V {
	res := make([]V, len(fm.data))
	for i := range fm.data {
		res[i] = fm.data[i].value
	}
	return res
}

func (fm *SmallMap[K, V]) Clear() {
	for i := range fm.data {
		fm.data[i] = entry[K, V]{}
	}
	fm.data = fm.data[:0]
	for i := range fm.index {
		fm.index[i] = 0
	}
}

func (fm *SmallMap[K, V]) Contains(key K) bool {
	_, found := fm.Get(key)
	return found
}
