package main

// Given an AST, interpret the
// operators in it and return
// a final value.
func generateAST(node *ASTNode, reg int, parentASTOp NodeType) int {
	switch node.op {
	case NodeIf:
		return genIFAST(node)
	case NodeGlue:
		// Do each child statement, and free the
		// registers after each child
		generateAST(node.left, NoReg, node.op)
		genfreeregs()
		generateAST(node.right, NoReg, node.op)
		genfreeregs()
		return NoReg
	}

	leftreg, rightreg := 0, 0
	// Get the left and right sub-tree values
	if node.left != nil {
		leftreg = generateAST(node.left, NoReg, node.op)
	}
	if node.right != nil {
		rightreg = generateAST(node.right, leftreg, node.op)
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
	case NodeEqual, NodeNotEqual, NodeLessThan, NodeGreaterThan, NodeLessThanOrEqual, NodeGreaterThanOrEqual:
		// If the parent AST node is an A_IF, generate a compare
		// followed by a jump. Otherwise, compare registers and
		// set one to 1 or 0 based on the comparison.
		if parentASTOp == NodeIf {
			return cgcompare_and_jump(node.op, leftreg, rightreg, reg)
		}
		return cgcompare_and_set(node.op, leftreg, rightreg)
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
	case NodePrint:
		// Print the left-child's value
		// and return no register
		genprintint(leftreg)
		genfreeregs()
		return NoReg
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

var id int

// Generate and return a new label number
func label() int {
	id++
	return id
}

// Generate the code for an IF statement
// and an optional ELSE clause
func genIFAST(node *ASTNode) int {
	Lfalse, Lend := 0, 0
	// Generate two labels: one for the
	// false compound statement, and one
	// for the end of the overall IF statement.
	// When there is no ELSE clause, Lfalse _is_
	// the ending label!
	Lfalse = label()
	if node.right != nil {
		Lend = label()
	}
	// Generate the condition code followed
	// by a zero jump to the false label.
	// We cheat by sending the Lfalse label as a register.
	generateAST(node.left, Lfalse, node.op)
	genfreeregs()
	// Generate the true compound statement
	generateAST(node.middle, NoReg, node.op)
	genfreeregs()
	// If there is an optional ELSE clause,
	// generate the jump to skip to the end
	if node.right != nil {
		cgjump(Lend)
	}
	// Now the false label
	cglabel(Lfalse)
	// Optional ELSE clause: generate the
	// false compound statement and the
	// end label
	if node.right != nil {
		generateAST(node.right, NoReg, node.op)
		genfreeregs()
		cglabel(Lend)
	}
	return NoReg
}
