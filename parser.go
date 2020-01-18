package main

import (
	"fmt"
)

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

	// Debug: Print what we are about to do
	if node.op == NodeIntLiteral {
		fmt.Printf("int %d\n", node.value)
	} else {
		fmt.Printf("%d %s %d\n", leftval, ASTop[node.op], rightval)
	}
	switch node.op {
	case NodeAdd:
		return (leftval + rightval)
	case NodeSubtract:
		return (leftval - rightval)
	case NodeMultiply:
		return (leftval * rightval)
	case NodeDivide:
		return (leftval / rightval)
	case NodeIntLiteral:
		return (node.value)
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
func binexpr() *ASTNode {
	// Get the integer literal on the left.
	// Fetch the next token at the same time.
	left := primary()
	// If no tokens left, return just the left node
	if CurrentToken.token == TokenEOF {
		return left
	}
	// Convert the token into a node type
	nodetype := arithop(CurrentToken.token)
	// Get the next token in
	scan(CurrentToken)
	// Recursively get the right-hand tree
	right := binexpr()
	// Now build a tree with both sub-trees
	return NewASTNode(nodetype, left, right, 0)
}
