package ungo

import (
	"fmt"
	"slices"
)

type Tuple struct {
	elements []any
}

func NewTuple(elements ...any) *Tuple {
	return &Tuple{elements: elements}
}

func (t *Tuple) ToSlice() []any {
	return t.elements
}

func (t *Tuple) String() string {
	result := "("
	for i, element := range t.elements {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%v", element)
	}
	result += ")"
	return result
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

func (t *Tuple) Filter(f func(any) bool) *Tuple {
	var newElements []any
	for _, element := range t.elements {
		if f(element) {
			newElements = append(newElements, element)
		}
	}
	return &Tuple{elements: newElements}
}

func (t *Tuple) Reduce(f func(any, any) any) any {
	result := t.elements[0]
	for _, element := range t.elements[1:] {
		result = f(result, element)
	}
	return result
}

func (t *Tuple) ForEach(f func(any)) {
	for _, element := range t.elements {
		f(element)
	}
}

func (t *Tuple) Contains(element any) bool {
	return slices.Contains(t.elements, element)
}

func (t *Tuple) IndexOf(element any) int {
	for i, e := range t.elements {
		if e == element {
			return i
		}
	}
	return -1
}

func (t *Tuple) LastIndexOf(element any) int {
	for i := len(t.elements) - 1; i >= 0; i-- {
		if t.elements[i] == element {
			return i
		}
	}
	return -1
}
