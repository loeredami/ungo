package ungo

type Exception[T any] struct {
	Value T
	Error error
}

func NewException[T any](value T, err error) Exception[T] {
	return Exception[T]{Value: value, Error: err}
}

func Try[T any](fn func() (T, error)) Exception[T] {
	value, err := fn()
	return NewException(value, err)
}

func (e Exception[T]) Catch(fn func(error) T) T {
	if e.Error != nil {
		return fn(e.Error)
	}
	return e.Value
}

func ExceptToResult[T any](fn func() (T, error)) Result[T] {
	res := Try(fn)
	return Result[T]{value: res.Value, err: res.Error}
}
