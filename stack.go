package ungo

import "fmt"

type stackNode[T any] struct {
	val  T
	last *stackNode[T]
}

type Stack[T any] struct {
	counter int
	last    *stackNode[T]
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		counter: 0,
		last:    nil,
	}
}

func (s *Stack[T]) String() string {
	if s.last == nil {
		return "Stack(empty)"
	}
	return fmt.Sprintf("Stack(last: %v)", s.last)
}

func (s *Stack[T]) Push(value T) {
	nLast := &stackNode[T]{
		val:  value,
		last: s.last,
	}
	s.last = nLast
	s.counter++
}

func (s *Stack[T]) Pop() Optional[T] {
	if s.last == nil {
		return None[T]()
	}

	value := s.last.val
	s.last = s.last.last
	s.counter--

	return Some(value)
}

func (s *Stack[T]) Peek() Optional[T] {
	if s.last == nil {
		return None[T]()
	}

	return Some(s.last.val)
}

func (s *Stack[T]) Size() int {
	return s.counter
}

func (s *Stack[T]) Clear() {
	for s.Pop().HasValue() {
	}
}
