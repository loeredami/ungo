package ungo

type Blackboard[K comparable] struct {
	data map[K]any
	subs MultiMap[K, func(any)]
}

func NewBlackboard[K comparable]() *Blackboard[K] {
	return &Blackboard[K]{
		data: make(map[K]any),
		subs: *NewMultiMap[K, func(any)](),
	}
}

func (b *Blackboard[K]) Write(key K, value any) {
	b.data[key] = value
	for _, sub := range b.subs.Get(key) {
		sub(value)
	}
}

func (b *Blackboard[K]) Watch(key K, sub func(any)) {
	b.subs.Add(key, sub)
}
