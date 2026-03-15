package ungo

type Branded[T any, Brand any] struct {
	value T
}

func NewBranded[T any, Brand any](value T) Branded[T, Brand] {
	return Branded[T, Brand]{value: value}
}
