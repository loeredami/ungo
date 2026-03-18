package ungo

import "fmt"

type Set[T comparable] FastMap[T, struct{}]

func NewSet[T comparable](elements ...T) Set[T] {
	s := Set[T]{}
	for _, elem := range elements {
		s.Add(elem)
	}
	return s
}

func (s *Set[T]) String() string {
	return fmt.Sprintf("Set(%v)", s.ToSlice())
}

func (s *Set[T]) Add(element T) {
	(*FastMap[T, struct{}])(s).Set(element, struct{}{})
}

func (s *Set[T]) Remove(element T) {
	(*FastMap[T, struct{}])(s).Delete(element)
}

func (s Set[T]) Contains(element T) bool {
	return (*FastMap[T, struct{}])(&s).Contains(element)
}

func (s Set[T]) Size() int {
	return (*FastMap[T, struct{}])(&s).Size()
}

func (s *Set[T]) Clear() {
	(*FastMap[T, struct{}])(s).Clear()
}

func (s *Set[T]) ToSlice() []T {
	return (*FastMap[T, struct{}])(s).Keys()
}

func SetFromSlice[T comparable](slice []T) Set[T] {
	s := Set[T]{}
	for _, elem := range slice {
		s.Add(elem)
	}
	return s
}
