package ungo

import "os"

type Codec[T any] interface {
	Type() string
	Save(*T, *os.File) error
	Load(*os.File, *T) error
	Encode(v T, dst []byte) (int, error)
	Decode(src []byte, v *T) error
	MaxSize() int
	Marshal(v T) ([]byte, error)
	Unmarshal(src []byte, v *T) error
}
