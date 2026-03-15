package ungo

type Either[T1, T2 any] struct {
	Left  *T1
	Right *T2
}

func NewEither[T1, T2 any](left T1, right T2) Either[T1, T2] {
	return Either[T1, T2]{Left: &left, Right: &right}
}

func (e Either[T1, T2]) IsLeft() bool {
	return e.Left != nil
}

func (e Either[T1, T2]) IsRight() bool {
	return e.Right != nil
}

func (e Either[T1, T2]) LeftValue() T1 {
	return *e.Left
}

func (e Either[T1, T2]) RightValue() T2 {
	return *e.Right
}

func (e Either[T1, T2]) LeftValueOrDefault(defaultValue T1) T1 {
	if e.Left != nil {
		return *e.Left
	}
	return defaultValue
}

func (e Either[T1, T2]) RightValueOrDefault(defaultValue T2) T2 {
	if e.Right != nil {
		return *e.Right
	}
	return defaultValue
}
