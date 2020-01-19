package main

var CurrentToken *Token = &Token{}

// Node Type
type NodeType int

// Abstract Syntax Tree Node Types
const (
	NodeAdd NodeType = iota
	NodeSubtract
	NodeMultiply
	NodeDivide
	NodeIntLiteral
)

// Abstract Syntax Tree structure
type ASTNode struct {
	op          NodeType
	left, right *ASTNode
	value       int
}

// Build and return a generic AST node
func NewASTNode(op NodeType, left, right *ASTNode, value int) *ASTNode {
	return &ASTNode{
		op:    op,
		left:  left,
		right: right,
		value: value,
	}
}

// Make an AST leaf node
func NewLeafASTNode(op NodeType, value int) *ASTNode {
	return NewASTNode(op, nil, nil, value)
}

// Make a unary AST node: only one child
func NewUnaryASTNode(op NodeType, left *ASTNode, value int) *ASTNode {
	return NewASTNode(op, left, nil, value)
}

// List of AST operators
var ASTop = []string{"+", "-", "*", "/"}

// Given an AST, interpret the
// operators in it and return
// a final value.
func interpretAST(node *ASTNode) int {
	leftval, rightval := 0, 0

	// Get the left and right sub-tree values
	if node.left != nil {
		leftval = interpretAST(node.left)
	}
	if node.right != nil {
		rightval = interpretAST(node.right)
	}

	switch node.op {
	case NodeAdd:
		return leftval + rightval
	case NodeSubtract:
		return leftval - rightval
	case NodeMultiply:
		return leftval * rightval
	case NodeDivide:
		return leftval / rightval
	case NodeIntLiteral:
		return node.value
	default:
		fatal("unknown AST operator %d\n", node.op)
		return 0
	}
}

// Parse a primary factor and return an
// AST node representing it.
func primary() *ASTNode {
	// For an INTLIT token, make a leaf AST node for it
	// and scan in the next token. Otherwise, a syntax error
	// for any other token type.
	switch CurrentToken.token {
	case TokenIntLiteral:
		n := NewLeafASTNode(NodeIntLiteral, CurrentToken.value)
		scan(CurrentToken)
		return n
	default:
		fatal("syntax error on line %d\n", Line)
		return nil
	}
}

// Convert a token into an AST operation.
func arithop(t TokenType) NodeType {
	switch t {
	case TokenPlus:
		return NodeAdd
	case TokenMinus:
		return NodeSubtract
	case TokenStar:
		return NodeMultiply
	case TokenSlash:
		return NodeDivide
	default:
		fatal("unknown token in arithop() on line %d\n", Line)
		return 0
	}
}

// Return an AST tree whose root is a binary operator
func binexpr(previousTokenPrecedence int) *ASTNode {
	// Get the integer literal on the left.
	// Fetch the next token at the same time.
	left := primary()
	tokenType := CurrentToken.token
	// If no tokens left, return just the left node
	if tokenType == TokenEOF {
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
		// Join that sub-tree with ours. Convert the token
		// into an AST operation at the same time.
		left = NewASTNode(arithop(tokenType), left, right, 0)
		// Update the details of the current token.
		tokenType = CurrentToken.token
		// If no tokens left, return just the left node
		if tokenType == TokenEOF {
			return left
		}
	}
	// Return the tree we have when the precedence
	// is the same or lower
	return left
}

// Operator precedence for each token
var OperatorPrecedence = map[TokenType]int{
	TokenEOF:        0,
	TokenPlus:       10,
	TokenMinus:      10,
	TokenStar:       20,
	TokenSlash:      20,
	TokenIntLiteral: 0,
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
