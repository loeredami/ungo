package ungo

import (
	"testing"
	"time"
)

func TestPromise(t *testing.T) {
	p := NewPromise(func() int {
		time.Sleep(time.Second)
		return 42
	})
	var result int
	p.Then(func(i int) {
		result = i
	})

	result2 := p.Await()
	if result != 0 { // since await somehow is faster than Then callback
		t.Errorf("expected 0, got %d", result)
	}
	if result2 != 42 {
		t.Errorf("expected 42, got %d", result2)
	}

}
