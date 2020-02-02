package main

const NoReg = -1

type NodeType int

const (
	NodeNone NodeType = iota
	NodeVoid
	NodeChar
	NodeInt
)

type StructuralNodeType int

const (
	NodeVariable StructuralNodeType = iota
	NodeFunction
)

// Op Type
type OpType int

// Abstract Syntax Tree Op Types
const (
	_ OpType = iota
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide

	OpIntLiteral

	OpEqual
	OpNotEqual
	OpLessThan
	OpLessThanOrEqual
	OpGreaterThan
	OpGreaterThanOrEqual

	OpGlue
	OpIf
	OpWhile
	OpFunction

	OpIdent
	OpLvIdent
	OpAssign
	OpPrint
	OpWiden
)

// Abstract Syntax Tree structure
type ASTNode struct {
	op                  OpType
	t                   NodeType
	left, middle, right *ASTNode
	value               int
}

// Build and return a generic AST node
func NewASTNode(op OpType, t NodeType, left, middle, right *ASTNode, value int) *ASTNode {
	return &ASTNode{
		op:     op,
		t:      t,
		left:   left,
		middle: middle,
		right:  right,
		value:  value,
	}
}

// Make an AST leaf node
func NewLeafASTNode(op OpType, t NodeType, value int) *ASTNode {
	return NewASTNode(op, t, nil, nil, nil, value)
}

// Make a unary AST node: only one child
func NewUnaryASTNode(op OpType, t NodeType, left *ASTNode, value int) *ASTNode {
	return NewASTNode(op, t, left, nil, nil, value)
}
