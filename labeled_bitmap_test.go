package ungo

import "testing"

func TestLabeledBitmap(t *testing.T) {
	bitmap := NewLabeledBitmap([]string{"one", "two"})

	bitmap.SetLabel("one")

	if !bitmap.Get("one") {
		t.Errorf("expected 'one' to be set")
	}

	bitmap.ClearLabel("one")

	if bitmap.Get("one") {
		t.Errorf("expected 'one' to be cleared")
	}
}
