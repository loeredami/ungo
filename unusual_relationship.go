package ungo

// UnusualRelationship defines how two UnusualNums relate to each other.
type UnusualRelationship int8

const (
	UnusualDisjoint    UnusualRelationship = iota // No shared possibilities
	UnusualOverlapping                            // Some shared possibilities
	UnusualIdentical                              // Exactly the same set
)

func CheckUnusualRelationship(u1, u2 *UnusualNum) UnusualRelationship {
	matches := 0
	slice1 := u1.ToSlice()
	slice2 := u2.ToSlice()

	for _, v1 := range slice1 {
		if u2.Contains(v1) {
			matches++
		}
	}

	if matches == 0 {
		return UnusualDisjoint
	}
	if matches == len(slice1) && matches == len(slice2) {
		return UnusualIdentical
	}
	return UnusualOverlapping
}
