package ungo

import "testing"

func TestSignal(t *testing.T) {
	is := NewSignal[int]()

	var observed int

	is.Watch(func(i int) {
		observed = i
	})

	is.Set(42)

	if observed != 42 {
		t.Errorf("expected observed to be 42, got %d", observed)
	}
}
