package ungo

type Dequeue[T any] struct {
	elements []T
}

func (d *Dequeue[T]) PushFront(t T) {
	d.elements = append([]T{t}, d.elements...)
}

func (d *Dequeue[T]) PushBack(t T) {
	d.elements = append(d.elements, t)
}

func (d *Dequeue[T]) PopFront() Optional[T] {
	if len(d.elements) == 0 {
		return EmptyOptional[T]()
	}
	t := d.elements[0]
	d.elements = d.elements[1:]
	return MakeOptional(t)
}

func (d *Dequeue[T]) PopBack() Optional[T] {
	if len(d.elements) == 0 {
		return EmptyOptional[T]()
	}
	t := d.elements[len(d.elements)-1]
	d.elements = d.elements[:len(d.elements)-1]
	return MakeOptional(t)
}

func (d *Dequeue[T]) Front() Optional[T] {
	if len(d.elements) == 0 {
		return EmptyOptional[T]()
	}
	return MakeOptional(d.elements[0])
}

func (d *Dequeue[T]) Back() Optional[T] {
	if len(d.elements) == 0 {
		return EmptyOptional[T]()
	}
	return MakeOptional(d.elements[len(d.elements)-1])
}

func (d *Dequeue[T]) Len() int {
	return len(d.elements)
}
