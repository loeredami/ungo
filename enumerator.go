package ungo

type Enumerator[T any] struct {
	Next func() Optional[T]
}

func (e Enumerator[T]) Filter(predicate func(T) bool) Enumerator[T] {
	return Enumerator[T]{
		Next: func() Optional[T] {
			for {
				val := e.Next()
				if !val.HasValue() {
					return val
				}
				if predicate(val.Value()) {
					return val
				}
			}
		},
	}
}

func (e Enumerator[T]) Map(mapper func(T) T) Enumerator[T] {
	return Enumerator[T]{
		Next: func() Optional[T] {
			for {
				val := e.Next()
				if !val.HasValue() {
					return val
				}
				return Optional[T]{value: mapper(val.Value()), valid: true}
			}
		},
	}
}

func (e Enumerator[T]) ForEach(consumer func(T)) {
	for {
		val := e.Next()
		if !val.HasValue() {
			return
		}
		consumer(val.Value())
	}
}

func (e Enumerator[T]) Collect() []T {
	var result []T
	for {
		val := e.Next()
		if !val.HasValue() {
			return result
		}
		result = append(result, val.Value())
	}
}

func (e Enumerator[T]) Reduce(reducer func(T, T) T, initial T) T {
	result := initial
	for {
		val := e.Next()
		if !val.HasValue() {
			return result
		}
		result = reducer(result, val.Value())
	}
}
