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
