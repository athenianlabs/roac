package main

// Given two primitive types, return true if they are compatible,
// false otherwise. Also return either zero or an OpWiden
// operation if one has to be widened to match the other.
// If onlyRight is true, only widen left to right.
func typeCompatible(left, right NodeType, onlyRight bool) (*OpType, *OpType, bool) {
	// Voids not compatible with anything
	if left == NodeVoid || right == NodeVoid {
		return nil, nil, false
	}
	// Same types, they are compatible
	if left == right {
		return nil, nil, true
	}
	// Get the sizes for each type
	leftSize := genprimsize(left)
	rightSize := genprimsize(right)
	// Types with zero size are not compatible with anything
	if leftSize == 0 || rightSize == 0 {
		return nil, nil, false
	}
	// Widen types as required
	if leftSize < rightSize {
		t := OpWiden
		return &t, nil, true
	}
	if rightSize < leftSize {
		if onlyRight {
			return nil, nil, false
		}
		t := OpWiden
		return nil, &t, true
	}
	// Anything remaining is the same size and thus compatible
	return nil, nil, true
}
