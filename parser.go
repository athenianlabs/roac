package main

var CurrentToken = &Token{}

// Operator precedence for each token
var OperatorPrecedence = map[TokenType]int{
	TokenEOF:                0,
	TokenPlus:               10,
	TokenMinus:              10,
	TokenStar:               20,
	TokenSlash:              20,
	TokenIntLiteral:         0,
	TokenEqual:              30,
	TokenNotEqual:           30,
	TokenLessThan:           40,
	TokenLessThanOrEqual:    40,
	TokenGreaterThan:        40,
	TokenGreaterThanOrEqual: 40,
}

// Check that we have a binary operator and
// return its precedence.
func OpPrecedence(t TokenType) int {
	prec := OperatorPrecedence[t]
	if prec == 0 {
		fatal("syntax error on line %d, token %d\n", Line, t)
	}
	return prec
}

// Return an AST tree whose root is a binary operator
func binexpr(previousTokenPrecedence int) *ASTNode {
	// Get the integer literal on the left.
	// Fetch the next token at the same time.
	left := primary()
	tokenType := CurrentToken.token
	// If no tokens left, return just the left node
	if tokenType == TokenSemicolon || tokenType == TokenRightParen {
		return left
	}
	// While the precedence of this token is
	// more than that of the previous token precedence
	for OpPrecedence(tokenType) > previousTokenPrecedence {
		// Fetch in the next integer literal
		scan(CurrentToken)
		// Recursively call binexpr() with the
		// precedence of our token to build a sub-tree
		right := binexpr(OperatorPrecedence[tokenType])
		// Ensure the two types are compatible.
		leftOp, rightOp, ok := typeCompatible(left.t, right.t, false)
		if !ok {
			fatal("incompatible types\n")
		}
		// Widen either side if required. type vars are A_WIDEN now
		if leftOp != nil {
			left = NewUnaryASTNode(*leftOp, right.t, left, 0)
		}
		if rightOp != nil {
			right = NewUnaryASTNode(*rightOp, left.t, right, 0)
		}
		// Join that sub-tree with ours. Convert the token
		// into an AST operation at the same time.
		left = NewASTNode(arithop(tokenType), left.t, left, nil, right, 0)
		// Update the details of the current token.
		tokenType = CurrentToken.token
		// If no tokens left, return just the left node
		if tokenType == TokenSemicolon || tokenType == TokenRightParen {
			return left
		}
	}
	// Return the tree we have when the precedence
	// is the same or lower
	return left
}

// Parse a primary factor and return an
// AST node representing it.
func primary() *ASTNode {
	var node *ASTNode
	// For an INTLIT token, make a leaf AST node for it
	// and scan in the next token. Otherwise, a syntax error
	// for any other token type.
	switch CurrentToken.token {
	case TokenIntLiteral:
		// For an INTLIT token, make a leaf AST node for it.
		// Make it a P_CHAR if it's within the P_CHAR range
		if CurrentToken.value >= 0 && CurrentToken.value < 256 {
			node = NewLeafASTNode(OpIntLiteral, NodeChar, CurrentToken.value)
		} else {
			node = NewLeafASTNode(OpIntLiteral, NodeInt, CurrentToken.value)
		}
	case TokenIdent:
		sym := GetSymbolByString(Text)
		node = NewLeafASTNode(OpIdent, sym.t, sym.id)
	default:
		fatal("syntax error on line %d\n", Line)
		return nil
	}
	scan(CurrentToken)
	return node
}

// Convert a token into an AST operation.
func arithop(t TokenType) OpType {
	switch t {
	case TokenPlus:
		return OpAdd
	case TokenMinus:
		return OpSubtract
	case TokenStar:
		return OpMultiply
	case TokenSlash:
		return OpDivide
	case TokenEqual:
		return OpEqual
	case TokenNotEqual:
		return OpNotEqual
	case TokenLessThan:
		return OpLessThan
	case TokenLessThanOrEqual:
		return OpLessThanOrEqual
	case TokenGreaterThan:
		return OpGreaterThan
	case TokenGreaterThanOrEqual:
		return OpGreaterThanOrEqual
	default:
		fatal("unknown token in arithop() on line %d\n", Line)
		return 0
	}
}
