package ungo

import "fmt"

type Optional[T any] struct {
	value T
	valid bool
}

func (o Optional[T]) String() string {
	if !o.valid {
		return "Optional(empty)"
	}
	return fmt.Sprintf("Optional(%v)", o.value)
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

func (o Optional[T]) IfPresent(f func(T)) {
	if o.valid {
		f(o.value)
	}
}

func (o Optional[T]) IfAbsent(f func(*T)) {
	if !o.valid {
		f(&o.value)
	}
}

func (o Optional[T]) OrElse(value T) T {
	if !o.valid {
		return value
	}
	return o.value
}
