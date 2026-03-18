package ungo

import (
	"fmt"
	"math"
)

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | complex64 | complex128
}

// Builds on a weird thought experiment, read more in unusual_num.md
type UnusualNum struct {
	v Set[complex128]
}

func NewUnusualNum() *UnusualNum {
	return &UnusualNum{v: Set[complex128]{}}
}

func (u *UnusualNum) String() string {
	if u.v.Size() == 0 {
		return "0"
	}
	if u.v.Size() == 1 {
		return fmt.Sprintf("%v", u.v.ToSlice()[0])
	}
	result := "{"
	for _, ni := range u.v.ToSlice() {
		result += fmt.Sprintf("%v ", ni)
	}
	result = result[:len(result)-1] + "}"
	return result
}

func (u *UnusualNum) AddPossibility(ni complex128) {
	u.v.Add(ni)
}

func AddUnknownPossibility[T Number](u *UnusualNum, num T) {
	switch v := any(num).(type) {
	case int:
		u.AddPossibility(complex(float64(v), 0))
	case int8:
		u.AddPossibility(complex(float64(v), 0))
	case int16:
		u.AddPossibility(complex(float64(v), 0))
	case int32:
		u.AddPossibility(complex(float64(v), 0))
	case int64:
		u.AddPossibility(complex(float64(v), 0))
	case uint:
		u.AddPossibility(complex(float64(v), 0))
	case uint8:
		u.AddPossibility(complex(float64(v), 0))
	case uint16:
		u.AddPossibility(complex(float64(v), 0))
	case uint32:
		u.AddPossibility(complex(float64(v), 0))
	case uint64:
		u.AddPossibility(complex(float64(v), 0))
	case float32:
		u.AddPossibility(complex(float64(v), 0))
	case float64:
		u.AddPossibility(complex(float64(v), 0))
	case complex64:
		u.AddPossibility(complex(float64(real(v)), float64(imag(v))))
	case complex128:
		u.AddPossibility(v)
	}
}

func AddN[T Number](u *UnusualNum, num T) *UnusualNum {
	other := NewUnusualNum()
	AddUnknownPossibility[T](other, num)
	return u.Add(other)
}

func SubtractN[T Number](u *UnusualNum, num T) *UnusualNum {
	other := NewUnusualNum()
	AddUnknownPossibility[T](other, num)
	return u.Subtract(other)
}

func MultiplyN[T Number](u *UnusualNum, num T) *UnusualNum {
	other := NewUnusualNum()
	AddUnknownPossibility[T](other, num)
	return u.Multiply(other)
}

func DivideN[T Number](u *UnusualNum, num T) *UnusualNum {
	other := NewUnusualNum()
	AddUnknownPossibility[T](other, num)
	return u.Divide(other)
}

func (u *UnusualNum) Add(other *UnusualNum) *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		for _, nj := range other.v.ToSlice() {
			result.AddPossibility(ni + nj)
		}
	}
	return result
}

func (u *UnusualNum) Subtract(other *UnusualNum) *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		for _, nj := range other.v.ToSlice() {
			result.AddPossibility(ni - nj)
		}
	}
	return result
}

func (u *UnusualNum) Multiply(other *UnusualNum) *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		for _, nj := range other.v.ToSlice() {
			result.AddPossibility(ni * nj)
		}
	}
	return result
}

func (u *UnusualNum) Divide(other *UnusualNum) *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		for _, nj := range other.v.ToSlice() {
			if nj == 0 {
				result.AddPossibility(complex(0, math.Inf(1)))
				result.AddPossibility(complex(0, imag(ni)))
				result.AddPossibility(complex(0, 0))
				result.AddPossibility(complex(imag(ni), 0))
				result.AddPossibility(complex(math.Inf(1), 0))

				result.AddPossibility(complex(0, -math.Inf(1)))
				result.AddPossibility(complex(0, -imag(ni)))
				result.AddPossibility(complex(-imag(ni), 0))
				result.AddPossibility(complex(-math.Inf(1), 0))
				continue
			}
			result.AddPossibility(ni / nj)
		}
	}
	return result
}

// more mathematic functions

func (u *UnusualNum) Abs() *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		result.AddPossibility(complex(math.Abs(real(ni)), math.Abs(imag(ni))))
	}
	return result
}

func (u *UnusualNum) Log() *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		result.AddPossibility(complex(math.Log(math.Abs(real(ni))), math.Log(math.Abs(imag(ni)))))
	}
	return result
}

func (u *UnusualNum) Exp() *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		result.AddPossibility(complex(math.Exp(real(ni)), math.Exp(imag(ni))))
	}
	return result
}

func (u *UnusualNum) Sin() *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		result.AddPossibility(complex(math.Sin(real(ni)), math.Sin(imag(ni))))
	}
	return result
}

func (u *UnusualNum) Cos() *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		result.AddPossibility(complex(math.Cos(real(ni)), math.Cos(imag(ni))))
	}
	return result
}

func (u *UnusualNum) Tan() *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		result.AddPossibility(complex(math.Tan(real(ni)), math.Tan(imag(ni))))
	}
	return result
}

func (u *UnusualNum) Asin() *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		result.AddPossibility(complex(math.Asin(real(ni)), math.Asin(imag(ni))))
	}
	return result
}

func (u *UnusualNum) Acos() *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		result.AddPossibility(complex(math.Acos(real(ni)), math.Acos(imag(ni))))
	}
	return result
}

func (u *UnusualNum) Atan() *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		result.AddPossibility(complex(math.Atan(real(ni)), math.Atan(imag(ni))))
	}
	return result
}

func (u *UnusualNum) Sqrt() *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		result.AddPossibility(complex(math.Sqrt(real(ni)), math.Sqrt(imag(ni))))
	}
	return result
}

// utility

func (u *UnusualNum) ToSlice() []complex128 {
	return u.v.ToSlice()
}

func (u *UnusualNum) NumPossibilities() int {
	return u.v.Size()
}

func (u *UnusualNum) Contains(ni complex128) bool {
	return u.v.Contains(ni)
}

func (u *UnusualNum) Constraint(filter func(complex128) bool) *UnusualNum {
	result := NewUnusualNum()
	for _, ni := range u.v.ToSlice() {
		if filter(ni) {
			result.AddPossibility(ni)
		}
	}
	return result
}
