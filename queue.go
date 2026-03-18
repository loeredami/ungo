package ungo

type Queue[T any] struct {
	elements []T
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{elements: []T{}}
}

func (q *Queue[T]) Push(element T) {
	q.elements = append(q.elements, element)
}

func (q *Queue[T]) Pop() Optional[T] {
	if len(q.elements) == 0 {
		return EmptyOptional[T]()
	}
	element := q.elements[0]
	q.elements = q.elements[1:]
	return MakeOptional(element)
}

func (q *Queue[T]) Peek() Optional[T] {
	if len(q.elements) == 0 {
		return EmptyOptional[T]()
	}
	return MakeOptional(q.elements[0])
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.elements) == 0
}
