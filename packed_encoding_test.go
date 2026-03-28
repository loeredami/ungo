package ungo

import (
	"bytes"
	"os"
	"testing"
)

func TestPackedEncoding_FullCycle(t *testing.T) {
	filename := "test_compact.bin"
	defer os.Remove(filename)

	pe := NewPackedEncoding()

	typeID4Bit := uint16(1)
	typeID12Bit := uint16(2)
	typeID128Bit := uint16(3)

	pe.RegisterType(typeID4Bit, 4)
	pe.RegisterType(typeID12Bit, 12)
	pe.RegisterType(typeID128Bit, 128)

	val1 := uint8(0x0A)
	if err := Add(pe, typeID4Bit, val1); err != nil {
		t.Fatalf("Failed to add 4-bit val: %v", err)
	}

	val2 := uint16(0x0FFF)
	if err := Add(pe, typeID12Bit, val2); err != nil {
		t.Fatalf("Failed to add 12-bit val: %v", err)
	}

	val3 := [16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10, 0x11, 0x12}
	if err := Add(pe, typeID128Bit, val3); err != nil {
		t.Fatalf("Failed to add 128-bit val: %v", err)
	}

	if len(pe.Data) != 18 {
		t.Errorf("Expected 18 bytes of packed data, got %d", len(pe.Data))
	}

	if err := pe.WriteToFile(filename); err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}

	pe2 := NewPackedEncoding()
	if err := pe2.ReadFromFile(filename); err != nil {
		t.Fatalf("Failed to read from file: %v", err)
	}

	res1 := Get[uint8](pe2, 0)
	if res1 != 0x0A {
		t.Errorf("Value 1 mismatch: expected 0x0A, got 0x%X", res1)
	}

	res2 := Get[uint16](pe2, 1)
	if res2 != 0x0FFF {
		t.Errorf("Value 2 mismatch: expected 0x0FFF, got 0x%X", res2)
	}

	res3 := Get[[16]byte](pe2, 2)
	if !bytes.Equal(res3[:], val3[:]) {
		t.Errorf("Value 3 mismatch: expected %v, got %v", val3, res3)
	}
}

func TestPackedEncoding_OverlappingAlignment(t *testing.T) {
	pe := NewPackedEncoding()
	type3Bit := uint16(1)
	pe.RegisterType(type3Bit, 3)

	Add(pe, type3Bit, uint8(7))
	Add(pe, type3Bit, uint8(7))
	Add(pe, type3Bit, uint8(7))

	if len(pe.Data) != 2 {
		t.Errorf("Expected 2 bytes for 9 bits, got %d", len(pe.Data))
	}

	if Get[uint8](pe, 0) != 7 || Get[uint8](pe, 1) != 7 || Get[uint8](pe, 2) != 7 {
		t.Error("Recovered bit-aligned values do not match")
	}
}
