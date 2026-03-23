package ungo

import "fmt"

type queueNode[T any] struct {
	val  T
	next *queueNode[T]
}

type Queue[T any] struct {
	counter int
	current *queueNode[T]
	last    *queueNode[T]
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{
		counter: 0,
		current: nil,
		last:    nil,
	}
}

func (q *Queue[T]) String() string {
	if q.current == nil {
		return "Queue(Empty)"
	}
	return fmt.Sprintf("Queue(%v)", q.current.val)
}

func (q *Queue[T]) Push(element T) {
	if q.current == nil {
		q.current = &queueNode[T]{
			val:  element,
			next: nil,
		}
		q.last = q.current
		q.counter++
		return
	}

	node := &queueNode[T]{
		val:  element,
		next: nil,
	}
	q.last.next = node
	q.last = node
	q.counter++
}

func (q *Queue[T]) Pop() Optional[T] {
	if q.current == nil {
		return None[T]()
	}

	val := q.current.val

	q.current = q.current.next
	if q.current == nil {
		q.last = nil
	}

	q.counter--

	return Some(val)
}

func (q *Queue[T]) Peek() Optional[T] {
	if q.current == nil {
		return None[T]()
	}
	return Some(q.current.val)
}

func (q *Queue[T]) IsEmpty() bool {
	return q.current == nil
}
