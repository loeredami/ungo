package ungo

import "testing"

func TestLazy(t *testing.T) {
	lazy := NewLazy(func() int {
		return 25
	})

	value := lazy.Value()

	if value != 25 {
		t.Errorf("expected value to be 25, got %v", value)
	}
}
