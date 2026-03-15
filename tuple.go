package ungo

type Tuple struct {
	elements []any
}

func NewTuple(elements ...any) *Tuple {
	return &Tuple{elements: elements}
}

func (t *Tuple) Get(index int) any {
	return t.elements[index]
}

func (t *Tuple) Len() int {
	return len(t.elements)
}

func (t *Tuple) Map(f func(any) any) *Tuple {
	newElements := make([]any, len(t.elements))
	for i, element := range t.elements {
		newElements[i] = f(element)
	}
	return &Tuple{elements: newElements}
}
