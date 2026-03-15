package ungo

import "sync"

type FlyweightFactory[K comparable, V any] struct {
	pool map[K]*V
	mu   sync.RWMutex
}

func NewFlyweightFactory[K comparable, V any]() *FlyweightFactory[K, V] {
	return &FlyweightFactory[K, V]{pool: make(map[K]*V)}
}

func (f *FlyweightFactory[K, V]) Get(key K, generator func() V) *V {
	f.mu.RLock()
	if instance, ok := f.pool[key]; ok {
		f.mu.RUnlock()
		return instance
	}
	f.mu.RUnlock()

	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check after lock
	if instance, ok := f.pool[key]; ok {
		return instance
	}

	val := generator()
	f.pool[key] = &val
	return &val
}

func (f *FlyweightFactory[K, V]) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.pool = make(map[K]*V)
}

func (f *FlyweightFactory[K, V]) Len() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.pool)
}

func (f *FlyweightFactory[K, V]) Keys() []K {
	f.mu.RLock()
	defer f.mu.RUnlock()
	keys := make([]K, 0, len(f.pool))
	for k := range f.pool {
		keys = append(keys, k)
	}
	return keys
}

func (f *FlyweightFactory[K, V]) Values() []*V {
	f.mu.RLock()
	defer f.mu.RUnlock()
	values := make([]*V, 0, len(f.pool))
	for _, v := range f.pool {
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
	}, 0, len(f.pool))
	for k, v := range f.pool {
		entries = append(entries, struct {
			Key   K
			Value *V
		}{Key: k, Value: v})
	}
	return entries
}
