package main

// Parse one or more global declarations, either
// variables or functions
func globalDeclarations() {
	for {
		// We have to read past the type and identifier
		// to see either a '(' for a function declaration
		// or a ',' or ';' for a variable declaration.
		// Text is filled in by the ident() call.
		t := parseType()
		ident()
		if CurrentToken.token == TokenLeftParen {
			// Parse the function declaration and
			// generate the assembly code for it
			tree := functionDeclaration(t)
			generateAST(tree, NoReg, 0)
		} else {
			// Parse the global variable declaration
			varDeclaration(t)
		}
		// Stop when we have reached EOF
		if CurrentToken.token == TokenEOF {
			return
		}
	}
}

// Parse the declaration of a list of variables.
// The identifier has been scanned & we have the type
func varDeclaration(t NodeType) {
	for {
		// Text now has the identifier's name.
		// Add it as a known identifier
		// and generate its space in assembly
		sym := AddSymbol(Text, t, NodeVariable, 0)
		genglobsym(sym)
		// If the next token is a semicolon,
		// skip it and return.
		if CurrentToken.token == TokenSemicolon {
			scan(CurrentToken)
			return
		}
		// If the next token is a comma, skip it,
		// get the identifier and loop back
		if CurrentToken.token == TokenComma {
			scan(CurrentToken)
			ident()
			continue
		}
		fatal("Missing , or ; after identifier\n")
	}
}

// Parse the declaration of a simplistic function.
// The identifier has been scanned & we have the type
func functionDeclaration(t NodeType) *ASTNode {
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
