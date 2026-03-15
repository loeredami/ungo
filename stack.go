package ungo

type Stack[T any] struct {
	elements []T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(value T) {
	s.elements = append(s.elements, value)
}

func (s *Stack[T]) Pop() Optional[T] {
	if len(s.elements) == 0 {
		return EmptyOptional[T]()
	}
	lastIndex := len(s.elements) - 1
	value := s.elements[lastIndex]
	s.elements = s.elements[:lastIndex]
	return MakeOptional(value)
}

func (s *Stack[T]) Peek() Optional[T] {
	if len(s.elements) == 0 {
		return EmptyOptional[T]()
	}
	return MakeOptional(s.elements[len(s.elements)-1])
}

func (s *Stack[T]) Size() int {
	return len(s.elements)
}

func (s *Stack[T]) Clear() {
	s.elements = make([]T, 0)
}
