package ungo

import (
	"testing"
	"time"
)

func TestDispatcher(t *testing.T) {
	d := NewDispatcher[string]()

	var holder string
	holder = "not working"

	fn := func(s string) {
		holder = s
	}
	d.Subscribe(Handler[string](&fn))

	d.Emit("working")

	time.Sleep(time.Second)
	if holder != "working" {
		t.Errorf("holder is %s, expected working", holder)
	}
}
