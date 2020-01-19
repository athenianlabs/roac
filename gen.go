package main

// Given an AST, interpret the
// operators in it and return
// a final value.
func generateAST(node *ASTNode) int {
	leftreg, rightreg := 0, 0
	// Get the left and right sub-tree values
	if node.left != nil {
		leftreg = generateAST(node.left)
	}
	if node.right != nil {
		rightreg = generateAST(node.right)
	}
	switch node.op {
	case NodeAdd:
		return cgadd(leftreg, rightreg)
	case NodeSubtract:
		return cgsub(leftreg, rightreg)
	case NodeMultiply:
		return cgmul(leftreg, rightreg)
	case NodeDivide:
		return cgdiv(leftreg, rightreg)
	case NodeIntLiteral:
		return cgload(node.value)
	default:
		fatal("unknown AST operator %d\n", node.op)
		return 0
	}
}

func generatecode(n *ASTNode) {
	reg := 0
	cgpreamble()
	reg = generateAST(n)
	cgprintint(reg)
	cgpostamble()
}
