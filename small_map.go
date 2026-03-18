package ungo

import (
	"unsafe"
)

const (
	hEmpty    uint8 = 0
	hOccupied uint8 = 1 << 7
)

type SmallMap[K comparable, V any] struct {
	keys     []K
	values   []V
	metadata []uint8
	size     int
	mask     uintptr
}

func NewSmallMap[K comparable, V any](capacity int) *SmallMap[K, V] {
	realCap := 1
	for realCap < capacity {
		realCap <<= 1
	}

	// We allocate a bit of extra padding at the end.
	// This prevents "wrap-around" logic inside the hot loops.
	padding := 8
	return &SmallMap[K, V]{
		keys:     make([]K, realCap+padding),
		values:   make([]V, realCap+padding),
		metadata: make([]uint8, realCap+padding),
		mask:     uintptr(realCap - 1),
	}
}

func (fm *SmallMap[K, V]) fastHash(key K) uintptr {
	ptr := unsafe.Pointer(&key)
	size := unsafe.Sizeof(key)

	if size == 8 {
		u := *(*uint64)(ptr)
		u ^= u >> 33
		u *= 0xff51afd7ed558ccd
		u ^= u >> 33
		return uintptr(u)
	}

	h := uint64(14695981039346656037)
	b := unsafe.Slice((*byte)(ptr), size)
	for _, x := range b {
		h ^= uint64(x)
		h *= 1099511628211
	}
	return uintptr(h)
}

func (fm *SmallMap[K, V]) Set(key K, value V) {
	h := fm.fastHash(key)
	idx := h & fm.mask
	tag := uint8(h&0x7F) | hOccupied

	// Optimized loop: minimal branching
	for {
		m := fm.metadata[idx]
		if m == hEmpty {
			fm.keys[idx] = key
			fm.values[idx] = value
			fm.metadata[idx] = tag
			fm.size++
			return
		}
		if m == tag && fm.keys[idx] == key {
			fm.values[idx] = value
			return
		}

		idx++
		// If we hit the absolute end of the allocated slice,
		// then and only then do we wrap back to 0.
		if idx >= uintptr(len(fm.metadata)) {
			idx = 0
		}
	}
}

func (fm *SmallMap[K, V]) Get(key K) (V, bool) {
	h := fm.fastHash(key)
	idx := h & fm.mask
	tag := uint8(h&0x7F) | hOccupied

	for {
		m := fm.metadata[idx]
		if m == hEmpty {
			var zero V
			return zero, false
		}
		if m == tag && fm.keys[idx] == key {
			return fm.values[idx], true
		}

		idx++
		if idx >= uintptr(len(fm.metadata)) {
			idx = 0
		}
	}
}

func (fm *SmallMap[K, V]) Size() int {
	return fm.size
}

func (fm *SmallMap[K, V]) ForEach(f func(key K, value V)) {
	// Optimization: we only iterate up to the real capacity + padding
	for i := 0; i < len(fm.metadata); i++ {
		if fm.metadata[i]&hOccupied != 0 {
			f(fm.keys[i], fm.values[i])
		}
	}
}

func (fm *SmallMap[K, V]) Keys() []K {
	res := make([]K, 0, fm.size)
	for i := 0; i < len(fm.metadata); i++ {
		if fm.metadata[i]&hOccupied != 0 {
			res = append(res, fm.keys[i])
		}
	}
	return res
}

func (fm *SmallMap[K, V]) Values() []V {
	res := make([]V, 0, fm.size)
	for i := 0; i < len(fm.metadata); i++ {
		if fm.metadata[i]&hOccupied != 0 {
			res = append(res, fm.values[i])
		}
	}
	return res
}

func (fm *SmallMap[K, V]) Clear() {
	var zeroK K
	var zeroV V
	for i := range fm.metadata {
		fm.metadata[i] = hEmpty
		fm.keys[i] = zeroK
		fm.values[i] = zeroV
	}
	fm.size = 0
}

func (fm *SmallMap[K, V]) Contains(key K) bool {
	_, found := fm.Get(key)
	return found
}
