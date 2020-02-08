package main

// Parse the declaration of a variable
func varDeclaration() {
	// Get the type of the variable, then the identifier
	t := parseType()
	ident()
	sym := AddSymbol(Text, t, NodeVariable, 0)
	genglobsym(sym)
	semi()
}

// Parse the declaration of a simplistic function
func functionDeclaration() *ASTNode {
	// Get the type of the variable, then the identifier
	t := parseType()
	ident()

	// Get a label-id for the end label, add the function
	// to the symbol table, and set the Functionid global
	// to the function's symbol-id
	sym := AddSymbol(Text, t, NodeFunction, label())
	FunctionId = sym.id

	// Scan in the parentheses
	lparen()
	rparen()

	// Get the AST tree for the compound statement
	tree := compoundStatement()

	// If the function type isn't P_VOID, check that
	// the last AST operation in the compound statement
	// was a return statement
	if t != NodeVoid {
		finalstmt := tree
		if tree.op == OpGlue {
			finalstmt = tree.right
		}
		if finalstmt == nil || finalstmt.op != OpReturn {
			fatal("No return for function with non-void type\n")
		}
	}
	// Return an A_FUNCTION node which has the function's nameslot
	// and the compound statement sub-tree
	return NewUnaryASTNode(OpFunction, t, tree, sym.id)
}
