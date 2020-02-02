package main

// Parse the current token and
// return a primitive type enum value
func parseType(t TokenType) NodeType {
	switch t {
	case TokenChar:
		return NodeChar
	case TokenInt:
		return NodeInt
	case TokenVoid:
		return NodeVoid
	default:
		fatal("Illegal type, token %d\n", t)
	}
	return 0
}

// Parse the declaration of a variable
func varDeclaration() {
	// Get the type of the variable, then the identifier
	t := parseType(CurrentToken.token)
	scan(CurrentToken)
	ident()
	sym := AddSymbol(Text, t, NodeVariable)
	genglobsym(sym)
	semi()
}

// Parse the declaration of a simplistic function
func functionDeclaration() *ASTNode {
	// Find the 'void', the identifier, and the '(' ')'.
	// For now, do nothing with them
	match(TokenVoid, "void")
	ident()
	sym := AddSymbol(Text, NodeVoid, NodeFunction)
	lparen()
	rparen()

	// Get the AST tree for the compound statement
	tree := compoundStatement()
	// Return an A_FUNCTION node which has the function's nameslot
	// and the compound statement sub-tree
	return NewUnaryASTNode(OpFunction, NodeVoid, tree, sym.id)
}
