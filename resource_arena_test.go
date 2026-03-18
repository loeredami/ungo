package ungo

import (
	"fmt"
	"testing"
)

type DisposableString struct {
	value string
}

func (ds *DisposableString) Dispose() {
	fmt.Println(ds.value, "has been disposed")
}

func TestResourceAerna(t *testing.T) {
	arena := NewResourceArena()

	arena.Track(&DisposableString{value: "A String"})
	arena.Track(&DisposableString{value: "Another String"})

	arena.Melt()
}
