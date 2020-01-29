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

	NodeEqual
	NodeNotEqual
	NodeLessThan
	NodeLessThanOrEqual
	NodeGreaterThan
	NodeGreaterThanOrEqual

	NodeGlue
	NodeIf
	NodeWhile

	NodeIdent
	NodeLvIdent
	NodeAssign
	NodePrint
)

// Abstract Syntax Tree structure
type ASTNode struct {
	op                  NodeType
	left, middle, right *ASTNode
	value               int
}

// Build and return a generic AST node
func NewASTNode(op NodeType, left, middle, right *ASTNode, value int) *ASTNode {
	return &ASTNode{
		op:     op,
		left:   left,
		middle: middle,
		right:  right,
		value:  value,
	}
}

// Make an AST leaf node
func NewLeafASTNode(op NodeType, value int) *ASTNode {
	return NewASTNode(op, nil, nil, nil, value)
}

// Make a unary AST node: only one child
func NewUnaryASTNode(op NodeType, left *ASTNode, value int) *ASTNode {
	return NewASTNode(op, left, nil, nil, value)
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
		node = NewLeafASTNode(NodeIntLiteral, CurrentToken.value)
	case TokenIdent:
		id, ok := GetSymbolIDByString(Text)
		if !ok {
			fatal("unknown variable %s\n", Text)
		}
		node = NewLeafASTNode(NodeIdent, id)
	default:
		fatal("syntax error on line %d\n", Line)
		return nil
	}
	scan(CurrentToken)
	return node
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
	case TokenEqual:
		return NodeEqual
	case TokenNotEqual:
		return NodeNotEqual
	case TokenLessThan:
		return NodeLessThan
	case TokenLessThanOrEqual:
		return NodeLessThanOrEqual
	case TokenGreaterThan:
		return NodeGreaterThan
	case TokenGreaterThanOrEqual:
		return NodeGreaterThanOrEqual
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
		// Join that sub-tree with ours. Convert the token
		// into an AST operation at the same time.
		left = NewASTNode(arithop(tokenType), left, nil, right, 0)
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
