package ungo

import "unsafe"

type Zone struct {
	data       []byte
	fixed_size int
}

type ZoneWrapper[T any] struct {
	zone *Zone
}

func NewEmptyZone(size int) *Zone {
	return &Zone{data: make([]byte, size), fixed_size: size}
}

func ZoneAlloc[T any](amount uintptr) ZoneWrapper[T] {
	var t T
	buf := make([]byte, unsafe.Sizeof(t)*amount)
	zone := NewEmptyZone(int(amount * unsafe.Sizeof(t)))
	copy(zone.data, buf)
	return ZoneWrapper[T]{zone: zone}
}

func (w ZoneWrapper[T]) Zone() *Zone {
	return w.zone
}

func (w ZoneWrapper[T]) Data() []byte {
	return w.zone.data
}

func (w ZoneWrapper[T]) Free() {
	*w.zone = Zone{}
}

func (w ZoneWrapper[T]) Get() []T {
	var t T
	buf := make([]byte, unsafe.Sizeof(t)*uintptr(w.zone.fixed_size))
	copy(buf, w.zone.data)
	return ReinterpretCast[[]byte, []T](buf)
}

func (w ZoneWrapper[T]) GetRef() []*T {
	positions := []uintptr{}
	for i := uintptr(0); i < uintptr(w.zone.fixed_size); i++ {
		positions = append(positions, i*unsafe.Sizeof(*new(T)))
	}
	pointers := make([]*T, len(positions))
	for i, pos := range positions {
		pointers[i] = ReinterpretCast[*byte, *T](&w.zone.data[pos])
	}
	return pointers
}

func (w ZoneWrapper[T]) GetRefAt(index uintptr) *T {
	return ReinterpretCast[*byte, *T](&w.zone.data[index*unsafe.Sizeof(*new(T))])
}
