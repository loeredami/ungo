package ungo

type Registry[T any] SmallMap[string, Lazy[T]]

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{}
}

func (r *Registry[T]) Register(name string, value Lazy[T]) {
	(*SmallMap[string, Lazy[T]])(r).Set(name, value)
}

func (r *Registry[T]) Get(name string) Optional[T] {
	val, ok := (*SmallMap[string, Lazy[T]])(r).Get(name)
	if !ok {
		return EmptyOptional[T]()
	}
	return MakeOptional(val.Value())
}

func (r *Registry[T]) Unregister(name string) {
	(*SmallMap[string, Lazy[T]])(r).Delete(name)
}

func (r *Registry[T]) Clear() {
	(*SmallMap[string, Lazy[T]])(r).Clear()
}

func (r *Registry[T]) Count() int {
	return (*SmallMap[string, Lazy[T]])(r).Size()
}

func (r *Registry[T]) Keys() []string {
	return (*SmallMap[string, Lazy[T]])(r).Keys()
}

func (r *Registry[T]) IfPresent(name string, consumer func(T)) {
	val, ok := (*SmallMap[string, Lazy[T]])(r).Get(name)
	if ok {
		consumer(val.Value())
	}
}

func (r *Registry[T]) IfPresentOrElse(name string, consumer func(T), orElse func()) {
	val, ok := (*SmallMap[string, Lazy[T]])(r).Get(name)
	if ok {
		consumer(val.Value())
	} else {
		orElse()
	}
}

func (r *Registry[T]) IfAbsent(name string, supplier func() T) {
	_, ok := (*SmallMap[string, Lazy[T]])(r).Get(name)
	if !ok {
		(*SmallMap[string, Lazy[T]])(r).Set(name, NewLazy(supplier))
	}
}
