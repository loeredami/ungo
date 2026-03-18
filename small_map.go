package ungo

import (
	"unsafe"
)

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
	d := fm.data
	n := len(d)

	// Optimization: Manual Loop Unrolling (4x)
	// This reduces branch instructions and allows the CPU to
	// parallelize the comparisons.
	i := 0
	for ; i <= n-4; i += 4 {
		if d[i].key == key {
			d[i].value = value
			return
		}
		if d[i+1].key == key {
			d[i+1].value = value
			return
		}
		if d[i+2].key == key {
			d[i+2].value = value
			return
		}
		if d[i+3].key == key {
			d[i+3].value = value
			return
		}
	}

	// Handle remaining elements
	for ; i < n; i++ {
		if d[i].key == key {
			d[i].value = value
			return
		}
	}

	fm.data = append(d, entry[K, V]{key: key, value: value})
}

func (fm *SmallMap[K, V]) Get(key K) (V, bool) {
	d := fm.data
	n := len(d)

	// Optimization: Pointer arithmetic to avoid slice bounds checks (BCE)
	// We Prove to the compiler that we are staying within limits.
	if n > 0 {
		ptr := uintptr(unsafe.Pointer(&d[0]))
		size := unsafe.Sizeof(entry[K, V]{})

		i := 0
		// Unrolled search for Get
		for ; i <= n-4; i += 4 {
			e0 := (*entry[K, V])(unsafe.Pointer(ptr + uintptr(i)*size))
			if e0.key == key {
				return e0.value, true
			}

			e1 := (*entry[K, V])(unsafe.Pointer(ptr + uintptr(i+1)*size))
			if e1.key == key {
				return e1.value, true
			}

			e2 := (*entry[K, V])(unsafe.Pointer(ptr + uintptr(i+2)*size))
			if e2.key == key {
				return e2.value, true
			}

			e3 := (*entry[K, V])(unsafe.Pointer(ptr + uintptr(i+3)*size))
			if e3.key == key {
				return e3.value, true
			}
		}

		for ; i < n; i++ {
			e := (*entry[K, V])(unsafe.Pointer(ptr + uintptr(i)*size))
			if e.key == key {
				return e.value, true
			}
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
			d[lastIdx] = entry[K, V]{} // Zero out for GC
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
	d := fm.data
	for i := range d {
		d[i] = entry[K, V]{}
	}
	fm.data = d[:0]
}

func (fm *SmallMap[K, V]) Contains(key K) bool {
	// Contains benefits from the same Get logic but returns early
	_, found := fm.Get(key)
	return found
}
