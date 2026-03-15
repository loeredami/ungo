package ungo

import "testing"

func TestFlyweight(t *testing.T) {
	flyweight := NewFlyweightFactory[string, int]()

	value := *flyweight.Get("key", func() int { return 42 })

	if value != 42 {
		t.Errorf("expected value to be 42, got %d", value)
	}
}
