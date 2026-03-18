package ungo

type Constraint[T comparable] struct {
	possibleValues Set[T]
}

func NewConstraint[T comparable]() Constraint[T] {
	return Constraint[T]{}
}

func (c Constraint[T]) Add(value T) {
	c.possibleValues.Add(value)
}

func (c Constraint[T]) Remove(value T) {
	c.possibleValues.Remove(value)
}

func (c Constraint[T]) Narrow(f func(T) bool) Constraint[T] {
	result := Constraint[T]{}
	result.possibleValues = Set[T]{}
	for _, v := range c.possibleValues.ToSlice() {
		if f(v) {
			result.possibleValues.Add(v)
		}
	}
	return result
}

func (c Constraint[T]) Values() []T {
	return c.possibleValues.ToSlice()
}

func (c Constraint[T]) IsNothing() bool {
	return c.possibleValues.Size() == 0
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
