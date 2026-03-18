package ungo

import "testing"

func TestQueue(t *testing.T) {
	q := NewQueue[int]()
	q.Push(42)
	if q.Pop().OrElse(0) != 42 {
		t.Errorf("expected 42, got %d", q.Pop())
	}
	if q.Pop().OrElse(0) != 0 {
		t.Errorf("expected 0, got %d", q.Pop())
	}

	q.Push(42)
	q.Push(43)
	if q.Pop().OrElse(0) != 42 {
		t.Errorf("expected 42, got %d", q.Pop())
	}
	if q.Pop().OrElse(0) != 43 {
		t.Errorf("expected 43, got %d", q.Pop())
	}
}
