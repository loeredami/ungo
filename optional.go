package ungo

type Optional[T any] struct {
	value T
	valid bool
}

func (o Optional[T]) HasValue() bool {
	return o.valid
}

func (o Optional[T]) Value() T {
	if !o.valid {
		var zero T
		return zero
	}
	return o.value
}

func MakeOptional[T any](value T) Optional[T] {
	return Optional[T]{value: value, valid: true}
}

func EmptyOptional[T any]() Optional[T] {
	return Optional[T]{valid: false}
}
