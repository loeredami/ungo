package ungo

import (
	"fmt"
	"math"
)

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64 | ~complex64 | ~complex128
}

// Builds on a weird thought experiment, read more in `unusual_num.md`
type UnusualNum struct {
	v Set[complex128]
}

func NewUnusualNum() *UnusualNum {
	return &UnusualNum{v: NewSet[complex128](0xFFFF)}
}

func (u *UnusualNum) ForceFirst() complex128 {
	if u.v.Size() == 0 {
		return 0
	}
	return u.v.ToSlice()[0]
}

func (u *UnusualNum) ForceFirstReal() float64 {
	if u.v.Size() == 0 {
		return 0
	}
	return real(u.v.ToSlice()[0])
}

func (u *UnusualNum) ForceFirstImag() float64 {
	if u.v.Size() == 0 {
		return 0
	}
	return imag(u.v.ToSlice()[0])
}

func (u *UnusualNum) ForceLast() complex128 {
	if u.v.Size() == 0 {
		return 0
	}
	return u.v.ToSlice()[u.v.Size()-1]
}

func (u *UnusualNum) ForceLastReal() float64 {
	if u.v.Size() == 0 {
		return 0
	}
	return real(u.v.ToSlice()[u.v.Size()-1])
}

func (u *UnusualNum) ForceLastImag() float64 {
	if u.v.Size() == 0 {
		return 0
	}
	return imag(u.v.ToSlice()[u.v.Size()-1])
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
	AddUnknownPossibility(other, num)
	return u.Add(other)
}

func SubtractN[T Number](u *UnusualNum, num T) *UnusualNum {
	other := NewUnusualNum()
	AddUnknownPossibility(other, num)
	return u.Subtract(other)
}

func MultiplyN[T Number](u *UnusualNum, num T) *UnusualNum {
	other := NewUnusualNum()
	AddUnknownPossibility(other, num)
	return u.Multiply(other)
}

func DivideN[T Number](u *UnusualNum, num T) *UnusualNum {
	other := NewUnusualNum()
	AddUnknownPossibility(other, num)
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

// UnusualEvaluatePredicate returns QuantBoolTrue if all possibilities match,
// QuantBoolFalse if none match, and QuantBoolMaybe if some match.
func UnusualEvaluatePredicate(u *UnusualNum, predicate func(complex128) bool) UnsureBool {
	trueCount := 0
	possibilities := u.ToSlice()
	total := len(possibilities)

	for _, p := range possibilities {
		if predicate(p) {
			trueCount++
		}
	}

	if trueCount == total {
		return QuantBoolTrue
	}
	if trueCount == 0 {
		return QuantBoolFalse
	}
	return QuantBoolMaybe
}

// UnusualComparePredicate performs a pairwise comparison of all possibilities between two numbers.
// If the condition is true for all pairs, it returns QuantBoolTrue.
// If it's false for all pairs, it returns QuantBoolFalse.
// Otherwise, it returns QuantBoolMaybe.
func UnusualComparePredicate(u1, u2 *UnusualNum, predicate func(complex128, complex128) bool) UnsureBool {
	trueSeen := false
	falseSeen := false

	for _, v1 := range u1.ToSlice() {
		for _, v2 := range u2.ToSlice() {
			if predicate(v1, v2) {
				trueSeen = true
			} else {
				falseSeen = true
			}
			if trueSeen && falseSeen {
				return QuantBoolMaybe
			}
		}
	}

	if trueSeen {
		return QuantBoolTrue
	}
	return QuantBoolFalse
}

// Certainty returns a value from 0.0 to 1.0.
// 1.0 means there is only one possibility (total certainty).
// 0.0 means the set is empty or infinitely large (though here it's limited by slice size).
func (u *UnusualNum) Certainty() float64 {
	size := u.NumPossibilities()
	if size <= 1 {
		return 1.0
	}
	return 1.0 / float64(size)
}

// CouldBecome returns True if u2 is a possible future state of u1
// (meaning u2 is a subset of u1, or u1 contains the "seed" for u2).
func CouldBecome(u1, u2 *UnusualNum) UnsureBool {
	if u1.NumPossibilities() == 0 {
		return QuantBoolFalse
	}

	allInParent := true
	for _, v2 := range u2.ToSlice() {
		if !u1.Contains(v2) {
			allInParent = false
			break
		}
	}

	if allInParent {
		return QuantBoolTrue
	}

	// If some are in parent and some aren't, it's a "Maybe" transition
	return QuantBoolMaybe
}

// Intersect returns a new UnusualNum containing only the possibilities
// that exist in both u and other.
func (u *UnusualNum) Intersect(other *UnusualNum) *UnusualNum {
	result := NewUnusualNum()
	for _, val := range u.ToSlice() {
		if other.Contains(val) {
			result.AddPossibility(val)
		}
	}
	return result
}

// Split divides the possibilities into two sets: those that satisfy
// the predicate and those that do not.
func (u *UnusualNum) Split(predicate func(complex128) bool) (match, noMatch *UnusualNum) {
	match = NewUnusualNum()
	noMatch = NewUnusualNum()

	for _, val := range u.ToSlice() {
		if predicate(val) {
			match.AddPossibility(val)
		} else {
			noMatch.AddPossibility(val)
		}
	}
	return
}

// Span returns the maximum distance between any two points in the set.
// A Span of 0 means all possibilities are identical (certainty).
func (u *UnusualNum) Span() float64 {
	maxDist := 0.0
	vals := u.ToSlice()
	for i := 0; i < len(vals); i++ {
		for j := i + 1; j < len(vals); j++ {
			// Calculate Euclidean distance between complex points
			diff := vals[i] - vals[j]
			dist := math.Sqrt(real(diff)*real(diff) + imag(diff)*imag(diff))
			if dist > maxDist {
				maxDist = dist
			}
		}
	}
	return maxDist
}

// Centroid returns the average of all current possibilities.
// This represents the "expected value" if every possibility is weighted equally.
func (u *UnusualNum) Centroid() complex128 {
	size := u.v.Size()
	if size == 0 {
		return 0
	}
	var sum complex128
	for _, val := range u.v.ToSlice() {
		sum += val
	}
	return sum / complex(float64(size), 0)
}

// Prune removes all possibilities further than 'radius' from the 'target' point.
func (u *UnusualNum) Prune(target complex128, radius float64) {
	for _, val := range u.v.ToSlice() {
		diff := val - target
		dist := math.Sqrt(real(diff)*real(diff) + imag(diff)*imag(diff))
		if dist > radius {
			u.v.Remove(val)
		}
	}
}

// Hull returns the min/max bounds of all possibilities in the complex plane.
func (u *UnusualNum) Hull() (min, max complex128) {
	vals := u.v.ToSlice()
	if len(vals) == 0 {
		return 0, 0
	}

	minR, minI := real(vals[0]), imag(vals[0])
	maxR, maxI := real(vals[0]), imag(vals[0])

	for _, v := range vals[1:] {
		minR = math.Min(minR, real(v))
		maxR = math.Max(maxR, real(v))
		minI = math.Min(minI, imag(v))
		maxI = math.Max(maxI, imag(v))
	}
	return complex(minR, minI), complex(maxR, maxI)
}

// Quantize rounds all possibilities to the nearest multiple of 'step'.
// This effectively collapses similar realities into a single point.
func (u *UnusualNum) Quantize(step float64) {
	newSet := NewSet[complex128]()
	for _, val := range u.v.ToSlice() {
		r := math.Round(real(val)/step) * step
		i := math.Round(imag(val)/step) * step
		newSet.Add(complex(r, i))
	}
	u.v = newSet
}
