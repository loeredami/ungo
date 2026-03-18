package ungo

type Registry[T any] FastMap[string, Lazy[T]]

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{}
}

func (r *Registry[T]) Register(name string, value Lazy[T]) {
	(*FastMap[string, Lazy[T]])(r).Set(name, value)
}

func (r *Registry[T]) Get(name string) Optional[T] {
	val, ok := (*FastMap[string, Lazy[T]])(r).Get(name)
	if !ok {
		return EmptyOptional[T]()
	}
	return MakeOptional(val.Value())
}

func (r *Registry[T]) Unregister(name string) {
	(*FastMap[string, Lazy[T]])(r).Delete(name)
}

func (r *Registry[T]) Clear() {
	(*FastMap[string, Lazy[T]])(r).Clear()
}

func (r *Registry[T]) Count() int {
	return (*FastMap[string, Lazy[T]])(r).Size()
}

func (r *Registry[T]) Keys() []string {
	return (*FastMap[string, Lazy[T]])(r).Keys()
}

func (r *Registry[T]) IfPresent(name string, consumer func(T)) {
	val, ok := (*FastMap[string, Lazy[T]])(r).Get(name)
	if ok {
		consumer(val.Value())
	}
}

func (r *Registry[T]) IfPresentOrElse(name string, consumer func(T), orElse func()) {
	val, ok := (*FastMap[string, Lazy[T]])(r).Get(name)
	if ok {
		consumer(val.Value())
	} else {
		orElse()
	}
}

func (r *Registry[T]) IfAbsent(name string, supplier func() T) {
	_, ok := (*FastMap[string, Lazy[T]])(r).Get(name)
	if !ok {
		(*FastMap[string, Lazy[T]])(r).Set(name, NewLazy(supplier))
	}
}
