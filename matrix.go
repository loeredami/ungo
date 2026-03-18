package ungo

type Matrix[T Number] struct {
	vectors []*Vector[T]
}

func (m *Matrix[T]) Add(m2 *Matrix[T]) *Matrix[T] {
	result := &Matrix[T]{}
	bigger_matrix := m
	if len(m2.vectors) > len(m.vectors) {
		bigger_matrix = m2
	}
	for i := range bigger_matrix.vectors {
		if i < len(m.vectors) {
			result.vectors = append(result.vectors, m.vectors[i].Add(m2.vectors[i]))
		} else {
			result.vectors = append(result.vectors, bigger_matrix.vectors[i])
		}
	}
	return result
}

func (m *Matrix[T]) Sub(m2 *Matrix[T]) *Matrix[T] {
	result := &Matrix[T]{}
	bigger_matrix := m
	if len(m2.vectors) > len(m.vectors) {
		bigger_matrix = m2
	}
	for i := range bigger_matrix.vectors {
		if i < len(m.vectors) {
			result.vectors = append(result.vectors, m.vectors[i].Sub(m2.vectors[i]))
		} else {
			result.vectors = append(result.vectors, bigger_matrix.vectors[i])
		}
	}
	return result
}

func (m *Matrix[T]) Mul(m2 *Matrix[T]) *Matrix[T] {
	result := &Matrix[T]{}
	for i := range m.vectors {
		result.vectors = append(result.vectors, m.vectors[i].Mul(m2.vectors[i]))
	}
	return result
}

func (m *Matrix[T]) Div(m2 *Matrix[T]) *Matrix[T] {
	result := &Matrix[T]{}
	for i := range m.vectors {
		result.vectors = append(result.vectors, m.vectors[i].Div(m2.vectors[i]))
	}
	return result
}

func (m *Matrix[T]) Transpose() *Matrix[T] {
	result := &Matrix[T]{}
	for i := range m.vectors[0].arr {
		row := &Vector[T]{}
		for j := range m.vectors {
			row.arr = append(row.arr, m.vectors[j].arr[i])
		}
		result.vectors = append(result.vectors, row)
	}
	return result
}
