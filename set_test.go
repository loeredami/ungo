package ungo

import "testing"

func TestSet(t *testing.T) {
	s := SetFromSlice([]int{52, 32})
	s.Add(1)
	s.Add(2)
	s.Add(3)
	if !s.Contains(2) {
		t.Errorf("expected set to contain 2")
	}
	if !s.Contains(32) {
		t.Errorf("expected set to contain 32")
	}
	if s.Contains(4) {
		t.Errorf("expected set to not contain 4")
	}

	s.Remove(2)
	if s.Contains(2) {
		t.Errorf("expected set to not contain 2 after remove")
	}

	s.Clear()
}
