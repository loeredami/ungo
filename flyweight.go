package ungo

import "sync"

type FlyweightFactory[K comparable, V any] struct {
	pool FastMap[K, *V]
	mu   sync.RWMutex
}

func NewFlyweightFactory[K comparable, V any]() *FlyweightFactory[K, V] {
	return &FlyweightFactory[K, V]{pool: FastMap[K, *V]{}}
}

func (f *FlyweightFactory[K, V]) Get(key K, generator func() V) *V {
	f.mu.RLock()
	if instance, ok := f.pool.Get(key); ok {
		f.mu.RUnlock()
		return instance
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check after lock
	if instance, ok := f.pool.Get(key); ok {
		return instance
	}

	val := generator()
	f.pool.Set(key, &val)
	return &val
}

func (f *FlyweightFactory[K, V]) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.pool = FastMap[K, *V]{}
}

func (f *FlyweightFactory[K, V]) Len() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.pool.Size()
}

func (f *FlyweightFactory[K, V]) Keys() []K {
	f.mu.RLock()
	defer f.mu.RUnlock()
	keys := make([]K, 0, f.pool.Size())
	for _, k := range f.pool.indexes {
		keys = append(keys, k)
	}
	return keys
}

func (f *FlyweightFactory[K, V]) Values() []*V {
	f.mu.RLock()
	defer f.mu.RUnlock()
	values := make([]*V, 0, f.pool.Size())
	for _, v := range f.pool.values {
		values = append(values, v)
	}
	return values
}

func (f *FlyweightFactory[K, V]) Entries() []struct {
	Key   K
	Value *V
} {
	f.mu.RLock()
	defer f.mu.RUnlock()
	entries := make([]struct {
		Key   K
		Value *V
	}, 0, f.pool.Size())
	for i := range f.pool.indexes {
		entries = append(entries, struct {
			Key   K
			Value *V
		}{Key: f.pool.indexes[i], Value: f.pool.values[i]})
	}
	return entries
}
