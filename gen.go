package main

// Given an AST, interpret the
// operators in it and return
// a final value.
func generateAST(node *ASTNode, reg int) int {
	leftreg, rightreg := 0, 0
	// Get the left and right sub-tree values
	if node.left != nil {
		leftreg = generateAST(node.left, -1)
	}
	if node.right != nil {
		rightreg = generateAST(node.right, leftreg)
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
		return cgloadint(node.value)
	case NodeIdent:
		name, _ := GetSymbolByID(node.value)
		return cgloadglob(name)
	case NodeLvIdent:
		name, _ := GetSymbolByID(node.value)
		return cgstorglob(reg, name)
	case NodeAssign:
		// The work has already been done, return the result
		return rightreg
	default:
		fatal("unknown AST operator %d\n", node.op)
		return 0
	}
}

func genpreamble() {
	cgpreamble()
}

func genpostamble() {
	cgpostamble()
}

func genfreeregs() {
	freeall_registers()
}

func genprintint(reg int) {
	cgprintint(reg)
}

func genglobsym(s string) {
	cgglobsym(s)
}
