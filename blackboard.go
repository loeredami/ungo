package ungo

type Blackboard[K comparable] struct {
	data FastMap[K, any]
	subs MultiMap[K, func(any)]
}

func NewBlackboard[K comparable]() *Blackboard[K] {
	return &Blackboard[K]{
		data: FastMap[K, any]{},
		subs: *NewMultiMap[K, func(any)](),
	}
}

func (b *Blackboard[K]) Write(key K, value any) {
	b.data.Set(key, value)
	for _, sub := range b.subs.Get(key) {
		sub(value)
	}
}

func (b *Blackboard[K]) Watch(key K, sub func(any)) {
	b.subs.Add(key, sub)
}

func (b *Blackboard[K]) Read(key K) any {
	val, ok := b.data.Get(key)
	if ok {
		return val
	}
	return nil
}

func (b *Blackboard[K]) Delete(key K) {
	b.data.Delete(key)
}

func (b *Blackboard[K]) Clear() {
	b.data.Clear()
}
