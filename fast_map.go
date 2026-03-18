package ungo

import (
	"hash/maphash"
	"unsafe"
)

type entry[K comparable, V any] struct {
	key      K
	value    V
	occupied bool
}

type FastMap[K comparable, V any] struct {
	buckets []entry[K, V]
	size    int
	mask    uint64
	seed    maphash.Seed
}

func NewFastMap[K comparable, V any](capacity int) *FastMap[K, V] {
	realCap := 1
	for realCap < capacity {
		realCap <<= 1
	}

	return &FastMap[K, V]{
		buckets: make([]entry[K, V], realCap),
		mask:    uint64(realCap - 1),
		seed:    maphash.MakeSeed(),
	}
}

func (fm *FastMap[K, V]) hash(key K) uint64 {
	var h maphash.Hash
	h.SetSeed(fm.seed)

	kSize := unsafe.Sizeof(key)
	ptr := unsafe.Pointer(&key)
	b := unsafe.Slice((*byte)(ptr), kSize)

	h.Write(b)
	return h.Sum64()
}

func (fm *FastMap[K, V]) Set(key K, value V) {
	idx := fm.hash(key) & fm.mask
	for {
		if !fm.buckets[idx].occupied {
			fm.buckets[idx].key = key
			fm.buckets[idx].value = value
			fm.buckets[idx].occupied = true
			fm.size++
			return
		}
		if fm.buckets[idx].key == key {
			fm.buckets[idx].value = value
			return
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *FastMap[K, V]) Get(key K) (V, bool) {
	idx := fm.hash(key) & fm.mask
	for {
		if !fm.buckets[idx].occupied {
			var zero V
			return zero, false
		}
		if fm.buckets[idx].key == key {
			return fm.buckets[idx].value, true
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *FastMap[K, V]) Delete(key K) {
	idx := fm.hash(key) & fm.mask
	for {
		if !fm.buckets[idx].occupied {
			return
		}
		if fm.buckets[idx].key == key {
			fm.buckets[idx].occupied = false
			fm.size--
			fm.rehashCluster(idx)
			return
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *FastMap[K, V]) rehashCluster(hole uint64) {
	i := hole
	for {
		i = (i + 1) & fm.mask
		if !fm.buckets[i].occupied {
			return
		}
		// Re-insert this element
		k, v := fm.buckets[i].key, fm.buckets[i].value
		fm.buckets[i].occupied = false
		fm.size--
		fm.Set(k, v)
	}
}

func (fm *FastMap[K, V]) Size() int {
	return fm.size
}

func (fm *FastMap[K, V]) ForEach(f func(key K, value V)) {
	for i := range fm.buckets {
		if fm.buckets[i].occupied {
			f(fm.buckets[i].key, fm.buckets[i].value)
		}
	}
}

func (fm *FastMap[K, V]) Keys() []K {
	keys := make([]K, 0, fm.size)
	for i := range fm.buckets {
		if fm.buckets[i].occupied {
			keys = append(keys, fm.buckets[i].key)
		}
	}
	return keys
}

func (fm *FastMap[K, V]) Values() []V {
	vals := make([]V, 0, fm.size)
	for i := range fm.buckets {
		if fm.buckets[i].occupied {
			vals = append(vals, fm.buckets[i].value)
		}
	}
	return vals
}

func (fm *FastMap[K, V]) Clear() {
	for i := range fm.buckets {
		fm.buckets[i].occupied = false
		// Clear values to help GC
		var zeroK K
		var zeroV V
		fm.buckets[i].key = zeroK
		fm.buckets[i].value = zeroV
	}
	fm.size = 0
}

func (fm *FastMap[K, V]) Contains(key K) bool {
	_, found := fm.Get(key)
	return found
}
