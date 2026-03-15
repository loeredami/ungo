package ungo

import "testing"

func TestDefault(t *testing.T) {
	d := NewDefault("Hello")

	hello := d.Pass(EmptyOptional[string]())
	not_hello := d.Pass(MakeOptional("World"))

	if hello != "Hello" {
		t.Errorf("Expected 'Hello', got %v", hello)
	}
	if not_hello != "World" {
		t.Errorf("Expected 'World', got %v", not_hello)
	}
}
