package ungo

type UnsureBool int8

const (
	QuantBoolTrue  UnsureBool = 1
	QuantBoolFalse UnsureBool = 0
	QuantBoolMaybe UnsureBool = -1
)

func (q UnsureBool) Bool() bool {
	return q != QuantBoolFalse
}

func (q UnsureBool) Maybe() bool {
	return q == QuantBoolMaybe
}

func (q UnsureBool) True() bool {
	return q == QuantBoolTrue
}

func (q UnsureBool) False() bool {
	return q == QuantBoolFalse
}

func (q UnsureBool) Ensure(b bool) bool {
	if q == QuantBoolMaybe {
		return b
	}
	return q.Bool()
}

func (q UnsureBool) Not() UnsureBool {
	switch q {
	case QuantBoolTrue:
		return QuantBoolFalse
	case QuantBoolFalse:
		return QuantBoolTrue
	default:
		return QuantBoolMaybe
	}
}

func (q UnsureBool) And(other UnsureBool) UnsureBool {
	if q == QuantBoolFalse || other == QuantBoolFalse {
		return QuantBoolFalse
	}
	if q == QuantBoolTrue && other == QuantBoolTrue {
		return QuantBoolTrue
	}
	return QuantBoolMaybe
}

func (q UnsureBool) Or(other UnsureBool) UnsureBool {
	if q == QuantBoolTrue || other == QuantBoolTrue {
		return QuantBoolTrue
	}
	if q == QuantBoolFalse && other == QuantBoolFalse {
		return QuantBoolFalse
	}
	return QuantBoolMaybe
}

func (q UnsureBool) String() string {
	switch q {
	case QuantBoolTrue:
		return "True"
	case QuantBoolFalse:
		return "False"
	default:
		return "Maybe"
	}
}

func FromBool(b bool) UnsureBool {
	if b {
		return QuantBoolTrue
	}
	return QuantBoolFalse
}

func IfQ[T any](condition UnsureBool, trueValue T, falseValue T) T {
	if condition.Bool() {
		return trueValue
	}
	return falseValue
}

// All returns True if every input is True, False if any is False, and Maybe otherwise.
func UnsureAll(bools ...UnsureBool) UnsureBool {
	res := QuantBoolTrue
	for _, b := range bools {
		res = res.And(b)
	}
	return res
}

// Any returns True if any input is True, False if all are False, and Maybe otherwise.
func UnsureAny(bools ...UnsureBool) UnsureBool {
	res := QuantBoolFalse
	for _, b := range bools {
		res = res.Or(b)
	}
	return res
}

// TruthDensity returns a float between 0 and 1 representing
// how much of the "possibility space" is currently True.
func TruthDensity(bools ...UnsureBool) float64 {
	if len(bools) == 0 {
		return 0
	}
	var sum float64
	for _, b := range bools {
		switch b {
		case QuantBoolTrue:
			sum += 1.0
		case QuantBoolMaybe:
			sum += 0.5
		case QuantBoolFalse:
			sum += 0.0
		}
	}
	return sum / float64(len(bools))
}
