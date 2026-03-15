package ungo

import "testing"

func TestEither(t *testing.T) {
	e := Either[int, string]{}
	if e.Left != nil || e.Right != nil {
		t.Errorf("either should be empty")
	}
	e.Left = new(int)
	if e.Left == nil {
		t.Errorf("either left should not be nil")
	}
	e.Right = new(string)
	if e.Right == nil {
		t.Errorf("either right should not be nil")
	}

	is_left := e.IsLeft()
	if !is_left {
		t.Errorf("either should be left")
	}

	is_right := e.IsRight()
	if !is_right {
		t.Errorf("either should be right")
	}

	// Test with right value
	e.Left = nil
	e.Right = new(string)
	is_left = e.IsLeft()
	if is_left {
		t.Errorf("either should not be left")
	}
	is_right = e.IsRight()
	if !is_right {
		t.Errorf("either should be right")
	}

	// Test with nil value
	e.Left = nil
	e.Right = nil
	is_left = e.IsLeft()
	if is_left {
		t.Errorf("either should not be left")
	}
	is_right = e.IsRight()
	if is_right {
		t.Errorf("either should not be right")
	}
}
