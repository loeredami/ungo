package ungo

import "testing"

func TestPipeline(t *testing.T) {
	p := NewPipeSequence(
		PipeFunc[int](func(f int) int {
			return f * 2
		}),
		PipeFunc[int](func(f int) int {
			return f + 1
		}),
	)
	result := p.Run(3)
	if result != 7 {
		t.Errorf("expected result 7, got %d", result)
	}
}
