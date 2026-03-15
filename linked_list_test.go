package ungo

import "testing"

func TestLinkedList(t *testing.T) {
	list := ListOf(1, 2, 3)

	idx0 := list.Get(0).Value()
	idx1 := list.Get(1).Value()
	idx2 := list.Get(2).Value()
	idx3 := list.Get(3)

	if idx0 != 1 {
		t.Errorf("expected 1, got %d", idx0)
	}
	if idx1 != 2 {
		t.Errorf("expected 2, got %d", idx1)
	}
	if idx2 != 3 {
		t.Errorf("expected 3, got %d", idx2)
	}
	if idx3.HasValue() {
		t.Errorf("expected no value, got %d", idx3.Value())
	}

	size := list.Size()
	if size != 3 {
		t.Errorf("expected size 3, got %d", size)
	}

	slice := list.ToSlice()
	if len(slice) != 3 {
		t.Errorf("expected slice length 3, got %d", len(slice))
	}
	if slice[0] != 1 {
		t.Errorf("expected slice[0] to be 1, got %d", slice[0])
	}
	if slice[1] != 2 {
		t.Errorf("expected slice[1] to be 2, got %d", slice[1])
	}
	if slice[2] != 3 {
		t.Errorf("expected slice[2] to be 3, got %d", slice[2])
	}

}
