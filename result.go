package ungo

import "fmt"

type Result[T any] struct {
	value T
	err   error
}

func (e Result[T]) String() string {
	if e.err != nil {
		return fmt.Sprintf("Result(Err: %v)", e.err)
	}
	return fmt.Sprintf("Result(Ok: %v)", e.value)
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

func (e Result[T]) OnSuccess(f func(T)) {
	if e.err == nil {
		f(e.value)
	}
}

func (e Result[T]) OnError(f func(error)) {
	if e.err != nil {
		f(e.err)
	}
}

func VFail[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func VSuccess[T any](value T) Result[T] {
	return Result[T]{value: value}
}
