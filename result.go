package ungo

type Result[T any] struct {
	value T
	err   error
}

func (e Result[T]) Success() bool {
	return e.err == nil
}

func (e Result[T]) Value() T {
	if e.err != nil {
		var zero T
		return zero
	}
	return e.value
}

func (e Result[T]) Error() error {
	return e.err
}

func VFail[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func VSuccess[T any](value T) Result[T] {
	return Result[T]{value: value}
}
