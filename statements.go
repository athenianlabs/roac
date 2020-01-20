package main

func printStatement() {
	// Match a 'print' as the first token
	match(TokenPrint, "print")
	// Parse the following expression and
	// generate the assembly code
	tree := binexpr(0)
	reg := generateAST(tree, -1)
	genprintint(reg)
	genfreeregs()
	// Match the following semicolon
	// and stop if we are at EOF
	semi()
}

// Parse the declaration of a variable
func varDeclaration() {
	// Ensure we have an 'int' token followed by an identifier
	// and a semicolon. Text now has the identifier's name.
	// Add it as a known identifier
	match(TokenInt, "int")
	ident()
	AddSymbol(Text)
	genglobsym(Text)
	semi()
}

func assignmentStatement() {
	// struct ASTnode *left, *right, *tree;
	// int id;

	// Ensure we have an identifier
	ident()

	id, exists := GetSymbolIDByString(Text)
	if !exists {
		fatal("undeclared variable %s\n", Text)
	}
	right := NewLeafASTNode(NodeLvIdent, id)

	// Ensure we have an equals sign
	match(TokenEquals, "=")

	// Parse the following expression
	left := binexpr(0)

	// Make an assignment AST tree
	tree := NewASTNode(NodeAssign, left, right, 0)

	// Generate the assembly code for the assignment
	generateAST(tree, -1)
	genfreeregs()

	// Match the following semicolon
	semi()
}

// Parse one or more statements
func statements() {
	for {
		switch CurrentToken.token {
		case TokenPrint:
			printStatement()
			break
		case TokenInt:
			varDeclaration()
			break
		case TokenIdent:
			assignmentStatement()
			break
		case TokenEOF:
			return
		default:
			fatal("syntax error, token %d", CurrentToken.token)
		}
	}
}
