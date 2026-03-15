package ungo

import (
	"fmt"
	"testing"
)

func TestException(t *testing.T) {
	result := Try(func() (string, error) {
		return "hello", nil
	})
	if result.Error != nil {
		t.Errorf("expected no error, got %v", result.Error)
	}
	if result.Value != "hello" {
		t.Errorf("expected value 'hello', got %v", result.Value)
	}

	failed := Try(func() (int, error) {
		return 0, fmt.Errorf("error message")
	})

	fail_result := failed.Catch(func(err error) int {
		return -1
	})

	if fail_result != -1 {
		t.Errorf("expected catch to return -1, got %v", fail_result)
	}
}
