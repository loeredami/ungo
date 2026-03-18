package ungo

type UnusualVector struct {
	nums []UnusualNum
}

func MakeUnusualVector(nums ...UnusualNum) *UnusualVector {
	return &UnusualVector{
		nums: nums,
	}
}

func (v *UnusualVector) Add(v2 *UnusualVector) {
	bigger := v
	if len(v2.nums) > len(v.nums) {
		bigger = v2
	}
	for i, num := range bigger.nums {
		if i < len(v.nums) {
			v.nums[i].Add(&num)
		} else {
			v.nums = append(v.nums, num)
		}
	}
}

func (v *UnusualVector) Subtract(v2 *UnusualVector) {
	bigger := v
	if len(v2.nums) > len(v.nums) {
		bigger = v2
	}
	for i, num := range bigger.nums {
		if i < len(v.nums) {
			v.nums[i].Subtract(&num)
		} else {
			v.nums = append(v.nums, num)
		}
	}
}

func (v *UnusualVector) Multiply(v2 *UnusualVector) {
	bigger := v
	if len(v2.nums) > len(v.nums) {
		bigger = v2
	}
	for i, num := range bigger.nums {
		if i < len(v.nums) {
			v.nums[i].Multiply(&num)
		} else {
			v.nums = append(v.nums, num)
		}
	}
}

func (v *UnusualVector) Divide(v2 *UnusualVector) {
	bigger := v
	if len(v2.nums) > len(v.nums) {
		bigger = v2
	}
	for i, num := range bigger.nums {
		if i < len(v.nums) {
			v.nums[i].Divide(&num)
		} else {
			v.nums = append(v.nums, num)
		}
	}
}

func (v *UnusualVector) Dot(v2 *UnusualVector) *UnusualNum {
	result := NewUnusualNum()
	for i, num := range v.nums {
		if i < len(v2.nums) {
			result.Add(&num)
		}
	}
	return result
}

func (v *UnusualVector) Length() *UnusualNum {
	result := NewUnusualNum()
	for _, num := range v.nums {
		result.Add(&num)
	}
	return result
}

func (v *UnusualVector) ForEach(f func(*UnusualNum)) *UnusualVector {
	for i := range v.nums {
		f(&v.nums[i])
	}
	return v
}

func (v *UnusualVector) Set(index int, num UnusualNum) {
	if index < 0 || index >= len(v.nums) {
		v.nums = append(v.nums, num)
	} else {
		v.nums[index] = num
	}
}

func (v *UnusualVector) At(index int) *UnusualNum {
	if index < 0 || index >= len(v.nums) {
		return NewUnusualNum()
	}
	return &v.nums[index]
}

func (v *UnusualVector) Normalize() *UnusualVector {
	length := v.Length()
	v.ForEach(func(num *UnusualNum) {
		num.Divide(length)
	})
	return v
}

func (v *UnusualVector) Clone() *UnusualVector {
	clone := &UnusualVector{nums: make([]UnusualNum, len(v.nums))}
	copy(clone.nums, v.nums)
	return clone
}
