package main

// Parse the current token and
// return a primitive type enum value
func parseType() NodeType {
	nt := NodeType(0)
	switch CurrentToken.token {
	case TokenChar:
		nt = NodeChar
	case TokenInt:
		nt = NodeInt
	case TokenLong:
		nt = NodeLong
	case TokenVoid:
		nt = NodeVoid
	default:
		fatal("Illegal type, token %d\n", CurrentToken.token)
	}
	// Scan in one or more further '*' tokens
	// and determine the correct pointer type
	for {
		scan(CurrentToken)
		if CurrentToken.token != TokenStar {
			break
		}
		nt = pointerTo(nt)
	}
	// We leave with the next token already scanned
	return nt
}

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

// Given a primitive type, return
// the type which is a pointer to it
func pointerTo(t NodeType) NodeType {
	switch t {
	case NodeVoid:
		return NodeVoidPointer
	case NodeChar:
		return NodeCharPointer
	case NodeInt:
		return NodeIntPointer
	case NodeLong:
		return NodeLongPointer
	default:
		fatal("unrecognized in pointerTo %v\n", t)
	}
	return 0
}

// Given a primitive pointer type, return
// the type which it points to
func valueAt(t NodeType) NodeType {
	switch t {
	case NodeVoidPointer:
		return NodeVoid
	case NodeCharPointer:
		return NodeChar
	case NodeIntPointer:
		return NodeInt
	case NodeLongPointer:
		return NodeLong
	default:
		fatal("unrecognized in valueAt %v\n", t)
	}
	return 0
}
