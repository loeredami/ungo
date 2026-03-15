package ungo

type Handler[T any] func(T)

type Dispatcher[T any] struct {
	handlers []Handler[T]
}

func NewDispatcher[T any]() *Dispatcher[T] {
	return &Dispatcher[T]{handlers: []Handler[T]{}}
}

func (d *Dispatcher[T]) Subscribe(handler Handler[T]) {
	d.handlers = append(d.handlers, handler)
}

func (d *Dispatcher[T]) Emit(event T) {
	for _, handler := range d.handlers {
		go handler(event)
	}
}
