package ungo

type Branded[T any, Brand any] struct {
	value T
}

func NewBranded[T any, Brand any](value T) Branded[T, Brand] {
	return Branded[T, Brand]{value: value}
}

func (b Branded[T, Brand]) Value() T {
	return b.value
}

func (b Branded[T, Brand]) Brand() Brand {
	var brand Brand
	return brand
}
