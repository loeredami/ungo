package ungo

type Signal[T any] struct {
	value     T
	observers []func(T)
}

func NewSignal[T any]() *Signal[T] {
	return &Signal[T]{}
}

func (s *Signal[T]) Set(val T) {
	s.value = val
	for _, obs := range s.observers {
		obs(val)
	}
}

func (s *Signal[T]) Get() T {
	return s.value
}

func (s *Signal[T]) Watch(sub func(T)) {
	s.observers = append(s.observers, sub)
}

// SignalCompute creates a "Derived" entity that updates when the source changes
func SignalCompute[T, U any](source *Signal[T], fn func(T) U) *Signal[U] {
	derived := &Signal[U]{value: fn(source.Get())}
	source.observers = append(source.observers, func(val T) {
		derived.Set(fn(val))
	})
	return derived
}
