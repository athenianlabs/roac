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

// Parse a single statement
// and return its AST
func singleStatement() *ASTNode {
	switch CurrentToken.token {
	case TokenPrint:
		return printStatement()
	case TokenInt:
		varDeclaration()
		return nil // No AST generated here
	case TokenIdent:
		return assignmentStatement()
	case TokenIf:
		return ifStatement()
	case TokenWhile:
		return whileStatement()
	case TokenFor:
		return forStatement()
	default:
		fatal("Syntax error, token %d\n", CurrentToken.token)
	}
	return nil
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
	return NewASTNode(NodeAssign, left, nil, right, 0)
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

// Parse a WHILE statement
// and return its AST
func whileStatement() *ASTNode {
	// Ensure we have 'while' '('
	match(TokenWhile, "while")
	lparen()
	// Parse the following expression
	// and the ')' following. Ensure
	// the tree's operation is a comparison.
	condAST := binexpr(0)
	if condAST.op < NodeEqual || condAST.op > NodeGreaterThanOrEqual {
		fatal("bad comparison operator")
	}
	rparen()
	// Get the AST for the compound statement
	bodyAST := compoundStatement()
	// Build and return the AST for this statement
	return NewASTNode(NodeWhile, condAST, nil, bodyAST, 0)
}

// Parse a FOR statement
// and return its AST
func forStatement() *ASTNode {
	// Ensure we have 'for' '('
	match(TokenFor, "for")
	lparen()
	// Get the pre_op statement and the ';'
	preopAST := singleStatement()
	semi()
	// Get the condition and the ';'
	condAST := binexpr(0)
	if condAST.op < NodeEqual || condAST.op > NodeGreaterThanOrEqual {
		fatal("Bad comparison operator")
	}
	semi()
	// Get the post_op statement and the ')'
	postopAST := singleStatement()
	rparen()
	// Get the compound statement which is the body
	bodyAST := compoundStatement()
	// For now, all four sub-trees have to be non-NULL.
	// Later on, we'll change the semantics for when some are missing
	// Glue the compound statement and the postop tree
	tree := NewASTNode(NodeGlue, bodyAST, nil, postopAST, 0)
	// Make a WHILE loop with the condition and this new body
	tree = NewASTNode(NodeWhile, condAST, nil, tree, 0)
	// And glue the preop tree to the A_WHILE tree
	return NewASTNode(NodeGlue, preopAST, nil, tree, 0)
}

// Parse the declaration of a simplistic function
func functionDeclaration() *ASTNode {
	// Find the 'void', the identifier, and the '(' ')'.
	// For now, do nothing with them
	match(TokenVoid, "void")
	ident()
	nameslot := AddSymbol(Text)
	lparen()
	rparen()

	// Get the AST tree for the compound statement
	tree := compoundStatement()
	// Return an A_FUNCTION node which has the function's nameslot
	// and the compound statement sub-tree
	return NewUnaryASTNode(NodeFunction, tree, nameslot)
}

// Parse a compound statement
// and return its AST
func compoundStatement() *ASTNode {
	var tree, left *ASTNode
	// Require a left curly bracket
	lbrace()
	for {
		// Parse a single statement
		tree = singleStatement()
		// Some statements must be followed by a semicolon
		if tree != nil && (tree.op == NodePrint || tree.op == NodeAssign) {
			semi()
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
			// When we hit a right curly bracket,
			// skip past it and return the AST
			if CurrentToken.token == TokenRightBrace {
				rbrace()
				return left
			}
		}
	}
}
