package ungo

type LinkedList[T any] struct {
	head *listNode[T]
	tail *listNode[T]
	size int
}

type listNode[T any] struct {
	value T
	next  *listNode[T]
}

func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{}
}

func (ll *LinkedList[T]) Add(value T) {
	newNode := &listNode[T]{value: value}
	if ll.head == nil {
		ll.head = newNode
		ll.tail = newNode
	} else {
		ll.tail.next = newNode
		ll.tail = newNode
	}
	ll.size++
}

func (ll *LinkedList[T]) Get(index int) Optional[T] {
	if index < 0 || index >= ll.size {
		return EmptyOptional[T]()
	}
	current := ll.head
	for range index {
		current = current.next
	}
	return MakeOptional(current.value)
}

func (ll *LinkedList[T]) Set(index int, value T) {
	if index < 0 || index >= ll.size {
		return
	}
	current := ll.head
	for range index {
		current = current.next
	}
	current.value = value
}

func (ll *LinkedList[T]) Size() int {
	return ll.size
}

func (ll *LinkedList[T]) Clear() {
	for current := ll.head; current != nil; {
		next := current.next
		current.next = nil // Help GC
		current = next
	}
	ll.head = nil
	ll.tail = nil
	ll.size = 0
}

func ListOf[T any](elements ...T) *LinkedList[T] {
	ll := NewLinkedList[T]()
	for _, elem := range elements {
		ll.Add(elem)
	}
	return ll
}

func ListFromSlice[T any](slice []T) *LinkedList[T] {
	ll := NewLinkedList[T]()
	for _, elem := range slice {
		ll.Add(elem)
	}
	return ll
}

func (ll *LinkedList[T]) ToSlice() []T {
	slice := make([]T, 0, ll.size)
	for current := ll.head; current != nil; current = current.next {
		slice = append(slice, current.value)
	}
	return slice
}

func (ll *LinkedList[T]) Remove(index int) {
	if index < 0 || index >= ll.size {
		return
	}
	if ll.size == 1 {
		ll.head = nil
		ll.tail = nil
		ll.size = 0
		return
	}
	current := ll.head
	for range index {
		current = current.next
	}
	current.next = current.next.next
	ll.size--
}
