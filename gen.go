package main

// Given an AST, interpret the
// operators in it and return
// a final value.
func generateAST(node *ASTNode, reg int, parentASTOp OpType) int {
	switch node.op {
	case OpIf:
		return genIFAST(node)
	case OpWhile:
		return genWHILE(node)
	case OpGlue:
		// Do each child statement, and free the
		// registers after each child
		generateAST(node.left, NoReg, node.op)
		genfreeregs()
		generateAST(node.right, NoReg, node.op)
		genfreeregs()
		return NoReg
	case OpFunction:
		// Generate the function's preamble before the code
		sym := GetSymbolByID(node.value)
		cgfuncpreamble(sym.name)
		generateAST(node.left, NoReg, node.op)
		cgfuncpostamble(sym)
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
	case OpAdd:
		return cgadd(leftreg, rightreg)
	case OpSubtract:
		return cgsub(leftreg, rightreg)
	case OpMultiply:
		return cgmul(leftreg, rightreg)
	case OpDivide:
		return cgdiv(leftreg, rightreg)
	case OpEqual, OpNotEqual, OpLessThan, OpGreaterThan, OpLessThanOrEqual, OpGreaterThanOrEqual:
		// If the parent AST node is an A_IF, generate a compare
		// followed by a jump. Otherwise, compare registers and
		// set one to 1 or 0 based on the comparison.
		if parentASTOp == OpIf || parentASTOp == OpWhile {
			return cgcompare_and_jump(node.op, leftreg, rightreg, reg)
		}
		return cgcompare_and_set(node.op, leftreg, rightreg)
	case OpIntLiteral:
		return cgloadint(node.value)
	case OpIdent:
		sym := GetSymbolByID(node.value)
		return cgloadglob(sym)
	case OpLvIdent:
		sym := GetSymbolByID(node.value)
		return cgstorglob(reg, sym)
	case OpAssign:
		// The work has already been done, return the result
		return rightreg
	case OpPrint:
		// Print the left-child's value
		// and return no register
		genprintint(leftreg)
		genfreeregs()
		return NoReg
	case OpWiden:
		// Widen the child's type to the parent's type
		return cgwiden(leftreg, node.left.t, node.t)
	case OpReturn:
		sym := GetSymbolByID(FunctionId)
		cgreturn(leftreg, sym)
		return NoReg
	case OpFunctionCall:
		sym := GetSymbolByID(node.value)
		return cgcall(leftreg, sym)
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

func genglobsym(s *Symbol) {
	cgglobsym(s)
}

func genprimsize(t NodeType) int {
	return cgprimsize(t)
}

var currentLabelId int

// Generate and return a new label number
func label() int {
	currentLabelId++
	return currentLabelId
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

// Generate the code for a WHILE statement
// and an optional ELSE clause
func genWHILE(n *ASTNode) int {
	Lstart, Lend := 0, 0
	// Generate the start and end labels
	// and output the start label
	Lstart = label()
	Lend = label()
	cglabel(Lstart)
	// Generate the condition code followed
	// by a jump to the end label.
	// We cheat by sending the Lfalse label as a register.
	generateAST(n.left, Lend, n.op)
	genfreeregs()
	// Generate the compound statement for the body
	generateAST(n.right, NoReg, n.op)
	genfreeregs()
	// Finally output the jump back to the condition,
	// and the end label
	cgjump(Lstart)
	cglabel(Lend)
	return (NoReg)
}
