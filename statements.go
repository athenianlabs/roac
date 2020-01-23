package main

func printStatement() *ASTNode {
	// Match a 'print' as the first token
	match(TokenPrint, "print")
	// Parse the following expression and
	// generate the assembly code
	tree := binexpr(0)
	tree = NewUnaryASTNode(NodePrint, tree, 0)
	// Match the following semicolon
	// and stop if we are at EOF
	semi()
	return tree
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

func assignmentStatement() *ASTNode {
	// Ensure we have an identifier
	ident()
	id, exists := GetSymbolIDByString(Text)
	if !exists {
		fatal("undeclared variable %s\n", Text)
	}
	right := NewLeafASTNode(NodeLvIdent, id)
	// Ensure we have an equals sign
	match(TokenAssign, "=")
	// Parse the following expression
	left := binexpr(0)
	// Make an assignment AST tree
	tree := NewASTNode(NodeAssign, left, nil, right, 0)
	// Match the following semicolon
	semi()
	return tree
}

// Parse an IF statement including
// any optional ELSE clause
// and return its AST
func ifStatement() *ASTNode {
	// Ensure we have 'if' '('
	match(TokenIf, "if")
	lparen()
	// Parse the following expression
	// and the ')' following. Ensure
	// the tree's operation is a comparison.
	condAST := binexpr(0)
	if condAST.op < NodeEqual || condAST.op > NodeGreaterThanOrEqual {
		fatal("bad comparison operator\n")
	}
	rparen()
	// Get the AST for the compound statement
	trueAST := compoundStatement()
	// If we have an 'else', skip it
	// and get the AST for the compound statement
	var falseAST *ASTNode
	if CurrentToken.token == TokenElse {
		scan(CurrentToken)
		falseAST = compoundStatement()
	}
	// Build and return the AST for this statement
	return NewASTNode(NodeIf, condAST, trueAST, falseAST, 0)
}

// Parse a compound statement
// and return its AST
func compoundStatement() *ASTNode {
	var tree, left *ASTNode
	// Require a left curly bracket
	lbrace()

	for {
		switch CurrentToken.token {
		case TokenPrint:
			tree = printStatement()
		case TokenInt:
			varDeclaration()
			tree = nil // No AST generated her
		case TokenIdent:
			tree = assignmentStatement()
		case TokenIf:
			tree = ifStatement()
		case TokenRightBrace:
			// When we hit a right curly bracket, skip past it and return the AST
			rbrace()
			return left
		default:
			fatal("syntax error, token %d\n", CurrentToken.token)
		}

		// For each new tree, either save it in left
		// if left is empty, or glue the left and the
		// new tree together
		if tree != nil {
			if left == nil {
				left = tree
			} else {
				left = NewASTNode(NodeGlue, left, nil, tree, 0)
			}
		}
	}
}
