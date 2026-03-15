package ungo

import "testing"

func TestMultiMap(t *testing.T) {
	m := NewMultiMap[int, string]()

	m.Add(1, "one")
	m.Add(2, "two")
	m.Add(1, "three")

	values := m.Get(1)
	if len(values) != 2 {
		t.Errorf("expected 2 values for key 1, got %d", len(values))
	}
	if values[0] != "one" || values[1] != "three" {
		t.Errorf("expected values [one, three], got %v", values)
	}

}
