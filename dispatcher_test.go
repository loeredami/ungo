package ungo

import (
	"testing"
	"time"
)

func TestDispatcher(t *testing.T) {
	d := NewDispatcher[string]()

	var holder string
	holder = "not working"

	d.Subscribe(Handler[string](func(s string) {
		holder = s
	}))

	d.Emit("working")

	time.Sleep(time.Second)
	if holder != "working" {
		t.Errorf("holder is %s, expected working", holder)
	}
}
