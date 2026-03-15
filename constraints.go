package ungo

type Constraint[T comparable] struct {
	possibleValues Set[T]
}

func (c Constraint[T]) Narrow(f func(T) bool) Constraint[T] {
	result := Constraint[T]{}
	result.possibleValues = make(Set[T])
	for v := range c.possibleValues {
		if f(v) {
			result.possibleValues[v] = struct{}{}
		}
	}
	return result
}

func (c Constraint[T]) Values() []T {
	return c.possibleValues.ToSlice()
}

func (c Constraint[T]) IsNothing() bool {
	return len(c.possibleValues) == 0
}

func (c Constraint[T]) IsJust() bool {
	return !c.IsNothing()
}

func (c Constraint[T]) EnforceFirst() T {
	if c.IsNothing() {
		var zero T
		return zero
	}
	return c.possibleValues.ToSlice()[0]
}

func (c Constraint[T]) EnforceLast() T {
	if c.IsNothing() {
		var zero T
		return zero
	}
	slice := c.possibleValues.ToSlice()
	return slice[len(slice)-1]
}
