package ungo

type Vector[T Number] struct {
	arr []T
}

func MakeVector[T Number](nums ...T) Vector[T] {
	return Vector[T]{arr: nums}
}

func (v *Vector[T]) Add(v2 *Vector[T]) *Vector[T] {
	bigger := v
	if len(v2.arr) > len(v.arr) {
		bigger = v2
	}
	for i := range bigger.arr {
		if i < len(v.arr) {
			v.arr[i] += bigger.arr[i]
		} else {
			v.arr = append(v.arr, bigger.arr[i])
		}
	}
	return bigger
}

func (v *Vector[T]) Sub(v2 *Vector[T]) *Vector[T] {
	bigger := v
	if len(v2.arr) > len(v.arr) {
		bigger = v2
	}
	for i := range bigger.arr {
		if i < len(v.arr) {
			v.arr[i] -= bigger.arr[i]
		} else {
			v.arr = append(v.arr, -bigger.arr[i])
		}
	}
	return bigger
}

func (v *Vector[T]) Mul(v2 *Vector[T]) *Vector[T] {
	bigger := v
	if len(v2.arr) > len(v.arr) {
		bigger = v2
	}
	for i := range bigger.arr {
		if i < len(v.arr) {
			v.arr[i] *= bigger.arr[i]
		} else {
			v.arr = append(v.arr, bigger.arr[i])
		}
	}
	return bigger
}

func (v *Vector[T]) Div(v2 *Vector[T]) *Vector[T] {
	bigger := v
	if len(v2.arr) > len(v.arr) {
		bigger = v2
	}
	for i := range bigger.arr {
		if i < len(v.arr) {
			v.arr[i] /= bigger.arr[i]
		} else {
			v.arr = append(v.arr, bigger.arr[i])
		}
	}
	return bigger
}

func (v *Vector[T]) Dot(v2 *Vector[T]) *T {
	bigger := v
	if len(v2.arr) > len(v.arr) {
		bigger = v2
	}
	result := T(0)
	for i := range bigger.arr {
		if i < len(v.arr) {
			result += v.arr[i] * bigger.arr[i]
		} else {
			result += bigger.arr[i]
		}
	}
	return &result
}

func (v *Vector[T]) Length() *T {
	result := T(0)
	for _, num := range v.arr {
		result += num * num
	}
	return &result
}

func (v *Vector[T]) Normalize() *Vector[T] {
	length := v.Length()
	for i := range v.arr {
		v.arr[i] /= *length
	}
	return v
}

func (v *Vector[T]) Clone() *Vector[T] {
	result := &Vector[T]{arr: make([]T, len(v.arr))}
	copy(result.arr, v.arr)
	return result
}

func (v *Vector[T]) Set(index int, value T) {
	if index >= len(v.arr) {
		v.arr = append(v.arr, make([]T, index-len(v.arr)+1)...)
	}
	v.arr[index] = value
}

func (v *Vector[T]) At(index int) *T {
	if index >= len(v.arr) {
		return nil
	}
	return &v.arr[index]
}

func (v *Vector[T]) ForEach(f func(T) T) *Vector[T] {
	for i := range v.arr {
		v.arr[i] = f(v.arr[i])
	}
	return v
}
