package ungo

import (
	"testing"
)

func TestBlackboard(t *testing.T) {
	b := NewBlackboard[string]()

	var value string
	value = "unset"

	b.Watch("key", func(a any) {
		value = a.(string)
	})

	b.Write("key", "value")

	if value != "value" {
		t.Errorf("expected value to be 'value', got %v", value)
	}
}
