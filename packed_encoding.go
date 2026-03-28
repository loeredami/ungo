package ungo

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"unsafe"
)

type PackedType struct {
	ID       uint16
	SizeBits uint64
}

type PackedHeaderEntry struct {
	TypeID   uint16
	Location uint64
}

type PackedEncoding struct {
	Types  map[uint16]PackedType
	Header []PackedHeaderEntry
	Data   []byte
}

func NewPackedEncoding() *PackedEncoding {
	return &PackedEncoding{
		Types: make(map[uint16]PackedType),
	}
}

func (pe *PackedEncoding) RegisterType(id uint16, bits uint64) {
	pe.Types[id] = PackedType{ID: id, SizeBits: bits}
}

func Add[T any](pe *PackedEncoding, typeID uint16, value T) error {
	t, ok := pe.Types[typeID]
	if !ok {
		return fmt.Errorf("type %d not registered", typeID)
	}

	var nextLoc uint64
	if len(pe.Header) > 0 {
		last := pe.Header[len(pe.Header)-1]
		nextLoc = last.Location + pe.Types[last.TypeID].SizeBits
	}

	pe.Header = append(pe.Header, PackedHeaderEntry{
		TypeID:   typeID,
		Location: nextLoc,
	})

	sizeT := unsafe.Sizeof(value)
	valBytes := unsafe.Slice((*byte)(unsafe.Pointer(&value)), sizeT)

	pe.writeBits(nextLoc, t.SizeBits, valBytes)
	return nil
}

func Get[T any](pe *PackedEncoding, index int) T {
	entry := pe.Header[index]
	sizeBits := pe.Types[entry.TypeID].SizeBits

	raw := pe.getRawBits(entry.Location, sizeBits)

	var result T
	sizeT := unsafe.Sizeof(result)
	resBytes := unsafe.Slice((*byte)(unsafe.Pointer(&result)), sizeT)

	// Copy the extracted bits into the memory of T
	copy(resBytes, raw)
	return result
}

func (pe *PackedEncoding) writeBits(offset uint64, sizeBits uint64, src []byte) {
	neededBytes := int((offset + sizeBits + 7) / 8)
	if len(pe.Data) < neededBytes {
		newBuf := make([]byte, neededBytes)
		copy(newBuf, pe.Data)
		pe.Data = newBuf
	}

	for i := uint64(0); i < sizeBits; i++ {
		destBitPos := offset + i
		srcBitPos := i

		if (src[srcBitPos/8] & (1 << (srcBitPos % 8))) != 0 {
			pe.Data[destBitPos/8] |= (1 << (destBitPos % 8))
		} else {
			pe.Data[destBitPos/8] &= ^(1 << (destBitPos % 8))
		}
	}
}

func (pe *PackedEncoding) getRawBits(offset uint64, sizeBits uint64) []byte {
	resBytes := (sizeBits + 7) / 8
	result := make([]byte, resBytes)

	for i := uint64(0); i < sizeBits; i++ {
		srcBitPos := offset + i
		destBitPos := i

		if (pe.Data[srcBitPos/8] & (1 << (srcBitPos % 8))) != 0 {
			result[destBitPos/8] |= (1 << (destBitPos % 8))
		}
	}
	return result
}

// --- File I/O ---

func (pe *PackedEncoding) WriteToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// 1. Types
	binary.Write(f, binary.LittleEndian, uint32(len(pe.Types)))
	for _, t := range pe.Types {
		binary.Write(f, binary.LittleEndian, t.ID)
		binary.Write(f, binary.LittleEndian, t.SizeBits)
	}

	// 2. Header
	binary.Write(f, binary.LittleEndian, uint32(len(pe.Header)))
	for _, h := range pe.Header {
		binary.Write(f, binary.LittleEndian, h.TypeID)
		binary.Write(f, binary.LittleEndian, h.Location)
	}

	// 3. Data
	binary.Write(f, binary.LittleEndian, uint32(len(pe.Data)))
	_, err = f.Write(pe.Data)
	return err
}

func (pe *PackedEncoding) ReadFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var count uint32

	// 1. Types
	if err := binary.Read(f, binary.LittleEndian, &count); err != nil {
		return err
	}
	pe.Types = make(map[uint16]PackedType)
	for i := uint32(0); i < count; i++ {
		var t PackedType
		binary.Read(f, binary.LittleEndian, &t.ID)
		binary.Read(f, binary.LittleEndian, &t.SizeBits)
		pe.Types[t.ID] = t
	}

	// 2. Header
	if err := binary.Read(f, binary.LittleEndian, &count); err != nil {
		return err
	}
	pe.Header = make([]PackedHeaderEntry, count)
	for i := uint32(0); i < count; i++ {
		binary.Read(f, binary.LittleEndian, &pe.Header[i].TypeID)
		binary.Read(f, binary.LittleEndian, &pe.Header[i].Location)
	}

	// 3. Data
	if err := binary.Read(f, binary.LittleEndian, &count); err != nil {
		return err
	}
	pe.Data = make([]byte, count)
	_, err = io.ReadFull(f, pe.Data)
	return err
}
