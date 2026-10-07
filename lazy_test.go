package ungo

import "testing"

func TestLazy(t *testing.T) {
	lazy := NewLazy(func() int {
		return 25
	}) // TODO: we should check if the function gets called again when requesting Value multiple times, ideally should only be called once.

	value := lazy.Value()

	if value != 25 {
		t.Errorf("expected value to be 25, got %v", value)
	}
}
