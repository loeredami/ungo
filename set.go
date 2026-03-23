package ungo

import "fmt"

type Set[T comparable] SmallMap[T, struct{}]

func NewSet[T comparable](elements ...T) Set[T] {
	s := Set[T](*NewSmallMap[T, struct{}](0xFFF))
	for _, elem := range elements {
		s.Add(elem)
	}
	return s
}

func (s *Set[T]) String() string {
	return fmt.Sprintf("Set(%v)", s.ToSlice())
}

func (s *Set[T]) Add(element T) {
	(*SmallMap[T, struct{}])(s).Set(element, struct{}{})
}

func (s *Set[T]) Remove(element T) {
	(*SmallMap[T, struct{}])(s).Delete(element)
}

func (s Set[T]) Contains(element T) bool {
	return (*SmallMap[T, struct{}])(&s).Contains(element)
}

func (s Set[T]) Size() int {
	return (*SmallMap[T, struct{}])(&s).Size()
}

func (s *Set[T]) Clear() {
	(*SmallMap[T, struct{}])(s).Clear()
}

func (s *Set[T]) ToSlice() []T {
	return (*SmallMap[T, struct{}])(s).Keys()
}

func SetFromSlice[T comparable](slice []T) Set[T] {
	s := Set[T](*NewSmallMap[T, struct{}](0xFFF))
	for _, elem := range slice {
		s.Add(elem)
	}
	return s
}
