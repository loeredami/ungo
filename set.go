package ungo

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](elements ...T) Set[T] {
	s := make(Set[T])
	for _, elem := range elements {
		s[elem] = struct{}{}
	}
	return s
}

func (s Set[T]) Add(element T) {
	s[element] = struct{}{}
}

func (s Set[T]) Remove(element T) {
	delete(s, element)
}

func (s Set[T]) Contains(element T) bool {
	_, exists := s[element]
	return exists
}

func (s Set[T]) Size() int {
	return len(s)
}

func (s Set[T]) Clear() {
	for key := range s {
		delete(s, key)
	}
}

func (s Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s))
	for key := range s {
		slice = append(slice, key)
	}
	return slice
}

func SetFromSlice[T comparable](slice []T) Set[T] {
	s := make(Set[T])
	for _, elem := range slice {
		s[elem] = struct{}{}
	}
	return s
}
