package ungo

type Default[T any] struct {
	value T
}

func NewDefault[T any](value T) *Default[T] {
	return &Default[T]{value: value}
}

func (d *Default[T]) Pass(value Optional[T]) T {
	if value.HasValue() {
		return value.Value()
	}
	return d.value
}

func (d *Default[T]) Value() T {
	return d.value
}
