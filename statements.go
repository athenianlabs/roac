package main

// Parse one or more statements
func statements() {
	for {
		// Match a 'print' as the first token
		match(TokenPrint, "print")
		// Parse the following expression and
		// generate the assembly code
		tree := binexpr(0)
		reg := generateAST(tree)
		genprintint(reg)
		genfreeregs()
		// Match the following semicolon
		// and stop if we are at EOF
		semi()
		if CurrentToken.token == TokenEOF {
			return
		}
	}
}
