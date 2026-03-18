package ungo

import (
	"hash/maphash"
	"unsafe"
)

const (
	empty    uint8 = 0
	occupied uint8 = 1
)

type SmallMap[K comparable, V any] struct {
	keys     []K
	values   []V
	metadata []uint8
	size     int
	mask     uint64
	seed     maphash.Seed
	hash     maphash.Hash
}

func NewSmallMap[K comparable, V any](capacity int) *SmallMap[K, V] {
	realCap := 1
	for realCap < capacity {
		realCap <<= 1
	}

	sm := &SmallMap[K, V]{
		keys:     make([]K, realCap),
		values:   make([]V, realCap),
		metadata: make([]uint8, realCap),
		mask:     uint64(realCap - 1),
		seed:     maphash.MakeSeed(),
	}
	sm.hash.SetSeed(sm.seed)
	return sm
}

func (fm *SmallMap[K, V]) getHash(key K) uint64 {
	fm.hash.Reset()
	kSize := unsafe.Sizeof(key)
	ptr := unsafe.Pointer(&key)
	b := unsafe.Slice((*byte)(ptr), kSize)

	fm.hash.Write(b)
	return fm.hash.Sum64()
}

func (fm *SmallMap[K, V]) Set(key K, value V) {
	idx := fm.getHash(key) & fm.mask

	keys := fm.keys
	meta := fm.metadata

	for {
		if meta[idx] == empty {
			keys[idx] = key
			fm.values[idx] = value
			meta[idx] = occupied
			fm.size++
			return
		}
		if keys[idx] == key {
			fm.values[idx] = value
			return
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *SmallMap[K, V]) Get(key K) (V, bool) {
	idx := fm.getHash(key) & fm.mask
	keys := fm.keys
	meta := fm.metadata

	for {
		if meta[idx] == empty {
			var zero V
			return zero, false
		}
		if keys[idx] == key {
			return fm.values[idx], true
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *SmallMap[K, V]) Delete(key K) {
	idx := fm.getHash(key) & fm.mask
	keys := fm.keys
	meta := fm.metadata

	for {
		if meta[idx] == empty {
			return
		}
		if keys[idx] == key {
			meta[idx] = empty
			var zeroK K
			var zeroV V
			keys[idx] = zeroK
			fm.values[idx] = zeroV
			fm.size--
			fm.rehashCluster(idx)
			return
		}
		idx = (idx + 1) & fm.mask
	}
}

func (fm *SmallMap[K, V]) rehashCluster(hole uint64) {
	i := hole
	for {
		i = (i + 1) & fm.mask
		if fm.metadata[i] == empty {
			return
		}

		k, v := fm.keys[i], fm.values[i]
		fm.metadata[i] = empty

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
	meta := fm.metadata
	keys := fm.keys
	vals := fm.values
	for i := range meta {
		if meta[i] == occupied {
			f(keys[i], vals[i])
		}
	}
}

func (fm *SmallMap[K, V]) Keys() []K {
	res := make([]K, 0, fm.size)
	meta := fm.metadata
	keys := fm.keys
	for i := range meta {
		if meta[i] == occupied {
			res = append(res, keys[i])
		}
	}
	return res
}

func (fm *SmallMap[K, V]) Values() []V {
	res := make([]V, 0, fm.size)
	meta := fm.metadata
	vals := fm.values
	for i := range meta {
		if meta[i] == occupied {
			res = append(res, vals[i])
		}
	}
	return res
}

func (fm *SmallMap[K, V]) Clear() {
	var zeroK K
	var zeroV V
	for i := range fm.metadata {
		fm.metadata[i] = empty
		fm.keys[i] = zeroK
		fm.values[i] = zeroV
	}
	fm.size = 0
}

func (fm *SmallMap[K, V]) Contains(key K) bool {
	_, found := fm.Get(key)
	return found
}
