package ungo

type Specification[T any] interface {
	IsSatisfiedBy(T) bool
}

type ListSpec[T any] struct {
	specs []Specification[T]
}

func (s ListSpec[T]) IsSatisfiedBy(value T) bool {
	for _, spec := range s.specs {
		if !spec.IsSatisfiedBy(value) {
			return false
		}
	}
	return true
}

type FuncSpec[T any] struct {
	fn func(T) bool
}

func (s FuncSpec[T]) IsSatisfiedBy(value T) bool {
	return s.fn(value)
}

type AndSpec[T any] struct {
	spec1 Specification[T]
	spec2 Specification[T]
}

func (s AndSpec[T]) IsSatisfiedBy(value T) bool {
	return s.spec1.IsSatisfiedBy(value) && s.spec2.IsSatisfiedBy(value)
}

type OrSpec[T any] struct {
	spec1 Specification[T]
	spec2 Specification[T]
}

func (s OrSpec[T]) IsSatisfiedBy(value T) bool {
	return s.spec1.IsSatisfiedBy(value) || s.spec2.IsSatisfiedBy(value)
}

type NotSpec[T any] struct {
	spec Specification[T]
}

func (s NotSpec[T]) IsSatisfiedBy(value T) bool {
	return !s.spec.IsSatisfiedBy(value)
}
