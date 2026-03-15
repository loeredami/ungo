package ungo

import "testing"

func TestBitmap(t *testing.T) {
	b := NewBitmap(2) // 2 bytes
	b.Set(0)          // set bit 0
	b.Set(8)          // set bit 8

	if !b.Test(0) {
		t.Errorf("bit 0 should be set")
	}
	if !b.Test(8) {
		t.Errorf("bit 8 should be set")
	}

	if b.Test(1) {
		t.Errorf("bit 1 should not be set")
	}

	if b.Test(9) {
		t.Errorf("bit 9 should not be set")
	}
}

func TestBitmap_Dumps(t *testing.T) {
	b := NewBitmap(2) // 2 bytes
	b.Set(0)          // set bit 0
	b.Set(8)          // set bit 8

	dump := make([]byte, 2)
	b.Dump(dump)
	if len(dump) != 2 {
		t.Errorf("dump should be 2 bytes")
	}
	if dump[0] != 1 {
		t.Errorf("dump[0] should be 1")
	}
	if dump[1] != 1 {
		t.Errorf("dump[1] should be 1")
	}
}
