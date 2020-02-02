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
	// Widen NodeChars to NodeInts as required
	if left == NodeChar && right == NodeInt {
		t := OpWiden
		return &t, nil, true
	}
	if left == NodeInt && right == NodeChar {
		if onlyRight {
			return nil, nil, false
		}
		t := OpWiden
		return nil, &t, true
	}
	// Anything remaining is compatible
	return nil, nil, true
}
