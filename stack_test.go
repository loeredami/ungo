package ungo

import "testing"

func TestStack(t *testing.T) {
	s := NewStack[int]()

	s.Push(42)

	top := s.Pop()

	if top.Value() != 42 {
		t.Errorf("expected top to be 42, got %d", top.Value())
	}

	empty := s.Pop()
	if empty.Value() != 0 {
		t.Errorf("expected empty stack to return 0, got %d", empty.Value())
	}
}
