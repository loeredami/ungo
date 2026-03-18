package ungo

import (
	"unsafe"
)

// We use a "Metadata" byte to store a 7-bit hash fragment.
// This allows us to skip 99% of key comparisons.
const (
	hEmpty    uint8 = 0
	hOccupied uint8 = 1 << 7
)

type SmallMap[K comparable, V any] struct {
	keys     []K
	values   []V
	metadata []uint8 // Bit 7: occupied, Bits 0-6: hash fragment
	size     int
	mask     uintptr
}

func NewSmallMap[K comparable, V any](capacity int) *SmallMap[K, V] {
	realCap := 1
	for realCap < capacity {
		realCap <<= 1
	}

	return &SmallMap[K, V]{
		keys:     make([]K, realCap),
		values:   make([]V, realCap),
		metadata: make([]uint8, realCap),
		mask:     uintptr(realCap - 1),
	}
}

// fastHash avoids the maphash object entirely for fixed-size types.
// It uses a simple but very fast Wyhash-inspired mixer.
func (fm *SmallMap[K, V]) fastHash(key K) uintptr {
	ptr := unsafe.Pointer(&key)
	size := unsafe.Sizeof(key)

	// Fast path for 8-byte types (int64, float64, pointers)
	if size == 8 {
		u := *(*uint64)(ptr)
		u ^= u >> 33
		u *= 0xff51afd7ed558ccd
		u ^= u >> 33
		return uintptr(u)
	}

	// Fallback for other sizes (uses a basic FNV-1a inline)
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

	for {
		m := fm.metadata[idx]
		if m == hEmpty {
			fm.keys[idx] = key
			fm.values[idx] = value
			fm.metadata[idx] = tag
			fm.size++
			return
		}
		// If the 7-bit hash fragment matches, then check the actual key.
		// This avoids expensive equality checks for different keys.
		if m == tag && fm.keys[idx] == key {
			fm.values[idx] = value
			return
		}
		idx = (idx + 1) & fm.mask
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
		idx = (idx + 1) & fm.mask
	}
}

func (fm *SmallMap[K, V]) Delete(key K) {
	h := fm.fastHash(key)
	idx := h & fm.mask
	tag := uint8(h&0x7F) | hOccupied

	for {
		m := fm.metadata[idx]
		if m == hEmpty {
			return
		}
		if m == tag && fm.keys[idx] == key {
			fm.metadata[idx] = hEmpty
			var zeroK K
			var zeroV V
			fm.keys[idx] = zeroK
			fm.values[idx] = zeroV
			fm.size--
			fm.rehashCluster(uint64(idx))
			return
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *SmallMap[K, V]) rehashCluster(hole uint64) {
	i := hole
	for {
		i = (i + 1) & uint64(fm.mask)
		if fm.metadata[i] == hEmpty {
			return
		}

		k, v := fm.keys[i], fm.values[i]
		fm.metadata[i] = hEmpty

		var zeroK K
		var zeroV V
		fm.keys[i] = zeroK
		fm.values[i] = zeroV

		fm.size--
		fm.Set(k, v)
	}
}

func (fm *SmallMap[K, V]) Size() int {
	return fm.size
}

func (fm *SmallMap[K, V]) ForEach(f func(key K, value V)) {
	for i, m := range fm.metadata {
		if m&hOccupied != 0 {
			f(fm.keys[i], fm.values[i])
		}
	}
}

func (fm *SmallMap[K, V]) Keys() []K {
	res := make([]K, 0, fm.size)
	for i, m := range fm.metadata {
		if m&hOccupied != 0 {
			res = append(res, fm.keys[i])
		}
	}
	return res
}

func (fm *SmallMap[K, V]) Values() []V {
	res := make([]V, 0, fm.size)
	for i, m := range fm.metadata {
		if m&hOccupied != 0 {
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
