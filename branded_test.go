package ungo

import (
	"reflect"
	"testing"
)

func TestBranded(t *testing.T) {
	b := NewBranded[string, struct{}]("Super Sonic Racing!")

	// same value, different brand
	b2 := NewBranded[string, struct{ o2 string }]("Super Sonic Racing!")

	if b.value != "Super Sonic Racing!" {
		t.Errorf("expected value to be 'Super Sonic Racing!', got %v", b.value)
	}

	if b2.value != "Super Sonic Racing!" {
		t.Errorf("expected value to be 'Super Sonic Racing!', got %v", b2.value)
	}

	if reflect.TypeOf(b) == reflect.TypeOf(b2) {
		t.Errorf("expected different types, got %v and %v", reflect.TypeOf(b), reflect.TypeOf(b2))
	}
}
