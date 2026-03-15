package ungo

type Registry[T any] struct {
	implementations map[string]RegistryEntry[T]
}

type RegistryEntry[T any] interface {
	Value(registry *Registry[T]) T
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{
		implementations: make(map[string]RegistryEntry[T]),
	}
}

func (r *Registry[T]) Register(name string, entry RegistryEntry[T]) {
	r.implementations[name] = entry
}

func (r *Registry[T]) Get(name string) T {
	return r.implementations[name].Value(r)
}

func (r *Registry[T]) GetOrDefault(name string, def T) T {
	if entry, ok := r.implementations[name]; ok {
		return entry.Value(r)
	}
	return def
}

func (r *Registry[T]) Names() []string {
	names := make([]string, 0, len(r.implementations))
	for name := range r.implementations {
		names = append(names, name)
	}
	return names
}
