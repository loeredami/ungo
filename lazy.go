package ungo

type Lazy[T any] struct {
	initializer func() T
	value       Optional[T]
}

func NewLazy[T any](initializer func() T) Lazy[T] {
	return Lazy[T]{initializer: initializer}
}

func (l Lazy[T]) Value() T {
	if !l.value.HasValue() {
		l.value = Optional[T]{value: l.initializer(), valid: true}
	}
	return l.value.Value()
}
