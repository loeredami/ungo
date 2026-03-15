package ungo

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
)

type UUID string

func NewUUID() UUID {
	b := make([]byte, 16)
	rand.Read(b)
	return UUID(hex.EncodeToString(b))
}

func (u UUID) String() string {
	return string(u)
}

func (u UUID) Bytes() []byte {
	return []byte(u)
}

func (u UUID) Equal(other UUID) bool {
	return u == other
}

func (u UUID) Compare(other UUID) int {
	return bytes.Compare([]byte(u), []byte(other))
}

func (u UUID) Less(other UUID) bool {
	return u.Compare(other) < 0
}

func (u UUID) Greater(other UUID) bool {
	return u.Compare(other) > 0
}

func (u UUID) Between(lower, upper UUID) bool {
	return u.Less(upper) && u.Greater(lower)
}

func (u UUID) Raw() []uint64 {
	b := []byte(u)
	return []uint64{
		((uint64(b[0]) << 56) | (uint64(b[1]) << 48) | (uint64(b[2]) << 40) | (uint64(b[3]) << 32) | (uint64(b[4]) << 24) | (uint64(b[5]) << 16) | (uint64(b[6]) << 8) | uint64(b[7])),
		((uint64(b[8]) << 56) | (uint64(b[9]) << 48) | (uint64(b[10]) << 40) | (uint64(b[11]) << 32) | (uint64(b[12]) << 24) | (uint64(b[13]) << 16) | (uint64(b[14]) << 8) | uint64(b[15])),
	}
}

func UUIDFromRaw(raw []uint64) UUID {
	b := make([]byte, 16)
	b[0] = byte((raw[0] >> 56) & 0xff)
	b[1] = byte((raw[0] >> 48) & 0xff)
	b[2] = byte((raw[0] >> 40) & 0xff)
	b[3] = byte((raw[0] >> 32) & 0xff)
	b[4] = byte((raw[0] >> 24) & 0xff)
	b[5] = byte((raw[0] >> 16) & 0xff)
	b[6] = byte((raw[0] >> 8) & 0xff)
	b[7] = byte(raw[0] & 0xff)
	b[8] = byte((raw[1] >> 56) & 0xff)
	b[9] = byte((raw[1] >> 48) & 0xff)
	b[10] = byte((raw[1] >> 40) & 0xff)
	b[11] = byte((raw[1] >> 32) & 0xff)
	b[12] = byte((raw[1] >> 24) & 0xff)
	b[13] = byte((raw[1] >> 16) & 0xff)
	b[14] = byte((raw[1] >> 8) & 0xff)
	b[15] = byte(raw[1] & 0xff)
	return UUID(hex.EncodeToString(b))
}
