package main

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
		if tree != nil && (tree.op == OpPrint || tree.op == OpAssign) {
			semi()
		}
		// For each new tree, either save it in left
		// if left is empty, or glue the left and the
		// new tree together
		if tree != nil {
			if left == nil {
				left = tree
			} else {
				left = NewASTNode(OpGlue, NodeNone, left, nil, tree, 0)
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

// Parse a single statement
// and return its AST
func singleStatement() *ASTNode {
	switch CurrentToken.token {
	case TokenPrint:
		return printStatement()
	case TokenChar, TokenInt:
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

func printStatement() *ASTNode {
	// Match a 'print' as the first token
	match(TokenPrint, "print")
	// Parse the following expression
	tree := binexpr(0)
	// Ensure the two types are compatible.
	_, rightOp, ok := typeCompatible(NodeInt, tree.t, false)
	if !ok {
		fatal("incompatible types\n")
	}
	// Widen the tree if required.
	if rightOp != nil {
		tree = NewUnaryASTNode(*rightOp, NodeInt, tree, 0)
	}
	// Make an print AST tree
	tree = NewUnaryASTNode(OpPrint, NodeNone, tree, 0)
	// Return the AST
	return tree
}

func assignmentStatement() *ASTNode {
	// Ensure we have an identifier
	ident()
	sym := GetSymbolByString(Text)
	right := NewLeafASTNode(OpLvIdent, sym.t, sym.id)
	// Ensure we have an equals sign
	match(TokenAssign, "=")
	// Parse the following expression
	left := binexpr(0)
	// Ensure the two types are compatible.
	leftOp, _, ok := typeCompatible(left.t, right.t, true)
	if !ok {
		fatal("incompatible types\n")
	}
	// Widen the left if required.
	if leftOp != nil {
		left = NewUnaryASTNode(*leftOp, right.t, left, 0)
	}
	// Make an assignment AST tree
	return NewASTNode(OpAssign, NodeInt, left, nil, right, 0)
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
	if condAST.op < OpEqual || condAST.op > OpGreaterThanOrEqual {
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
	return NewASTNode(OpIf, NodeNone, condAST, trueAST, falseAST, 0)
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
	if condAST.op < OpEqual || condAST.op > OpGreaterThanOrEqual {
		fatal("bad comparison operator")
	}
	rparen()
	// Get the AST for the compound statement
	bodyAST := compoundStatement()
	// Build and return the AST for this statement
	return NewASTNode(OpWhile, NodeNone, condAST, nil, bodyAST, 0)
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
	if condAST.op < OpEqual || condAST.op > OpGreaterThanOrEqual {
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
	tree := NewASTNode(OpGlue, NodeNone, bodyAST, nil, postopAST, 0)
	// Make a WHILE loop with the condition and this new body
	tree = NewASTNode(OpWhile, NodeNone, condAST, nil, tree, 0)
	// And glue the preop tree to the A_WHILE tree
	return NewASTNode(OpGlue, NodeNone, preopAST, nil, tree, 0)
}
