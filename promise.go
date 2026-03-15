package ungo

import "time"

type Promise[T any] struct {
	Ch chan T
}

func NewPromise[T any](fn func() T) Promise[T] {
	p := Promise[T]{Ch: make(chan T)}
	go func() {
		p.Ch <- fn()
	}()
	return p
}

func (p Promise[T]) Then(fn func(T)) {
	go func() {
		v := <-p.Ch
		fn(v)
	}()
}
func (p Promise[T]) Resolve(value T) {
	p.Ch <- value
}

func (p Promise[T]) Reject(err error) {
	p.Ch <- err.(T)
}

func (p Promise[T]) Await() T {
	return <-p.Ch
}

func (p Promise[T]) Timeout(timeout time.Duration) Optional[T] {
	select {
	case v := <-p.Ch:
		return Optional[T]{value: v, valid: true}
	case <-time.After(timeout):
		return Optional[T]{valid: false}
	}
}
