package ungo

type Lazy[T any] struct {
	initializer func() T
	value       Optional[T]
	onInit      []func(*T)
}

func NewLazy[T any](initializer func() T) Lazy[T] {
	return Lazy[T]{initializer: initializer}
}

func (l *Lazy[T]) OnInit(f func(*T)) {
	l.onInit = append(l.onInit, f)
}

func (l *Lazy[T]) ClearOnInit() {
	l.onInit = make([]func(*T), 0)
}

func (l Lazy[T]) Value() T {
	if !l.value.HasValue() {
		l.value = Optional[T]{value: l.initializer(), valid: true}
		for _, f := range l.onInit {
			f(&l.value.value)
		}
	}
	return l.value.Value()
}

func (l Lazy[T]) IsInitialized() bool {
	return l.value.HasValue()
}
